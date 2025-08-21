package withdraw

import (
	"context"
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FufillWithdrawRequestInput struct {
	ID uint `path:"id" json:"id" description:"Withdraw request ID"`
}

type FufillWithdrawRequestResponse struct {
	response.Response
}

func FufillWithdrawRequest(c *gin.Context, in *FufillWithdrawRequestInput) (*FufillWithdrawRequestResponse, error) {
	dbCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var record models.WithdrawRecord

	if err := config.GetDB().WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.WithdrawRecord{}).Where("id = ?", in.ID).First(&record).Error; err != nil {
			return err
		}

		if record.Status != models.WithdrawStatusPending {
			return response.NewValidationErrorResponse("id", "withdraw request is not pending")
		}

		if err := tx.Model(&models.WithdrawRecord{}).Where("id = ?", in.ID).Update("status", models.WithdrawStatusSuccess).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &FufillWithdrawRequestResponse{}, nil
}
