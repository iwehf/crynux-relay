package models

import (
	"context"
	"crynux_relay/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type NativeTokenListener struct {
	gorm.Model
	Network        string    `json:"network" gorm:"not null"`
	LastBlockNum   uint64    `json:"last_block_num" gorm:"not null;default:0"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"not null"`
}

func GetNativeTokenListener(ctx context.Context, db *gorm.DB, network string) (*NativeTokenListener, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	appConfig := config.GetConfig()
	blockchain, exists := appConfig.Blockchains[network]
	if !exists {
		return nil, errors.New("blockchain not found")
	}
	var listener NativeTokenListener
	err := db.WithContext(dbCtx).Model(&NativeTokenListener{}).Attrs(&NativeTokenListener{
		LastBlockNum:   blockchain.StartBlockNum,
		LastUpdateTime: time.Now(),
		Network: network,
	}).FirstOrCreate(&listener).Error

	if err != nil {
		return nil, err
	}

	return &listener, nil
}
