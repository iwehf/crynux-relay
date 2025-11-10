package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeScore struct {
	gorm.Model
	NodeAddress  string    `json:"node_address" gorm:"uniqueIndex:idx_node_time;type:string;size:191;not null"`
	Time         time.Time `json:"time" gorm:"uniqueIndex:idx_node_time"`
	ProbWeight   float64   `json:"prob_weight"`
	StakingScore float64   `json:"staking_score"`
	QOSScore     float64   `json:"qos_score"`
}
