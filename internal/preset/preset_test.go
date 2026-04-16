package preset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	// Create a temporary presets file
	tmpDir := t.TempDir()
	presetsPath := filepath.Join(tmpDir, "presets.json")
	content := `{
		"node": {
			"image": "docker.io/library/node:20-alpine",
			"description": "Node.js 20 environment",
			"category": "Develop",
			"icon": "code"
		},
		"redis": {
			"image": "docker.io/library/redis:7-alpine",
			"description": "Redis Cache",
			"category": "Database",
			"icon": "database"
		}
	}`
	if err := os.WriteFile(presetsPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	manager, err := NewManager(presetsPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if len(manager.GetAll()) != 2 {
		t.Errorf("Expected 2 presets, got %d", len(manager.GetAll()))
	}

	p, ok := manager.Find("node")
	if !ok {
		t.Error("Failed to find 'node' preset")
	}
	if p.Category != "Develop" {
		t.Errorf("Expected category 'Develop', got '%s'", p.Category)
	}
	if p.Icon != "code" {
		t.Errorf("Expected icon 'code', got '%s'", p.Icon)
	}
}

func TestAutoDetect(t *testing.T) {
	// Setup mock manager with standard keys
	manager := &Manager{
		presets: map[string]Preset{
			"node":   {Category: "Develop"},
			"python": {Category: "Develop"},
			"go":     {Category: "Develop"},
		},
	}

	tmpDir := t.TempDir()

	// Test Node.js detection
	pkgJson := filepath.Join(tmpDir, "package.json")
	os.WriteFile(pkgJson, []byte("{}"), 0644)
	res, ok := manager.AutoDetect(tmpDir)
	if !ok || res != "node" {
		t.Errorf("Expected 'node' detection, got '%s', ok=%v", res, ok)
	}
	os.Remove(pkgJson)

	// Test Go detection
	goMod := filepath.Join(tmpDir, "go.mod")
	os.WriteFile(goMod, []byte("module test"), 0644)
	res, ok = manager.AutoDetect(tmpDir)
	if !ok || res != "go" {
		t.Errorf("Expected 'go' detection, got '%s', ok=%v", res, ok)
	}
}
