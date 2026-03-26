package relayaccount

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/api/v1/validate"
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/utils"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GetRelayAccountEventLogsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
}

type GetRelayAccountEventLogsInputWithSignature struct {
	GetRelayAccountEventLogsInput
	Timestamp int64  `query:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `query:"signature" json:"signature" description:"Signature" validate:"required"`
}

type RelayAccountEventLog struct {
	ID        uint                         `json:"id"`
	CreatedAt uint64                       `json:"created_at"`
	Address   string                       `json:"address"`
	Amount    string                       `json:"amount"`
	Type      models.RelayAccountEventType `json:"type"`
	Payload   string                       `json:"payload"`
}

type GetRelayAccountEventLogsResponse struct {
	response.Response
	Data []RelayAccountEventLog `json:"data"`
}

func getEventPayload(eventType models.RelayAccountEventType, reason string) string {
	if eventType != models.RelayAccountEventTypeDeposit {
		return "{}"
	}

	reasons := strings.SplitN(reason, "-", 3)
	if len(reasons) != 3 || reasons[0] != strconv.Itoa(int(models.RelayAccountEventTypeDeposit)) {
		return "{}"
	}
	payload := map[string]string{
		"tx_hash": reasons[1],
		"network": reasons[2],
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(payloadJSON)
}

func GetRelayAccountEventLogs(c *gin.Context, in *GetRelayAccountEventLogsInputWithSignature) (*GetRelayAccountEventLogsResponse, error) {
	match, address, err := validate.ValidateSignature(in.GetRelayAccountEventLogsInput, in.Timestamp, in.Signature)
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

	if err := config.GetDB().WithContext(dbCtx).
		Model(&models.RelayAccountEvent{}).
		Where("id > ?", in.StartID).
		Where("status != ?", models.RelayAccountEventStatusInvalid).
		Order("id ASC").
		Limit(in.Limit).
		Find(&events).Error; err != nil {
		return nil, err
	}

	logs := make([]RelayAccountEventLog, 0, len(events))
	if len(events) > 0 {
		appConfig := config.GetConfig()
		invalidIDs := make([]uint, 0)
		for _, event := range events {
			if event.Status == models.RelayAccountEventStatusPending {
				break
			}

			if event.Status != models.RelayAccountEventStatusProcessed {
				continue
			}

			valid := utils.VerifyMAC([]byte(event.Reason), appConfig.MAC.SecretKey, event.MAC)
			if !valid {
				invalidIDs = append(invalidIDs, event.ID)
				continue
			}

			logs = append(logs, RelayAccountEventLog{
				ID:        event.ID,
				CreatedAt: uint64(event.CreatedAt.Unix()),
				Address:   event.Address,
				Amount:    event.Amount.String(),
				Type:      event.Type,
				Payload:   getEventPayload(event.Type, event.Reason),
			})
		}

		if len(invalidIDs) > 0 {
			if err := config.GetDB().WithContext(dbCtx).
				Model(&models.RelayAccountEvent{}).
				Where("id IN (?)", invalidIDs).
				Update("status", models.RelayAccountEventStatusInvalid).Error; err != nil {
				return nil, err
			}
		}
	}

	return &GetRelayAccountEventLogsResponse{
		Data: logs,
	}, nil
}
