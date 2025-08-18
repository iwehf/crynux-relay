package blockchain

import (
	"context"
	"crynux_relay/config"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// GetCredits retrieves the credits balance for a given address
func GetCredits(ctx context.Context, addr common.Address) (*big.Int, error) {
	creditsContractInstance := GetCreditsContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return creditsContractInstance.GetCredits(opts, addr)
}

// GetAllCreditAddresses retrieves all addresses that have credits
func GetAllCreditAddresses(ctx context.Context) ([]common.Address, error) {
	creditsContractInstance := GetCreditsContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return creditsContractInstance.GetAllCreditAddresses(opts)
}

// GetAllCredits retrieves all credit addresses and their corresponding credit amounts
func GetAllCredits(ctx context.Context) ([]common.Address, []*big.Int, error) {
	creditsContractInstance := GetCreditsContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return creditsContractInstance.GetAllCredits(opts)
}

// GetOwner retrieves the owner address of the Credits contract
func GetOwner(ctx context.Context) (common.Address, error) {
	creditsContractInstance := GetCreditsContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return creditsContractInstance.Owner(opts)
}

// BuyCredits purchases credits for the caller (msg.sender)
func BuyCredits(ctx context.Context, amount *big.Int) (string, error) {
	creditsContractInstance := GetCreditsContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	// Set the value to the amount being purchased (payable function)
	auth.Value = amount

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

	tx, err := creditsContractInstance.BuyCredits(auth, amount)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// BuyCreditsFor purchases credits for a specified address
func BuyCreditsFor(ctx context.Context, addr common.Address, amount *big.Int) (string, error) {
	creditsContractInstance := GetCreditsContractInstance()

	appConfig := config.GetConfig()
	address := common.HexToAddress(appConfig.Blockchain.Account.Address)
	privkey := appConfig.Blockchain.Account.PrivateKey

	txMutex.Lock()
	defer txMutex.Unlock()

	auth, err := GetAuth(ctx, address, privkey)
	if err != nil {
		return "", err
	}

	// Set the value to the amount being purchased (payable function)
	auth.Value = amount

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

	tx, err := creditsContractInstance.BuyCreditsFor(auth, addr, amount)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// SetStakingAddress sets the staking address (only callable by owner)
func SetStakingAddress(ctx context.Context, stakingAddress common.Address) (string, error) {
	creditsContractInstance := GetCreditsContractInstance()

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

	tx, err := creditsContractInstance.SetStakingAddress(auth, stakingAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}
