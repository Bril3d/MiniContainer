package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [CONTAINER]",
	Short: "Start a stopped container",
	Long:  "Restart a previously stopped container by name or ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()
		name := args[0]

		if err := podman.Start(name); err != nil {
			color.Red("✗ Failed to start container '%s': %v", name, err)
			os.Exit(1)
		}

		fmt.Printf("Container '%s' started successfully.\n", name)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
