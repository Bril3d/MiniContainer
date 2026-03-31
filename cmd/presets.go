package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	ps "github.com/Bril3d/minicontainer/internal/preset"
	"github.com/spf13/cobra"
)

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "List available presets",
	Long:  "Show all pre-configured container environments that can be launched with `mini run <preset>`.",
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := ps.NewManager(ps.GetDefaultPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading presets: %v\n", err)
			os.Exit(1)
		}

		names := mgr.List()
		sort.Strings(names)

		if len(names) == 0 {
			fmt.Println("No presets available.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PRESET\tIMAGE\tDESCRIPTION")
		fmt.Fprintln(w, "──────\t─────\t───────────")

		for _, name := range names {
			p, _ := mgr.Find(name)
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, p.Image, p.Description)
		}
		w.Flush()

		fmt.Printf("\nRun a preset: mini run <preset>\n")
	},
}

func init() {
	rootCmd.AddCommand(presetsCmd)
}
