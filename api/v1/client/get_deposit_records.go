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

type GetDepositRecordsInput struct {
	Address  string  `path:"address" description:"The address of the user"`
	Page     uint    `query:"page" description:"The page number" default:"1"`
	PageSize uint    `query:"page_size" description:"The page size" default:"10"`
	Network  *string `query:"network" description:"The network of the deposit"`
}

type DepositRecord struct {
	ID        uint   `json:"id"`
	CreatedAt uint64 `json:"created_at"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Network   string `json:"network"`
	TxHash    string `json:"tx_hash"`
}

type GetDepositRecordsData struct {
	Total          int64           `json:"total" description:"The total number of deposit records"`
	DepositRecords []DepositRecord `json:"deposit_records" description:"The deposit records"`
}

type GetDepositRecordsResponse struct {
	response.Response
	Data *GetDepositRecordsData `json:"data" description:"The data of the deposit records"`
}

func GetDepositRecords(c *gin.Context, in *GetDepositRecordsInput) (*GetDepositRecordsResponse, error) {
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

	dbi := db.WithContext(dbCtx).Model(&models.DepositRecord{}).Where("address = ?", address)

	if in.Network != nil {
		dbi = dbi.Where("network = ?", *in.Network)
	}

	var total int64
	if err := dbi.Count(&total).Error; err != nil {
		return nil, err
	}

	var depositRecords []models.DepositRecord
	if err := dbi.Order("id DESC").Offset(int((in.Page - 1) * in.PageSize)).Limit(int(in.PageSize)).Find(&depositRecords).Error; err != nil {
		return nil, err
	}

	results := make([]DepositRecord, len(depositRecords))
	for i, record := range depositRecords {
		results[i] = DepositRecord{
			ID:        record.ID,
			CreatedAt: uint64(record.CreatedAt.Unix()),
			Address:   record.Address,
			Amount:    record.Amount.String(),
			Network:   record.Network,
			TxHash:    record.TxHash,
		}
	}

	return &GetDepositRecordsResponse{
		Data: &GetDepositRecordsData{
			Total:          total,
			DepositRecords: results,
		},
	}, nil
}
