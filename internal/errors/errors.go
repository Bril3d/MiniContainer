package errors

import (
	"fmt"
	"strings"
)

// Humanize returns a user-friendly error message for common container runtime errors.
func Humanize(err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "port is already allocated"), strings.Contains(msg, "bind: address already in use"):
		return fmt.Errorf("port conflict: the requested host port is already in use by another process")
	case strings.Contains(msg, "already in use by container"):
		return fmt.Errorf("naming conflict: a container with that name already exists")
	case strings.Contains(msg, "image not found"), strings.Contains(msg, "no such image"), strings.Contains(msg, "image not known"):
		return fmt.Errorf("image not found: please check the image name or specify a registry (e.g. docker.io/...)")
	case strings.Contains(msg, "executable file not found in $PATH"), strings.Contains(msg, "OCI runtime exec failed: exec failed: container_linux.go"):
		return fmt.Errorf("command not found: the specified command does not exist inside the container")
	case strings.Contains(msg, "cgroups v2"):
		return fmt.Errorf("system limitation: resource monitoring requires cgroups v2 in rootless mode")
	case strings.Contains(msg, "failed to connect"), strings.Contains(msg, "connect: no such file or directory"), strings.Contains(msg, "Cannot connect to Podman"):
		return fmt.Errorf("connection error: MiniContainer could not reach Podman. Please ensure Podman is running (e.g., 'podman machine start')")
	case strings.Contains(msg, "Permission denied"), strings.Contains(msg, "EPERM"):
		return fmt.Errorf("permission denied: try running with elevated privileges or check your Podman rootless configuration")
	case strings.Contains(msg, "Temporary failure in name resolution"):
		return fmt.Errorf("network error: container cannot resolve DNS. Check your host's internet connection or VPN settings")
	}

	return err
}
