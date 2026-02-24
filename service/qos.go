package service

import (
	"context"
	"crynux_relay/config"
	"crynux_relay/models"
	"database/sql"
	"math"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	TASK_SCORE_REWARDS [3]uint64 = [3]uint64{10, 5, 2}
	nodeQoSScorePool   NodeQosScorePool = NodeQosScorePool{
		pool: make(map[string][]uint64),
	}
)

func getTaskQosScore(order int) uint64 {
	return TASK_SCORE_REWARDS[order]
}

type NodeQosScorePool struct {
	mu   sync.RWMutex
	pool map[string][]uint64
}

func getNodeTaskQosScore(node *models.Node, qos uint64) (float64, error) {
	poolSize := config.GetConfig().QoS.ScorePoolSize
	if poolSize == 0 {
		poolSize = 50
	}

	nodeQoSScorePool.mu.RLock()
	qosScorePool, ok := nodeQoSScorePool.pool[node.Address]
	nodeQoSScorePool.mu.RUnlock()
	if !ok {
		qosScorePool = make([]uint64, 0)
		if node.QOSScore > 0 {
			for i := 0; i < int(poolSize)-1; i++ {
				qosScorePool = append(qosScorePool, uint64(node.QOSScore))
			}
		}
	}
	qosScorePool = append(qosScorePool, qos)
	if len(qosScorePool) > int(poolSize) {
		qosScorePool = qosScorePool[1:]
	}

	nodeQoSScorePool.mu.Lock()
	nodeQoSScorePool.pool[node.Address] = qosScorePool
	nodeQoSScorePool.mu.Unlock()
	var sum uint64 = 0
	for _, score := range qosScorePool {
		sum += score
	}
	return float64(sum) / float64(len(qosScorePool)), nil
}

// ShouldPermanentKickout returns true if the node's QoS score is below the
// configured kickout threshold and the QoS score pool has enough samples.
func ShouldPermanentKickout(node *models.Node) bool {
	cfg := config.GetConfig().QoS
	poolSize := cfg.ScorePoolSize

	nodeQoSScorePool.mu.RLock()
	qosScorePool, ok := nodeQoSScorePool.pool[node.Address]
	nodeQoSScorePool.mu.RUnlock()

	// Only kick out if we have enough samples
	if !ok || uint64(len(qosScorePool)) < poolSize {
		return false
	}

	return node.QOSScore < cfg.KickoutThreshold
}

// getEffectiveHealth computes the current effective health from health fields
// using exponential decay recovery toward 1.0.
// If HealthUpdatedAt is not set, the node is considered fully healthy.
func getEffectiveHealth(healthBase float64, healthUpdatedAt sql.NullTime) float64 {
	if !healthUpdatedAt.Valid {
		return 1.0
	}
	cfg := config.GetConfig().QoS
	tau := cfg.RecoveryTauMinutes
	if tau <= 0 {
		tau = 30.0
	}

	elapsed := time.Since(healthUpdatedAt.Time).Minutes()
	if elapsed < 0 {
		elapsed = 0
	}

	// H_effective(t) = H_base + (1 - H_base) * (1 - exp(-(t - t_base) / tau))
	hEffective := healthBase + (1.0-healthBase)*(1.0-math.Exp(-elapsed/tau))
	if hEffective > 1.0 {
		hEffective = 1.0
	}
	return hEffective
}

// ApplyHealthPenalty is called when a task times out. It multiplies the
// current effective health by the configured penalty factor.
func ApplyHealthPenalty(ctx context.Context, db *gorm.DB, node *models.Node) error {
	cfg := config.GetConfig().QoS
	hEffective := getEffectiveHealth(node.HealthBase, node.HealthUpdatedAt)
	hNew := hEffective * cfg.PenaltyFactor

	return node.Update(ctx, db, map[string]interface{}{
		"health_base":       hNew,
		"health_updated_at": sql.NullTime{Time: time.Now(), Valid: true},
	})
}

// ApplyHealthBoost is called on successful task completion. It adds the
// configured success boost to the current effective health, capped at 1.0.
func ApplyHealthBoost(ctx context.Context, db *gorm.DB, node *models.Node) error {
	cfg := config.GetConfig().QoS
	hEffective := getEffectiveHealth(node.HealthBase, node.HealthUpdatedAt)
	hNew := hEffective + cfg.SuccessBoost
	if hNew > 1.0 {
		hNew = 1.0
	}

	return node.Update(ctx, db, map[string]interface{}{
		"health_base":       hNew,
		"health_updated_at": sql.NullTime{Time: time.Now(), Valid: true},
	})
}

// CalculateQosScore returns the node's current QoS score (0 to 1),
// combining long-term performance and short-term reliability.
// Returns 0 if the node should be hard-excluded.
func CalculateQosScore(qosScore float64, healthBase float64, healthUpdatedAt sql.NullTime) float64 {
	h := getEffectiveHealth(healthBase, healthUpdatedAt)
	cfg := config.GetConfig().QoS
	if h < cfg.ExcludeThreshold {
		return 0 // hard exclusion
	}
	// globalMaxQosScore is defined in selecting_prob.go, but visible here since it's the same package
	qosLong := qosScore / globalMaxQosScore
	if qosLong == 0 {
		qosLong = 0.5 // default for new nodes
	}
	return qosLong * h
}
