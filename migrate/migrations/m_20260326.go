package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260326(db *gorm.DB) *gormigrate.Gormigrate {
	type DepositRecord struct {
		LocalStatus int8 `gorm:"column:local_status"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260326",
			Migrate: func(tx *gorm.DB) error {
				m := tx.Migrator()
				if err := m.AddColumn(&DepositRecord{}, "LocalStatus"); err != nil {
					return err
				}
				if err := m.CreateIndex(&DepositRecord{}, "LocalStatus"); err != nil {
					return err
				}
				return tx.Table("deposit_records").
					Where("local_status = 0").
					Update("local_status", 1).Error
			},
			Rollback: func(tx *gorm.DB) error {
				m := tx.Migrator()
				if err := m.DropIndex(&DepositRecord{}, "LocalStatus"); err != nil {
					return err
				}
				if err := m.DropColumn(&DepositRecord{}, "local_status"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
