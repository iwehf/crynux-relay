package service

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/models"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

func CreateCredits(ctx context.Context, db *gorm.DB, address string, amount *big.Int, network string) (*models.CreditsRecord, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	record := &models.CreditsRecord{
		Address: address,
		Amount:  models.BigInt{Int: *amount},
		Network: network,
	}

	if err := db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		blockchainTransaction, err := blockchain.QueueCreateCredits(dbCtx, tx, common.HexToAddress(address), amount, network)
		if err != nil {
			return err
		}
		record.BlockchainTransactionID = blockchainTransaction.ID
		
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return record, nil
}

