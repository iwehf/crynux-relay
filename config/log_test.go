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
