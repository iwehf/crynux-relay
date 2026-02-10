package service

import (
	"context"
	"crynux_relay/config"
	"crynux_relay/models"
	"database/sql"
	"math"
	"time"

	"gorm.io/gorm"
)

// GetEffectiveHealth computes the current effective health from HealthBase and
// HealthUpdatedAt using exponential decay recovery toward 1.0.
// If HealthUpdatedAt is not set, the node is considered fully healthy.
func GetEffectiveHealth(node *models.Node) float64 {
	if !node.HealthUpdatedAt.Valid {
		return 1.0
	}
	cfg := config.GetConfig().NodeHealth
	tau := cfg.RecoveryTauMinutes
	if tau <= 0 {
		tau = 30.0
	}

	elapsed := time.Since(node.HealthUpdatedAt.Time).Minutes()
	if elapsed < 0 {
		elapsed = 0
	}

	// H_effective(t) = H_base + (1 - H_base) * (1 - exp(-(t - t_base) / tau))
	hBase := node.HealthBase
	hEffective := hBase + (1.0-hBase)*(1.0-math.Exp(-elapsed/tau))
	if hEffective > 1.0 {
		hEffective = 1.0
	}
	return hEffective
}

// ApplyHealthPenalty is called when a task times out. It multiplies the
// current effective health by the configured penalty factor.
func ApplyHealthPenalty(ctx context.Context, db *gorm.DB, node *models.Node) error {
	cfg := config.GetConfig().NodeHealth
	hEffective := GetEffectiveHealth(node)
	hNew := hEffective * cfg.PenaltyFactor

	now := time.Now()
	node.HealthBase = hNew
	node.HealthUpdatedAt = sql.NullTime{Time: now, Valid: true}

	return node.Update(ctx, db, map[string]interface{}{
		"health_base":       hNew,
		"health_updated_at": sql.NullTime{Time: now, Valid: true},
	})
}

// ApplyHealthBoost is called on successful task completion. It adds the
// configured success boost to the current effective health, capped at 1.0.
func ApplyHealthBoost(ctx context.Context, db *gorm.DB, node *models.Node) error {
	cfg := config.GetConfig().NodeHealth
	hEffective := GetEffectiveHealth(node)
	hNew := hEffective + cfg.SuccessBoost
	if hNew > 1.0 {
		hNew = 1.0
	}

	now := time.Now()
	node.HealthBase = hNew
	node.HealthUpdatedAt = sql.NullTime{Time: now, Valid: true}

	return node.Update(ctx, db, map[string]interface{}{
		"health_base":       hNew,
		"health_updated_at": sql.NullTime{Time: now, Valid: true},
	})
}
