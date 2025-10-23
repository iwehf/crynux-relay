package service

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/config"
	"crynux_relay/models"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// StartNativeTokenListener starts the native token transfer listener
func StartNativeTokenListener(ctx context.Context) {
	appConfig := config.GetConfig()

	// Start the listener goroutine
	for network := range appConfig.Blockchains {
		go func(network string) {
			if err := runNativeTokenListener(ctx, config.GetDB(), network); err != nil {
				log.Errorf("Native token listener failed: %v", err)
			}
		}(network)
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
	processedBlock := listener.LastBlockNum
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
			break
		}
		processedBlock = blockNum
	}

	// Update listener status
	if err := db.Model(&listener).Updates(map[string]interface{}{
		"last_block_num":   processedBlock,
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
			return err
		}
	}

	return nil
}

func processBuyQuotaTransaction(ctx context.Context, db *gorm.DB, tx *types.Transaction, client *blockchain.BlockchainClient) error {

	// Check if transaction is successful
	receipt, err := client.RpcClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt of %s, network: %s, error: %w", tx.Hash().Hex(), client.Network, err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil
	}

	// Check if already processed
	event, err := models.GetTaskQuotaBoughtEvent(ctx, db, tx.Hash().Hex(), client.Network)
	if err != nil {
		return err
	}
	if event != nil {
		return nil
	}

	// Get sender address (need to recover from signature)
	from, err := types.Sender(types.LatestSignerForChainID(client.ChainID), tx)
	if err != nil {
		return fmt.Errorf("failed to get sender address of %s, network: %s, error: %w", tx.Hash().Hex(), client.Network, err)
	}

	// Call BuyTaskQuota to add quota for the sender
	commitFunc, err := buyTaskQuota(ctx, db, tx.Hash().Hex(), from.Hex(), tx.Value(), client.Network)
	if err != nil {
		log.Errorf("Failed to buy task quota for %s, network: %s, error: %v", from.Hex(), client.Network, err)
		return err
	}

	// Execute quota update
	if err := commitFunc(); err != nil {
		log.Errorf("Failed to buy task quota for %s, network: %s, error: %v", from.Hex(), client.Network, err)
		return err
	}
	return nil
}

func processBuyTaskFeeTransaction(ctx context.Context, db *gorm.DB, tx *types.Transaction, client *blockchain.BlockchainClient) error {
	// Check if transaction is successful
	receipt, err := client.RpcClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt of %s, network: %s, error: %w", tx.Hash().Hex(), client.Network, err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil
	}

	// Check if already processed
	event, err := models.GetTaskFeeBoughtEvent(ctx, db, tx.Hash().Hex(), client.Network)
	if err != nil {
		return err
	}
	if event != nil {
		return nil
	}

	// Get sender address (need to recover from signature)
	from, err := types.Sender(types.LatestSignerForChainID(client.ChainID), tx)
	if err != nil {
		return fmt.Errorf("failed to get sender address of %s, network: %s, error: %w", tx.Hash().Hex(), client.Network, err)
	}

	// Call BuyTaskQuota to add quota for the sender
	commitFunc, err := buyTaskFee(ctx, db, tx.Hash().Hex(), from.Hex(), tx.Value(), client.Network)
	if err != nil {
		log.Errorf("Failed to buy task fee for %s, network: %s, error: %v", from.Hex(), client.Network, err)
		return err
	}

	// Execute quota update
	if err := commitFunc(); err != nil {
		log.Errorf("Failed to buy task fee for %s, network: %s, error: %v", from.Hex(), client.Network, err)
		return err
	}
	return nil
}

// processTransaction processes a single transaction
func processTransaction(ctx context.Context, db *gorm.DB, tx *types.Transaction, client *blockchain.BlockchainClient) error {
	// Only process native token transfers (to field is not empty and data field is empty)
	if tx.To() == nil || len(tx.Data()) > 0 {
		return nil
	}

	appConfig := config.GetConfig()

	// Check if transfer is to the target address
	if strings.EqualFold(tx.To().Hex(), appConfig.BuyQuota.Address) {
		return processBuyQuotaTransaction(ctx, db, tx, client)
	} else if strings.EqualFold(tx.To().Hex(), appConfig.BuyTaskFee.Address) {
		return processBuyTaskFeeTransaction(ctx, db, tx, client)
	}

	return nil
}
