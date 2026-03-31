package cmd

import (
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rmiCmd = &cobra.Command{
	Use:   "rmi [image]",
	Short: "Remove a local image",
	Long:  "Delete a locally pulled container image by name or ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		image := args[0]
		color.White("Removing image %s...", image)

		err := podman.RemoveImage(image)
		if err != nil {
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Image removed: %s", image)
	},
}

func init() {
	rootCmd.AddCommand(rmiCmd)
}
