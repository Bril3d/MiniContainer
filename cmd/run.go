package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ps "github.com/Bril3d/minicontainer/internal/preset"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/Bril3d/minicontainer/internal/ui"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	runName        string
	runPorts       []string
	runVolumes     []string
	runEnvVars     []string
	runInteractive bool
)

var runCmd = &cobra.Command{
	Use:   "run [image] [command...]",
	Short: "Run a container from an image or preset",
	Long: `Start a new container from a container image or a preset name.

Examples:
  mini run docker.io/library/alpine sleep 300
  mini run --name myapp --port 8080:80 docker.io/library/nginx
  mini run -i docker.io/library/alpine sh
  mini run python    (uses preset)`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		// Load presets
		presetMgr, err := ps.NewManager(ps.GetDefaultPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load presets: %v\n", err)
		}

		var image string
		var presetName string
		containerCmd := []string{}

		if len(args) > 0 {
			image = args[0]
			containerCmd = args[1:]
		} else {
			// Auto-detection
			if presetMgr != nil {
				if detected, ok := presetMgr.AutoDetect("."); ok {
					color.Green("✓ Auto-detected project type, using preset: %s", detected)
					image = detected
					presetName = detected
				}
			}
		}

		if image == "" {
			color.Red("✗ No project detected and no image/preset specified.")
			fmt.Fprintln(os.Stderr)
			if presetMgr != nil {
				fmt.Fprintf(os.Stderr, "Available presets: %s\n", strings.Join(presetMgr.List(), ", "))
			}
			fmt.Fprintf(os.Stderr, "\nUsage: mini run [preset|image] [command...] [flags]\n")
			os.Exit(1)
		}

		var activePreset ps.Preset
		hasPreset := false
		if presetMgr != nil {
			if p, ok := presetMgr.Find(image); ok {
				activePreset = p
				hasPreset = true
				presetName = image
				image = p.Image // Use image from preset
				if len(containerCmd) == 0 && p.Cmd != "" {
					containerCmd = []string{p.Cmd}
				}
			}
		}

		// Build env map
		envMap := make(map[string]string)
		if hasPreset {
			for k, v := range activePreset.Env {
				envMap[k] = v
			}
		}
		// CLI env flags override preset env
		for _, e := range runEnvVars {
			parts := splitEnvVar(e)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}

		// Build merged ports and volumes
		allPorts := runPorts
		if hasPreset {
			allPorts = append(activePreset.Ports, runPorts...)
		}

		allVolumes := runVolumes
		if hasPreset {
			allVolumes = append(activePreset.Volumes, runVolumes...)
		}

		opts := rt.RunOptions{
			Image:       image,
			Name:        runName,
			Ports:       allPorts,
			Volumes:     allVolumes,
			Env:         envMap,
			Cmd:         containerCmd,
			Detach:      !runInteractive,
			Interactive: runInteractive,
		}

		// Auto-resolve port conflicts for ALL runs (not just presets)
		for i, pStr := range opts.Ports {
			hostPort := parseHostPort(pStr)
			if hostPort > 0 && !rt.IsPortAvailable(hostPort) {
				newPort := rt.FindAvailablePort(hostPort + 1)
				if newPort == -1 {
					color.Red("✗ Error: Could not find an available port in range %d-%d", hostPort+1, hostPort+100)
					os.Exit(1)
				}
				// Rebuild the port string with the new host port
				opts.Ports[i] = replaceHostPort(pStr, newPort)
				color.Yellow("⚠ Port %d is in use. Auto-resolved to %d", hostPort, newPort)
			}
		}

		if runInteractive {
			color.Cyan("Starting interactive container from %s...", image)
			_, err := podman.Run(opts)
			if err != nil {
				color.Red("✗ Error: %v", err)
				os.Exit(1)
			}
			return
		}

		msg := fmt.Sprintf("Starting container from %s", opts.Image)
		if presetName != "" {
			msg = fmt.Sprintf("Starting container from %s (preset: %s)", opts.Image, presetName)
		}
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " " + msg + "..."
		s.Start()

		containerID, err := podman.Run(opts)
		s.Stop()
		fmt.Fprint(os.Stderr, "\033[2K\r") // Full line clear

		if err != nil {
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}

		shortID := containerID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}
		color.Green("✓ Container started: %s", shortID)

		// Auto-open browser for web presets
		if hasPreset && activePreset.Web && len(opts.Ports) > 0 {
			hostPort := parseHostPort(opts.Ports[0])
			if hostPort > 0 {
				url := fmt.Sprintf("http://localhost:%d", hostPort)
				color.Cyan("🔗 Opening browser: %s", url)
				_ = ui.OpenBrowser(url)
			}
		}
	},
}

// parseHostPort extracts the host port from a port mapping string like "8080:80" or "8080:80/tcp".
func parseHostPort(portStr string) int {
	// Format: hostPort:containerPort or hostPort:containerPort/proto
	parts := strings.SplitN(portStr, ":", 2)
	if len(parts) < 2 {
		return 0
	}
	port, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return port
}

// replaceHostPort replaces the host port in a port mapping string.
func replaceHostPort(portStr string, newHostPort int) string {
	parts := strings.SplitN(portStr, ":", 2)
	if len(parts) < 2 {
		return portStr
	}
	return fmt.Sprintf("%d:%s", newHostPort, parts[1])
}

// splitEnvVar splits "KEY=VALUE" into ["KEY", "VALUE"].
func splitEnvVar(s string) []string {
	for i, c := range s {
		if c == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func init() {
	runCmd.Flags().SetInterspersed(false)
	runCmd.Flags().StringVar(&runName, "name", "", "Assign a name to the container")
	runCmd.Flags().StringSliceVarP(&runPorts, "port", "p", nil, "Publish port (e.g. 8080:80)")
	runCmd.Flags().StringSliceVarP(&runVolumes, "volume", "v", nil, "Bind mount volume (e.g. ./:/app)")
	runCmd.Flags().StringSliceVarP(&runEnvVars, "env", "e", nil, "Set environment variable (e.g. KEY=VALUE)")
	runCmd.Flags().BoolVarP(&runInteractive, "interactive", "i", false, "Run in interactive mode (attach TTY)")
	rootCmd.AddCommand(runCmd)
}
