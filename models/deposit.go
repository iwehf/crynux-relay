package models

import "gorm.io/gorm"

type DepositRecord struct {
	gorm.Model
	Address string `json:"address" gorm:"not null;index"`
	Amount  BigInt `json:"amount" gorm:"not null"`
	Network string `json:"network" gorm:"not null;index"`
	TxHash  string `json:"tx_hash" gorm:"not null;index"`
}
