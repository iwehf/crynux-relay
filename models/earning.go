package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStakingEarning struct {
	gorm.Model
	UserAddress string    `json:"user_address" gorm:"uniqueIndex:idx_user_node_time;type:string;size:191;not null"`
	NodeAddress string    `json:"node_address" gorm:"uniqueIndex:idx_user_node_time;type:string;size:191;not null"`
	Earning     BigInt    `json:"earning" gorm:"type:string;size:191;not null"`
	Time        time.Time `json:"time" gorm:"uniqueIndex:idx_user_node_time"`
}

type NodeEarning struct {
	gorm.Model
	NodeAddress      string    `json:"node_address" gorm:"uniqueIndex:idx_node_time;type:string;size:191;not null"`
	OperatorEarning  BigInt    `json:"operator_earning" gorm:"type:string;size:191;not null"`
	DelegatorEarning BigInt    `json:"delegator_earning" gorm:"type:string;size:191;not null"`
	Time             time.Time `json:"time" gorm:"uniqueIndex:idx_node_time"`
}
