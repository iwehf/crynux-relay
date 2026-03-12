package relayaccount

import (
	"context"
	"crynux_relay/api/v1/middleware"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskFeeRecordsInput struct {
	Address  string `path:"address" description:"The address of the user"`
	Page     uint   `query:"page" description:"The page number" default:"1"`
	PageSize uint   `query:"page_size" description:"The page size" default:"10"`
}

type TaskFeeRecord struct {
	ID        uint                         `json:"id"`
	CreatedAt uint64                       `json:"created_at"`
	Address   string                       `json:"address"`
	Amount    string                       `json:"amount"`
	Type      models.RelayAccountEventType `json:"type"`
	Reason    string                       `json:"reason"`
}

type TaskFeeRecordsData struct {
	Total   int64           `json:"total" description:"The total number of records"`
	Records []TaskFeeRecord `json:"records" description:"The task fee records"`
}

type TaskFeeRecordsResponse struct {
	response.Response
	Data *TaskFeeRecordsData `json:"data" description:"The task fee records data"`
}

func GetTaskFeeLedgerRecords(c *gin.Context, in *TaskFeeRecordsInput) (*TaskFeeRecordsResponse, error) {
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}

	address := middleware.GetUserAddress(c)
	if address != in.Address {
		validationErr := response.NewValidationErrorResponse("address", "Address mismatch")
		return nil, validationErr
	}

	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := config.GetDB()
	eventTypes := []models.RelayAccountEventType{
		models.RelayAccountEventTypeTaskPayment,
		models.RelayAccountEventTypeTaskIncome,
		models.RelayAccountEventTypeTaskRefund,
	}

	dbi := db.WithContext(dbCtx).
		Model(&models.RelayAccountEvent{}).
		Where("address = ?", address).
		Where("type IN (?)", eventTypes).
		Where("status = ?", models.RelayAccountEventStatusProcessed)

	var total int64
	if err := dbi.Count(&total).Error; err != nil {
		return nil, err
	}

	var events []models.RelayAccountEvent
	if err := dbi.Order("id DESC").
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&events).Error; err != nil {
		return nil, err
	}

	results := make([]TaskFeeRecord, len(events))
	for i, event := range events {
		results[i] = TaskFeeRecord{
			ID:        event.ID,
			CreatedAt: uint64(event.CreatedAt.Unix()),
			Address:   event.Address,
			Amount:    event.Amount.String(),
			Type:      event.Type,
			Reason:    event.Reason,
		}
	}

	return &TaskFeeRecordsResponse{
		Data: &TaskFeeRecordsData{
			Total:   total,
			Records: results,
		},
	}, nil
}
