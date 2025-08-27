package withdraw

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/api/v1/validate"
	"crynux_relay/config"
	"crynux_relay/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FufillWithdrawRequestInput struct {
	ID uint `path:"id" json:"id" description:"Withdraw request ID"`
}

type FufillWithdrawRequestInputWithSignature struct {
	FufillWithdrawRequestInput
	Timestamp int64  `form:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `form:"signature" json:"signature" description:"Signature" validate:"required"`
}

type FufillWithdrawRequestResponse struct {
	response.Response
}

func FufillWithdrawRequest(c *gin.Context, in *FufillWithdrawRequestInputWithSignature) (*FufillWithdrawRequestResponse, error) {
	match, address, err := validate.ValidateSignature(in.FufillWithdrawRequestInput, in.Timestamp, in.Signature)

	if err != nil || !match {
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	if address != config.GetConfig().Wallet.Address {
		validationErr := response.NewValidationErrorResponse("address", "Invalid address")
		return nil, validationErr
	}

	if err := service.FufillWithdraw(c.Request.Context(), config.GetDB(), in.ID); err != nil {
		if errors.Is(err, service.ErrWithdrawRequestNotPending) || errors.Is(err, service.ErrWithdrawRequestNotProcessedLocally) {
			validationErr := response.NewValidationErrorResponse("id", err.Error())
			return nil, validationErr
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			validationErr := response.NewValidationErrorResponse("id", "Withdraw request not found")
			return nil, validationErr
		}
		return nil, response.NewExceptionResponse(err)
	}

	return &FufillWithdrawRequestResponse{}, nil
}
