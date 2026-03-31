package cmd

import (
	"fmt"
	"os"
	"time"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/briandowns/spinner"
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
		
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Stopping container %s...", id)
		s.Start()

		err := podman.Stop(id)
		s.Stop()
		fmt.Fprint(os.Stderr, "\r\033[2K") // Clear spinner line

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
