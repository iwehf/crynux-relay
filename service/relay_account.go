package service

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/utils"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ErrInsufficientRelayAccount = errors.New("insufficient relay account balance")

type relayAccountCacheType struct {
	accounts map[string]*big.Int
	mu       sync.RWMutex
}

var relayAccountCache = &relayAccountCacheType{
	accounts: make(map[string]*big.Int),
}

func InitRelayAccountCache(ctx context.Context, db *gorm.DB) error {
	for {
		events, err := getPendingRelayAccountEvents(ctx, db, 50)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			break
		}
		if err := processPendingRelayAccountEvents(ctx, db, events); err != nil {
			return err
		}
	}
	return nil
}

func StartRelayAccountSync(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := syncRelayAccountsToDB(ctx, db); err != nil {
				log.Errorf("Failed to sync relay accounts: %v", err)
			}
		}
	}
}

func getPendingRelayAccountEvents(ctx context.Context, db *gorm.DB, limit int) ([]models.RelayAccountEvent, error) {
	var events []models.RelayAccountEvent
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := db.WithContext(dbCtx).
		Where("status = ?", models.RelayAccountEventStatusPending).
		Order("id").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func validatePendingRelayAccountEvents(ctx context.Context, db *gorm.DB, events []models.RelayAccountEvent) ([]models.RelayAccountEvent, []models.RelayAccountEvent, error) {
	invalidEvents := make([]models.RelayAccountEvent, 0)
	candidateEvents := make([]models.RelayAccountEvent, 0)
	taskIDSet := make(map[string]struct{})
	withdrawIDSet := make(map[uint]struct{})
	txHashes := make([]string, 0)
	networks := make([]string, 0)

	for _, event := range events {
		reasons := strings.Split(event.Reason, "-")
		eventTypeStr := fmt.Sprintf("%d", event.Type)
		if len(reasons) < 2 || reasons[0] != eventTypeStr {
			invalidEvents = append(invalidEvents, event)
			continue
		}
		switch event.Type {
		case models.RelayAccountEventTypeDeposit:
			if len(reasons) != 3 {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			txHashes = append(txHashes, reasons[1])
			networks = append(networks, reasons[2])
		case models.RelayAccountEventTypeTaskPayment, models.RelayAccountEventTypeTaskIncome, models.RelayAccountEventTypeDaoTaskShare, models.RelayAccountEventTypeTaskRefund:
			if len(reasons) != 2 {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			taskIDSet[reasons[1]] = struct{}{}
		case models.RelayAccountEventTypeWithdraw, models.RelayAccountEventTypeWithdrawRefund, models.RelayAccountEventTypeWithdrawFeeIncome:
			if len(reasons) != 2 {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			withdrawID, err := strconv.ParseUint(reasons[1], 10, 64)
			if err != nil {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			withdrawIDSet[uint(withdrawID)] = struct{}{}
		default:
			invalidEvents = append(invalidEvents, event)
			continue
		}
		candidateEvents = append(candidateEvents, event)
	}

	taskIDs := make([]string, 0, len(taskIDSet))
	for taskID := range taskIDSet {
		taskIDs = append(taskIDs, taskID)
	}
	withdrawIDs := make([]uint, 0, len(withdrawIDSet))
	for withdrawID := range withdrawIDSet {
		withdrawIDs = append(withdrawIDs, withdrawID)
	}

	var tasks []models.InferenceTask
	if len(taskIDs) > 0 {
		dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := db.WithContext(dbCtx).
			Where("task_id_commitment IN (?)", taskIDs).
			Find(&tasks).Error; err != nil {
			return nil, nil, err
		}
	}
	taskMap := make(map[string]models.InferenceTask)
	for _, task := range tasks {
		taskMap[task.TaskIDCommitment] = task
	}

	var withdrawRecords []models.WithdrawRecord
	if len(withdrawIDs) > 0 {
		dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := db.WithContext(dbCtx).
			Where("id IN (?)", withdrawIDs).
			Find(&withdrawRecords).Error; err != nil {
			return nil, nil, err
		}
	}
	withdrawMap := make(map[uint]models.WithdrawRecord)
	for _, record := range withdrawRecords {
		withdrawMap[record.ID] = record
	}

	validTxHashes := make(map[string]struct{})
	for i, txHash := range txHashes {
		network := networks[i]
		client, err := blockchain.GetBlockchainClient(network)
		if err != nil {
			return nil, nil, err
		}
		receipt, err := client.RpcClient.TransactionReceipt(ctx, common.HexToHash(txHash))
		if errors.Is(err, ethereum.NotFound) {
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			continue
		}
		validTxHashes[txHash] = struct{}{}
	}

	validEvents := make([]models.RelayAccountEvent, 0, len(candidateEvents))
	for _, event := range candidateEvents {
		reasons := strings.Split(event.Reason, "-")
		switch event.Type {
		case models.RelayAccountEventTypeDeposit:
			if _, ok := validTxHashes[reasons[1]]; ok {
				validEvents = append(validEvents, event)
			} else {
				invalidEvents = append(invalidEvents, event)
			}
		case models.RelayAccountEventTypeTaskPayment:
			task, ok := taskMap[reasons[1]]
			if !ok || event.Amount.Int.Cmp(&task.TaskFee.Int) != 0 {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			validEvents = append(validEvents, event)
		case models.RelayAccountEventTypeTaskIncome, models.RelayAccountEventTypeDaoTaskShare:
			task, ok := taskMap[reasons[1]]
			if !ok {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			if task.Status != models.TaskEndSuccess && task.Status != models.TaskEndGroupSuccess && task.Status != models.TaskEndGroupRefund {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			validEvents = append(validEvents, event)
		case models.RelayAccountEventTypeTaskRefund:
			task, ok := taskMap[reasons[1]]
			if !ok || event.Amount.Int.Cmp(&task.TaskFee.Int) != 0 {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			if task.Status != models.TaskEndGroupRefund && task.Status != models.TaskEndAborted {
				invalidEvents = append(invalidEvents, event)
				continue
			}
			validEvents = append(validEvents, event)
		case models.RelayAccountEventTypeWithdraw, models.RelayAccountEventTypeWithdrawRefund, models.RelayAccountEventTypeWithdrawFeeIncome:
			withdrawID, _ := strconv.ParseUint(reasons[1], 10, 64)
			if _, ok := withdrawMap[uint(withdrawID)]; ok {
				validEvents = append(validEvents, event)
			} else {
				invalidEvents = append(invalidEvents, event)
			}
		}
	}

	return validEvents, invalidEvents, nil
}

func mergeRelayAccountEvents(events []models.RelayAccountEvent) map[string]*big.Int {
	merged := make(map[string]*big.Int)
	for _, event := range events {
		if _, exists := merged[event.Address]; !exists {
			merged[event.Address] = big.NewInt(0)
		}
		switch event.Type {
		case models.RelayAccountEventTypeTaskPayment, models.RelayAccountEventTypeWithdraw:
			merged[event.Address].Sub(merged[event.Address], &event.Amount.Int)
		default:
			merged[event.Address].Add(merged[event.Address], &event.Amount.Int)
		}
	}
	return merged
}

func processPendingRelayAccountEvents(ctx context.Context, db *gorm.DB, events []models.RelayAccountEvent) error {
	validEvents, invalidEvents, err := validatePendingRelayAccountEvents(ctx, db, events)
	if err != nil {
		return err
	}

	if len(invalidEvents) > 0 {
		invalidIDs := make([]uint, 0, len(invalidEvents))
		invalidWithdrawEventIDs := make([]uint, 0)
		for _, event := range invalidEvents {
			invalidIDs = append(invalidIDs, event.ID)
			if event.Type == models.RelayAccountEventTypeWithdraw {
				invalidWithdrawEventIDs = append(invalidWithdrawEventIDs, event.ID)
			}
		}
		if err := db.Model(&models.RelayAccountEvent{}).
			Where("id IN (?)", invalidIDs).
			Update("status", models.RelayAccountEventStatusInvalid).Error; err != nil {
			return err
		}
		if len(invalidWithdrawEventIDs) > 0 {
			if err := db.Model(&models.WithdrawRecord{}).
				Where("relay_account_event_id IN (?)", invalidWithdrawEventIDs).
				Where("local_status = ?", models.WithdrawLocalStatusPending).
				Update("local_status", models.WithdrawLocalStatusInvalid).Error; err != nil {
				return err
			}
		}
	}

	mergedAccounts := mergeRelayAccountEvents(validEvents)
	addresses := make([]string, 0, len(mergedAccounts))
	for address := range mergedAccounts {
		addresses = append(addresses, address)
	}
	eventIDs := make([]uint, 0, len(validEvents))
	for _, event := range validEvents {
		eventIDs = append(eventIDs, event.ID)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	appConfig := config.GetConfig()
	eventMACs := make(map[uint]string, len(validEvents))
	for _, event := range validEvents {
		eventMACs[event.ID] = utils.GenerateMAC([]byte(event.Reason), appConfig.MAC.SecretKey)
	}

	return db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var existedAccounts []models.RelayAccount
		if err := tx.Model(&models.RelayAccount{}).Where("address IN (?)", addresses).Find(&existedAccounts).Error; err != nil {
			return err
		}

		existedAddresses := make([]string, 0, len(existedAccounts))
		existedMap := make(map[string]*models.RelayAccount)
		for i, account := range existedAccounts {
			existedMap[account.Address] = &existedAccounts[i]
			existedAddresses = append(existedAddresses, account.Address)
		}

		newAccounts := make([]models.RelayAccount, 0)
		for address, amount := range mergedAccounts {
			if account, ok := existedMap[address]; ok {
				account.Balance.Int.Add(&account.Balance.Int, amount)
			} else {
				newAccounts = append(newAccounts, models.RelayAccount{
					Address: address,
					Balance: models.BigInt{Int: *new(big.Int).Set(amount)},
				})
			}
		}

		if len(newAccounts) > 0 {
			if err := tx.CreateInBatches(&newAccounts, 100).Error; err != nil {
				return err
			}
		}

		if len(existedMap) > 0 {
			var cases string
			for _, account := range existedMap {
				cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", account.Address, account.Balance.String())
			}
			if err := tx.Model(&models.RelayAccount{}).
				Where("address IN (?)", existedAddresses).
				Update("balance", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
				return err
			}
		}

		var macCases string
		for eventID, mac := range eventMACs {
			macCases += fmt.Sprintf(" WHEN id = %d THEN '%s'", eventID, mac)
		}
		if err := tx.Model(&models.RelayAccountEvent{}).
			Where("id IN (?)", eventIDs).
			Where("status = ?", models.RelayAccountEventStatusPending).
			Updates(map[string]interface{}{
				"mac":    gorm.Expr("CASE" + macCases + " END"),
				"status": models.RelayAccountEventStatusProcessed,
			}).Error; err != nil {
			return err
		}

		withdrawEventIDs := make([]uint, 0)
		for _, event := range validEvents {
			if event.Type == models.RelayAccountEventTypeWithdraw {
				withdrawEventIDs = append(withdrawEventIDs, event.ID)
			}
		}
		if len(withdrawEventIDs) > 0 {
			var records []models.WithdrawRecord
			if err := tx.Model(&models.WithdrawRecord{}).
				Where("relay_account_event_id IN (?)", withdrawEventIDs).
				Where("local_status = ?", models.WithdrawLocalStatusPending).
				Find(&records).Error; err != nil {
				return err
			}

			if len(records) > 0 {
				recordIDs := make([]uint, 0, len(records))
				var recordMACCases string
				for _, record := range records {
					recordIDs = append(recordIDs, record.ID)
					recordMAC := utils.GenerateMAC([]byte(record.MACString()), appConfig.MAC.SecretKey)
					recordMACCases += fmt.Sprintf(" WHEN id = %d THEN '%s'", record.ID, recordMAC)
				}

				if err := tx.Model(&models.WithdrawRecord{}).
					Where("id IN (?)", recordIDs).
					Where("local_status = ?", models.WithdrawLocalStatusPending).
					Updates(map[string]interface{}{
						"mac":          gorm.Expr("CASE" + recordMACCases + " END"),
						"local_status": models.WithdrawLocalStatusProcessed,
					}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func syncRelayAccountsToDB(ctx context.Context, db *gorm.DB) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for {
				events, err := getPendingRelayAccountEvents(ctx, db, 50)
				if err != nil {
					log.Errorf("Failed to load pending relay account events: %v", err)
					break
				}
				if len(events) == 0 {
					break
				}
				if err := processPendingRelayAccountEvents(ctx, db, events); err != nil {
					log.Errorf("Failed to process pending relay account events: %v", err)
					break
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
				continue
			}
		}
	}
}

func getRelayAccountFromCache(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	relayAccountCache.mu.RLock()
	balance, exists := relayAccountCache.accounts[address]
	relayAccountCache.mu.RUnlock()
	if exists {
		return balance, nil
	}

	relayAccountCache.mu.Lock()
	defer relayAccountCache.mu.Unlock()
	if balance, exists := relayAccountCache.accounts[address]; exists {
		return balance, nil
	}

	var dbAccount models.RelayAccount
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).
		Where("address = ?", address).
		Attrs(models.RelayAccount{Balance: models.BigInt{Int: *big.NewInt(0)}}).
		FirstOrInit(&dbAccount).Error; err != nil {
		return nil, err
	}
	relayAccountCache.accounts[address] = &dbAccount.Balance.Int
	return relayAccountCache.accounts[address], nil
}

func createRelayAccountEvent(ctx context.Context, db *gorm.DB, eventType models.RelayAccountEventType, reason, address string, amount *big.Int) error {
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return db.WithContext(dbCtx).Create(&models.RelayAccountEvent{
		Reason:    reason,
		Address:   address,
		Amount:    models.BigInt{Int: *new(big.Int).Set(amount)},
		CreatedAt: time.Now(),
		Status:    models.RelayAccountEventStatusPending,
		Type:      eventType,
	}).Error
}

func depositRelayAccount(ctx context.Context, db *gorm.DB, txHash, address string, amount *big.Int, network string) (func() error, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		reason := fmt.Sprintf("%d-%s-%s", models.RelayAccountEventTypeDeposit, txHash, network)
		if err := createRelayAccountEvent(ctx, tx, models.RelayAccountEventTypeDeposit, reason, address, amount); err != nil {
			return err
		}
		record := &models.DepositRecord{
			Address: address,
			Amount:  models.BigInt{Int: *new(big.Int).Set(amount)},
			Network: network,
			TxHash:  txHash,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	balance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		balance.Add(balance, amountCopy)
		return nil
	}, nil
}

func chargeTaskFromRelayAccount(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int) (func() error, error) {
	balance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	if balance.Cmp(amount) < 0 {
		return nil, ErrInsufficientRelayAccount
	}
	reason := fmt.Sprintf("%d-%s", models.RelayAccountEventTypeTaskPayment, taskIDCommitment)
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeTaskPayment, reason, address, amount); err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		if balance.Cmp(amountCopy) < 0 {
			return ErrInsufficientRelayAccount
		}
		balance.Sub(balance, amountCopy)
		return nil
	}, nil
}

func refundTaskPaymentToRelayAccount(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int) (func() error, error) {
	reason := fmt.Sprintf("%d-%s", models.RelayAccountEventTypeTaskRefund, taskIDCommitment)
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeTaskRefund, reason, address, amount); err != nil {
		return nil, err
	}
	balance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		balance.Add(balance, amountCopy)
		return nil
	}, nil
}

func sendTaskIncome(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int, taskType models.TaskType) (func() error, error) {
	appConfig := config.GetConfig()
	daoTaskShare := big.NewInt(0).Mul(amount, big.NewInt(0).SetUint64(appConfig.Dao.TaskFeeSharePercent))
	daoTaskShare.Div(daoTaskShare, big.NewInt(100))
	nodeIncome := big.NewInt(0).Sub(amount, daoTaskShare)

	nodeReason := fmt.Sprintf("%d-%s", models.RelayAccountEventTypeTaskIncome, taskIDCommitment)
	daoReason := fmt.Sprintf("%d-%s", models.RelayAccountEventTypeDaoTaskShare, taskIDCommitment)
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeTaskIncome, nodeReason, address, nodeIncome); err != nil {
		return nil, err
	}
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeDaoTaskShare, daoReason, appConfig.Dao.TaskFeeShareAddress, daoTaskShare); err != nil {
		return nil, err
	}

	incentive, _ := utils.WeiToEther(nodeIncome).Float64()
	if err := addNodeIncentive(ctx, db, address, incentive, taskType); err != nil {
		return nil, err
	}

	nodeBalance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	daoBalance, err := getRelayAccountFromCache(ctx, db, appConfig.Dao.TaskFeeShareAddress)
	if err != nil {
		return nil, err
	}

	nodeIncomeCopy := new(big.Int).Set(nodeIncome)
	daoTaskShareCopy := new(big.Int).Set(daoTaskShare)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		nodeBalance.Add(nodeBalance, nodeIncomeCopy)
		daoBalance.Add(daoBalance, daoTaskShareCopy)
		return nil
	}, nil
}

func chargeWithdrawFromRelayAccount(ctx context.Context, db *gorm.DB, withdrawID uint, address string, amount *big.Int) (uint, func() error, error) {
	balance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return 0, nil, err
	}
	if balance.Cmp(amount) < 0 {
		return 0, nil, ErrInsufficientRelayAccount
	}
	reason := fmt.Sprintf("%d-%d", models.RelayAccountEventTypeWithdraw, withdrawID)
	event := &models.RelayAccountEvent{
		Reason:    reason,
		Address:   address,
		Amount:    models.BigInt{Int: *new(big.Int).Set(amount)},
		CreatedAt: time.Now(),
		Status:    models.RelayAccountEventStatusPending,
		Type:      models.RelayAccountEventTypeWithdraw,
	}
	if err := db.Create(event).Error; err != nil {
		return 0, nil, err
	}

	amountCopy := new(big.Int).Set(amount)
	return event.ID, func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		if balance.Cmp(amountCopy) < 0 {
			return ErrInsufficientRelayAccount
		}
		balance.Sub(balance, amountCopy)
		return nil
	}, nil
}

func fulfillWithdrawFeeIncome(ctx context.Context, db *gorm.DB, withdrawID uint, feeAddress string, fee *big.Int) (func() error, error) {
	reason := fmt.Sprintf("%d-%d", models.RelayAccountEventTypeWithdrawFeeIncome, withdrawID)
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeWithdrawFeeIncome, reason, feeAddress, fee); err != nil {
		return nil, err
	}
	balance, err := getRelayAccountFromCache(ctx, db, feeAddress)
	if err != nil {
		return nil, err
	}
	feeCopy := new(big.Int).Set(fee)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		balance.Add(balance, feeCopy)
		return nil
	}, nil
}

func rejectWithdrawToRelayAccount(ctx context.Context, db *gorm.DB, withdrawID uint, address string, amount *big.Int) (func() error, error) {
	reason := fmt.Sprintf("%d-%d", models.RelayAccountEventTypeWithdrawRefund, withdrawID)
	if err := createRelayAccountEvent(ctx, db, models.RelayAccountEventTypeWithdrawRefund, reason, address, amount); err != nil {
		return nil, err
	}
	balance, err := getRelayAccountFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	return func() error {
		relayAccountCache.mu.Lock()
		defer relayAccountCache.mu.Unlock()
		balance.Add(balance, amountCopy)
		return nil
	}, nil
}

func GetRelayAccountBalance(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	return getRelayAccountFromCache(ctx, db, address)
}
