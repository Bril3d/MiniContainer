package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var buildTags []string
var dockerfilePath string

var buildCmd = &cobra.Command{
	Use:   "build [context]",
	Short: "Build an image from a Dockerfile",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()
		
		context := "."
		if len(args) > 0 {
			context = args[0]
		}

		opts := rt.BuildOptions{
			Context:    context,
			Dockerfile: dockerfilePath,
			Tags:       buildTags,
		}

		fmt.Printf("Building image in %s...\n", context)
		if err := podman.Build(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Build successful!")
	},
}

func init() {
	buildCmd.Flags().StringSliceVarP(&buildTags, "tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringVarP(&dockerfilePath, "file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	rootCmd.AddCommand(buildCmd)
}
