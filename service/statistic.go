package service

import (
	"context"
	"crynux_relay/models"
	"errors"
	"math/big"
	"time"

	"gorm.io/gorm"
)

func addNodeIncentive(ctx context.Context, db *gorm.DB, nodeAddress string, incentive float64, taskType models.TaskType) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	t := time.Now().UTC().Truncate(24 * time.Hour)
	nodeIncentive := models.NodeIncentive{Time: t, NodeAddress: nodeAddress}
	if err := db.WithContext(ctx).Model(&nodeIncentive).Where(&nodeIncentive).First(&nodeIncentive).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if nodeIncentive.ID > 0 {
		nodeIncentive.Incentive += incentive
		nodeIncentive.TaskCount += 1
		if taskType == models.TaskTypeSD {
			nodeIncentive.SDTaskCount += 1
		} else if taskType == models.TaskTypeLLM {
			nodeIncentive.LLMTaskCount += 1
		} else if taskType == models.TaskTypeSDFTLora {
			nodeIncentive.SDFTLoraTaskCount += 1
		}
		if err := db.WithContext(dbCtx).Save(&nodeIncentive).Error; err != nil {
			return err
		}
	} else {
		nodeIncentive.Incentive = incentive
		nodeIncentive.TaskCount = 1
		if taskType == models.TaskTypeSD {
			nodeIncentive.SDTaskCount = 1
		} else if taskType == models.TaskTypeLLM {
			nodeIncentive.LLMTaskCount = 1
		} else if taskType == models.TaskTypeSDFTLora {
			nodeIncentive.SDFTLoraTaskCount = 1
		}
		if err := db.WithContext(dbCtx).Create(&nodeIncentive).Error; err != nil {
			return err
		}
	}
	return nil
}

func addNodeEarning(ctx context.Context, db *gorm.DB, nodeAddress string, operatorEarning, delegatorEarning *big.Int) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	t := time.Now().UTC().Truncate(24 * time.Hour)
	if err := db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var nodeEarning models.NodeEarning
		// find today's node earning
		if err := tx.Model(&models.NodeEarning{}).Where("node_address = ?", nodeAddress).Where("time = ?", t).First(&nodeEarning).Error; err != nil {
			// if today's node earning doesn't exist, find yesterday's node earning
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var oldNodeEarning models.NodeEarning
				if err := tx.Model(&models.NodeEarning{}).Where("node_address = ?", nodeAddress).Where("time = ?", t.Add(-24*time.Hour)).First(&oldNodeEarning).Error; err != nil {
					// if yesterdays node earning doesn't exist, create a new node earning
					if errors.Is(err, gorm.ErrRecordNotFound) {
						nodeEarning = models.NodeEarning{
							NodeAddress:      nodeAddress,
							Time:             t,
							OperatorEarning:  models.BigInt{Int: *big.NewInt(0)},
							DelegatorEarning: models.BigInt{Int: *big.NewInt(0)},
						}
					} else {
						return err
					}
				}
				nodeEarning.Time = t
				nodeEarning.NodeAddress = nodeAddress
				nodeEarning.OperatorEarning = oldNodeEarning.OperatorEarning
				nodeEarning.DelegatorEarning = oldNodeEarning.DelegatorEarning
			} else {
				return err
			}
		}
		nodeEarning.OperatorEarning = models.BigInt{Int: *big.NewInt(0).Add(&nodeEarning.OperatorEarning.Int, operatorEarning)}
		nodeEarning.DelegatorEarning = models.BigInt{Int: *big.NewInt(0).Add(&nodeEarning.DelegatorEarning.Int, delegatorEarning)}
		if err := tx.Save(&nodeEarning).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func addUserStakingEarning(ctx context.Context, db *gorm.DB, userAddress, nodeAddress string, earning *big.Int) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	t := time.Now().UTC().Truncate(24 * time.Hour)
	if err := db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var userStakingEarning models.UserStakingEarning
		// find today's user staking earning
		if err := tx.Model(&models.UserStakingEarning{}).Where("user_address = ?", userAddress).Where("node_address = ?", nodeAddress).Where("time = ?", t).First(&userStakingEarning).Error; err != nil {
			// if today's user staking earning doens't exist, find yesterday's node earning
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var oldUserStakingEarning models.UserStakingEarning
				if err := tx.Model(&models.UserStakingEarning{}).Where("user_address = ?", userAddress).Where("node_address = ?", nodeAddress).Where("time = ?", t.Add(-24 * time.Hour)).First(&oldUserStakingEarning).Error; err != nil {
					// if yesterday's node earning doesn't exist, create a new user staking earning
					if errors.Is(err, gorm.ErrRecordNotFound) {
						userStakingEarning = models.UserStakingEarning{
							UserAddress: userAddress,
							NodeAddress: nodeAddress,
							Time: t,
							Earning: models.BigInt{Int: *big.NewInt(0)},
						}
					} else {
						return err
					}
				}
				userStakingEarning.Time = t
				userStakingEarning.UserAddress = userAddress
				userStakingEarning.NodeAddress = nodeAddress
				userStakingEarning.Earning = oldUserStakingEarning.Earning
			} else {
				return err
			}
		}
		userStakingEarning.Earning = models.BigInt{Int: *big.NewInt(0).Add(&userStakingEarning.Earning.Int, earning)}
		if err := tx.Save(&userStakingEarning).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}