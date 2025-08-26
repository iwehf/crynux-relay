package taskfee

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetTaskFeeLogsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
}

type TaskFeeLog struct {
	ID               uint               `json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	Address          string             `json:"address"`
	TaskFee          string             `json:"task_fee"`
}

type GetTaskFeeLogsResponse struct {
	response.Response
	Data []TaskFeeLog `json:"data"`
}

func GetTaskFeeLogs(c *gin.Context, in *GetTaskFeeLogsInput) (*GetTaskFeeLogsResponse, error) {
	var events []models.TaskFeeEvent
	
	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.TaskFeeEvent{}).Where("id > ?", in.StartID).Order("id ASC").Limit(in.Limit).Find(&events).Error; err != nil {
		return nil, err
	}

	// Filter out events with non-consecutive IDs
	logs := make([]TaskFeeLog, 0, len(events))
	if len(events) > 0 {
		lastID := events[0].ID

		for i, event := range events {
			if i > 0 && event.ID != lastID+1 {
				break
			}

			lastID = event.ID

			logs = append(logs, TaskFeeLog{
				ID: event.ID,
				CreatedAt: event.CreatedAt,
				Address: event.Address,
				TaskFee: event.TaskFee.String(),
			})
		}
	}

	return &GetTaskFeeLogsResponse{
		Data: logs,
	}, nil
}