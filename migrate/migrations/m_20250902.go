package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250902(db *gorm.DB) *gormigrate.Gormigrate {
	type CreditsRecord struct {
		gorm.Model
		Address                 string `json:"address" gorm:"index"`
		Amount                  string `json:"amount" gorm:"type:string;size:255"`
		Network                 string `json:"network" gorm:"index"`
		BlockchainTransactionID uint   `json:"blockchain_transaction_id" gorm:"index;not null"`
	}

	type BlockchainTransaction struct {
		LastRetryAt       sql.NullTime   `json:"last_retry_at" gorm:"null"`
		RetryTransactionID sql.NullInt64 `json:"retry_transaction_id" gorm:"null;index"`
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250902",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().CreateTable(&CreditsRecord{}); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&BlockchainTransaction{}, "RetryTransactionID"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&BlockchainTransaction{}, "LastRetryAt"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable(&CreditsRecord{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&BlockchainTransaction{}, "RetryTransactionID"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&BlockchainTransaction{}, "LastRetryAt"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
