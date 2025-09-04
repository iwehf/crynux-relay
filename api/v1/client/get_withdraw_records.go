package client

import (
	"context"
	"crynux_relay/api/v1/middleware"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetWithdrawRecordsInput struct {
	Address  string                 `path:"address" description:"The address of the user"`
	Page     uint                   `query:"page" description:"The page number" default:"1"`
	PageSize uint                   `query:"page_size" description:"The page size" default:"10"`
	Network  *string                `query:"network" description:"The network of the withdraw"`
	Status   *models.WithdrawStatus `query:"status" description:"The status of the withdraw"`
}

type WithdrawRecord struct {
	ID             uint                  `json:"id"`
	Address        string                `json:"address"`
	BenefitAddress string                `json:"benefit_address"`
	Amount         string                `json:"amount"`
	Network        string                `json:"network"`
	Status         models.WithdrawStatus `json:"status"`
	TxHash         string                `json:"tx_hash"`
}

type GetWithdrawRecordsData struct {
	Total           int64            `json:"total" description:"The total number of withdraw records"`
	WithdrawRecords []WithdrawRecord `json:"withdraw_records" description:"The withdraw records"`
}

type GetWithdrawRecordsResponse struct {
	response.Response
	Data *GetWithdrawRecordsData `json:"data" description:"The data of the withdraw records"`
}

func GetWithdrawRecords(c *gin.Context, in *GetWithdrawRecordsInput) (*GetWithdrawRecordsResponse, error) {
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

	dbi := db.WithContext(dbCtx).Model(&models.WithdrawRecord{}).Where("address = ?", address)

	if in.Network != nil {
		dbi = dbi.Where("network = ?", *in.Network)
	}

	if in.Status != nil {
		dbi = dbi.Where("status = ?", *in.Status)
	}

	var total int64
	if err := dbi.Count(&total).Error; err != nil {
		return nil, err
	}

	var withdrawRecords []models.WithdrawRecord
	if err := dbi.Order("id DESC").Offset(int((in.Page - 1) * in.PageSize)).Limit(int(in.PageSize)).Find(&withdrawRecords).Error; err != nil {
		return nil, err
	}

	results := make([]WithdrawRecord, len(withdrawRecords))
	for i, record := range withdrawRecords {
		var txHash string
		if record.TxHash.Valid {
			txHash = record.TxHash.String
		}
		results[i] = WithdrawRecord{
			ID:             record.ID,
			Address:        record.Address,
			BenefitAddress: record.BenefitAddress,
			Amount:         record.Amount.String(),
			Network:        record.Network,
			Status:         record.Status,
			TxHash:         txHash,
		}
	}

	return &GetWithdrawRecordsResponse{
		Data: &GetWithdrawRecordsData{
			Total:           total,
			WithdrawRecords: results,
		},
	}, nil
}
