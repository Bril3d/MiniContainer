package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	statsWatch bool
	statsJSON  bool
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display a live stream of container(s) resource usage statistics",
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		for {
			stats, err := podman.Stats()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if statsJSON {
				data, _ := json.MarshalIndent(stats, "", "  ")
				fmt.Println(string(data))
				return // JSON doesn't support watching in this simple implementation
			}

			// Clear screen if watching
			if statsWatch {
				fmt.Print("\033[H\033[2J") // Clear screen
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "CONTAINER ID\tNAME\tCPU %\tMEM USAGE / LIMIT\tMEM %\tNET I/O")

			for _, s := range stats {
				id := s.ID
				if len(id) > 12 {
					id = id[:12]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", id, s.Name, s.CPUPerc, s.MemUsage, s.MemPerc, s.NetIO)
			}
			w.Flush()

			if !statsWatch {
				break
			}

			time.Sleep(2 * time.Second)
		}
	},
}

func init() {
	statsCmd.Flags().BoolVarP(&statsWatch, "watch", "w", false, "Continuously watch statistics")
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(statsCmd)
}
