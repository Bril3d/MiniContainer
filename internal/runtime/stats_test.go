package runtime

import (
	"encoding/json"
	"testing"
)

// Simplified struct matching what's in cmd/stats.go or expected by GUI
func TestStatsJSONParsing(t *testing.T) {
	// Mock raw Podman stats output
	rawOutput := `[{"block_io":"0B / 0B","cpu_perc":"0.00%","id":"123","mem_perc":"0.10%","mem_usage":"10MB / 100MB","name":"test-node","net_io":"1KB / 1KB","pids":"1"}]`
	
	var stats []ContainerStats
	err := json.Unmarshal([]byte(rawOutput), &stats)
	if err != nil {
		t.Fatalf("Failed to parse mock JSON: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("Expected 1 stat entry, got %d", len(stats))
	}

	if stats[0].Name != "test-node" {
		t.Errorf("Expected name 'test-node', got '%s'", stats[0].Name)
	}
}
