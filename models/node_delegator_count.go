package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeDelegatorCount struct {
	gorm.Model
	NodeAddress string    `json:"node_address" gorm:"uniqueIndex:idx_node_time;type:string;size:191;not null"`
	Count       uint64    `json:"count"`
	Time        time.Time `json:"time" gorm:"uniqueIndex:idx_node_time"`
}
