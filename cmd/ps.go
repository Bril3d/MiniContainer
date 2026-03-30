package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/spf13/cobra"
)

var psJSON bool

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Long:  "List all containers (running and stopped).",
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		containers, err := podman.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		if psJSON {
			data, _ := json.MarshalIndent(containers, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(containers) == 0 {
			fmt.Println("No containers running — start one with `mini run <preset>`")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTAINER ID\tNAME\tIMAGE\tSTATUS\tPORTS")
		fmt.Fprintln(w, "────────────\t────\t─────\t──────\t─────")

		for _, c := range containers {
			id := c.ID
			if len(id) > 12 {
				id = id[:12]
			}

			ports := ""
			for i, p := range c.Ports {
				if i > 0 {
					ports += ", "
				}
				if p.HostPort > 0 {
					ports += fmt.Sprintf("%d->%d/%s", p.HostPort, p.ContainerPort, p.Protocol)
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", id, c.Name, c.Image, c.Status, ports)
		}
		w.Flush()
	},
}

func init() {
	psCmd.Flags().BoolVar(&psJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(psCmd)
}
