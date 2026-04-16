package minifile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	// Try loading global .env
	globalEnv := o.loadEnvFile(m.BaseDir)

	for name, svc := range m.Services {
		imageName := svc.Image
		if svc.Build != nil {
			imageName = fmt.Sprintf("%s-%s:latest", m.Project, name)
			
			// Resolve relative build context
			ctx := svc.Build.Context
			if !filepath.IsAbs(ctx) {
				ctx = filepath.Join(m.BaseDir, ctx)
			}

			err := o.runtime.Build(rt.BuildOptions{
				Context:    ctx,
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
		
		// Merge environments: .env (global) < Minifile (per-service)
		env := make(map[string]string)
		for k, v := range globalEnv {
			env[k] = v
		}
		for k, v := range svc.Environment {
			env[k] = v
		}

		// Resolve relative volume paths
		volumes := make([]string, len(svc.Volumes))
		for i, v := range svc.Volumes {
			volumes[i] = o.resolveVolumePath(m.BaseDir, v)
		}

		opts := rt.RunOptions{
			Image:  imageName,
			Name:   containerName,
			Ports:  svc.Ports,
			Env:    env,
			Volumes: volumes,
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

func (o *Orchestrator) resolveVolumePath(baseDir, vol string) string {
	parts := strings.Split(vol, ":")
	if len(parts) < 2 {
		return vol // Not a bind mount or malformed
	}

	hostPath := parts[0]
	if !filepath.IsAbs(hostPath) {
		hostPath = filepath.Join(baseDir, hostPath)
	}

	// Reconstruct with mode if present
	return strings.Join(append([]string{hostPath}, parts[1:]...), ":")
}

func (o *Orchestrator) loadEnvFile(dir string) map[string]string {
	env := make(map[string]string)
	path := filepath.Join(dir, ".env")
	
	f, err := os.Open(path)
	if err != nil {
		return env // No .env file, just return empty
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return env
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
