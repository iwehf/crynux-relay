package models

import (
	"context"
	"gorm.io/gorm"
)

type CreditsRecord struct {
	gorm.Model
	Address                 string `json:"address" gorm:"index"`
	Amount                  BigInt `json:"amount" gorm:"type:string;size:255"`
	Network                 string `json:"network" gorm:"index"`
	BlockchainTransactionID uint   `json:"blockchain_transaction_id" gorm:"index"`
}

func GetCreditsRecordByID(ctx context.Context, db *gorm.DB, id uint) (*CreditsRecord, error) {
	var record CreditsRecord
	if err := db.WithContext(ctx).Model(&CreditsRecord{}).Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}