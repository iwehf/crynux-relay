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

func GetBenefitAddress(ctx context.Context, nodeAddress common.Address, network string) (common.Address, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return common.Address{}, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	ba, err := client.BenefitAddressContractInstance.GetBenefitAddress(opts, nodeAddress)
	if err != nil {
		return common.Address{}, err
	}
	if ba.Big().Cmp(big.NewInt(0)) == 0 {
		return nodeAddress, nil
	} else {
		return ba, nil
	}
}

func SetBenefitAddress(ctx context.Context, benefitAddress common.Address, network string) (string, error) {
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

	tx, err := client.BenefitAddressContractInstance.SetBenefitAddress(auth, benefitAddress)
	if err != nil {
		err = client.processSendingTxError(err)
		return "", err
	}

	client.IncrementNonce()
	return tx.Hash().Hex(), nil
}

// QueueSetBenefitAddress queues a set benefit address transaction to be sent later
func QueueSetBenefitAddress(ctx context.Context, db *gorm.DB, benefitAddress common.Address, network string) (*models.BlockchainTransaction, error) {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[network]
	if !ok {
		return nil, fmt.Errorf("network %s not found", network)
	}

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
		Network:     network,
		Type:        "BenefitAddress::setBenefitAddress",
		Status:      models.TransactionStatusPending,
		FromAddress: blockchain.Account.Address,
		ToAddress:   blockchain.Contracts.BenefitAddress,
		Value:       "0",
		Data: sql.NullString{
			String: dataStr,
			Valid:  true,
		},
		MaxRetries: blockchain.MaxRetries,
	}

	if err := transaction.Save(ctx, db); err != nil {
		return nil, err
	}

	return transaction, nil
}
