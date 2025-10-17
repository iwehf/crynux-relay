package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20251014(db *gorm.DB) *gormigrate.Gormigrate {
	type DepositRecord struct {
		gorm.Model
		Address string `json:"address" gorm:"not null;type:string;size:191;index"`
		Amount  string `json:"amount" gorm:"not null;type:string;size:191"`
		Network string `json:"network" gorm:"not null;type:string;size:191;index"`
		TxHash  string `json:"tx_hash" gorm:"not null;type:string;size:191;index"`
	}
	
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20251014",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().CreateTable(&DepositRecord{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable(&DepositRecord{}); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
