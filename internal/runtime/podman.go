package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Bril3d/minicontainer/internal/errors"
)

// PodmanRuntime implements ContainerRuntime using Podman CLI.
type PodmanRuntime struct {
	// useWSL indicates we should run podman through WSL on Windows
	useWSL    bool
	wslDistro string
}

// NewPodmanRuntime creates a new PodmanRuntime instance.
// On Windows, it auto-detects whether to use WSL to reach Podman.
func NewPodmanRuntime() *PodmanRuntime {
	p := &PodmanRuntime{}

	if runtime.GOOS == "windows" {
		// Try native podman first
		_, err := Exec("podman", "info", "--format", "{{.Host.Os}}")
		if err != nil {
			// Native podman connection failed — try via WSL
			_, wslErr := Exec("wsl", "-d", "podman-machine-default", "--", "podman", "--version")
			if wslErr == nil {
				p.useWSL = true
				p.wslDistro = "podman-machine-default"
			}
		}
	}

	return p
}

// podmanCmd returns the base command and prefix args for calling podman.
// On Windows with WSL fallback, this wraps the call through `wsl -d <distro> --`.
func (p *PodmanRuntime) podmanCmd() (string, []string) {
	if p.useWSL {
		return "wsl", []string{"-d", p.wslDistro, "--", "podman"}
	}
	return "podman", nil
}

// buildArgs prepends WSL prefix args (if needed) to the podman subcommand args.
func (p *PodmanRuntime) buildArgs(podmanArgs ...string) (string, []string) {
	cmd, prefix := p.podmanCmd()
	return cmd, append(prefix, podmanArgs...)
}

// Version returns the Podman version string.
func (p *PodmanRuntime) Version() (string, error) {
	cmd, args := p.buildArgs("--version")
	out, err := Exec(cmd, args...)
	if err != nil {
		return "", errors.Humanize(fmt.Errorf("podman is not installed or not reachable.\n\n  Install guide: https://podman.io/getting-started/installation\n\n  On Windows, make sure WSL2 is enabled and Podman is installed inside WSL."))
	}
	return out, nil
}

// Run starts a new container with the given options.
func (p *PodmanRuntime) Run(opts RunOptions) (ContainerID, error) {
	podmanArgs := []string{"run"}

	if opts.Detach {
		podmanArgs = append(podmanArgs, "-d")
		// Add -i (Interactive) and -t (TTY) by default for detached containers to keep them alive
		// if they don't have a long-running foreground process (e.g. Python REPL).
		podmanArgs = append(podmanArgs, "-i", "-t")
	}

	if opts.Interactive {
		podmanArgs = append(podmanArgs, "-it")
	}

	if opts.Name != "" {
		podmanArgs = append(podmanArgs, "--name", opts.Name)
	}

	for _, port := range opts.Ports {
		podmanArgs = append(podmanArgs, "-p", port)
	}

	for _, vol := range opts.Volumes {
		podmanArgs = append(podmanArgs, "-v", vol)
	}

	for key, val := range opts.Env {
		podmanArgs = append(podmanArgs, "-e", fmt.Sprintf("%s=%s", key, val))
	}

	podmanArgs = append(podmanArgs, opts.Image)
	podmanArgs = append(podmanArgs, opts.Cmd...)

	cmd, args := p.buildArgs(podmanArgs...)

	if opts.Interactive {
		err := ExecInteractive(cmd, args...)
		return "", err
	}

	out, err := Exec(cmd, args...)
	if err != nil {
		return "", errors.Humanize(fmt.Errorf("failed to start container: %w", err))
	}
	return strings.TrimSpace(cleanOutput(out)), nil
}

// containerInspect holds the subset of `podman inspect` we need.
type containerInspect struct {
	Config struct {
		Tty       bool     `json:"Tty"`
		OpenStdin bool     `json:"OpenStdin"`
		Cmd       []string `json:"Cmd"`
		Image     string   `json:"Image"`
	} `json:"Config"`
	Name            string          `json:"Name"`
	HostConfig      json.RawMessage `json:"HostConfig"`
}

// Start starts a stopped container by ID or name.
// If the container was created without TTY/Stdin (legacy), it recreates it
// with proper flags so interactive processes (python, node) stay alive.
func (p *PodmanRuntime) Start(id string) error {
	// First, try a normal start
	cmd, args := p.buildArgs("start", id)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to start container '%s': %w", id, err))
	}

	// Check if the container actually stayed running
	cmd2, args2 := p.buildArgs("inspect", "--format", "{{.State.Running}}", id)
	out, err := Exec(cmd2, args2...)
	if err != nil {
		return nil // started but can't inspect — don't error
	}

	running := strings.TrimSpace(cleanOutput(out))
	if running == "true" {
		return nil // container is running, all good
	}

	// Container exited immediately — likely missing TTY.
	// Inspect it to get config, then recreate with -it flags.
	cmd3, args3 := p.buildArgs("inspect", "--format", "json", id)
	inspectOut, err := Exec(cmd3, args3...)
	if err != nil {
		return nil // can't inspect, nothing more we can do
	}

	inspectOut = cleanJSON(inspectOut)
	var inspections []containerInspect
	if err := json.Unmarshal([]byte(inspectOut), &inspections); err != nil || len(inspections) == 0 {
		return nil
	}
	info := inspections[0]

	// Only attempt recreate if Tty was not set (legacy container)
	if info.Config.Tty {
		return nil // container had TTY but still exited — nothing we can do
	}

	// Get the container's image, name, ports, etc for recreation
	containerName := strings.TrimPrefix(info.Name, "/")
	image := info.Config.Image

	// Remove the old container
	_ = p.Remove(id, true)

	// Recreate with -it flags so interactive processes stay alive
	recreateArgs := []string{"run", "-d", "-i", "-t", "--name", containerName}
	recreateArgs = append(recreateArgs, image)
	recreateArgs = append(recreateArgs, info.Config.Cmd...)

	cmd4, args4 := p.buildArgs(recreateArgs...)
	_, err = Exec(cmd4, args4...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to recreate container '%s' with TTY: %w", id, err))
	}

	return nil
}

// Stop stops a running container by ID or name.
func (p *PodmanRuntime) Stop(id string) error {
	cmd, args := p.buildArgs("stop", id)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to stop container '%s': %w", id, err))
	}
	return nil
}

// Remove deletes a container. If force is true, removes even running containers.
func (p *PodmanRuntime) Remove(id string, force bool) error {
	podmanArgs := []string{"rm"}
	if force {
		podmanArgs = append(podmanArgs, "-f")
	}
	podmanArgs = append(podmanArgs, id)

	cmd, args := p.buildArgs(podmanArgs...)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to remove container '%s': %w", id, err))
	}
	return nil
}

// List returns all containers (running and stopped).
func (p *PodmanRuntime) List() ([]Container, error) {
	cmd, args := p.buildArgs("ps", "-a", "--format", "json")
	out, err := Exec(cmd, args...)
	if err != nil {
		return nil, errors.Humanize(fmt.Errorf("failed to list containers: %w", err))
	}

	if out == "" {
		return []Container{}, nil
	}

	out = cleanJSON(out)
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
	podmanArgs := []string{"logs"}
	if follow {
		podmanArgs = append(podmanArgs, "-f")
	}
	podmanArgs = append(podmanArgs, id)

	cmd, args := p.buildArgs(podmanArgs...)
	return ExecStream(cmd, args...)
}

// Stats returns resource usage for all running containers.
func (p *PodmanRuntime) Stats() ([]ContainerStats, error) {
	cmd, args := p.buildArgs("stats", "--no-stream", "--format", "json")
	out, err := Exec(cmd, args...)
	if err != nil {
		errMsg := err.Error()
		// Known Podman/WSL issues that should not crash the UI
		if strings.Contains(errMsg, "cgroups v2") ||
			strings.Contains(errMsg, "unknown FS magic") ||
			strings.Contains(errMsg, "netns") {
			return []ContainerStats{}, nil
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if out == "" {
		return []ContainerStats{}, nil
	}

	out = cleanJSON(out)
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
	cmd, args := p.buildArgs("pull", image)
	err := ExecStream(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to pull image '%s': %w", image, err))
	}
	return nil
}

// Images lists all locally available images.
func (p *PodmanRuntime) Images() ([]Image, error) {
	cmd, args := p.buildArgs("images", "--format", "json")
	out, err := Exec(cmd, args...)
	if err != nil {
		return nil, errors.Humanize(fmt.Errorf("failed to list images: %w", err))
	}

	if out == "" {
		return []Image{}, nil
	}

	out = cleanJSON(out)
	if out == "" || out == "[]" || out == "null" {
		return []Image{}, nil
	}

	var images []Image
	if err := json.Unmarshal([]byte(out), &images); err != nil {
		return nil, fmt.Errorf("failed to parse image list: %w", err)
	}

	// Post-process to ensure Repository and Tag are populated from Names if empty
	for i := range images {
		if (images[i].Repository == "" || images[i].Repository == "<none>") && len(images[i].Names) > 0 {
			fullName := images[i].Names[0]
			// Split by colon to get repo and tag
			parts := strings.Split(fullName, ":")
			if len(parts) > 1 {
				images[i].Repository = strings.Join(parts[:len(parts)-1], ":")
				images[i].Tag = parts[len(parts)-1]
			} else {
				images[i].Repository = fullName
				images[i].Tag = "latest"
			}
		}
	}

	return images, nil
}

// RemoveImage deletes a local image by name or ID.
func (p *PodmanRuntime) RemoveImage(image string) error {
	cmd, args := p.buildArgs("rmi", image)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to remove image '%s': %w", image, err))
	}
	return nil
}

// Pause pauses a running container.
func (p *PodmanRuntime) Pause(id string) error {
	cmd, args := p.buildArgs("pause", id)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to pause container '%s': %w", id, err))
	}
	return nil
}

// Unpause resumes a paused container.
func (p *PodmanRuntime) Unpause(id string) error {
	cmd, args := p.buildArgs("unpause", id)
	_, err := Exec(cmd, args...)
	if err != nil {
		return errors.Humanize(fmt.Errorf("failed to unpause container '%s': %w", id, err))
	}
	return nil
}

// Exec runs a command inside a running container.
func (p *PodmanRuntime) Exec(id string, cmdArgs []string, interactive bool) error {
	podmanArgs := []string{"exec"}
	if interactive {
		podmanArgs = append(podmanArgs, "-it")
	}
	podmanArgs = append(podmanArgs, id)
	podmanArgs = append(podmanArgs, cmdArgs...)

	cmd, args := p.buildArgs(podmanArgs...)

	if interactive {
		return ExecInteractive(cmd, args...)
	}

	return ExecStream(cmd, args...)
}

// cleanJSON strips leading diagnostic warnings or non-JSON lines often prepended by Podman/WSL.
func cleanJSON(out string) string {
	// Find the start of the JSON array or object
	idx := strings.IndexAny(out, "[{")
	if idx == -1 {
		return out
	}
	return out[idx:]
}

// cleanOutput handles stripping diagnostic warnings from standard command output (non-JSON).
func cleanOutput(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return ""
	}
	
	// Podman often prepends diagnostics like "time=..." or "level=...".
	// The actual result (like container ID) is usually on the LAST line.
	return lines[len(lines)-1]
}

// Compile-time check that PodmanRuntime implements ContainerRuntime.
var _ ContainerRuntime = (*PodmanRuntime)(nil)
