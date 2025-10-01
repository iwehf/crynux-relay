package utils

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func generateMACBytes(data []byte, secretKey string) []byte {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(data)
	return mac.Sum(nil)
}


func GenerateMAC(data []byte, secretKey string) string {
	return hexutil.Encode(generateMACBytes(data, secretKey))
}

func VerifyMAC(data []byte, secretKey string, mac string) bool {
	macBytes, err := hexutil.Decode(mac)
	if err != nil {
		return false
	}
	return hmac.Equal(macBytes, generateMACBytes(data, secretKey))
}