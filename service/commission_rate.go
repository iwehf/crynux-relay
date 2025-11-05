package service

import "sync"

var globalCommissionRateCache = NewCommissionRateCache()

type CommissionRateCache struct {
	sync.RWMutex
	commissionRates map[string]uint8
}

func (c *CommissionRateCache) set(nodeAddress string, rate uint8) {
	c.Lock()
	defer c.Unlock()

	c.commissionRates[nodeAddress] = rate
}

func (c *CommissionRateCache) get(nodeAddress string) uint8 {
	c.RLock()
	defer c.RUnlock()

	if rate, ok := c.commissionRates[nodeAddress]; ok {
		return rate
	}
	return 0
}

func NewCommissionRateCache() *CommissionRateCache {
	return &CommissionRateCache{
		commissionRates: make(map[string]uint8),
	}
}

func SetCommissionRate(nodeAddress string, rate uint8) {
	globalCommissionRateCache.set(nodeAddress, rate)
}

func GetCommissionRate(nodeAddress string) uint8 {
	return globalCommissionRateCache.get(nodeAddress)
}
