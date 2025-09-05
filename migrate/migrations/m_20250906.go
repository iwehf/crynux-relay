package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250906(db *gorm.DB) *gormigrate.Gormigrate {
	type WithdrawRecord struct {
		WithdrawalFee string `json:"withdrawal_fee" gorm:"not null;type:string;size:191"`
	}
	type TaskFeeEvent struct {
		TaskIDCommitment string `json:"task_id_commitment" gorm:"type:string;size:191;not null;uniqueIndex"`
		Reason           string `json:"reason" gorm:"type:string;size:191;not null;uniqueIndex"`
		Type             uint8  `json:"type" gorm:"not null;index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250906",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&WithdrawRecord{}, "WithdrawalFee"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&TaskFeeEvent{}, "TaskIDCommitment"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&TaskFeeEvent{}, "Reason"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&TaskFeeEvent{}, "Type"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&WithdrawRecord{}, "WithdrawalFee"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&TaskFeeEvent{}, "TaskIDCommitment"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&TaskFeeEvent{}, "Reason"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&TaskFeeEvent{}, "Type"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
