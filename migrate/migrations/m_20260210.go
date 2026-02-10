package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260210(db *gorm.DB) *gormigrate.Gormigrate {
	type Node struct {
		HealthBase      float64      `json:"health_base" gorm:"default:1.0"`
		HealthUpdatedAt sql.NullTime `json:"health_updated_at" gorm:"null;default:null"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260210",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&Node{}, "HealthBase"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&Node{}, "HealthUpdatedAt"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&Node{}, "HealthBase"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&Node{}, "HealthUpdatedAt"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
