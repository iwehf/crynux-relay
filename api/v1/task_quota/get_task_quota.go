package taskquota

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/service"

	"github.com/gin-gonic/gin"
)

type GetTaskQuotaInput struct {
	Address string `path:"address" json:"address" description:"Address of account"`
}

type GetTaskQuotaResponse struct {
	response.Response
	Data models.BigInt `json:"data"`
}

func GetTaskQuota(c *gin.Context, in *GetTaskQuotaInput) (*GetTaskQuotaResponse, error) {
	taskQuota, err := service.GetTaskQuota(c.Request.Context(), config.GetDB(), in.Address)
	if err != nil {
		return nil, err
	}
	return &GetTaskQuotaResponse{
		Data: models.BigInt{Int: *taskQuota},
	}, nil
}
