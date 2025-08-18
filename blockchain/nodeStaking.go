package blockchain

import (
	"context"
	"crynux_relay/blockchain/bindings"
	"crynux_relay/config"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// GetMinStakeAmount gets the minimum stake amount
func GetMinStakeAmount(ctx context.Context) (*big.Int, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return nodeStakingContractInstance.GetMinStakeAmount(opts)
}

// GetStakingInfo gets the staking information for a specific node
func GetStakingInfo(ctx context.Context, nodeAddress common.Address) (bindings.NodeStakingStakingInfo, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return nodeStakingContractInstance.GetStakingInfo(opts, nodeAddress)
}

// GetAllNodeAddresses gets all staked node addresses
func GetAllNodeAddresses(ctx context.Context) ([]common.Address, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return nodeStakingContractInstance.GetAllNodeAddresses(opts)
}

// GetNodeStakingOwner gets the contract owner address
func GetNodeStakingOwner(ctx context.Context) (common.Address, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return nodeStakingContractInstance.Owner(opts)
}

// Stake stakes tokens
func Stake(ctx context.Context, stakedAmount *big.Int) (string, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	// Set the stake amount as the transaction value
	auth.Value = stakedAmount

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := getLimiter().Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	nonce, err := getNonce(callCtx, address)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := nodeStakingContractInstance.Stake(auth, stakedAmount)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// Unstake removes the stake
func Unstake(ctx context.Context, nodeAddress common.Address) (string, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := getLimiter().Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	nonce, err := getNonce(callCtx, address)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := nodeStakingContractInstance.Unstake(auth, nodeAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// SetAdminAddress sets the admin address (owner only)
func SetAdminAddress(ctx context.Context, adminAddress common.Address) (string, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := getLimiter().Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	nonce, err := getNonce(callCtx, address)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := nodeStakingContractInstance.SetAdminAddress(auth, adminAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// SlashStaking slashes the node's stake (owner or admin only)
func SlashStaking(ctx context.Context, nodeAddress common.Address) (string, error) {
	nodeStakingContractInstance := GetNodeStakingContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := getLimiter().Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	nonce, err := getNonce(callCtx, address)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := nodeStakingContractInstance.SlashStaking(auth, nodeAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}
