package minifile

import (
	"fmt"
	"strings"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
)

// Orchestrator handles the lifecycle of a project defined by a Minifile.
type Orchestrator struct {
	runtime rt.ContainerRuntime
}

func NewOrchestrator(runtime rt.ContainerRuntime) *Orchestrator {
	return &Orchestrator{runtime: runtime}
}

// Up builds and starts all services defined in the Minifile.
func (o *Orchestrator) Up(m *Minifile) error {
	for name, svc := range m.Services {
		imageName := svc.Image
		if svc.Build != nil {
			imageName = fmt.Sprintf("%s-%s:latest", m.Project, name)
			err := o.runtime.Build(rt.BuildOptions{
				Context:    svc.Build.Context,
				Dockerfile: svc.Build.Dockerfile,
				Tags:       []string{imageName},
				Args:       svc.Build.Args,
			})
			if err != nil {
				return fmt.Errorf("build failed for %s: %w", name, err)
			}
		}

		if imageName == "" {
			continue
		}

		containerName := fmt.Sprintf("%s-%s", m.Project, name)
		opts := rt.RunOptions{
			Image:  imageName,
			Name:   containerName,
			Ports:  svc.Ports,
			Env:    svc.Environment,
			Cmd:    svc.Command,
			Detach: true,
		}

		_, err := o.runtime.Run(opts)
		if err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err)
		}
	}
	return nil
}

// Down stops and removes all containers matching the project services.
func (o *Orchestrator) Down(m *Minifile) error {
	containers, err := o.runtime.List()
	if err != nil {
		return err
	}

	for _, c := range containers {
		for serviceName := range m.Services {
			expectedPrefix := fmt.Sprintf("%s-%s", m.Project, serviceName)
			match := false
			for _, n := range c.Names {
				cleanName := strings.TrimPrefix(n, "/")
				if cleanName == expectedPrefix {
					match = true
					break
				}
			}

			if match {
				if c.State == "running" {
					_ = o.runtime.Stop(c.ID)
				}
				if err := o.runtime.Remove(c.ID, true); err != nil {
					return fmt.Errorf("failed to remove %s: %w", expectedPrefix, err)
				}
			}
		}
	}
	return nil
}
