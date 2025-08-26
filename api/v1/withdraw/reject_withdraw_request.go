package withdraw

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/service"

	"github.com/gin-gonic/gin"
)

type RejectWithdrawRequestInput struct {
	ID uint `path:"id" json:"id" description:"Withdraw request ID"`
}

type RejectWithdrawRequestResponse struct {
	response.Response
}

func RejectWithdrawRequest(c *gin.Context, in *RejectWithdrawRequestInput) (*RejectWithdrawRequestResponse, error) {
	if err := service.RejectWithdraw(c.Request.Context(), config.GetDB(), in.ID); err != nil {
		return nil, response.NewExceptionResponse(err)
	}
	return &RejectWithdrawRequestResponse{}, nil
}
