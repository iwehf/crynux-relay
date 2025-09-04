package models

import "gorm.io/gorm"

type WithdrawRecord struct {
	gorm.Model
	Address        string              `json:"address" gorm:"not null;index"`
	BenefitAddress string              `json:"benefit_address" gorm:"not null;index"`
	Amount         BigInt              `json:"amount" gorm:"not null"`
	Network        string              `json:"network" gorm:"not null;index"`
	Status         WithdrawStatus      `json:"status" gorm:"not null;default:0;index"`
	LocalStatus    WithdrawLocalStatus `json:"local_status" gorm:"not null;default:0;index"`
	TaskFeeEventID uint                `json:"task_fee_event_id" gorm:"not null"`
}

type WithdrawStatus int8

const (
	WithdrawStatusPending WithdrawStatus = iota
	WithdrawStatusSuccess
	WithdrawStatusFailed
)

type WithdrawLocalStatus int8

const (
	WithdrawLocalStatusPending WithdrawLocalStatus = iota
	WithdrawLocalStatusProcessed
	WithdrawLocalStatusInvalid
)
