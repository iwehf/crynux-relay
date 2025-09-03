package credits

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/api/v1/validate"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetCreditsRecordInput struct {
	ID uint `path:"id" json:"id" description:"Credits record ID"`
}

type GetCreditsRecordInputWithSignature struct {
	GetCreditsRecordInput
	Timestamp int64  `query:"timestamp" json:"timestamp" description:"Signature timestamp" validate:"required"`
	Signature string `query:"signature" json:"signature" description:"Signature" validate:"required"`
}

type CreditsRecordStatus uint8

const (
	CreditsRecordStatusPending CreditsRecordStatus = iota
	CreditsRecordStatusSuccess
	CreditsRecordStatusFailed
)

type CreditsRecord struct {
	ID        uint                `json:"id"`
	CreatedAt uint64              `json:"created_at"`
	Address   string              `json:"address"`
	Amount    string              `json:"amount"`
	Network   string              `json:"network"`
	Status    CreditsRecordStatus `json:"status"`
}

type GetCreditsRecordResponse struct {
	response.Response
	Data *CreditsRecord `json:"data" description:"The data of the credits record"`
}

func GetCreditsRecord(c *gin.Context, in *GetCreditsRecordInputWithSignature) (*GetCreditsRecordResponse, error) {
	match, address, err := validate.ValidateSignature(in.GetCreditsRecordInput, in.Timestamp, in.Signature)
	if err != nil || !match {
		validationErr := response.NewValidationErrorResponse("signature", "Invalid signature")
		return nil, validationErr
	}

	if address != config.GetConfig().Credits.Address {
		validationErr := response.NewValidationErrorResponse("address", "Invalid address")
		return nil, validationErr
	}

	record, status, err := func() (*models.CreditsRecord, CreditsRecordStatus, error) {
		db := config.GetDB()
		dbCtx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		var status CreditsRecordStatus
		record, err := models.GetCreditsRecordByID(dbCtx, db, in.ID)
		if err != nil {
			return nil, CreditsRecordStatusPending, err
		}

		transaction, err := models.GetTransactionByID(dbCtx, db, record.BlockchainTransactionID)
		if err != nil {
			return nil, CreditsRecordStatusPending, err
		}

		switch transaction.Status {
		case models.TransactionStatusPending, models.TransactionStatusSent:
			status = CreditsRecordStatusPending
		case models.TransactionStatusConfirmed:
			status = CreditsRecordStatusSuccess
		case models.TransactionStatusFailed:
			retryTransactions, err := models.GetRetryTransactionsByID(dbCtx, db, record.BlockchainTransactionID)
			if err != nil {
				return nil, status, err
			}
			if len(retryTransactions) > 0 {
				retryTransaction := retryTransactions[len(retryTransactions)-1]
				switch retryTransaction.Status {
				case models.TransactionStatusPending, models.TransactionStatusSent:
					status = CreditsRecordStatusPending
				case models.TransactionStatusConfirmed:
					status = CreditsRecordStatusSuccess
				case models.TransactionStatusFailed:
					if retryTransaction.RetryCount >= retryTransaction.MaxRetries {
						status = CreditsRecordStatusFailed
					} else {
						status = CreditsRecordStatusPending
					}
				}
			} else {
				status = CreditsRecordStatusPending
			}
		}
		return record, status, nil
	}()
	if err != nil {
		return nil, response.NewExceptionResponse(err)
	}

	return &GetCreditsRecordResponse{
		Data: &CreditsRecord{
			ID:        record.ID,
			CreatedAt: uint64(record.CreatedAt.Unix()),
			Address:   record.Address,
			Amount:    record.Amount.String(),
			Network:   record.Network,
			Status:    status,
		},
	}, nil
}
