package tools

import (
	"crynux_relay/blockchain"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// SignatureValidator handles signature validation for API key management
type SignatureValidator struct {
	verifier       *blockchain.SignatureVerifier
	timeoutSeconds int64
}

// NewSignatureValidator creates a new signature validator
func NewSignatureValidator(timeoutSeconds int64) *SignatureValidator {
	return &SignatureValidator{
		verifier:       blockchain.NewSignatureVerifier(),
		timeoutSeconds: timeoutSeconds,
	}
}

// ValidateTimestamp validates if the timestamp is within acceptable range
func (sv *SignatureValidator) ValidateTimestamp(timestamp int64) error {
	now := time.Now().Unix()

	// Allow some clock drift - accept timestamps from 1 minute ago to 1 minute in the future
	const clockDriftSeconds = 60

	if timestamp < now-sv.timeoutSeconds-clockDriftSeconds {
		return errors.New("signature timestamp too old")
	}

	if timestamp > now+clockDriftSeconds {
		return errors.New("signature timestamp too far in the future")
	}

	return nil
}

// GenerateConnectWalletMessage generates the standard message for wallet connection
func (sv *SignatureValidator) GenerateConnectWalletMessage(address string, timestamp int64) string {
	return sv.verifier.GenerateSignMessage("Connect Wallet", address, timestamp)
}

// GenerateConnectWalletMessage generates the standard message for wallet connection
func (sv *SignatureValidator) GenerateWithdrawMessage(address string, amount *big.Int, benefitAddress string, network string, timestamp int64) string {
	message := fmt.Sprintf("Withdraw %s from %s to %s on %s", amount.String(), address, benefitAddress, network)
	return sv.verifier.GenerateSignMessage(message, address, timestamp)
}

// ValidateAndRecoverAddress validates signature and returns the recovered address
func (sv *SignatureValidator) ValidateAndRecoverAddress(message, signature string, timestamp int64) (string, error) {
	// Validate timestamp
	if err := sv.ValidateTimestamp(timestamp); err != nil {
		return "", err
	}

	// Validate signature format
	if err := sv.verifier.ValidateSignatureFormat(signature); err != nil {
		return "", err
	}

	// Recover address from signature
	recoveredAddress, err := sv.verifier.RecoverAddress(message, signature)
	if err != nil {
		return "", err
	}

	return recoveredAddress, nil
}

// Example usage for testing signature functionality
func (sv *SignatureValidator) GenerateTestSignature(message, privateKey string) (string, error) {
	return sv.verifier.SignMessage(message, privateKey)
}

// Global signature validator instance (can be configured via config)
var DefaultSignatureValidator = NewSignatureValidator(60) // 1 minute timeout

// GenerateConnectWalletMessage generates the standard message for wallet connection
func GenerateConnectWalletMessage(address string, timestamp int64) string {
	return DefaultSignatureValidator.GenerateConnectWalletMessage(address, timestamp)
}

// GenerateWithdrawMessage generates the standard message for withdraw
func GenerateWithdrawMessage(address string, amount *big.Int, benefitAddress string, network string, timestamp int64) string {
	return DefaultSignatureValidator.GenerateWithdrawMessage(address, amount, benefitAddress, network, timestamp)
}

// ValidateAndRecover validates signature and returns recovered address with default settings
func ValidateAndRecover(message, signature string, timestamp int64) (string, error) {
	return DefaultSignatureValidator.ValidateAndRecoverAddress(message, signature, timestamp)
}
