package minifile

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Parse reads a Minifile from the given path and returns the parsed struct.
func Parse(path string) (*Minifile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Minifile: %w", err)
	}

	var m Minifile
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse Minifile YAML: %w", err)
	}

	// Basic validation
	if m.Version == "" {
		m.Version = "1.0"
	}

	return &m, nil
}

// FindLookUp looks for a file named "Minifile" in the current or parent directories.
func FindLookUp() (string, error) {
	// For now, just look in the current directory
	if _, err := os.Stat("Minifile"); err == nil {
		return "Minifile", nil
	}
	return "", fmt.Errorf("Minifile not found")
}
