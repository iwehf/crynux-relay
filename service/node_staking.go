package service

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/models"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// StakeAmountUpdater is responsible for periodically updating node staking amounts
type StakeAmountUpdater struct {
	db        *gorm.DB
	interval  time.Duration
	stopChan  chan struct{}
	isRunning bool
	mutex     sync.RWMutex
}

// NewStakeAmountUpdater creates a new StakeAmountUpdater instance
func NewStakeAmountUpdater(db *gorm.DB) *StakeAmountUpdater {
	return &StakeAmountUpdater{
		db:       db,
		interval: time.Minute,
		stopChan: make(chan struct{}),
	}
}

// Start starts the staking amount updater
func (updater *StakeAmountUpdater) Start(ctx context.Context) {
	updater.mutex.Lock()
	if updater.isRunning {
		updater.mutex.Unlock()
		return
	}
	updater.isRunning = true
	updater.mutex.Unlock()

	log.Info("StakeAmountUpdater: starting stake amount updater")

	go updater.run(ctx)
}

func (updater *StakeAmountUpdater) run(ctx context.Context) {
	ticker := time.NewTicker(updater.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("StakeAmountUpdater: stopping due to context cancellation")
			return
		case <-updater.stopChan:
			log.Info("StakeAmountUpdater: stopping due to stop signal")
			return
		case <-ticker.C:
			updater.updateStakeAmounts(ctx)
		}
	}
}

// Stop stops the staking amount updater
func (updater *StakeAmountUpdater) Stop() {
	updater.mutex.Lock()
	defer updater.mutex.Unlock()

	if !updater.isRunning {
		return
	}

	close(updater.stopChan)
	updater.isRunning = false
	log.Info("StakeAmountUpdater: stopped")
}

// IsRunning checks if the updater is currently running
func (updater *StakeAmountUpdater) IsRunning() bool {
	updater.mutex.RLock()
	defer updater.mutex.RUnlock()
	return updater.isRunning
}

// updateStakeAmounts updates staking amounts for all active nodes
func (updater *StakeAmountUpdater) updateStakeAmounts(ctx context.Context) {
	log.Info("StakeAmountUpdater: starting to update stake amounts")

	// Get all nodes that haven't quit
	nodes, err := updater.getActiveNodes(ctx)
	if err != nil {
		log.Errorf("StakeAmountUpdater: failed to get active nodes: %v", err)
		return
	}

	if len(nodes) == 0 {
		log.Info("StakeAmountUpdater: no active nodes to update")
		return
	}

	log.Infof("StakeAmountUpdater: found %d active nodes to update", len(nodes))

	// Concurrently update node staking amounts
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrency to avoid too many blockchain calls
	successCount := int64(0)
	errorCount := int64(0)

	for _, node := range nodes {
		wg.Add(1)
		go func(node models.Node) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := updater.updateNodeStakeAmount(ctx, &node); err != nil {
				log.Errorf("StakeAmountUpdater: failed to update node %s stake amount: %v", node.Address, err)
				atomic.AddInt64(&errorCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(node)
	}

	wg.Wait()

	log.Infof("StakeAmountUpdater: finished updating stake amounts. Success: %d, Errors: %d",
		successCount, errorCount)
}

// getActiveNodes gets all active nodes (nodes that haven't quit)
func (updater *StakeAmountUpdater) getActiveNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node

	// Use timeout context
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Query all nodes whose status is not NodeStatusQuit
	if err := updater.db.WithContext(dbCtx).
		Where("status != ?", models.NodeStatusQuit).
		Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// updateNodeStakeAmount updates the staking amount for a single node
func (updater *StakeAmountUpdater) updateNodeStakeAmount(ctx context.Context, node *models.Node) error {
	// Use timeout context
	blockchainCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get latest staking info from blockchain
	stakingInfo, err := blockchain.GetStakingInfo(blockchainCtx, common.HexToAddress(node.Address), node.Network)
	if err != nil {
		log.Warnf("StakeAmountUpdater: failed to get staking info for node %s: %v", node.Address, err)
		return err
	}

	// Calculate total staking amount
	totalStakeAmount := new(big.Int).Add(stakingInfo.StakedBalance, stakingInfo.StakedCredits)

	// Skip update if staking amount hasn't changed
	if totalStakeAmount.Cmp(&node.StakeAmount.Int) == 0 {
		log.Debugf("StakeAmountUpdater: node %s stake amount unchanged (%s)",
			node.Address, totalStakeAmount.String())
		return nil
	}

	log.Infof("StakeAmountUpdater: updating node %s stake amount from %s to %s",
		node.Address, node.StakeAmount.Int.String(), totalStakeAmount.String())

	// Update staking amount in database
	dbCtx, dbCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dbCancel()

	if err := updater.db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(node).Update("stake_amount", models.BigInt{Int: *totalStakeAmount}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.NetworkNodeData{}).Where("address = ?", node.Address).Update("staking", models.BigInt{Int: *totalStakeAmount}).Error; err != nil {
			return err
		}
		if err := emitEvent(ctx, tx, &models.NodeStakingEvent{NodeAddress: node.Address, StakingAmount: models.BigInt{Int: *totalStakeAmount}}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Errorf("StakeAmountUpdater: failed to update database for node %s: %v", node.Address, err)
		return err
	}

	// Update value in memory
	node.StakeAmount = models.BigInt{Int: *totalStakeAmount}
	UpdateMaxStaking(totalStakeAmount)

	log.Infof("StakeAmountUpdater: successfully updated node %s stake amount to %s",
		node.Address, totalStakeAmount.String())

	return nil
}
