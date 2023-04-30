package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "📄 list packages installed with pack",
	Run:   List,
}

func List(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"pack package", "related pacman package"})
	mp := ReadMapping()
	for pack, pacman := range mp {
		t.AppendRow(table.Row{pack, pacman})
	}
	t.Render()
}
