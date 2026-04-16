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

		for name, svc := range m.Services {
			fmt.Printf("📦 Service: %s\n", name)

			imageName := svc.Image
			// If build is defined, build it first
			if svc.Build != nil {
				imageName = fmt.Sprintf("%s-%s:latest", m.Project, name)
				fmt.Printf("  🛠️  Building image: %s...\n", imageName)
				err := podman.Build(rt.BuildOptions{
					Context:    svc.Build.Context,
					Dockerfile: svc.Build.Dockerfile,
					Tags:       []string{imageName},
					Args:       svc.Build.Args,
				})
				if err != nil {
					fmt.Printf("  ❌ Build failed for %s: %v\n", name, err)
					continue
				}
			}

			if imageName == "" {
				fmt.Printf("  ⚠️  Skipping %s: no image or build defined\n", name)
				continue
			}

			// Run the container
			containerName := fmt.Sprintf("%s-%s", m.Project, name)
			
			opts := rt.RunOptions{
				Image:  imageName,
				Name:   containerName,
				Ports:  svc.Ports,
				Env:    svc.Environment,
				Cmd:    svc.Command,
				Detach: true,
			}

			id, err := podman.Run(opts)
			if err != nil {
				fmt.Printf("  ❌ Failed to start %s: %v\n", name, err)
				continue
			}
			fmt.Printf("  ✅ Started! ID: %s\n", id[:12])
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
