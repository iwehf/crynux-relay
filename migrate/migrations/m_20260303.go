package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260303(db *gorm.DB) *gormigrate.Gormigrate {
	type NetworkNodeData struct {
		HealthBase      float64      `json:"health_base" gorm:"default:1.0"`
		HealthUpdatedAt sql.NullTime `json:"health_updated_at" gorm:"null;default:null"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260303",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&NetworkNodeData{}, "HealthBase"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&NetworkNodeData{}, "HealthUpdatedAt"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&NetworkNodeData{}, "HealthUpdatedAt"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&NetworkNodeData{}, "HealthBase"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
