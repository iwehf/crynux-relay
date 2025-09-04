package incentive

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetNodeDailyIncentiveInput struct {
	Address  string `path:"address" json:"address" validate:"required"`
	Page     int    `query:"page" json:"page" default:"1"`
	PageSize int    `query:"page_size" json:"page_size" default:"10"`
}

type NodeDailyIncentive struct {
	Timestamp uint64  `json:"timestamp"`
	Amount    float64 `json:"amount"`
}

type GetNodeDailyIncentiveOutput struct {
	response.Response
	Data []NodeDailyIncentive `json:"data"`
}

func GetNodeDailyIncentive(c *gin.Context, input *GetNodeDailyIncentiveInput) (*GetNodeDailyIncentiveOutput, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 10
	}

	offset := (input.Page - 1) * input.PageSize
	limit := input.PageSize
	
	var nodeIncentives []models.NodeIncentive

	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.NodeIncentive{}).Where("node_address = ?", input.Address).Order("time DESC").Offset(offset).Limit(limit).Find(&nodeIncentives).Error; err != nil {
		return nil, response.NewExceptionResponse(err)
	}

	result := make([]NodeDailyIncentive, len(nodeIncentives))
	for i, nodeIncentive := range nodeIncentives {
		result[i] = NodeDailyIncentive{
			Timestamp: uint64(nodeIncentive.Time.Unix()),
			Amount: nodeIncentive.Incentive,
		}
	}

	return &GetNodeDailyIncentiveOutput{Data: result}, nil
}