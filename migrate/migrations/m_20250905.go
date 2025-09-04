package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250905(db *gorm.DB) *gormigrate.Gormigrate {
	type WithdrawRecord struct {
		TxHash sql.NullString `json:"tx_hash" gorm:"null"`
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250905",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&WithdrawRecord{}, "TxHash"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&WithdrawRecord{}, "TxHash"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
