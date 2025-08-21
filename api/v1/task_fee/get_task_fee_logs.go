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

type GetTaskFeeLogsResponse struct {
	response.Response
	Data []models.TaskFeeEvent `json:"data"`
}

func GetTaskFeeLogs(c *gin.Context, in *GetTaskFeeLogsInput) (*GetTaskFeeLogsResponse, error) {
	var events []models.TaskFeeEvent
	
	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.TaskFeeEvent{}).Where("id > ?", in.StartID).Order("id ASC").Limit(in.Limit).Find(&events).Error; err != nil {
		return nil, err
	}

	// Filter out events with non-consecutive IDs
	if len(events) > 0 {
		lastID := events[0].ID
		end := 1
		
		for ; end < len(events); end++ {
			if events[end].ID == lastID+1 {
				lastID = events[end].ID
			} else {
				break
			}
		}
		events = events[:end]
	}
	
	return &GetTaskFeeLogsResponse{
		Data: events,
	}, nil
}