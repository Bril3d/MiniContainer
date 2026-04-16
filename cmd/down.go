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

		// Get all containers to find project matches
		containers, err := podman.List()
		if err != nil {
			fmt.Printf("❌ Failed to list containers: %v\n", err)
			os.Exit(1)
		}

		for _, c := range containers {
			// Basic name matching: project-service
			// This could be more robust by using labels in the future
			for serviceName := range m.Services {
				expectedPrefix := fmt.Sprintf("%s-%s", m.Project, serviceName)
				match := false
				for _, n := range c.Names {
					if n == expectedPrefix || n == "/"+expectedPrefix {
						match = true
						break
					}
				}

				if match {
					fmt.Printf("🗑️  Removing %s (%s)...\n", expectedPrefix, c.ID[:12])
					
					// Stop first
					if c.State == "running" {
						_ = podman.Stop(c.ID)
					}
					
					// Remove - boolean force = true
					if err := podman.Remove(c.ID, true); err != nil {
						fmt.Printf("  ❌ Failed to remove %s: %v\n", expectedPrefix, err)
					} else {
						fmt.Printf("  ✅ Removed\n")
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
