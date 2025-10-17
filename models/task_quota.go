package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaskQuotaEventStatus int8

const (
	TaskQuotaEventStatusPending TaskQuotaEventStatus = iota
	TaskQuotaEventStatusProcessed
	TaskQuotaEventStatusInvalid
)

type TaskQuotaType int8

const (
	TaskQuotaTypeBought TaskQuotaType = iota
	TaskQuotaTypeSpent
	TaskQuotaTypeRefunded
)

type TaskQuota struct {
	gorm.Model
	Address string `json:"address" gorm:"uniqueIndex"`
	Quota   BigInt `json:"quota" gorm:"type:string;size:255"`
}

type TaskQuotaEvent struct {
	ID            uint                 `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time            `json:"created_at" gorm:"not null"`
	Address       string               `json:"address" gorm:"not null;index"`
	Quota         BigInt               `json:"quota" gorm:"not null"`
	Status        TaskQuotaEventStatus `json:"status" gorm:"not null;default:0;index"`
	TaskQuotaType TaskQuotaType        `json:"task_quota_type" gorm:"not null;index"`
	Reason        string               `json:"reason" gorm:"not null;uniqueIndex"`
}

func GetTaskQuotaBoughtEvent(ctx context.Context, db *gorm.DB, txHash, network string) (*TaskQuotaEvent, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var event TaskQuotaEvent
	if err := db.WithContext(dbCtx).Where("reason = ?", fmt.Sprintf("%d-%s-%s", TaskQuotaTypeBought, txHash, network)).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}