package cmd

import (
	"fmt"
	"os"

	"github.com/Bril3d/minicontainer/internal/minifile"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers from Minifile",
	Long:  `Parses the Minifile and removes all containers associated with the project.`,
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

		fmt.Printf("🛑 Stopping project: %s\n", m.Project)
		podman := rt.NewPodmanRuntime()
		orch := minifile.NewOrchestrator(podman)

		if err := orch.Down(m); err != nil {
			fmt.Printf("❌ Shutdown failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ All project containers removed.")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
