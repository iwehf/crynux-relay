package blockchain

import (
	"context"
	"crynux_relay/blockchain/bindings"
	"crynux_relay/config"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"crynux_relay/models"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
)

// GetCredits retrieves the credits balance for a given address
func GetCredits(ctx context.Context, addr common.Address, network string) (*big.Int, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return nil, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return client.CreditsContractInstance.GetCredits(opts, addr)
}

// GetAllCreditAddresses retrieves all addresses that have credits
func GetAllCreditAddresses(ctx context.Context, network string) ([]common.Address, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return nil, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return client.CreditsContractInstance.GetAllCreditAddresses(opts)
}

// GetAllCredits retrieves all credit addresses and their corresponding credit amounts
func GetAllCredits(ctx context.Context, network string) ([]common.Address, []*big.Int, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return nil, nil, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return client.CreditsContractInstance.GetAllCredits(opts)
}


// BuyCreditsFor purchases credits for a specified address
func CreateCredits(ctx context.Context, addr common.Address, amount *big.Int, network string) (string, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return "", err
	}

	client.NonceMu.Lock()
	defer client.NonceMu.Unlock()

	nonce, err := client.GetNonce(ctx)
	if err != nil {
		return "", err
	}

	auth, err := client.GetAuth(ctx)
	if err != nil {
		return "", err
	}

	// Set the value to the amount being purchased (payable function)
	auth.Value = amount

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := client.Limiter.Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := client.CreditsContractInstance.CreateCredits(auth, addr, amount)
	if err != nil {
		err = client.processSendingTxError(err)
		return "", err
	}

	client.IncrementNonce()
	return tx.Hash().Hex(), nil
}

// SetStakingAddress sets the staking address (only callable by owner)
func SetStakingAddressForCredits(ctx context.Context, stakingAddress common.Address, network string) (string, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return "", err
	}

	client.NonceMu.Lock()
	defer client.NonceMu.Unlock()

	nonce, err := client.GetNonce(ctx)
	if err != nil {
		return "", err
	}

	auth, err := client.GetAuth(ctx)
	if err != nil {
		return "", err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := client.Limiter.Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := client.CreditsContractInstance.SetStakingAddress(auth, stakingAddress)
	if err != nil {
		err = client.processSendingTxError(err)
		return "", err
	}

	client.IncrementNonce()
	return tx.Hash().Hex(), nil
}

// SetAdminAddress sets the admin address (only callable by owner)
func SetAdminAddressForCredits(ctx context.Context, adminAddress common.Address, network string) (string, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return "", err
	}

	client.NonceMu.Lock()
	defer client.NonceMu.Unlock()

	nonce, err := client.GetNonce(ctx)
	if err != nil {
		return "", err
	}

	auth, err := client.GetAuth(ctx)
	if err != nil {
		return "", err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := client.Limiter.Wait(callCtx); err != nil {
		return "", err
	}
	auth.Context = callCtx
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := client.CreditsContractInstance.SetAdminAddress(auth, adminAddress)
	if err != nil {
		err = client.processSendingTxError(err)
		return "", err
	}

	client.IncrementNonce()
	return tx.Hash().Hex(), nil

}

// QueueCreateCredits queues a create credits transaction to be sent later
func QueueCreateCredits(ctx context.Context, db *gorm.DB, addr common.Address, amount *big.Int, network string) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[network]
	if !ok {
		return nil, fmt.Errorf("network %s not found", network)
	}

	abi, err := bindings.CreditsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("createCredits", addr, amount)
	if err != nil {
		return nil, err
	}
	dataStr := hexutil.Encode(data)

	transaction := &models.BlockchainTransaction{
		Network:     network,
		Type:        "Credits::createCredits",
		Status:      models.TransactionStatusPending,
		FromAddress: blockchain.Account.Address,
		ToAddress:   blockchain.Contracts.Credits,
		Value:       amount.String(),
		Data: sql.NullString{
			String: dataStr,
			Valid:  true,
		},
	}

	if err := transaction.Save(ctx, db); err != nil {
		return nil, err
	}

	return transaction, nil
}

// QueueSetStakingAddressForCredits queues a set staking address transaction to be sent later
func QueueSetStakingAddressForCredits(ctx context.Context, db *gorm.DB, stakingAddress common.Address, network string) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[network]
	if !ok {
		return nil, fmt.Errorf("network %s not found", network)
	}

	abi, err := bindings.CreditsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("setStakingAddress", stakingAddress)
	if err != nil {
		return nil, err
	}
	dataStr := hexutil.Encode(data)

	transaction := &models.BlockchainTransaction{
		Network:     network,
		Type:        "Credits::setStakingAddress",
		Status:      models.TransactionStatusPending,
		FromAddress: blockchain.Account.Address,
		ToAddress:   blockchain.Contracts.Credits,
		Value:       "0",
		Data: sql.NullString{
			String: dataStr,
			Valid:  true,
		},
	}

	if err := transaction.Save(ctx, db); err != nil {
		return nil, err
	}

	return transaction, nil
}

// QueueSetAdminAddressForCredits queues a set admin address transaction to be sent later
func QueueSetAdminAddressForCredits(ctx context.Context, db *gorm.DB, adminAddress common.Address, network string) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[network]
	if !ok {
		return nil, fmt.Errorf("network %s not found", network)
	}

	abi, err := bindings.CreditsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("setAdminAddress", adminAddress)
	if err != nil {
		return nil, err
	}
	dataStr := hexutil.Encode(data)

	transaction := &models.BlockchainTransaction{
		Network:     network,
		Type:        "Credits::setAdminAddress",
		Status:      models.TransactionStatusPending,
		FromAddress: blockchain.Account.Address,
		ToAddress:   blockchain.Contracts.Credits,
		Value:       "0",
		Data: sql.NullString{
			String: dataStr,
			Valid:  true,
		},
	}

	if err := transaction.Save(ctx, db); err != nil {
		return nil, err
	}

	return transaction, nil
}
