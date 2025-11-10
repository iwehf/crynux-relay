package service

import "sync"

var globalDelegatorShareCache = NewDelegatorShareCache()

type DelegatorShareCache struct {
	sync.RWMutex
	delegatorShares map[string]uint8
}

func (c *DelegatorShareCache) set(nodeAddress string, share uint8) {
	c.Lock()
	defer c.Unlock()

	c.delegatorShares[nodeAddress] = share
}

func (c *DelegatorShareCache) get(nodeAddress string) uint8 {
	c.RLock()
	defer c.RUnlock()

	if share, ok := c.delegatorShares[nodeAddress]; ok {
		return share
	}
	return 0
}

func NewDelegatorShareCache() *DelegatorShareCache {
	return &DelegatorShareCache{
		delegatorShares: make(map[string]uint8),
	}
}

func SetDelegatorShare(nodeAddress string, share uint8) {
	globalDelegatorShareCache.set(nodeAddress, share)
}

func GetDelegatorShare(nodeAddress string) uint8 {
	return globalDelegatorShareCache.get(nodeAddress)
}
