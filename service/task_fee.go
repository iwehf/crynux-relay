package service

import (
	"context"
	"crynux_relay/models"
	"math/big"

	"gorm.io/gorm"
)

var ErrInsufficientTaskFee = ErrInsufficientRelayAccount

func InitTaskFeeCache(ctx context.Context, db *gorm.DB) error {
	return InitRelayAccountCache(ctx, db)
}

func StartTaskFeeSync(ctx context.Context, db *gorm.DB) {
	StartRelayAccountSync(ctx, db)
}

func buyTaskFee(ctx context.Context, db *gorm.DB, txHash, address string, amount *big.Int, network string) (func() error, error) {
	return depositRelayAccount(ctx, db, txHash, address, amount, network)
}

func sendTaskFee(ctx context.Context, db *gorm.DB, taskIDCommitment, address string, amount *big.Int, taskType models.TaskType) (func() error, error) {
	return sendTaskIncome(ctx, db, taskIDCommitment, address, amount, taskType)
}

func withdrawTaskFee(ctx context.Context, db *gorm.DB, address string, amount *big.Int) (func() error, error) {
	return chargeTaskFromRelayAccount(ctx, db, "legacy-withdraw", address, amount)
}

func fulfillWithdrawTaskFee(ctx context.Context, db *gorm.DB, withdrawID uint, withdrawalFeeAddress string, withdrawalFee *big.Int) (func() error, error) {
	return fulfillWithdrawFeeIncome(ctx, db, withdrawID, withdrawalFeeAddress, withdrawalFee)
}

func rejectWithdrawTaskFee(ctx context.Context, db *gorm.DB, address string, amount *big.Int) (func() error, error) {
	return rejectWithdrawToRelayAccount(ctx, db, 0, address, amount)
}

func GetTaskFee(ctx context.Context, db *gorm.DB, address string) (*big.Int, error) {
	return GetRelayAccountBalance(ctx, db, address)
}
