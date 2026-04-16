package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Bril3d/minicontainer/internal/minifile"
	"github.com/Bril3d/minicontainer/internal/runtime"
)

func TestPhase3_VolumesAndEnv(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test; set INTEGRATION=true")
	}

	rt := runtime.NewPodmanRuntime()

	orch := minifile.NewOrchestrator(rt)

	// Create temp project directory
	tmpDir, err := ioutil.TempDir("", "mini-test-p3-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Create .env file
	envContent := "GLOBAL_VAR=hello_from_env\n# Comment\nFOO=bar"
	if err := ioutil.WriteFile(filepath.Join(tmpDir, ".env"), []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Create data file for volume mount
	dataDir := filepath.Join(tmpDir, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dataDir, "msg.txt"), []byte("persistence_active"), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Create Minifile
	// Note: Using alpine for simplicity. We map ./data to /mnt/data
	miniContent := `
version: "1.0"
project: p3test
services:
  web:
    image: alpine
    environment:
      SERVICE_VAR: "local_value"
    volumes:
      - "./data:/mnt/data"
    command: ["sleep", "3600"]
`
	miniPath := filepath.Join(tmpDir, "Minifile")
	if err := ioutil.WriteFile(miniPath, []byte(miniContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Run 'Up'
	m, err := minifile.Parse(miniPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := orch.Up(m); err != nil {
		t.Fatal(err)
	}
	defer orch.Down(m)

	containerName := "p3test-web"
	
	// Wait a bit for container to settle
	time.Sleep(2 * time.Second)

	// 5. Verify Environment Variables (Global .env + Minifile)
	out, err := rt.ExecWithOutput(containerName, []string{"env"})
	if err != nil {
		t.Fatal(err)
	}
	envOutput := string(out)
	if !strings.Contains(envOutput, "GLOBAL_VAR=hello_from_env") {
		t.Errorf("Global env var missing. Got: %s", envOutput)
	}
	if !strings.Contains(envOutput, "SERVICE_VAR=local_value") {
		t.Errorf("Service env var missing. Got: %s", envOutput)
	}

	// 6. Verify Volume Mount (Relative path resolution + Bind mount)
	out, err = rt.ExecWithOutput(containerName, []string{"cat", "/mnt/data/msg.txt"})
	if err != nil {
		t.Errorf("Failed to read volume content: %v", err)
	} else {
		if strings.TrimSpace(string(out)) != "persistence_active" {
			t.Errorf("Volume content mismatch. Got: %s", string(out))
		}
	}

	fmt.Println("Phase 3 integration test passed!")
}
