package credits

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/api/v1/validate"
	"crynux_relay/config"
	"crynux_relay/service"
	"math/big"

	"github.com/gin-gonic/gin"
)

type CreateCreditsInput struct {
	Address string `path:"address" json:"address" validate:"required"`
	Amount  string `json:"amount" validate:"required"`
	Network string `json:"network" validate:"required"`
}

type CreateCreditsInputWithSignature struct {
	CreateCreditsInput
	Timestamp int64  `json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `json:"signature" description:"Signature" validate:"required"`
}

type CreateCreditsData struct {
	RequestId uint `json:"request_id"`
}

type CreateCreditsResponse struct {
	response.Response
	Data *CreateCreditsData `json:"data" description:"The data of the create credits"`
}

func CreateCredits(c *gin.Context, in *CreateCreditsInputWithSignature) (*CreateCreditsResponse, error) {
	match, address, err := validate.ValidateSignature(in.CreateCreditsInput, in.Timestamp, in.Signature)

	if err != nil || !match {
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	if address != config.GetConfig().Credits.Address {
		validationErr := response.NewValidationErrorResponse("address", "Invalid address")
		return nil, validationErr
	}

	amount, ok := big.NewInt(0).SetString(in.Amount, 10)
	if !ok {
		return nil, response.NewValidationErrorResponse("amount", "Invalid amount")
	}

	record, err := service.CreateCredits(c.Request.Context(), config.GetDB(), in.Address, amount, in.Network)
	if err != nil {
		return nil, response.NewExceptionResponse(err)
	}

	return &CreateCreditsResponse{
		Data: &CreateCreditsData{
			RequestId: record.ID,
		},
	}, nil
}
