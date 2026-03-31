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
	case strings.Contains(msg, "port is already allocated"):
		return fmt.Errorf("port conflict: the requested host port is already in use by another process")
	case strings.Contains(msg, "already in use by container"):
		return fmt.Errorf("naming conflict: a container with that name already exists")
	case strings.Contains(msg, "image not found"):
		return fmt.Errorf("image not found: please check the image name or specify a registry (e.g. docker.io/...)")
	case strings.Contains(msg, "executable file not found in $PATH"):
		return fmt.Errorf("command not found: the specified command does not exist inside the container")
	case strings.Contains(msg, "cgroups v2"):
		return fmt.Errorf("system limitation: resource monitoring requires cgroups v2 in rootless mode")
	case strings.Contains(msg, "failed to connect"):
		return fmt.Errorf("connection error: could not connect to Podman. Is the service running?")
	}

	return err
}
