package tools_test

import (
	"crynux_relay/api/tools"
	"crynux_relay/blockchain"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	validator := tools.NewSignatureValidator(60)
	privateKey := "4c0883a69102937d6231471b5dbb6204fe5129617082794b75bc6dc66f6cb6b2"
	verifier := blockchain.NewSignatureVerifier()

	address, err := verifier.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to derive address: %v", err)
	}

	timestamp := time.Now().Unix()
	message := validator.GenerateConnectWalletMessage(address, timestamp)
	signature, err := validator.GenerateTestSignature(message, privateKey)
	if err != nil {
		t.Fatalf("failed to sign message: %v", err)
	}

	recoveredAddress, err := validator.ValidateAndRecoverAddress(message, signature, timestamp)
	if err != nil {
		t.Fatalf("failed to validate signature: %v", err)
	}
	t.Logf("Recovered address: %s", recoveredAddress)
	if recoveredAddress != address {
		t.Fatalf("recovered address does not match expected address")
	}
}
