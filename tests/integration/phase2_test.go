package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bril3d/minicontainer/internal/minifile"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
)

func TestRuntime_Build(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	if _, err := podman.Version(); err != nil {
		t.Skip("Podman not available")
	}

	tmpDir := t.TempDir()
	dockerfile := filepath.Join(tmpDir, "Dockerfile")
	content := "FROM docker.io/library/alpine:latest\nCMD [\"echo\", \"built-it\"]"
	os.WriteFile(dockerfile, []byte(content), 0644)

	tagName := fmt.Sprintf("test-build-%d:latest", time.Now().Unix())
	
	t.Logf("Building image %s...", tagName)
	err := podman.Build(rt.BuildOptions{
		Context:    tmpDir,
		Dockerfile: "Dockerfile",
		Tags:       []string{tagName},
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Clean up image
	// podman rmi tagName
}

func TestMinifile_Orchestration(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	if _, err := podman.Version(); err != nil {
		t.Skip("Podman not available")
	}

	orch := minifile.NewOrchestrator(podman)
	
	m := &minifile.Minifile{
		Project: "test-proj",
		Services: map[string]minifile.Service{
			"web": {
				Image: "docker.io/library/alpine:latest",
				Command: []string{"sleep", "300"},
			},
		},
	}

	t.Log("Testing Up...")
	if err := orch.Up(m); err != nil {
		t.Fatalf("Up failed: %v", err)
	}

	// Verify container exists
	containers, _ := podman.List()
	found := false
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "test-proj-web" || name == "/test-proj-web" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Container 'test-proj-web' not found after Up")
	}

	t.Log("Testing Down...")
	if err := orch.Down(m); err != nil {
		t.Fatalf("Down failed: %v", err)
	}

	// Verify removed
	containers, _ = podman.List()
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "test-proj-web" || name == "/test-proj-web" {
				t.Error("Container 'test-proj-web' still exists after Down")
			}
		}
	}
}
