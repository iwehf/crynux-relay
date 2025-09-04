package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250904(db *gorm.DB) *gormigrate.Gormigrate {
	type WithdrawRecord struct {
		TaskFeeEventID uint `json:"task_fee_event_id" gorm:"not null"`
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250904",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&WithdrawRecord{}, "TaskFeeEventID"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&WithdrawRecord{}, "TaskFeeEventID"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
