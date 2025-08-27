package client

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/api/tools"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ConnectWalletInput defines the input parameters for wallet connection
type ConnectWalletInput struct {
	Address   string `form:"address" json:"address" description:"Wallet address" validate:"required"`
	Signature string `form:"signature" json:"signature" description:"Signature" validate:"required"`
	Timestamp int64  `form:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
}

// ConnectWalletOutput defines the output structure for wallet connection
type ConnectWalletOutput struct {
	Token string `json:"token" description:"JWT token for authentication"`
	ExpiresAt int64 `json:"expires_at" description:"JWT token expiration time"`
}

// ConnectWalletResponse defines the API response structure
type ConnectWalletResponse struct {
	response.Response
	Data *ConnectWalletOutput `json:"data"`
}

// ConnectWallet handles wallet connection requests
// It verifies the signature and creates or retrieves a client associated with the wallet address
func ConnectWallet(c *gin.Context, in *ConnectWalletInput) (*ConnectWalletResponse, error) {
	// Generate the correct signature message format
	message := tools.GenerateConnectWalletMessage(in.Address, in.Timestamp)

	// Recover address from signature using the proper signature validation tools
	signerAddress, err := tools.ValidateAndRecover(message, in.Signature, in.Timestamp)
	if err != nil {
		log.Debugf("Error validating signature: %v", err)
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	// Verify that the signer address matches the provided address
	if signerAddress != in.Address {
		validationErr := response.NewValidationErrorResponse("address", "Signature address mismatch")
		return nil, validationErr
	}

	// Generate JWT token for the authenticated user
	token, exp, err := tools.GenerateToken(in.Address)
	if err != nil {
		log.Errorf("Error generating JWT token: %v", err)
		return nil, err
	}

	// Create response with token
	output := &ConnectWalletOutput{
		Token:     token,
		ExpiresAt: exp.Unix(),
	}

	return &ConnectWalletResponse{
		Data:     output,
	}, nil
}
