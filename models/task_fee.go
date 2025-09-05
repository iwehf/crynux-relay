package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskFee struct {
	gorm.Model
	Address string `json:"address" gorm:"uniqueIndex"`
	TaskFee BigInt `json:"task_fee" gorm:"type:string;size:255"`
}

type TaskFeeEventStatus int8

const (
	TaskFeeEventStatusPending TaskFeeEventStatus = iota
	TaskFeeEventStatusProcessed
	TaskFeeEventStatusInvalid
)

type TaskFeeEventType int8

const (
	TaskFeeEventTypeTask TaskFeeEventType = iota
	TaskFeeEventTypeDraw
	TaskFeeEventTypeWithdrawalFee
)

type TaskFeeEvent struct {
	ID        uint               `json:"id" gorm:"primarykey"`
	CreatedAt time.Time          `json:"created_at" gorm:"not null"`
	Address   string             `json:"address" gorm:"not null;index"`
	TaskFee   BigInt             `json:"task_fee" gorm:"not null"`
	Status    TaskFeeEventStatus `json:"status" gorm:"not null;default:0;index"`
	Reason    string             `json:"reason" gorm:"not null;uniqueIndex"`
	Type      TaskFeeEventType   `json:"type" gorm:"not null;index"`
}
