package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
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
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Image pulled: %s\n", image)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
