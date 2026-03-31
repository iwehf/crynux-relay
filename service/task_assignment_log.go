package service

import (
	"context"
	"crynux_relay/config"
	"crynux_relay/models"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func logTaskAssignmentEvent(ctx context.Context, task *models.InferenceTask, nodes []models.Node) {
	logger := config.GetTaskAssignmentLogger()
	if logger == nil {
		return
	}

	queueSizeLabel := "unknown"
	queueSize, err := models.GetQueuedTaskCount(ctx, config.GetDB())
	if err != nil {
		log.Errorf("TaskAssignment: get queued task count error: %v", err)
	} else {
		queueSizeLabel = fmt.Sprintf("%d", queueSize)
	}

	logger.Infof(
		"[TaskAssignment] [%s] [%s] [%d] [%s] [%s] task_id=%s",
		taskTypeTag(task.TaskType),
		taskVramRequirementLabel(task),
		len(nodes),
		taskAssignmentNodesLabel(nodes),
		queueSizeLabel,
		task.TaskIDCommitment,
	)
}

func taskVramRequirementLabel(task *models.InferenceTask) string {
	if len(task.RequiredGPU) > 0 {
		return fmt.Sprintf("%dGB", task.RequiredGPUVRAM)
	}
	return fmt.Sprintf("%dGB", task.MinVRAM)
}

func taskAssignmentNodesLabel(nodes []models.Node) string {
	if len(nodes) == 0 {
		return "none"
	}

	addresses := make([]string, 0, len(nodes))
	for _, node := range nodes {
		addresses = append(addresses, node.Address)
	}
	return strings.Join(addresses, ",")
}
