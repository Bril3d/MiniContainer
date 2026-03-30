package preset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Preset represents a pre-configured container environment.
type Preset struct {
	Image       string            `json:"image"`
	Description string            `json:"description"`
	Ports       []string          `json:"ports"`
	Volumes     []string          `json:"volumes"`
	Env         map[string]string `json:"env"`
	Cmd         string            `json:"cmd"`
}

// Manager handles loading and finding presets.
type Manager struct {
	presets map[string]Preset
}

// NewManager creates a new Manager and loads presets from the given path.
func NewManager(presetsPath string) (*Manager, error) {
	data, err := os.ReadFile(presetsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read presets file: %w", err)
	}

	var presets map[string]Preset
	if err := json.Unmarshal(data, &presets); err != nil {
		return nil, fmt.Errorf("failed to parse presets JSON: %w", err)
	}

	return &Manager{presets: presets}, nil
}

// Find retrieves a preset by name.
func (m *Manager) Find(name string) (Preset, bool) {
	p, ok := m.presets[name]
	return p, ok
}

// List returns all preset names.
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.presets))
	for name := range m.presets {
		names = append(names, name)
	}
	return names
}

// AutoDetect attempts to identify a project in the directory and returns the preset name.
func (m *Manager) AutoDetect(dir string) (string, bool) {
	// Node.js
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		if _, ok := m.presets["node"]; ok {
			return "node", true
		}
	}

	// Python
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		if _, ok := m.presets["python"]; ok {
			return "python", true
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil {
		if _, ok := m.presets["python"]; ok {
			return "python", true
		}
	}

	// Go
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		if _, ok := m.presets["go"]; ok {
			return "go", true
		}
	}

	return "", false
}

// GetDefaultPath returns the standard path for the presets file relative to the project root.
func GetDefaultPath() string {
	// For production, we might want this in a config dir, but for now, 
	// we assume it's in the project's presets/ folder.
	return filepath.Join("presets", "presets.json")
}
