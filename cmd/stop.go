package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
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
		fmt.Fprintf(os.Stderr, "Stopping container %s...\n", id)

		err := podman.Stop(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Container stopped: %s\n", id)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
