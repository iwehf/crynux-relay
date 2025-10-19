package migrations

import (
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/utils"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20251019(db *gorm.DB) *gormigrate.Gormigrate {
	
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20251019",
			Migrate: func(tx *gorm.DB) error {
				appConfig := config.GetConfig()
				var events []models.TaskFeeEvent
				if err := tx.Model(&models.TaskFeeEvent{}).Where("id > ?", 5142472).Where("mac = ?", "").Find(&events).Error; err != nil {
					return err
				}
				for _, event := range events {
					event.MAC = utils.GenerateMAC([]byte(event.Reason), appConfig.MAC.SecretKey)
					updates := map[string]interface{}{
						"mac": event.MAC,
						"status": models.TaskFeeEventStatusProcessed,
					}
					if err := tx.Model(&models.TaskFeeEvent{}).Where("id = ?", event.ID).Updates(updates).Error; err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})
}
