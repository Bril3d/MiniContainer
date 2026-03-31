package cmd

import (
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [container]",
	Short: "Stop a running container",
	Long:  "Stop a running container by ID or name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		id := args[0]
		color.White("Stopping container %s...", id)

		err := podman.Stop(id)
		if err != nil {
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Container stopped: %s", id)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
