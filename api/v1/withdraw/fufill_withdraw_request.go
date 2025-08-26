package withdraw

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/service"

	"github.com/gin-gonic/gin"
)

type FufillWithdrawRequestInput struct {
	ID uint `path:"id" json:"id" description:"Withdraw request ID"`
}

type FufillWithdrawRequestResponse struct {
	response.Response
}

func FufillWithdrawRequest(c *gin.Context, in *FufillWithdrawRequestInput) (*FufillWithdrawRequestResponse, error) {
	if err := service.FufillWithdraw(c.Request.Context(), config.GetDB(), in.ID); err != nil {
		return nil, response.NewExceptionResponse(err)
	}

	return &FufillWithdrawRequestResponse{}, nil
}
