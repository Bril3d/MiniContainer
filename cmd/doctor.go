package cmd

import (
	"fmt"
	"runtime"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system requirements and health",
	Long:  "Verify that Podman and all dependencies are correctly installed and configured.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 MiniContainer Health Check")
		fmt.Println("─────────────────────────────")

		// Check OS
		fmt.Printf("\n  OS:       %s/%s\n", runtime.GOOS, runtime.GOARCH)

		// Check Podman
		podman := rt.NewPodmanRuntime()
		version, err := podman.Version()
		if err != nil {
			fmt.Println("  Podman:   ✗ Not found")
			fmt.Printf("\n  %s\n", err)
		} else {
			fmt.Printf("  Podman:   ✓ %s\n", version)
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
