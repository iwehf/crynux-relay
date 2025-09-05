package service

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/config"
	"crynux_relay/models"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type TaskQuotaCache struct {
	taskQuotas map[string]*big.Int
	mu         sync.RWMutex
}

var taskQuotaCache = &TaskQuotaCache{
	taskQuotas: make(map[string]*big.Int),
}

func InitTaskQuotaCache(ctx context.Context, db *gorm.DB) error {
	for {
		events, err := getPendingTaskQuotaEvents(ctx, db, 50)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			break
		}

		if err := processPendingTaskQuotaEvents(ctx, db, events); err != nil {
			return err
		}
	}
	return nil

}

func getPendingTaskQuotaEvents(ctx context.Context, db *gorm.DB, limit int) ([]models.TaskQuotaEvent, error) {
	var events []models.TaskQuotaEvent
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := db.WithContext(dbCtx).Where("status = ?", models.TaskQuotaEventStatusPending).Order("id").Limit(limit).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func validatePendingTaskQuotaEvents(ctx context.Context, db *gorm.DB, events []models.TaskQuotaEvent) ([]models.TaskQuotaEvent, []models.TaskQuotaEvent, error) {
	invalidEvents := make([]models.TaskQuotaEvent, 0)
	candidateEvents := make([]models.TaskQuotaEvent, 0)

	var taskIDCommitments []string
	taskIDCommitmentMap := make(map[string]struct{})
	var txHashes []string
	var networks []string
	for _, event := range events {
		reasons := strings.Split(event.Reason, "-")
		if len(reasons) != 2 && len(reasons) != 3 {
			invalidEvents = append(invalidEvents, event)
			continue
		}
		eventType := reasons[0]
		if eventType != fmt.Sprintf("%d", event.TaskQuotaType) {
			invalidEvents = append(invalidEvents, event)
			continue
		}
		if event.TaskQuotaType == models.TaskQuotaTypeSpent || event.TaskQuotaType == models.TaskQuotaTypeRefunded {
			taskIDCommitment := reasons[1]
			if _, ok := taskIDCommitmentMap[taskIDCommitment]; !ok {
				taskIDCommitmentMap[taskIDCommitment] = struct{}{}
				taskIDCommitments = append(taskIDCommitments, taskIDCommitment)
			}
		} else {
			txHash := reasons[1]
			network := reasons[2]
			txHashes = append(txHashes, txHash)
			networks = append(networks, network)
		}
		candidateEvents = append(candidateEvents, event)
	}

	var tasks []models.InferenceTask
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).Where("task_id_commitment IN (?)", taskIDCommitments).Find(&tasks).Error; err != nil {
		return nil, nil, err
	}
	taskMap := make(map[string]models.InferenceTask)
	for _, task := range tasks {
		taskMap[task.TaskIDCommitment] = task
	}

	validTxHashesMap := make(map[string]struct{})
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
		validTxHashesMap[txHash] = struct{}{}
	}

	validEvents := make([]models.TaskQuotaEvent, 0)
	for _, event := range candidateEvents {
		reasons := strings.Split(event.Reason, "-")
		if event.TaskQuotaType == models.TaskQuotaTypeSpent || event.TaskQuotaType == models.TaskQuotaTypeRefunded {
			taskIDCommitment := reasons[1]
			if task, exists := taskMap[taskIDCommitment]; exists {
				if event.TaskQuotaType == models.TaskQuotaTypeSpent && event.Quota.Int.Cmp(&task.TaskFee.Int) == 0 {
					validEvents = append(validEvents, event)
				} else if event.TaskQuotaType == models.TaskQuotaTypeRefunded && event.Quota.Int.Cmp(&task.TaskFee.Int) == 0 && (task.Status == models.TaskEndGroupRefund || task.Status == models.TaskEndAborted) {
					validEvents = append(validEvents, event)
				} else {
					invalidEvents = append(invalidEvents, event)
				}
			} else {
				invalidEvents = append(invalidEvents, event)
			}
		} else {
			txHash := reasons[1]
			if _, exists := validTxHashesMap[txHash]; exists {
				validEvents = append(validEvents, event)
			} else {
				invalidEvents = append(invalidEvents, event)
			}
		}
	}

	return validEvents, invalidEvents, nil
}

func getTaskQuotaFromCache(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	taskQuotaCache.mu.RLock()
	taskQuota, exists := taskQuotaCache.taskQuotas[address]
	taskQuotaCache.mu.RUnlock()

	if exists {
		return taskQuota, nil
	}

	taskQuotaCache.mu.Lock()
	defer taskQuotaCache.mu.Unlock()

	if taskQuota, exists := taskQuotaCache.taskQuotas[address]; exists {
		return taskQuota, nil
	}

	var dbTaskQuota models.TaskQuota
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.WithContext(dbCtx).Where("address = ?", address).Attrs(models.TaskQuota{Quota: models.BigInt{Int: *big.NewInt(0)}}).FirstOrInit(&dbTaskQuota).Error; err != nil {
		return nil, err
	}

	taskQuotaCache.taskQuotas[address] = &dbTaskQuota.Quota.Int

	return taskQuotaCache.taskQuotas[address], nil
}

func StartTaskQuotaSync(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := syncTaskQuotasToDB(ctx, db); err != nil {
				log.Errorf("Failed to sync task quotas: %v", err)
			}
		}
	}
}

func mergeTaskQuotaEvents(events []models.TaskQuotaEvent) map[string]*big.Int {
	mergedTaskQuotas := make(map[string]*big.Int)
	for _, event := range events {
		if _, exists := mergedTaskQuotas[event.Address]; !exists {
			mergedTaskQuotas[event.Address] = big.NewInt(0)
		}
		if event.TaskQuotaType == models.TaskQuotaTypeBought || event.TaskQuotaType == models.TaskQuotaTypeRefunded {
			mergedTaskQuotas[event.Address].Add(mergedTaskQuotas[event.Address], &event.Quota.Int)
		} else {
			mergedTaskQuotas[event.Address].Sub(mergedTaskQuotas[event.Address], &event.Quota.Int)
		}
	}
	return mergedTaskQuotas
}

func processPendingTaskQuotaEvents(ctx context.Context, db *gorm.DB, events []models.TaskQuotaEvent) error {
	validEvents, invalidEvents, err := validatePendingTaskQuotaEvents(ctx, db, events)
	if err != nil {
		return err
	}

	if len(invalidEvents) > 0 {
		var invalidEventIDs []uint
		for _, event := range invalidEvents {
			invalidEventIDs = append(invalidEventIDs, event.ID)
		}
		if err := db.Model(&models.TaskQuotaEvent{}).Where("id IN (?)", invalidEventIDs).Update("status", models.TaskQuotaEventStatusInvalid).Error; err != nil {
			return err
		}
	}

	mergedTaskQuotas := mergeTaskQuotaEvents(validEvents)

	var eventIDs []uint
	for _, event := range validEvents {
		eventIDs = append(eventIDs, event.ID)
	}

	var addresses []string
	for address := range mergedTaskQuotas {
		addresses = append(addresses, address)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var existedTaskQuotas []models.TaskQuota
		if err := tx.Model(&models.TaskQuota{}).Where("address IN (?)", addresses).Find(&existedTaskQuotas).Error; err != nil {
			return err
		}

		existedAddresses := make([]string, len(existedTaskQuotas))
		existedTaskQuotasMap := make(map[string]*models.TaskQuota)
		for i, taskQuota := range existedTaskQuotas {
			existedTaskQuotasMap[taskQuota.Address] = &existedTaskQuotas[i]
			existedAddresses[i] = taskQuota.Address
		}

		var newTaskQuotas []models.TaskQuota
		for address, amount := range mergedTaskQuotas {
			if taskQuota, exists := existedTaskQuotasMap[address]; !exists {
				newTaskQuotas = append(newTaskQuotas, models.TaskQuota{Address: address, Quota: models.BigInt{Int: *amount}})
			} else {
				taskQuota.Quota.Int.Add(&taskQuota.Quota.Int, amount)
			}
		}

		if len(newTaskQuotas) > 0 {
			if err := tx.CreateInBatches(&newTaskQuotas, 100).Error; err != nil {
				return err
			}
		}

		if len(existedTaskQuotasMap) > 0 {
			var cases string
			for _, taskQuota := range existedTaskQuotasMap {
				cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", taskQuota.Address, taskQuota.Quota.String())
			}
			if err := tx.Model(&models.TaskQuota{}).Where("address IN (?)", existedAddresses).
				Update("quota", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&models.TaskQuotaEvent{}).Where("id IN (?)", eventIDs).Where("status = ?", models.TaskQuotaEventStatusPending).Update("status", models.TaskQuotaEventStatusProcessed).Error; err != nil {
			return err
		}

		return nil
	})
}

func syncTaskQuotasToDB(ctx context.Context, db *gorm.DB) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := func() error {
				for {
					events, err := getPendingTaskQuotaEvents(ctx, db, 50)
					if err != nil {
						return err
					}

					if len(events) == 0 {
						break
					}

					if err := processPendingTaskQuotaEvents(ctx, db, events); err != nil {
						return err
					}
				}
				return nil
			}()

			if err != nil {
				log.Errorf("Failed to sync task quotas: %v", err)
			}

			// Wait for 2 seconds before next iteration
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
				continue
			}
		}
	}
}

func BuyTaskQuota(ctx context.Context, db *gorm.DB, txHash, address string, amount *big.Int, network string) (func() error, error) {

	event := &models.TaskQuotaEvent{
		Reason:        fmt.Sprintf("%d-%s-%s", models.TaskQuotaTypeBought, txHash, network),
		Address:       address,
		Quota:         models.BigInt{Int: *new(big.Int).Set(amount)},
		CreatedAt:     time.Now(),
		Status:        models.TaskQuotaEventStatusPending,
		TaskQuotaType: models.TaskQuotaTypeBought,
	}

	if err := db.Create(event).Error; err != nil {
		return nil, err
	}

	taskQuota, err := getTaskQuotaFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	commitFunc := func() error {
		taskQuotaCache.mu.Lock()
		defer taskQuotaCache.mu.Unlock()
		taskQuota.Add(taskQuota, amountCopy)
		return nil
	}

	return commitFunc, nil
}

func SpendTaskQuota(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int) (func() error, error) {
	taskQuota, err := getTaskQuotaFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}

	if taskQuota.Cmp(amount) < 0 {
		return nil, fmt.Errorf("insufficient task quota when spend")
	}

	event := &models.TaskQuotaEvent{
		Reason:        fmt.Sprintf("%d-%s", models.TaskQuotaTypeSpent, taskIDCommitment),
		Address:       address,
		Quota:         models.BigInt{Int: *new(big.Int).Set(amount)},
		CreatedAt:     time.Now(),
		Status:        models.TaskQuotaEventStatusPending,
		TaskQuotaType: models.TaskQuotaTypeSpent,
	}

	if err := db.Create(event).Error; err != nil {
		return nil, err
	}

	amountCopy := new(big.Int).Set(amount)
	commitFunc := func() error {
		taskQuotaCache.mu.Lock()
		defer taskQuotaCache.mu.Unlock()
		if taskQuota.Cmp(amountCopy) < 0 {
			return fmt.Errorf("insufficient task quota when spend")
		}
		taskQuota.Sub(taskQuota, amountCopy)
		return nil
	}

	return commitFunc, nil

}

func RefundTaskQuota(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int) (func() error, error) {
	event := &models.TaskQuotaEvent{
		Reason:        fmt.Sprintf("%d-%s", models.TaskQuotaTypeRefunded, taskIDCommitment),
		Address:       address,
		Quota:         models.BigInt{Int: *new(big.Int).Set(amount)},
		CreatedAt:     time.Now(),
		Status:        models.TaskQuotaEventStatusPending,
		TaskQuotaType: models.TaskQuotaTypeRefunded,
	}

	if err := db.Create(event).Error; err != nil {
		return nil, err
	}

	taskQuota, err := getTaskQuotaFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	amountCopy := new(big.Int).Set(amount)
	commitFunc := func() error {
		taskQuotaCache.mu.Lock()
		defer taskQuotaCache.mu.Unlock()
		taskQuota.Add(taskQuota, amountCopy)
		return nil
	}

	return commitFunc, nil
}

func GetTaskQuota(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	return getTaskQuotaFromCache(ctx, db, address)
}

// StartNativeTokenListener starts the native token transfer listener
func StartNativeTokenListener(ctx context.Context) {
	appConfig := config.GetConfig()

	// Start the listener goroutine
	for network := range appConfig.Blockchains {
		go func() {
			if err := runNativeTokenListener(ctx, config.GetDB(), network); err != nil {
				log.Errorf("Native token listener failed: %v", err)
			}
		}()
	}

	log.Info("Native token listener started")
}

// runNativeTokenListener runs the native token transfer listener
func runNativeTokenListener(ctx context.Context, db *gorm.DB, network string) error {
	ticker := time.NewTicker(5 * time.Second) // Check for new blocks every 5 seconds
	defer ticker.Stop()

	client, err := blockchain.GetBlockchainClient(network)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := processNewBlocks(ctx, db, client); err != nil {
				log.Errorf("Failed to process new blocks: %v", err)
			}
		}
	}
}

// processNewBlocks processes new blocks
func processNewBlocks(ctx context.Context, db *gorm.DB, client *blockchain.BlockchainClient) error {
	// Get current block height
	latestBlock, err := client.RpcClient.BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	// Get listener status
	listener, err := models.GetNativeTokenListener(ctx, db, client.Network)
	if err != nil {
		return err
	}

	// If already at the latest block, skip
	if listener.LastBlockNum >= latestBlock.NumberU64() {
		return nil
	}

	// Process new blocks
	startBlock := listener.LastBlockNum + 1
	endBlock := latestBlock.NumberU64()

	// Limit the number of blocks processed each time to avoid long processing time
	if endBlock-startBlock > 10 {
		endBlock = startBlock + 10
	}

	log.Infof("Processing blocks from %d to %d", startBlock, endBlock)

	for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
		if err := processBlock(ctx, db, client, blockNum); err != nil {
			log.Errorf("Failed to process block %d: %v", blockNum, err)
			continue
		}
	}

	// Update listener status
	if err := db.Model(&listener).Updates(map[string]interface{}{
		"last_block_num":   endBlock,
		"last_update_time": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to update block listener: %w", err)
	}

	return nil
}

// processBlock processes a single block
func processBlock(ctx context.Context, db *gorm.DB, client *blockchain.BlockchainClient, blockNum uint64) error {
	block, err := client.RpcClient.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
	if err != nil {
		return fmt.Errorf("failed to get block %d: %w", blockNum, err)
	}

	// Check transactions in the block
	for _, tx := range block.Transactions() {
		if err := processTransaction(ctx, db, tx, client); err != nil {
			log.Errorf("Failed to process transaction %s: %v", tx.Hash().Hex(), err)
			continue
		}
	}

	return nil
}

// processTransaction processes a single transaction
func processTransaction(ctx context.Context, db *gorm.DB, tx *types.Transaction, client *blockchain.BlockchainClient) error {
	// Only process native token transfers (to field is not empty and data field is empty)
	if tx.To() == nil || len(tx.Data()) > 0 {
		return nil
	}

	// Check if transfer is to the target address
	if !strings.EqualFold(tx.To().Hex(), client.Address) {
		return nil
	}

	// Check if transaction is successful
	receipt, err := client.RpcClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil
	}

	// Check if already processed
	event, err := models.GetTaskQuotaBoughtEvent(ctx, db, tx.Hash().Hex())
	if err != nil {
		return err
	}
	if event != nil {
		return nil
	}

	// Get sender address (need to recover from signature)
	from, err := types.Sender(types.NewEIP155Signer(client.ChainID), tx)
	if err != nil {
		return fmt.Errorf("failed to get sender address: %w", err)
	}

	// Call BuyTaskQuota to add quota for the sender
	commitFunc, err := BuyTaskQuota(ctx, db, tx.Hash().Hex(), from.Hex(), tx.Value(), client.Network)
	if err != nil {
		log.Errorf("Failed to buy task quota for %s: %v", from.Hex(), err)
		return err
	}

	// Execute quota update
	if err := commitFunc(); err != nil {
		log.Errorf("Failed to buy task quota for %s: %v", from.Hex(), err)
		return err
	}

	log.Infof("Processed native token transfer: %s -> %s, amount: %s", from.Hex(), tx.To().Hex(), tx.Value().String())

	return nil
}
