package service

import (
	"math/big"
	"sync"
)

var globalUserStakingCache = newUserStakingCache()

type userStakingCache struct {
	sync.RWMutex
	nodeUserStakings map[string]map[string]*big.Int
	userStakeAmount  map[string]*big.Int
	nodeStakeAmount  map[string]*big.Int
}

func newUserStakingCache() *userStakingCache {
	return &userStakingCache{
		nodeUserStakings: make(map[string]map[string]*big.Int),
		userStakeAmount:  make(map[string]*big.Int),
		nodeStakeAmount:  make(map[string]*big.Int),
	}
}

func (c *userStakingCache) update(userAddress, nodeAddress string, amount *big.Int) {
	c.Lock()
	defer c.Unlock()

	oldAmount := big.NewInt(0)
	if userStakings, ok := c.nodeUserStakings[nodeAddress]; ok {
		if val, ok1 := userStakings[userAddress]; ok1 {
			oldAmount.Set(val)
		}
	}
	c.userStakeAmount[userAddress].Sub(c.userStakeAmount[userAddress], oldAmount)
	c.userStakeAmount[userAddress].Add(c.userStakeAmount[userAddress], amount)
	c.nodeStakeAmount[nodeAddress].Sub(c.nodeStakeAmount[nodeAddress], oldAmount)
	c.nodeStakeAmount[nodeAddress].Add(c.nodeStakeAmount[nodeAddress], amount)
	c.nodeUserStakings[nodeAddress][userAddress] = big.NewInt(0).Set(amount)
}

func (c *userStakingCache) unstake(userAddress, nodeAddress string) {
	c.Lock()
	defer c.Unlock()

	oldAmount := big.NewInt(0)
	if userStakings, ok := c.nodeUserStakings[nodeAddress]; ok {
		if val, ok1 := userStakings[userAddress]; ok1 {
			oldAmount.Set(val)
			delete(c.nodeUserStakings[nodeAddress], userAddress)
		}
		if len(c.nodeUserStakings[nodeAddress]) == 0 {
			delete(c.nodeUserStakings, nodeAddress)
		}
	}
	c.userStakeAmount[userAddress].Sub(c.userStakeAmount[userAddress], oldAmount)
	c.nodeStakeAmount[nodeAddress].Sub(c.nodeStakeAmount[nodeAddress], oldAmount)
}

func (c *userStakingCache) removeNode(nodeAddress string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.nodeUserStakings[nodeAddress]; ok {
		for userAddress, amount := range c.nodeUserStakings[nodeAddress] {
			c.userStakeAmount[userAddress].Sub(c.userStakeAmount[userAddress], amount)
		}
		delete(c.nodeUserStakings, nodeAddress)
		c.nodeStakeAmount[nodeAddress] = big.NewInt(0)
	}
}

func (c *userStakingCache) getUserStakeAmount(userAddress string) *big.Int {
	c.RLock()
	defer c.RUnlock()
	res := big.NewInt(0)
	if val, ok := c.userStakeAmount[userAddress]; ok {
		res.Set(val)
	}
	return res
}

func (c *userStakingCache) getNodeStakeAmount(nodeAddress string) *big.Int {
	c.RLock()
	defer c.RUnlock()
	res := big.NewInt(0)
	if val, ok := c.nodeStakeAmount[nodeAddress]; ok {
		res.Set(val)
	}
	return res
}

func (c *userStakingCache) getUserStakingsOfNode(nodeAddress string) (map[string]*big.Int, *big.Int) {
	c.RLock()
	defer c.RUnlock()

	userStakeAmount := c.userStakeAmount[nodeAddress]
	if res, ok := c.nodeUserStakings[nodeAddress]; ok {
		return res, userStakeAmount
	} else {
		return make(map[string]*big.Int), userStakeAmount
	}
}

func UpdateUserStaking(userAddress, nodeAddress string, amount *big.Int) {
	globalUserStakingCache.update(userAddress, nodeAddress, amount)
}

func UnstakeUserStaking(userAddress, nodeAddress string) {
	globalUserStakingCache.unstake(userAddress, nodeAddress)
}

func RemoveNodeUserStaking(nodeAddress string) {
	globalUserStakingCache.removeNode(nodeAddress)
}

func GetUserStakeAmountOfUser(userAddress string) *big.Int {
	return globalUserStakingCache.getUserStakeAmount(userAddress)
}

func GetUserStakeAmountOfNode(nodeAddress string) *big.Int {
	return globalUserStakingCache.getNodeStakeAmount(nodeAddress)
}

func GetUserStakingsOfNode(nodeAddress string) (map[string]*big.Int, *big.Int) {
	return globalUserStakingCache.getUserStakingsOfNode(nodeAddress)
}
