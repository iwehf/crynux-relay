package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250826(db *gorm.DB) *gormigrate.Gormigrate {
	type WithdrawRecord struct {
		LocalStatus    int `json:"local_status" gorm:"not null;default:0;index"`
	}
	
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250826",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().AddColumn(&WithdrawRecord{}, "LocalStatus")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&WithdrawRecord{}, "LocalStatus")
			},
		},
	})
}
