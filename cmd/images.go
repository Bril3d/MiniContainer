package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var imagesJSON bool

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List pulled images",
	Long:  "List all locally available container images.",
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		images, err := podman.Images()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		if imagesJSON {
			data, _ := json.MarshalIndent(images, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(images) == 0 {
			fmt.Println("No images pulled yet — try `mini pull alpine`")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tSIZE")
		fmt.Fprintln(w, "──────────\t───\t────────\t────")

		for _, img := range images {
			id := img.ID
			if len(id) > 12 {
				id = id[:12]
			}

			repo := img.Repository
			tag := img.Tag

			// Podman puts the full name:tag in Names[]. Parse repo and tag from it.
			if len(img.Names) > 0 {
				name := img.Names[0]
				if idx := strings.LastIndex(name, ":"); idx > 0 {
					repo = name[:idx]
					tag = name[idx+1:]
				} else {
					repo = name
					tag = "latest"
				}
			}

			if tag == "" {
				tag = "<none>"
			}

			sizeMB := float64(img.Size) / 1024 / 1024

			fmt.Fprintf(w, "%s\t%s\t%s\t%.1f MB\n", repo, tag, id, sizeMB)
		}
		w.Flush()
	},
}

func init() {
	imagesCmd.Flags().BoolVar(&imagesJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(imagesCmd)
}
