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

type RejectWithdrawRequestInput struct {
	ID uint `path:"id" json:"id" description:"Withdraw request ID"`
}

type RejectWithdrawRequestResponse struct {
	response.Response
}

func RejectWithdrawRequest(c *gin.Context, in *RejectWithdrawRequestInput) (*RejectWithdrawRequestResponse, error) {
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

		if err := tx.Model(&models.WithdrawRecord{}).Where("id = ?", in.ID).Update("status", models.WithdrawStatusFailed).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &RejectWithdrawRequestResponse{}, nil
}
