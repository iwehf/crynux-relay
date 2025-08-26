package withdraw

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetWithdrawRequestsInput struct {
	StartID uint                   `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int                    `query:"limit" json:"limit" description:"Limit"`
}

type WithdrawRecord struct {
	ID             uint                  `json:"id"`
	Address        string                `json:"address"`
	BenefitAddress string                `json:"benefit_address"`
	Amount         string                `json:"amount"`
	Network        string                `json:"network"`
	Status         models.WithdrawStatus `json:"status"`
}

type GetWithdrawRequestsResponse struct {
	response.Response
	Data []WithdrawRecord `json:"data"`
}

func GetWithdrawRequests(c *gin.Context, in *GetWithdrawRequestsInput) (*GetWithdrawRequestsResponse, error) {
	var records []models.WithdrawRecord

	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := config.GetDB().WithContext(dbCtx).Model(&models.WithdrawRecord{}).Where("id > ?", in.StartID).Order("id ASC").Limit(in.Limit).Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]WithdrawRecord, 0, len(records))
	if len(records) > 0 {
		lastID := records[0].ID
		for i, record := range records {
			if i > 0 && record.ID != lastID+1 {
				break
			}

			lastID = record.ID

			if record.LocalStatus != models.WithdrawLocalStatusProcessed {
				break
			}
	
			results = append(results, WithdrawRecord{
				ID:             record.ID,
				Address:        record.Address,
				BenefitAddress: record.BenefitAddress,
				Amount:         record.Amount.String(),
				Network:        record.Network,
				Status:         record.Status,
			})
		}
	}

	return &GetWithdrawRequestsResponse{
		Data: results,
	}, nil
}
