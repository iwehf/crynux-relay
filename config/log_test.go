package config

import (
	"path/filepath"
	"testing"
)

func TestGetNodeHealthLogPath_DefaultPath(t *testing.T) {
	tests := []string{"", "stdout", "stderr"}
	for _, output := range tests {
		if got := getNodeHealthLogPath(output); got != defaultNodeHealthLogPath {
			t.Fatalf("expected default path for output %q, got %q", output, got)
		}
	}
}

func TestGetNodeHealthLogPath_BasedOnMainLogDir(t *testing.T) {
	got := getNodeHealthLogPath("data/logs/relay.log")
	if got != filepath.Join("data", "logs", "node_health.log") {
		t.Fatalf("expected node health log path in same directory, got %q", got)
	}
}

func TestGetTaskAssignmentLogPath_DefaultPath(t *testing.T) {
	tests := []string{"", "stdout", "stderr"}
	for _, output := range tests {
		if got := getTaskAssignmentLogPath(output); got != defaultTaskAssignmentLogPath {
			t.Fatalf("expected default path for output %q, got %q", output, got)
		}
	}
}

func TestGetTaskAssignmentLogPath_BasedOnMainLogDir(t *testing.T) {
	got := getTaskAssignmentLogPath("data/logs/relay.log")
	if got != filepath.Join("data", "logs", "task_assignment.log") {
		t.Fatalf("expected task assignment log path in same directory, got %q", got)
	}
}

func TestInitNodeHealthLogger_DefaultDisabled(t *testing.T) {
	nodeHealthLogger = nil

	cfg := &AppConfig{}

	initNodeHealthLogger(cfg)

	if nodeHealthLogger != nil {
		t.Fatal("expected node health logger to be nil by default")
	}
}

func TestInitNodeHealthLogger_Disabled(t *testing.T) {
	nodeHealthLogger = nil

	enabled := false
	cfg := &AppConfig{}
	cfg.Log.Features.NodeHealthEnabled = &enabled

	initNodeHealthLogger(cfg)

	if nodeHealthLogger != nil {
		t.Fatal("expected node health logger to be nil when disabled")
	}
}

func TestInitNodeHealthLogger_Enabled(t *testing.T) {
	nodeHealthLogger = nil

	enabled := true
	cfg := &AppConfig{}
	cfg.Log.Features.NodeHealthEnabled = &enabled

	initNodeHealthLogger(cfg)

	if nodeHealthLogger == nil {
		t.Fatal("expected node health logger to be initialized when enabled")
	}
}

func TestInitTaskAssignmentLogger_Disabled(t *testing.T) {
	taskAssignmentLogger = nil

	cfg := &AppConfig{}
	cfg.Log.Features.TaskAssignmentEnabled = false

	initTaskAssignmentLogger(cfg)

	if taskAssignmentLogger != nil {
		t.Fatal("expected task assignment logger to be nil when disabled")
	}
}

func TestInitTaskAssignmentLogger_Enabled(t *testing.T) {
	taskAssignmentLogger = nil

	cfg := &AppConfig{}
	cfg.Log.Features.TaskAssignmentEnabled = true

	initTaskAssignmentLogger(cfg)

	if taskAssignmentLogger == nil {
		t.Fatal("expected task assignment logger to be initialized when enabled")
	}
}
