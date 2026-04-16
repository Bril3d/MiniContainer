package cmd

import (
	"fmt"
	"os"

	"github.com/Bril3d/minicontainer/internal/minifile"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers from Minifile",
	Long:  `Parses the Minifile in the current directory and starts all defined services.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := minifile.FindLookUp()
		if err != nil {
			fmt.Printf("❌ %s\n", err)
			os.Exit(1)
		}

		m, err := minifile.Parse(path)
		if err != nil {
			fmt.Printf("❌ Error parsing %s: %v\n", path, err)
			os.Exit(1)
		}

		fmt.Printf("🚀 Starting project: %s\n", m.Project)
		podman := rt.NewPodmanRuntime()
		orch := minifile.NewOrchestrator(podman)

		if err := orch.Up(m); err != nil {
			fmt.Printf("❌ Orchestration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ All services are up!")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
