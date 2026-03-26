package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260324(db *gorm.DB) *gormigrate.Gormigrate {
	type DepositRecord struct {
		RelayAccountEventID uint `gorm:"column:relay_account_event_id;index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260324",
			Migrate: func(tx *gorm.DB) error {
				m := tx.Migrator()
				if err := m.AddColumn(&DepositRecord{}, "RelayAccountEventID"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				m := tx.Migrator()
				if err := m.DropColumn(&DepositRecord{}, "relay_account_event_id"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
