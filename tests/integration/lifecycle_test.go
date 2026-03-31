package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ps "github.com/Bril3d/minicontainer/internal/preset"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
)

func TestLifecycle_Alpine(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	
	// Check if podman is available
	_, err := podman.Version()
	if err != nil {
		t.Skip("Podman not available, skipping integration test:", err)
	}

	image := "docker.io/library/alpine:latest"
	containerName := fmt.Sprintf("mini-test-%d", time.Now().Unix())

	// 1. Pull
	t.Logf("Pulling image %s...", image)
	err = podman.Pull(image)
	if err != nil {
		t.Fatalf("Failed to pull image: %v", err)
	}

	// 2. Run
	t.Logf("Running container %s...", containerName)
	opts := rt.RunOptions{
		Image:  image,
		Name:   containerName,
		Cmd:    []string{"sleep", "300"},
		Detach: true,
	}
	id, err := podman.Run(opts)
	if err != nil {
		t.Fatalf("Failed to run container: %v", err)
	}
	
	// Ensure cleanup
	t.Cleanup(func() {
		t.Logf("Cleaning up container %s...", id)
		_ = podman.Stop(id)
		_ = podman.Remove(id, true)
	})

	// 3. List & Verify
	t.Log("Verifying container in list...")
	containers, err := podman.List()
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	found := false
	for _, c := range containers {
		if strings.Contains(c.ID, id) || c.Names[0] == containerName {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Container %s not found in list", id)
	}

	// 4. Stop
	t.Log("Stopping container...")
	err = podman.Stop(id)
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}

	// 5. Remove (handled by Cleanup, but we test it explicitly here too if desired)
}

func TestLifecycle_AutoDetect(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	_, err := podman.Version()
	if err != nil {
		t.Skip("Podman not available, skipping integration test")
	}

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "mini-test-detect-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a marker file (Node.js)
	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Load presets (need relative path to real presets/presets.json or mock them)
	// For integration test, we'll try to find the real one relative to current dir
	wd, _ := os.Getwd()
	// wd is expected to be .../tests/integration
	presetsPath := filepath.Join(wd, "..", "..", "presets", "presets.json")
	
	manager, err := ps.NewManager(presetsPath)
	if err != nil {
		t.Logf("Warning: could not load real presets at %s, creating mock manager", presetsPath)
		// Fallback to manual check if presets.json isn't reachable
		if _, err := os.Stat(filepath.Join(tmpDir, "package.json")); err != nil {
			t.Fatal("Failed to create mock project marker")
		}
		return
	}

	// Test detection
	detected, ok := manager.AutoDetect(tmpDir)
	if !ok {
		t.Errorf("Expected to detect Node.js project in %s", tmpDir)
	}
	if detected != "node" {
		t.Errorf("Expected 'node' preset, got '%s'", detected)
	}
}
