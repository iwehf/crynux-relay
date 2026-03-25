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

type GetRelayAccountLogsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
}

type GetRelayAccountLogsInputWithSignature struct {
	GetRelayAccountLogsInput
	Timestamp int64  `query:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `query:"signature" json:"signature" description:"Signature" validate:"required"`
}

type RelayAccountLog struct {
	ID        uint                         `json:"id"`
	CreatedAt uint64                       `json:"created_at"`
	Address   string                       `json:"address"`
	Amount    string                       `json:"amount"`
	Type      models.RelayAccountEventType `json:"type"`
}

type GetRelayAccountLogsResponse struct {
	response.Response
	Data []RelayAccountLog `json:"data"`
}

func GetRelayAccountLogs(c *gin.Context, in *GetRelayAccountLogsInputWithSignature) (*GetRelayAccountLogsResponse, error) {
	match, address, err := validate.ValidateSignature(in.GetRelayAccountLogsInput, in.Timestamp, in.Signature)

	if err != nil || !match {
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	if address != config.GetConfig().Withdraw.RelayWalletAddress {
		validationErr := response.NewValidationErrorResponse("address", "Invalid address")
		return nil, validationErr
	}

	var events []models.RelayAccountEvent

	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.RelayAccountEvent{}).Where("id > ?", in.StartID).Where("status != ?", models.RelayAccountEventStatusInvalid).Order("id ASC").Limit(in.Limit).Find(&events).Error; err != nil {
		return nil, err
	}

	// Filter out events with non-consecutive IDs
	logs := make([]RelayAccountLog, 0, len(events))
	if len(events) > 0 {
		appConfig := config.GetConfig()
		invalidIDs := make([]uint, 0)
		for _, event := range events {
			if event.Status == models.RelayAccountEventStatusPending {
				break
			}

			if event.Status == models.RelayAccountEventStatusProcessed {
				valid := utils.VerifyMAC([]byte(event.Reason), appConfig.MAC.SecretKey, event.MAC)
				if !valid {
					invalidIDs = append(invalidIDs, event.ID)
					continue
				}
				logs = append(logs, RelayAccountLog{
					ID:        event.ID,
					CreatedAt: uint64(event.CreatedAt.Unix()),
					Address:   event.Address,
					Amount:    event.Amount.String(),
					Type:      event.Type,
				})
			}
		}

		if len(invalidIDs) > 0 {
			if err := config.GetDB().WithContext(dbCtx).Model(&models.RelayAccountEvent{}).Where("id IN (?)", invalidIDs).Update("status", models.RelayAccountEventStatusInvalid).Error; err != nil {
				return nil, err
			}
		}
	}

	return &GetRelayAccountLogsResponse{
		Data: logs,
	}, nil
}
