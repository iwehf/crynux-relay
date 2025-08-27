package migrations

import (
	"database/sql"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250821(db *gorm.DB) *gormigrate.Gormigrate {
	type NativeTokenListener struct {
		gorm.Model
		Network        string    `json:"network" gorm:"not null"`
		LastBlockNum   uint64    `json:"last_block_num" gorm:"not null;default:0"`
		LastUpdateTime time.Time `json:"last_update_time" gorm:"not null"`
	}
	type BlockchainTransaction struct {
		gorm.Model
		Network           string         `json:"network" gorm:"type:string;size:191;index;not null"`
		Type              string         `json:"type" gorm:"type:string;size:191;index;not null"`
		Status            uint8          `json:"status" gorm:"index;not null;default:0"`
		FromAddress       string         `json:"from_address" gorm:"type:string;size:191;not null"`
		ToAddress         string         `json:"to_address" gorm:"type:string;size:191;null"`
		Value             string         `json:"value" gorm:"not null;default:'0'"`
		Data              sql.NullString `json:"data" gorm:"null"`
		TxHash            sql.NullString `json:"tx_hash" gorm:"type:string;size:191;null;uniqueIndex"`
		BlockNumber       sql.NullInt64  `json:"block_number" gorm:"null"`
		GasUsed           sql.NullInt64  `json:"gas_used" gorm:"null"`
		EffectiveGasPrice sql.NullString `json:"effective_gas_price" gorm:"type:string;size:191;null"`
		StatusMessage     sql.NullString `json:"status_message" gorm:"null"`
		RetryCount        uint8          `json:"retry_count" gorm:"not null;default:0"`
		MaxRetries        uint8          `json:"max_retries" gorm:"not null;default:0"`
		LastRetryAt       sql.NullTime   `json:"last_retry_at" gorm:"null"`
		NextRetryAt       sql.NullTime   `json:"next_retry_at" gorm:"null"`
		SentAt            sql.NullTime   `json:"sent_at" gorm:"null"`
		ConfirmedAt       sql.NullTime   `json:"confirmed_at" gorm:"null"`
		FailedAt          sql.NullTime   `json:"failed_at" gorm:"null"`
	}

	type Node struct {
		Network string `json:"network" gorm:"index"`
	}

	type TaskFee struct {
		gorm.Model
		Address string `json:"address" gorm:"type:string;size:191;uniqueIndex"`
		TaskFee string `json:"task_fee" gorm:"type:string;size:191;not null"`
	}

	type TaskFeeEvent struct {
		ID               uint      `json:"id" gorm:"primarykey"`
		CreatedAt        time.Time `json:"created_at" gorm:"not null"`
		TaskIDCommitment string    `json:"task_id_commitment" gorm:"type:string;size:191;not null;uniqueIndex"`
		Address          string    `json:"address" gorm:"type:string;size:191;not null;index"`
		TaskFee          string    `json:"task_fee" gorm:"type:string;size:191;not null"`
		Status           int8       `json:"status" gorm:"not null;default:0;index"`
	}

	type TaskQuota struct {
		gorm.Model
		Address string `json:"address" gorm:"type:string;size:191;uniqueIndex"`
		Quota   string `json:"quota" gorm:"type:string;size:191;not null"`
	}

	type TaskQuotaEvent struct {
		ID            uint      `json:"id" gorm:"primarykey"`
		CreatedAt     time.Time `json:"created_at" gorm:"not null"`
		Address       string    `json:"address" gorm:"type:string;size:191;not null;index"`
		Quota         string    `json:"quota" gorm:"type:string;size:191;not null"`
		Status        int8       `json:"status" gorm:"not null;default:0;index"`
		TaskQuotaType int8       `json:"task_quota_type" gorm:"not null;index"`
		Reason        string    `json:"reason" gorm:"type:string;size:191;not null;uniqueIndex"`
	}

	type WithdrawRecord struct {
		gorm.Model
		Address        string `json:"address" gorm:"type:string;size:191;not null;index"`
		BenefitAddress string `json:"benefit_address" gorm:"type:string;size:191;not null;index"`
		Amount         string `json:"amount" gorm:"type:string;size:191;not null"`
		Network        string `json:"network" gorm:"type:string;size:191;not null;index"`
		Status         int8    `json:"status" gorm:"not null;default:0;index"`
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250821",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&Node{}, "Network"); err != nil {
					return err
				}

				return tx.AutoMigrate(&NativeTokenListener{}, &BlockchainTransaction{}, &TaskFee{}, &TaskFeeEvent{}, &TaskQuota{}, &TaskQuotaEvent{}, &WithdrawRecord{})
			},
			Rollback: func(tx *gorm.DB) error {
				if tx.Migrator().HasColumn(&Node{}, "Network") {
					if err := tx.Migrator().DropColumn(&Node{}, "Network"); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&NativeTokenListener{}) {
					if err := tx.Migrator().DropTable(&NativeTokenListener{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&BlockchainTransaction{}) {
					if err := tx.Migrator().DropTable(&BlockchainTransaction{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&TaskFee{}) {
					if err := tx.Migrator().DropTable(&TaskFee{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&TaskFeeEvent{}) {
					if err := tx.Migrator().DropTable(&TaskFeeEvent{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&TaskQuota{}) {
					if err := tx.Migrator().DropTable(&TaskQuota{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&TaskQuotaEvent{}) {
					if err := tx.Migrator().DropTable(&TaskQuotaEvent{}); err != nil {
						return err
					}
				}

				if tx.Migrator().HasTable(&WithdrawRecord{}) {
					if err := tx.Migrator().DropTable(&WithdrawRecord{}); err != nil {
						return err
					}
				}

				return nil
			},
		},
	})
}
