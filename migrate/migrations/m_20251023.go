package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20251023(db *gorm.DB) *gormigrate.Gormigrate {
	type TaskFeeEvent struct {
		TaskIDCommitment string `json:"task_id_commitment" gorm:"type:string;size:191;not null;uniqueIndex"`
		Reason           string `json:"reason" gorm:"type:string;size:191;not null;uniqueIndex"`
		Type             uint8  `json:"type" gorm:"not null;index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20251023",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().CreateIndex(&TaskFeeEvent{}, "TaskIDCommitment"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&TaskFeeEvent{}, "Reason"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&TaskFeeEvent{}, "Type"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if tx.Migrator().HasIndex(&TaskFeeEvent{}, "TaskIDCommitment") {
					if err := tx.Migrator().DropIndex(&TaskFeeEvent{}, "TaskIDCommitment"); err != nil {
						return err
					}
				}
				if tx.Migrator().HasIndex(&TaskFeeEvent{}, "Reason") {
					if err := tx.Migrator().DropIndex(&TaskFeeEvent{}, "Reason"); err != nil {
						return err
					}
				}
				if tx.Migrator().HasIndex(&TaskFeeEvent{}, "Type") {
					if err := tx.Migrator().DropIndex(&TaskFeeEvent{}, "Type"); err != nil {
						return err
					}
				}
				return nil
			},
		},
	})
}
