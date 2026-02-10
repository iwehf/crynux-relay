package service

import (
	"crynux_relay/config"
	"crynux_relay/models"
	"sync"
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
	if poolSize == 0 {
		poolSize = 50
	}

	nodeQoSScorePool.mu.RLock()
	qosScorePool, ok := nodeQoSScorePool.pool[node.Address]
	nodeQoSScorePool.mu.RUnlock()

	// Only kick out if we have enough samples
	if !ok || uint64(len(qosScorePool)) < poolSize {
		return false
	}

	return node.QOSScore < cfg.KickoutThreshold
}
