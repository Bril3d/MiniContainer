package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
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
		fmt.Fprintf(os.Stderr, "Removing container %s...\n", id)

		err := podman.Remove(id, rmForce)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Container removed: %s\n", id)
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Force remove a running container")
	rootCmd.AddCommand(rmCmd)
}
