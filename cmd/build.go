package cmd

import (
	"fmt"
	"os"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	buildTags       []string
	buildFile       string
	buildArgs       []string
)

var buildCmd = &cobra.Command{
	Use:   "build [CONTEXT]",
	Short: "Build an image from a Dockerfile",
	Long:  `Create a new container image using a Dockerfile. By default, it looks for 'Dockerfile' in the current directory.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := "."
		if len(args) > 0 {
			context = args[0]
		}

		podman := rt.NewPodmanRuntime()

		// Parse build-args (key=value)
		argMap := make(map[string]string)
		// For now, simple implementation, can be expanded
		
		opts := rt.BuildOptions{
			Context:    context,
			Dockerfile: buildFile,
			Tags:       buildTags,
			Args:       argMap,
		}

		fmt.Printf("🏗️  Building image in context: %s\n", context)
		if err := podman.Build(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Build successful!")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringSliceVarP(&buildTags, "tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringVarP(&buildFile, "file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
}
