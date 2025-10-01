package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20251001(db *gorm.DB) *gormigrate.Gormigrate {
	type TaskFeeEvent struct {
		MAC string `json:"mac" gorm:"not null;type:string;size:191"`
	}

	type WithdrawRecord struct {
		MAC string `json:"mac" gorm:"not null;type:string;size:191"`
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20251001",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&TaskFeeEvent{}, "MAC"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&WithdrawRecord{}, "MAC"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&TaskFeeEvent{}, "MAC"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&WithdrawRecord{}, "MAC"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
