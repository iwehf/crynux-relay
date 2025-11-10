package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeStaking struct {
	gorm.Model
	NodeAddress      string    `json:"node_address" gorm:"uniqueIndex:idx_node_time;type:string;size:191;not null"`
	OperatorStaking  BigInt    `json:"operator_staking" gorm:"type:string;size:191;not null"`
	DelegatorStaking BigInt    `json:"delegator_staking" gorm:"type:string;size:191;not null"`
	Time             time.Time `json:"time" gorm:"uniqueIndex:idx_node_time"`
}
