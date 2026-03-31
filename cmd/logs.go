package cmd

import (
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	logsFollow bool
)

var logsCmd = &cobra.Command{
	Use:   "logs [container]",
	Short: "Fetch the logs of a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()
		containerID := args[0]

		err := podman.Logs(containerID, logsFollow)
		if err != nil {
			color.Red("✗ Error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Stream logs in real-time")
	rootCmd.AddCommand(logsCmd)
}
