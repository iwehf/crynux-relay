package service

import (
	"context"
	"crynux_relay/config"
	"crynux_relay/models"
	"fmt"
	"math/big"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type TaskFeeCache struct {
	taskFees map[string]*big.Int
	mu       sync.RWMutex
}

var taskFeeCache = &TaskFeeCache{
	taskFees: make(map[string]*big.Int),
}

func InitTaskFeeCache(ctx context.Context, db *gorm.DB) error {
	for {
		events, err := getPendingTaskFeeEvents(ctx, db, 50)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			break
		}

		if err := processPendingTaskFeeEvents(ctx, db, events); err != nil {
			return err
		}
	}
	for {
		events, err := getPendingWithdrawEvents(ctx, db, 50)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			break
		}

		if err := processPendingWithdrawEvents(ctx, db, events); err != nil {
			return err
		}
	}
	return nil
}

func getPendingTaskFeeEvents(ctx context.Context, db *gorm.DB, limit int) ([]models.TaskFeeEvent, error) {
	var events []models.TaskFeeEvent
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := db.WithContext(dbCtx).Where("status = ?", models.TaskFeeEventStatusPending).Order("id").Limit(limit).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func getPendingWithdrawEvents(ctx context.Context, db *gorm.DB, limit int) ([]models.WithdrawRecord, error) {
	var withdrawRecords []models.WithdrawRecord
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := db.WithContext(dbCtx).Where("local_status = ?", models.WithdrawLocalStatusPending).Order("id").Limit(limit).Find(&withdrawRecords).Error
	if err != nil {
		return nil, err
	}

	return withdrawRecords, nil
}

func validatePendingTaskFeeEvents(ctx context.Context, db *gorm.DB, events []models.TaskFeeEvent) ([]models.TaskFeeEvent, []models.TaskFeeEvent, error) {
	var taskIDCommitments []string
	taskFeeEventMap := make(map[string]models.TaskFeeEvent)
	for _, event := range events {
		taskIDCommitments = append(taskIDCommitments, event.TaskIDCommitment)
		taskFeeEventMap[event.TaskIDCommitment] = event
	}

	var tasks []models.InferenceTask
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).
		Where("task_id_commitment IN (?)", taskIDCommitments).
		Find(&tasks).Error; err != nil {
		return nil, nil, err
	}
	taskMap := make(map[string]models.InferenceTask)
	for _, task := range tasks {
		taskMap[task.TaskIDCommitment] = task
	}

	validEvents := make([]models.TaskFeeEvent, 0)
	invalidEvents := make([]models.TaskFeeEvent, 0)
	for taskIDCommitment, event := range taskFeeEventMap {
		if task, exists := taskMap[taskIDCommitment]; exists {
			if task.Status == models.TaskValidated || task.Status == models.TaskGroupValidated {
				continue
			} else if task.Status == models.TaskEndSuccess || task.Status == models.TaskEndGroupSuccess || task.Status == models.TaskEndGroupRefund {
				if event.TaskFee.Int.Cmp(&task.TaskFee.Int) > 0 {
					invalidEvents = append(invalidEvents, event)
					continue
				}
				validEvents = append(validEvents, event)
			} else {
				invalidEvents = append(invalidEvents, event)
			}
		} else {
			invalidEvents = append(invalidEvents, event)
		}
	}

	return validEvents, invalidEvents, nil
}

func validatePendingWithdrawEvents(ctx context.Context, db *gorm.DB, events []models.WithdrawRecord) ([]models.WithdrawRecord, []models.WithdrawRecord, error) {
	var addresses []string
	for _, event := range events {
		addresses = append(addresses, event.Address)
	}

	var taskFees []models.TaskFee
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).Where("address IN (?)", addresses).Find(&taskFees).Error; err != nil {
		return nil, nil, err
	}

	existedTaskFeeAddressMap := make(map[string]struct{})
	for _, taskFee := range taskFees {
		existedTaskFeeAddressMap[taskFee.Address] = struct{}{}
	}

	validEvents := make([]models.WithdrawRecord, 0)
	invalidEvents := make([]models.WithdrawRecord, 0)
	for _, event := range events {
		if _, exists := existedTaskFeeAddressMap[event.Address]; exists {
			validEvents = append(validEvents, event)
		} else {
			invalidEvents = append(invalidEvents, event)
		}
	}
	return validEvents, invalidEvents, nil
}

func getTaskFeeFromCache(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	taskFeeCache.mu.RLock()
	taskFee, exists := taskFeeCache.taskFees[address]
	taskFeeCache.mu.RUnlock()

	if exists {
		return taskFee, nil
	}

	taskFeeCache.mu.Lock()
	defer taskFeeCache.mu.Unlock()

	if taskFee, exists := taskFeeCache.taskFees[address]; exists {
		return taskFee, nil
	}

	var dbTaskFee models.TaskFee
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.WithContext(dbCtx).Where("address = ?", address).Attrs(models.TaskFee{TaskFee: models.BigInt{Int: *big.NewInt(0)}}).FirstOrInit(&dbTaskFee).Error; err != nil {
		return nil, err
	}

	taskFeeCache.taskFees[address] = &dbTaskFee.TaskFee.Int

	return taskFeeCache.taskFees[address], nil
}

func StartTaskFeeSync(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := syncTaskFeesToDB(ctx, db); err != nil {
				log.Errorf("Failed to sync task fees: %v", err)
			}
		}
	}
}

func mergeTaskFeeEvents(events []models.TaskFeeEvent) map[string]*big.Int {
	mergedTaskFees := make(map[string]*big.Int)
	for _, event := range events {
		if _, exists := mergedTaskFees[event.Address]; !exists {
			mergedTaskFees[event.Address] = big.NewInt(0).Set(&event.TaskFee.Int)
		} else {
			mergedTaskFees[event.Address].Add(mergedTaskFees[event.Address], &event.TaskFee.Int)
		}
	}
	return mergedTaskFees
}

func mergeWithdrawEvents(events []models.WithdrawRecord) map[string]*big.Int {
	mergedWithdrawFees := make(map[string]*big.Int)
	for _, event := range events {
		if _, exists := mergedWithdrawFees[event.Address]; !exists {
			mergedWithdrawFees[event.Address] = big.NewInt(0)
		}
		if event.Status == models.WithdrawStatusFailed {
			mergedWithdrawFees[event.Address].Add(mergedWithdrawFees[event.Address], &event.Amount.Int)
		} else {
			mergedWithdrawFees[event.Address].Sub(mergedWithdrawFees[event.Address], &event.Amount.Int)
		}
	}
	return mergedWithdrawFees
}

func processPendingTaskFeeEvents(ctx context.Context, db *gorm.DB, events []models.TaskFeeEvent) error {
	validEvents, invalidEvents, err := validatePendingTaskFeeEvents(ctx, db, events)
	if err != nil {
		return err
	}

	if len(invalidEvents) > 0 {
		var invalidEventIDs []uint
		for _, event := range invalidEvents {
			invalidEventIDs = append(invalidEventIDs, event.ID)
		}
		if err := db.Model(&models.TaskFeeEvent{}).Where("id IN (?)", invalidEventIDs).Update("status", models.TaskFeeEventStatusInvalid).Error; err != nil {
			return err
		}
	}
	mergedTaskFees := mergeTaskFeeEvents(validEvents)

	var eventIDs []uint
	for _, event := range validEvents {
		eventIDs = append(eventIDs, event.ID)
	}

	var addresses []string
	for address := range mergedTaskFees {
		addresses = append(addresses, address)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var existedTaskFees []models.TaskFee
		if err := tx.Model(&models.TaskFee{}).Where("address IN (?)", addresses).Find(&existedTaskFees).Error; err != nil {
			return err
		}

		existedAddresses := make([]string, len(existedTaskFees))
		existedTaskFeesMap := make(map[string]*models.TaskFee)
		for i, taskFee := range existedTaskFees {
			existedTaskFeesMap[taskFee.Address] = &existedTaskFees[i]
			existedAddresses[i] = taskFee.Address
		}

		var newTaskFees []models.TaskFee
		for address, amount := range mergedTaskFees {
			if taskFee, exists := existedTaskFeesMap[address]; !exists {
				newTaskFees = append(newTaskFees, models.TaskFee{Address: address, TaskFee: models.BigInt{Int: *amount}})
			} else {
				taskFee.TaskFee.Int.Add(&taskFee.TaskFee.Int, amount)
			}
		}

		if len(newTaskFees) > 0 {
			if err := tx.CreateInBatches(&newTaskFees, 100).Error; err != nil {
				return err
			}
		}

		if len(existedTaskFeesMap) > 0 {
			var cases string
			for _, taskFee := range existedTaskFeesMap {
				cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", taskFee.Address, taskFee.TaskFee.String())
			}
			if err := tx.Model(&models.TaskFee{}).Where("address IN (?)", existedAddresses).
				Update("task_fee", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&models.TaskFeeEvent{}).Where("id IN (?)", eventIDs).Where("status = ?", models.TaskFeeEventStatusPending).Update("status", models.TaskFeeEventStatusProcessed).Error; err != nil {
			return err
		}

		return nil
	})
}

func processPendingWithdrawEvents(ctx context.Context, db *gorm.DB, events []models.WithdrawRecord) error {
	validEvents, invalidEvents, err := validatePendingWithdrawEvents(ctx, db, events)
	if err != nil {
		return err
	}

	if len(invalidEvents) != 0 {
		var invalidEventIDs []uint
		for _, event := range invalidEvents {
			invalidEventIDs = append(invalidEventIDs, event.ID)
		}
		if err := db.Model(&models.WithdrawRecord{}).Where("id IN (?)", invalidEventIDs).Update("local_status", models.WithdrawLocalStatusInvalid).Error; err != nil {
			return err
		}
	}

	mergedWithdrawFees := mergeWithdrawEvents(validEvents)

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var eventIDs []uint
	for _, event := range validEvents {
		eventIDs = append(eventIDs, event.ID)
	}

	var addresses []string
	for address := range mergedWithdrawFees {
		addresses = append(addresses, address)
	}

	return db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var existedTaskFees []models.TaskFee
		if err := tx.Model(&models.TaskFee{}).Where("address IN (?)", addresses).Find(&existedTaskFees).Error; err != nil {
			return err
		}

		existedAddresses := make([]string, len(existedTaskFees))
		existedTaskFeesMap := make(map[string]*models.TaskFee)
		for i, taskFee := range existedTaskFees {
			existedTaskFeesMap[taskFee.Address] = &existedTaskFees[i]
			existedAddresses[i] = taskFee.Address
		}

		for address, amount := range mergedWithdrawFees {
			if taskFee, exists := existedTaskFeesMap[address]; exists {
				taskFee.TaskFee.Int.Add(&taskFee.TaskFee.Int, amount)
			}
		}

		var cases string
		for _, taskFee := range existedTaskFeesMap {
			cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", taskFee.Address, taskFee.TaskFee.String())
		}
		if err := tx.Model(&models.TaskFee{}).Where("address IN (?)", existedAddresses).
			Update("task_fee", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.WithdrawRecord{}).Where("id IN (?)", eventIDs).Where("local_status = ?", models.WithdrawLocalStatusPending).Update("local_status", models.WithdrawLocalStatusProcessed).Error; err != nil {
			return err
		}

		return nil
	})
}

func syncTaskFeesToDB(ctx context.Context, db *gorm.DB) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := func() error {
				for {
					events, err := getPendingTaskFeeEvents(ctx, db, 50)
					if err != nil {
						return err
					}

					if len(events) == 0 {
						break
					}

					if err := processPendingTaskFeeEvents(ctx, db, events); err != nil {
						return err
					}
				}
				return nil
			}()

			if err != nil {
				log.Errorf("Failed to sync task fee events: %v", err)
			}

			err = func() error {
				for {
					events, err := getPendingWithdrawEvents(ctx, db, 50)
					if err != nil {
						return err
					}

					if len(events) == 0 {
						break
					}

					if err := processPendingWithdrawEvents(ctx, db, events); err != nil {
						return err
					}
				}
				return nil
			}()

			if err != nil {
				log.Errorf("Failed to sync withdraw events: %v", err)
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

func sendTaskFee(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int) (func() error, error) {
	appConfig := config.GetConfig()
	daoFee := big.NewInt(0).Mul(amount, big.NewInt(0).SetUint64(appConfig.Dao.Percent))
	daoFee.Div(daoFee, big.NewInt(100))

	reward := big.NewInt(0).Sub(amount, daoFee)

	rewardEvent := &models.TaskFeeEvent{
		TaskIDCommitment: taskIDCommitment,
		Address:          address,
		TaskFee:          models.BigInt{Int: *new(big.Int).Set(reward)},
		CreatedAt:        time.Now(),
		Status:           models.TaskFeeEventStatusPending,
	}
	daoEvent := &models.TaskFeeEvent{
		TaskIDCommitment: taskIDCommitment,
		Address:          appConfig.Dao.Address,
		TaskFee:          models.BigInt{Int: *new(big.Int).Set(daoFee)},
		CreatedAt:        time.Now(),
		Status:           models.TaskFeeEventStatusPending,
	}
	events := []*models.TaskFeeEvent{rewardEvent, daoEvent}

	if err := db.Create(events).Error; err != nil {
		return nil, err
	}

	nodeTaskFee, err := getTaskFeeFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	daoTaskFee, err := getTaskFeeFromCache(ctx, db, appConfig.Dao.Address)
	if err != nil {
		return nil, err
	}
	commitFunc := func() error {
		taskFeeCache.mu.Lock()
		defer taskFeeCache.mu.Unlock()
		nodeTaskFee.Add(nodeTaskFee, reward)
		daoTaskFee.Add(daoTaskFee, daoFee)
		return nil
	}

	return commitFunc, nil
}

func withdrawTaskFee(ctx context.Context, db *gorm.DB, address string, amount *big.Int) (func() error, error) {
	taskFee, err := getTaskFeeFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}
	if taskFee.Cmp(amount) < 0 {
		return nil, fmt.Errorf("insufficient task fee when withdraw")
	}
	amountCopy := new(big.Int).Set(amount)
	commitFunc := func() error {
		taskFeeCache.mu.Lock()
		defer taskFeeCache.mu.Unlock()
		if taskFee.Cmp(amountCopy) < 0 {
			return fmt.Errorf("insufficient task fee when withdraw")
		}
		taskFee.Sub(taskFee, amountCopy)
		return nil
	}

	return commitFunc, nil
}

func rejectWithdrawTaskFee(ctx context.Context, db *gorm.DB, address string, amount *big.Int) (func() error, error) {
	taskFee, err := getTaskFeeFromCache(ctx, db, address)
	if err != nil {
		return nil, err
	}

	amountCopy := new(big.Int).Set(amount)
	commitFunc := func() error {
		taskFeeCache.mu.Lock()
		defer taskFeeCache.mu.Unlock()
		taskFee.Add(taskFee, amountCopy)
		return nil
	}

	return commitFunc, nil
}

func GetTaskFee(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	return getTaskFeeFromCache(ctx, db, address)
}
