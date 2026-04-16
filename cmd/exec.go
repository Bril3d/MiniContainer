package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	execInteractive bool
)

var execCmd = &cobra.Command{
	Use:   "exec [container] [command...]",
	Short: "Execute a command inside a running container",
	Long: `Run a process inside a running container.

Examples:
  mini exec (id|name) ls -l /
  mini exec -i (id|name) sh`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		containerID := args[0]
		commandArgs := args[1:]

		err := podman.Exec(containerID, commandArgs, execInteractive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&execInteractive, "interactive", "i", false, "Keep STDIN open and allocate a pseudo-TTY")
}
