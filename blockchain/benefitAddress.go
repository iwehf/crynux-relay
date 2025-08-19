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

func GetBenefitAddress(ctx context.Context, nodeAddress common.Address) (common.Address, error) {
	benefitAddressContractInstance := GetBenefitAddressContractInstance()
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	return benefitAddressContractInstance.GetBenefitAddress(opts, nodeAddress)
}

func SetBenefitAddress(ctx context.Context, benefitAddress common.Address) (string, error) {
	benefitAddressContractInstance := GetBenefitAddressContractInstance()

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

	tx, err := benefitAddressContractInstance.SetBenefitAddress(auth, benefitAddress)
	if err != nil {
		return "", err
	}

	addNonce(nonce)
	return tx.Hash().Hex(), nil
}

// QueueSetBenefitAddress queues a set benefit address transaction to be sent later
func QueueSetBenefitAddress(ctx context.Context, db *gorm.DB, benefitAddress common.Address) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	address := appConfig.Blockchain.Account.Address

	abi, err := bindings.BenefitAddressMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("setBenefitAddress", benefitAddress)
	if err != nil {
		return nil, err
	}
	dataStr := hexutil.Encode(data)

	transaction := &models.BlockchainTransaction{
		Type:        "BenefitAddress::setBenefitAddress",
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
