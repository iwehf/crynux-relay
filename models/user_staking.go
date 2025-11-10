package models

import "gorm.io/gorm"

type UserStaking struct {
	gorm.Model
	UserAddress string `json:"user_address" gorm:"index;type:string;size:191;not null"`
	NodeAddress string `json:"node_address" gorm:"index;type:string;size:191;not null"`
	Amount      BigInt `json:"amount" gorm:"type:string;size:191;not null"`
	Valid       bool   `json:"valid" gorm:"not null;index"`
	Network     string `json:"network" gorm:"not null;index;type:string;size:191"`
}
