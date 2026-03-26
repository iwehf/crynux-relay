package models

import "gorm.io/gorm"

type DepositRecord struct {
	gorm.Model
	Address             string             `json:"address" gorm:"not null;index"`
	Amount              BigInt             `json:"amount" gorm:"not null"`
	Network             string             `json:"network" gorm:"not null;index"`
	TxHash              string             `json:"tx_hash" gorm:"not null;index"`
	RelayAccountEventID uint               `json:"relay_account_event_id" gorm:"not null;default:0;index"`
	LocalStatus         DepositLocalStatus `json:"local_status" gorm:"not null;default:0;index"`
}

type DepositLocalStatus int8

const (
	DepositLocalStatusPending DepositLocalStatus = iota
	DepositLocalStatusProcessed
	DepositLocalStatusInvalid
)
