package blockchain

import (
	"context"
	"crynux_relay/config"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
