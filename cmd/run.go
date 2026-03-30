package cmd

import (
	"fmt"
	"os"
	"strings"

	ps "github.com/Bril3d/minicontainer/internal/preset"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
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
					fmt.Fprintf(os.Stderr, "✓ Auto-detected project type, using preset: %s\n", detected)
					image = detected
					presetName = detected
				}
			}
		}

		if image == "" {
			fmt.Fprintf(os.Stderr, "Error: No project detected and no image/preset specified.\n\n")
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

		desc := image
		if presetName != "" {
			desc = fmt.Sprintf("%s (preset: %s)", image, presetName)
		}

		if runInteractive {
			fmt.Fprintf(os.Stderr, "Starting interactive container from %s...\n", desc)
			_, err := podman.Run(opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			return
		}

		fmt.Fprintf(os.Stderr, "Starting container from %s...\n", desc)
		containerID, err := podman.Run(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		shortID := containerID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}
		fmt.Printf("✓ Container started: %s\n", shortID)
	},
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
	runCmd.Flags().StringVar(&runName, "name", "", "Assign a name to the container")
	runCmd.Flags().StringSliceVarP(&runPorts, "port", "p", nil, "Publish port (e.g. 8080:80)")
	runCmd.Flags().StringSliceVarP(&runVolumes, "volume", "v", nil, "Bind mount volume (e.g. ./:/app)")
	runCmd.Flags().StringSliceVarP(&runEnvVars, "env", "e", nil, "Set environment variable (e.g. KEY=VALUE)")
	runCmd.Flags().BoolVarP(&runInteractive, "interactive", "i", false, "Run in interactive mode (attach TTY)")
	rootCmd.AddCommand(runCmd)
}
