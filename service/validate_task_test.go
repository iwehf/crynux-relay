package service

import (
	"crynux_relay/models"
	"testing"
)

func TestAssignValidationGroupQosScoresAllTimeoutTasksDoNotContribute(t *testing.T) {
	tasks := []*models.InferenceTask{
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-a"),
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-b"),
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-c"),
	}

	assignValidationGroupQosScores(tasks)

	for i, task := range tasks {
		if task.QOSScore.Valid {
			t.Fatalf("task %d should have invalid qos score when the whole group timed out", i)
		}
		if shouldPersistValidationGroupTimeoutQos(task) {
			t.Fatalf("task %d timeout score should not be persisted when the whole group timed out", i)
		}
	}
}

func TestAssignValidationGroupQosScoresTwoTimeoutsStillPenalizeLongTerm(t *testing.T) {
	tasks := []*models.InferenceTask{
		newValidationGroupTask(models.TaskScoreReady, models.TaskAbortReasonNone, "node-a"),
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-b"),
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-c"),
	}

	assignValidationGroupQosScores(tasks)

	if !tasks[0].QOSScore.Valid || tasks[0].QOSScore.Int64 != 10 {
		t.Fatalf("finished task should receive score 10, got %+v", tasks[0].QOSScore)
	}
	for i := 1; i < len(tasks); i++ {
		if !tasks[i].QOSScore.Valid || tasks[i].QOSScore.Int64 != 0 {
			t.Fatalf("timeout task %d should receive valid zero score, got %+v", i, tasks[i].QOSScore)
		}
		if !shouldPersistValidationGroupTimeoutQos(tasks[i]) {
			t.Fatalf("timeout task %d zero score should be persisted to long-term qos", i)
		}
	}
}

func TestAssignValidationGroupQosScoresSingleTimeoutStillPenalizesLongTerm(t *testing.T) {
	tasks := []*models.InferenceTask{
		newValidationGroupTask(models.TaskScoreReady, models.TaskAbortReasonNone, "node-a"),
		newValidationGroupTask(models.TaskEndGroupRefund, models.TaskAbortReasonNone, "node-b"),
		newValidationGroupTask(models.TaskEndAborted, models.TaskAbortTimeout, "node-c"),
	}

	assignValidationGroupQosScores(tasks)

	if !tasks[0].QOSScore.Valid || tasks[0].QOSScore.Int64 != 10 {
		t.Fatalf("first finished task should receive score 10, got %+v", tasks[0].QOSScore)
	}
	if !tasks[1].QOSScore.Valid || tasks[1].QOSScore.Int64 != 5 {
		t.Fatalf("second finished task should receive score 5, got %+v", tasks[1].QOSScore)
	}
	if !tasks[2].QOSScore.Valid || tasks[2].QOSScore.Int64 != 0 {
		t.Fatalf("timeout task should receive valid zero score, got %+v", tasks[2].QOSScore)
	}
	if !shouldPersistValidationGroupTimeoutQos(tasks[2]) {
		t.Fatalf("single timeout task zero score should be persisted to long-term qos")
	}
}

func newValidationGroupTask(status models.TaskStatus, abortReason models.TaskAbortReason, selectedNode string) *models.InferenceTask {
	return &models.InferenceTask{
		Status:       status,
		AbortReason:  abortReason,
		SelectedNode: selectedNode,
	}
}
