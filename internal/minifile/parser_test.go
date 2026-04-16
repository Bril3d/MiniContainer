package minifile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	content := `
version: "1.0"
services:
  app:
    build:
      context: .
    ports:
      - "8080:80"
    environment:
      NODE_ENV: production
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: password
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "Minifile")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mf, err := Parse(path)
	if err != nil {
		t.Fatalf("Failed to parse Minifile: %v", err)
	}

	if mf.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", mf.Version)
	}

	app, ok := mf.Services["app"]
	if !ok {
		t.Fatal("Service 'app' missing")
	}
	if app.Build == nil || app.Build.Context != "." {
		t.Errorf("Expected build context '.', got %v", app.Build)
	}
	if len(app.Ports) != 1 || app.Ports[0] != "8080:80" {
		t.Errorf("Unexpected ports: %v", app.Ports)
	}

	db, ok := mf.Services["db"]
	if !ok {
		t.Fatal("Service 'db' missing")
	}
	if db.Image != "postgres:15-alpine" {
		t.Errorf("Expected image postgres:15-alpine, got %s", db.Image)
	}
}

func TestFindLookUp(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create Minifile in temp dir
	content := "version: 1.0"
	path := filepath.Join(tmpDir, "Minifile")
	os.WriteFile(path, []byte(content), 0644)

	// Change working directory to temp dir for lookup test
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	found, err := FindLookUp()
	if err != nil {
		t.Fatalf("FindLookUp failed: %v", err)
	}
	if filepath.Base(found) != "Minifile" {
		t.Errorf("Expected Minifile, got %s", found)
	}
}
