package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [image]",
	Short: "Pull a container image",
	Long:  "Download a container image from a registry.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		image := args[0]
		err := podman.Pull(image)
		if err != nil {
			fmt.Fprintln(os.Stderr)
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Image pulled: %s", image)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
