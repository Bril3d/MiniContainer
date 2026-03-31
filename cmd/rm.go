package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/Bril3d/minicontainer/internal/ui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rmForce bool

var rmCmd = &cobra.Command{
	Use:   "rm [container]",
	Short: "Remove a container",
	Long:  "Remove a stopped container. Use --force to remove a running container.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		id := args[0]

		if !rmForce {
			if !ui.Confirm(fmt.Sprintf("Are you sure you want to remove container %s?", id)) {
				color.Yellow("Aborted.")
				return
			}
		}

		color.White("Removing container %s...", id)

		err := podman.Remove(id, rmForce)
		if err != nil {
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Container removed: %s", id)
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Force remove a running container")
	rootCmd.AddCommand(rmCmd)
}
