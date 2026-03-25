package stats

import (
	"testing"
	"time"
)

func TestTaskFeeHistogramCacheHitAndExpire(t *testing.T) {
	taskFeeHistogramCacheLock.Lock()
	taskFeeHistogramCache = make(map[TaskTypeString]taskFeeHistogramCacheEntry)
	taskFeeHistogramCacheLock.Unlock()

	now := time.Now().UTC()
	setTaskFeeHistogramCache(AllTaskType, &GetTaskFeeHistogramData{
		TaskFees:   []string{"100", "200"},
		TaskCounts: []int64{1, 2},
	}, now)

	data, ok := getTaskFeeHistogramFromCache(AllTaskType, now.Add(30*time.Second))
	if !ok {
		t.Fatalf("expected cache hit before expiry")
	}
	if len(data.TaskFees) != 2 || data.TaskFees[0] != "100" || data.TaskCounts[1] != 2 {
		t.Fatalf("unexpected cache payload: %+v", data)
	}

	data.TaskFees[0] = "mutated"
	dataAgain, ok := getTaskFeeHistogramFromCache(AllTaskType, now.Add(30*time.Second))
	if !ok {
		t.Fatalf("expected second cache hit before expiry")
	}
	if dataAgain.TaskFees[0] != "100" {
		t.Fatalf("cached data should be immutable copy, got %q", dataAgain.TaskFees[0])
	}

	if _, ok := getTaskFeeHistogramFromCache(AllTaskType, now.Add(61*time.Second)); ok {
		t.Fatalf("expected cache miss after expiry")
	}
}

func TestPow10Big(t *testing.T) {
	if got := pow10Big(-1).String(); got != "1" {
		t.Fatalf("pow10Big(-1) got %s, want 1", got)
	}
	if got := pow10Big(0).String(); got != "1" {
		t.Fatalf("pow10Big(0) got %s, want 1", got)
	}
	if got := pow10Big(3).String(); got != "1000" {
		t.Fatalf("pow10Big(3) got %s, want 1000", got)
	}
}
