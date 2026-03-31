package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "0.1.0-dev"
var debug bool

var rootCmd = &cobra.Command{
	Use:   "mini",
	Short: "MiniContainer — Developer Environment Launcher",
	Long: `MiniContainer is a lightweight, fast container management tool.
Run dev environments instantly — no Docker complexity, no heavy setup.

Usage:
  mini run <preset|image>    Start a container
  mini ps                    List running containers
  mini stop <container>      Stop a container
  mini rm <container>        Remove a container
  mini images                List pulled images
  mini presets               List available presets`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rt.SetDebug(debug)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("MiniContainer v{{.Version}}\n")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show raw commands and execution timing")
}
