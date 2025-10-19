package taskfee

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/api/v1/validate"
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type GetTaskFeeLogsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
}

type GetTaskFeeLogsInputWithSignature struct {
	GetTaskFeeLogsInput
	Timestamp int64  `query:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `query:"signature" json:"signature" description:"Signature" validate:"required"`
}

type TaskFeeLog struct {
	ID        uint                    `json:"id"`
	CreatedAt uint64                  `json:"created_at"`
	Address   string                  `json:"address"`
	TaskFee   string                  `json:"task_fee"`
	Type      models.TaskFeeEventType `json:"type"`
}

type GetTaskFeeLogsResponse struct {
	response.Response
	Data []TaskFeeLog `json:"data"`
}

func GetTaskFeeLogs(c *gin.Context, in *GetTaskFeeLogsInputWithSignature) (*GetTaskFeeLogsResponse, error) {
	match, address, err := validate.ValidateSignature(in.GetTaskFeeLogsInput, in.Timestamp, in.Signature)

	if err != nil || !match {
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	if address != config.GetConfig().Withdraw.Address {
		validationErr := response.NewValidationErrorResponse("address", "Invalid address")
		return nil, validationErr
	}

	var events []models.TaskFeeEvent

	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.TaskFeeEvent{}).Where("id > ?", in.StartID).Where("status != ?", models.TaskFeeEventStatusInvalid).Order("id ASC").Limit(in.Limit).Find(&events).Error; err != nil {
		return nil, err
	}

	// Filter out events with non-consecutive IDs
	logs := make([]TaskFeeLog, 0, len(events))
	if len(events) > 0 {
		lastID := events[0].ID
		appConfig := config.GetConfig()
		invalidIDs := make([]uint, 0)
		for i, event := range events {
			if event.Status == models.TaskFeeEventStatusPending || (i > 0 && event.ID != lastID+1) {
				break
			}

			lastID = event.ID

			if event.Status == models.TaskFeeEventStatusProcessed {
				valid := utils.VerifyMAC([]byte(event.Reason), appConfig.MAC.SecretKey, event.MAC)
				if !valid {
					invalidIDs = append(invalidIDs, event.ID)
					continue
				}
				logs = append(logs, TaskFeeLog{
					ID:        event.ID,
					CreatedAt: uint64(event.CreatedAt.Unix()),
					Address:   event.Address,
					TaskFee:   event.TaskFee.String(),
					Type:      event.Type,
				})
			}
		}

		if len(invalidIDs) > 0 {
			if err := config.GetDB().WithContext(dbCtx).Model(&models.TaskFeeEvent{}).Where("id IN (?)", invalidIDs).Update("status", models.TaskFeeEventStatusInvalid).Error; err != nil {
				return nil, err
			}
		}
	}

	return &GetTaskFeeLogsResponse{
		Data: logs,
	}, nil
}
