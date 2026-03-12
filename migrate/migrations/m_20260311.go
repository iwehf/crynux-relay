package migrations

import (
	"fmt"
	"math/big"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func M20260311(db *gorm.DB) *gormigrate.Gormigrate {
	type RelayAccount struct {
		gorm.Model
		Address string `gorm:"uniqueIndex"`
		Balance string `gorm:"type:string;size:255"`
	}

	type RelayAccountEvent struct {
		ID        uint      `gorm:"primarykey"`
		CreatedAt time.Time `gorm:"not null"`
		Address   string    `gorm:"not null;index"`
		Amount    string    `gorm:"not null"`
		Status    int8      `gorm:"not null;default:0;index"`
		Reason    string    `gorm:"not null;uniqueIndex"`
		Type      int8      `gorm:"not null;index"`
		MAC       string    `gorm:"not null"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260311",
			Migrate: func(tx *gorm.DB) error {
				// Run the full migration in one transaction to guarantee atomicity.
				return tx.Transaction(func(tx *gorm.DB) error {
					// Read and validate the total task fee amount before data migration.
					calcTaskFeeTotal := func() (*big.Int, int, error) {
						type TaskFeeValueRow struct {
							TaskFee string `gorm:"column:task_fee"`
						}
						var taskFees []TaskFeeValueRow
						if err := tx.Table("task_fees").Select("task_fee").Find(&taskFees).Error; err != nil {
							return nil, 0, err
						}
						total := new(big.Int)
						for _, taskFee := range taskFees {
							value := new(big.Int)
							if taskFee.TaskFee != "" {
								if _, ok := value.SetString(taskFee.TaskFee, 10); !ok {
									return nil, 0, fmt.Errorf("invalid task_fees.task_fee: %q", taskFee.TaskFee)
								}
							}
							total.Add(total, value)
						}
						return total, len(taskFees), nil
					}

					beforeTotal, beforeCount, err := calcTaskFeeTotal()
					if err != nil {
						return err
					}

					// Aggregate all task quota rows by address using arbitrary-precision integers.
					type TaskQuotaRow struct {
						Address string `gorm:"column:address"`
						Quota   string `gorm:"column:quota"`
					}
					var quotas []TaskQuotaRow
					if err := tx.Table("task_quotas").Select("address, quota").Find(&quotas).Error; err != nil {
						return err
					}
					quotaSums := make(map[string]*big.Int, len(quotas))
					for _, quota := range quotas {
						q := new(big.Int)
						if _, ok := q.SetString(quota.Quota, 10); !ok {
							return fmt.Errorf("invalid task_quotas.quota for address %s: %q", quota.Address, quota.Quota)
						}
						if total, exists := quotaSums[quota.Address]; exists {
							total.Add(total, q)
							continue
						}
						quotaSums[quota.Address] = q
					}
					quotaTotal := new(big.Int)
					for _, quotaSum := range quotaSums {
						quotaTotal.Add(quotaTotal, quotaSum)
					}
					// This is the expected ledger total after applying all quotas.
					expectedFinalTotal := new(big.Int).Add(new(big.Int).Set(beforeTotal), quotaTotal)

					// Load existing account rows only for addresses that appear in task_quotas.
					type RelayAccountRow struct {
						ID      uint   `gorm:"column:id"`
						Address string `gorm:"column:address"`
						TaskFee string `gorm:"column:task_fee"`
					}
					relayAccountByAddress := make(map[string]RelayAccountRow, len(quotaSums))
					if len(quotaSums) > 0 {
						addresses := make([]string, 0, len(quotaSums))
						for address := range quotaSums {
							addresses = append(addresses, address)
						}
						var relayAccounts []RelayAccountRow
						if err := tx.Table("task_fees").Select("id, address, task_fee").Where("address IN (?)", addresses).Find(&relayAccounts).Error; err != nil {
							return err
						}
						for _, relayAccount := range relayAccounts {
							relayAccountByAddress[relayAccount.Address] = relayAccount
						}
					}
					type RelayAccountInsertRow struct {
						Address   string    `gorm:"column:address"`
						TaskFee   string    `gorm:"column:task_fee"`
						CreatedAt time.Time `gorm:"column:created_at"`
						UpdatedAt time.Time `gorm:"column:updated_at"`
					}
					newRelayAccounts := make([]RelayAccountInsertRow, 0)
					// Apply quota sums to existing rows and collect missing rows for bulk insert.
					for address, quotaSum := range quotaSums {
						if relayAccount, exists := relayAccountByAddress[address]; exists {
							balance := new(big.Int)
							if relayAccount.TaskFee != "" {
								if _, ok := balance.SetString(relayAccount.TaskFee, 10); !ok {
									return fmt.Errorf("invalid task_fees.task_fee for address %s: %q", address, relayAccount.TaskFee)
								}
							}
							balance.Add(balance, quotaSum)
							if err := tx.Table("task_fees").Where("id = ?", relayAccount.ID).Updates(map[string]interface{}{
								"task_fee": balance.String(),
							}).Error; err != nil {
								return err
							}
							continue
						}

						now := time.Now()
						newRelayAccounts = append(newRelayAccounts, RelayAccountInsertRow{
							Address:   address,
							TaskFee:   quotaSum.String(),
							CreatedAt: now,
							UpdatedAt: now,
						})
					}
					if len(newRelayAccounts) > 0 {
						if err := tx.Table("task_fees").Create(&newRelayAccounts).Error; err != nil {
							return err
						}
					}

					// Log a totals check for operational verification after data updates.
					afterTotal, afterCount, err := calcTaskFeeTotal()
					if err != nil {
						return err
					}
					log.WithFields(log.Fields{
						"migration":            "M20260311",
						"task_fee_rows_before": beforeCount,
						"task_fee_rows_after":  afterCount,
						"task_quota_rows":      len(quotas),
						"quota_addresses":      len(quotaSums),
						"before_total":         beforeTotal.String(),
						"quota_total":          quotaTotal.String(),
						"expected_after_total": expectedFinalTotal.String(),
						"actual_after_total":   afterTotal.String(),
						"totals_match":         afterTotal.Cmp(expectedFinalTotal) == 0,
					}).Info("relay account migration totals check")

					// Rename legacy tables and columns to the relay account naming.
					m := tx.Migrator()
					if err := m.RenameTable("task_fees", "relay_accounts"); err != nil {
						return err
					}
					if err := m.RenameTable("task_fee_events", "relay_account_events"); err != nil {
						return err
					}
					if err := m.RenameColumn("withdraw_records", "task_fee_event_id", "relay_account_event_id"); err != nil {
						return err
					}
					if err := m.RenameColumn("relay_accounts", "task_fee", "balance"); err != nil {
						return err
					}
					if err := m.RenameColumn("relay_account_events", "task_fee", "amount"); err != nil {
						return err
					}
					return nil
				})
			},
			Rollback: func(tx *gorm.DB) error {
				// Reverse table and column renames in the opposite order.
				m := tx.Migrator()
				if err := m.RenameColumn("withdraw_records", "relay_account_event_id", "task_fee_event_id"); err != nil {
					return err
				}
				if err := m.RenameColumn("relay_account_events", "amount", "task_fee"); err != nil {
					return err
				}
				if err := m.RenameColumn("relay_accounts", "balance", "task_fee"); err != nil {
					return err
				}
				if err := m.RenameTable("relay_account_events", "task_fee_events"); err != nil {
					return err
				}
				if err := m.RenameTable("relay_accounts", "task_fees"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
