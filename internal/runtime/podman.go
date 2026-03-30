package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PodmanRuntime implements ContainerRuntime using Podman CLI.
type PodmanRuntime struct{}

// NewPodmanRuntime creates a new PodmanRuntime instance.
func NewPodmanRuntime() *PodmanRuntime {
	return &PodmanRuntime{}
}

// Version returns the Podman version string.
func (p *PodmanRuntime) Version() (string, error) {
	out, err := Exec("podman", "--version")
	if err != nil {
		return "", fmt.Errorf("podman is not installed or not in PATH.\n\n  Install guide: https://podman.io/getting-started/installation\n\n  On Windows, make sure WSL2 is enabled and Podman is installed inside WSL.")
	}
	return out, nil
}

// Run starts a new container with the given options.
func (p *PodmanRuntime) Run(opts RunOptions) (ContainerID, error) {
	args := []string{"run"}

	if opts.Detach {
		args = append(args, "-d")
	}

	if opts.Interactive {
		args = append(args, "-it")
	}

	if opts.Name != "" {
		args = append(args, "--name", opts.Name)
	}

	for _, port := range opts.Ports {
		args = append(args, "-p", port)
	}

	for _, vol := range opts.Volumes {
		args = append(args, "-v", vol)
	}

	for key, val := range opts.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, val))
	}

	args = append(args, opts.Image)
	args = append(args, opts.Cmd...)

	if opts.Interactive {
		err := ExecInteractive("podman", args...)
		return "", err
	}

	out, err := Exec("podman", args...)
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}
	// Podman returns the full container ID
	return strings.TrimSpace(out), nil
}

// Stop stops a running container by ID or name.
func (p *PodmanRuntime) Stop(id string) error {
	_, err := Exec("podman", "stop", id)
	if err != nil {
		return fmt.Errorf("failed to stop container '%s': %w", id, err)
	}
	return nil
}

// Remove deletes a container. If force is true, removes even running containers.
func (p *PodmanRuntime) Remove(id string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)

	_, err := Exec("podman", args...)
	if err != nil {
		return fmt.Errorf("failed to remove container '%s': %w", id, err)
	}
	return nil
}

// List returns all containers (running by default).
func (p *PodmanRuntime) List() ([]Container, error) {
	out, err := Exec("podman", "ps", "-a", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	if out == "" || out == "[]" || out == "null" {
		return []Container{}, nil
	}

	var containers []Container
	if err := json.Unmarshal([]byte(out), &containers); err != nil {
		return nil, fmt.Errorf("failed to parse container list: %w", err)
	}

	return containers, nil
}

// Logs streams container logs. If follow is true, streams in real time.
func (p *PodmanRuntime) Logs(id string, follow bool) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, id)

	// Stream logs directly to user's terminal
	return ExecStream("podman", args...)
}

// Stats returns resource usage for all running containers.
func (p *PodmanRuntime) Stats() ([]ContainerStats, error) {
	out, err := Exec("podman", "stats", "--no-stream", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if out == "" || out == "[]" || out == "null" {
		return []ContainerStats{}, nil
	}

	var stats []ContainerStats
	if err := json.Unmarshal([]byte(out), &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return stats, nil
}

// Pull downloads a container image with live progress output.
func (p *PodmanRuntime) Pull(image string) error {
	fmt.Fprintf(os.Stderr, "Pulling image: %s\n", image)
	err := ExecStream("podman", "pull", image)
	if err != nil {
		return fmt.Errorf("failed to pull image '%s': %w", image, err)
	}
	return nil
}

// Images lists all locally available images.
func (p *PodmanRuntime) Images() ([]Image, error) {
	out, err := Exec("podman", "images", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	if out == "" || out == "[]" || out == "null" {
		return []Image{}, nil
	}

	var images []Image
	if err := json.Unmarshal([]byte(out), &images); err != nil {
		return nil, fmt.Errorf("failed to parse image list: %w", err)
	}

	return images, nil
}

// RemoveImage deletes a local image by name or ID.
func (p *PodmanRuntime) RemoveImage(image string) error {
	_, err := Exec("podman", "rmi", image)
	if err != nil {
		return fmt.Errorf("failed to remove image '%s': %w", image, err)
	}
	return nil
}

// Compile-time check that PodmanRuntime implements ContainerRuntime.
var _ ContainerRuntime = (*PodmanRuntime)(nil)
