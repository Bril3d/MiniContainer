package runtime

// ContainerID is the unique identifier assigned by the container runtime.
type ContainerID = string

// RunOptions holds all parameters for starting a new container.
type RunOptions struct {
	Image       string            // Container image (e.g. "docker.io/library/python:3.11")
	Name        string            // Optional container name
	Ports       []string          // Port mappings (e.g. "8080:80")
	Volumes     []string          // Volume mounts (e.g. "./:/app")
	Env         map[string]string // Environment variables
	Cmd         []string          // Command + args to run inside the container
	Detach      bool              // Run in background
	Interactive bool              // Attach stdin/stdout (TTY)
}

// Container represents a running or stopped container.
type Container struct {
	ID      string   `json:"ID"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	Status  string   `json:"Status"`
	State   string   `json:"State"`
	Ports   []Port   `json:"Ports"`
	Created int64    `json:"Created"`
	Command []string `json:"Command"`
}

// Port represents a container port mapping.
type Port struct {
	HostIP        string `json:"host_ip"`
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"`
}

// Image represents a locally pulled container image.
type Image struct {
	ID         string   `json:"ID"`
	Repository string   `json:"Repository"`
	Tag        string   `json:"Tag"`
	Size       int64    `json:"Size"`
	Created    int64    `json:"Created"`
	Names      []string `json:"Names"`
}

// ContainerStats holds resource usage data for a running container.
type ContainerStats struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	CPUPerc  string `json:"CPUPerc"`
	MemUsage string `json:"MemUsage"`
	MemPerc  string `json:"MemPerc"`
	NetIO    string `json:"NetIO"`
	BlockIO  string `json:"BlockIO"`
}

// ContainerRuntime defines the interface for container operations.
// This abstraction decouples MiniContainer from any specific runtime (Podman, containerd, etc.)
// allowing future support for alternative backends.
type ContainerRuntime interface {
	// Run starts a new container with the given options.
	Run(opts RunOptions) (ContainerID, error)

	// Start starts a stopped container by ID or name.
	Start(id string) error

	// Stop stops a running container by ID or name.
	Stop(id string) error

	// Remove deletes a container. If force is true, removes even running containers.
	Remove(id string, force bool) error

	// List returns all containers (running and stopped).
	List() ([]Container, error)

	// Logs streams container logs. If follow is true, streams in real time.
	Logs(id string, follow bool) error

	// Stats returns resource usage for all running containers.
	Stats() ([]ContainerStats, error)

	// Pull downloads a container image.
	Pull(image string) error

	// Images lists all locally available images.
	Images() ([]Image, error)

	// RemoveImage deletes a local image by name or ID.
	RemoveImage(image string) error

	// Pause pauses a running container.
	Pause(id string) error

	// Unpause resumes a paused container.
	Unpause(id string) error

	// Version returns the runtime engine version string.
	Version() (string, error)

	// Exec runs a command inside a running container.
	Exec(id string, cmd []string, interactive bool) error
}

// ExecOptions holds parameters for executing a command in a container.
type ExecOptions struct {
	Container   string   `json:"container"`
	Command     []string `json:"command"`
	Interactive bool     `json:"interactive"`
}
