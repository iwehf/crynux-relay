package admin

import (
	"bytes"
	"context"
	"crynux_relay/api/v2/response"
	"crynux_relay/config"
	"crynux_relay/models"
	"crynux_relay/service"
	"encoding/csv"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type nodeQosCSVRow struct {
	NodeAddress  string
	Card         string
	ProbWeight   float64
	StakingScore float64
	QosScore     float64
	QosLong      float64
	QosShort     float64
}

func ExportActiveNodesQosCSV(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var nodes []models.Node
	if err := config.GetDB().WithContext(ctx).
		Model(&models.Node{}).
		Where("status != ?", models.NodeStatusQuit).
		Find(&nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Message: err.Error()})
		return
	}

	maxStaking := service.GetMaxStaking()
	rows := make([]nodeQosCSVRow, 0, len(nodes))

	for _, node := range nodes {
		qosLong, qosShort, qosScore := service.CalculateQosComponents(node.QOSScore, node.HealthBase, node.HealthUpdatedAt)
		stakingScore, qosProb, probWeight := service.CalculateSelectingProb(&node.StakeAmount.Int, maxStaking, qosScore)

		rows = append(rows, nodeQosCSVRow{
			NodeAddress:  node.Address,
			Card:         node.GPUName + " + " + strconv.FormatUint(node.GPUVram, 10),
			ProbWeight:   probWeight,
			StakingScore: stakingScore,
			QosScore:     qosProb,
			QosLong:      qosLong,
			QosShort:     qosShort,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ProbWeight > rows[j].ProbWeight
	})

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"node address",
		"card (GPT Type + VRAM)",
		"prob_weight",
		"staking_score",
		"qos_score",
		"qos_long",
		"qos_short",
	}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Message: err.Error()})
		return
	}

	for _, row := range rows {
		if err := writer.Write([]string{
			row.NodeAddress,
			row.Card,
			formatFloat(row.ProbWeight),
			formatFloat(row.StakingScore),
			formatFloat(row.QosScore),
			formatFloat(row.QosLong),
			formatFloat(row.QosShort),
		}); err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{Message: err.Error()})
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Message: err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=nodes_qos.csv")
	c.String(http.StatusOK, buf.String())
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 8, 64)
}
