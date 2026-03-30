package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
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
		fmt.Fprintf(os.Stderr, "Removing image %s...\n", image)

		err := podman.RemoveImage(image)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Image removed: %s\n", image)
	},
}

func init() {
	rootCmd.AddCommand(rmiCmd)
}
