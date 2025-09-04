package tools_test

import (
	"crynux_relay/api/tools"
	"testing"
)

func TestSignature(t *testing.T) {
	
	signature := "0x3076ab25f2a00baac115a8a8e44dd915c871a955e89899d5fd70a65af49e4100187cfc721489722cc552c25bc30ce7f7562aaa759500b0a7d1e5d51d31e10ff31b"
	address := "0xbab3d0dfe3f631a8aca3c8e25633317a7d4c9dcb"
	timestamp := int64(1756973246)

	message := tools.GenerateConnectWalletMessage(address, timestamp)
	recoveredAddress, err := tools.ValidateAndRecover(message, signature, timestamp)
	if err != nil {
		t.Fatalf("Failed to validate signature: %v", err)
	}
	t.Logf("Recovered address: %s", recoveredAddress)
	if recoveredAddress != address {
		t.Fatalf("Recovered address does not match expected address")
	}
}
