package withdraw

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type GetWithdrawRequestsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
	Status  *models.WithdrawStatus `query:"status" json:"status" description:"Status"`
}

type GetWithdrawRequestsResponse struct {
	response.Response
	Data []models.WithdrawRecord `json:"data"`
}

func GetWithdrawRequests(c *gin.Context, in *GetWithdrawRequestsInput) (*GetWithdrawRequestsResponse, error) {
	var records []models.WithdrawRecord
	
	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbi := config.GetDB().WithContext(dbCtx).Model(&models.WithdrawRecord{}).Where("id > ?", in.StartID)
	if in.Status != nil {
		dbi = dbi.Where("status = ?", *in.Status)
	}

	if err := dbi.Order("id ASC").Limit(in.Limit).Find(&records).Error; err != nil {
		return nil, err
	}

	return &GetWithdrawRequestsResponse{
		Data: records,
	}, nil
}