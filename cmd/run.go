package cmd

import (
	"fmt"
	"os"

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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		image := args[0]
		containerCmd := args[1:]

		// Build env map from --env flags
		envMap := make(map[string]string)
		for _, e := range runEnvVars {
			parts := splitEnvVar(e)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}

		opts := rt.RunOptions{
			Image:       image,
			Name:        runName,
			Ports:       runPorts,
			Volumes:     runVolumes,
			Env:         envMap,
			Cmd:         containerCmd,
			Detach:      !runInteractive,
			Interactive: runInteractive,
		}

		if runInteractive {
			fmt.Fprintf(os.Stderr, "Starting interactive container from %s...\n", image)
			_, err := podman.Run(opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			return
		}

		fmt.Fprintf(os.Stderr, "Starting container from %s...\n", image)
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
