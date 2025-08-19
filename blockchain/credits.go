package blockchain

import (
	"context"
	"crynux_relay/blockchain/bindings"
	"crynux_relay/config"
	"database/sql"
	"math/big"
	"time"

	"crynux_relay/models"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
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

// BuyCreditsFor purchases credits for a specified address
func CreateCredits(ctx context.Context, addr common.Address, amount *big.Int) (string, error) {
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

	tx, err := creditsContractInstance.CreateCredits(auth, addr, amount)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// SetStakingAddress sets the staking address (only callable by owner)
func SetStakingAddressForCredits(ctx context.Context, stakingAddress common.Address) (string, error) {
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

// SetAdminAddress sets the admin address (only callable by owner)
func SetAdminAddressForCredits(ctx context.Context, adminAddress common.Address) (string, error) {
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

	tx, err := creditsContractInstance.SetAdminAddress(auth, adminAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil

}

// QueueCreateCredits queues a create credits transaction to be sent later
func QueueCreateCredits(ctx context.Context, db *gorm.DB, addr common.Address, amount *big.Int) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	address := appConfig.Blockchain.Account.Address

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
		Type:        "Credits::createCredits",
		Status:      models.TransactionStatusPending,
		FromAddress: address,
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
func QueueSetStakingAddressForCredits(ctx context.Context, db *gorm.DB, stakingAddress common.Address) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	address := appConfig.Blockchain.Account.Address

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
		Type:        "Credits::setStakingAddress",
		Status:      models.TransactionStatusPending,
		FromAddress: address,
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
func QueueSetAdminAddressForCredits(ctx context.Context, db *gorm.DB, adminAddress common.Address) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	address := appConfig.Blockchain.Account.Address

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
		Type:        "Credits::setAdminAddress",
		Status:      models.TransactionStatusPending,
		FromAddress: address,
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
