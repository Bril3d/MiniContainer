package minifile

// Minifile represents the YAML configuration for a project.
type Minifile struct {
	Version  string             `yaml:"version"`
	Project  string             `yaml:"project"`
	Services map[string]Service `yaml:"services"`
}

// Service represents a single container service in the Minifile.
type Service struct {
	Image       string            `yaml:"image,omitempty"`
	Build       *BuildConfig      `yaml:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
}

// BuildConfig holds instructions for building a service image.
type BuildConfig struct {
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
}
