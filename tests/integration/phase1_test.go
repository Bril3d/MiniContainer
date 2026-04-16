package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
)

func TestRuntime_Restart(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	if _, err := podman.Version(); err != nil {
		t.Skip("Podman not available")
	}

	image := "docker.io/library/alpine:latest"
	name := fmt.Sprintf("test-restart-%d", time.Now().Unix())

	id, err := podman.Run(rt.RunOptions{
		Image:  image,
		Name:   name,
		Cmd:    []string{"sleep", "300"},
		Detach: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer podman.Remove(id, true)

	t.Log("Restarting container...")
	if err := podman.Restart(id); err != nil {
		t.Errorf("Failed to restart container: %v", err)
	}

	// Verify it's still running
	containers, _ := podman.List()
	found := false
	for _, c := range containers {
		if strings.Contains(c.ID, id) && strings.Contains(c.Status, "Up") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Container not running after restart")
	}
}

func TestRuntime_Exec(t *testing.T) {
	podman := rt.NewPodmanRuntime()
	if _, err := podman.Version(); err != nil {
		t.Skip("Podman not available")
	}

	image := "docker.io/library/alpine:latest"
	name := fmt.Sprintf("test-exec-%d", time.Now().Unix())

	id, err := podman.Run(rt.RunOptions{
		Image:  image,
		Name:   name,
		Cmd:    []string{"sleep", "300"},
		Detach: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer podman.Remove(id, true)

	t.Log("Executing command in container...")
	// Note: Our current Exec implementation uses ExecStream/ExecInteractive which hooks to os.Stdout
	// For testing, we just check if it returns without error
	err = podman.Exec(id, []string{"echo", "hello-from-test"}, false)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}
}
