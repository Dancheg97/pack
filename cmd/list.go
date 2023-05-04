package cmd

import (
	"os"
	"strings"

	"fmnx.io/dev/pack/system"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "📄 list packages installed with pack",
	Run:     List,
}

func List(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"pacman package", "version", "pack link"})
	pkgs := GetPacmanPackages()
	mp := ReadMapping()
	revmp := ReverseMapping(mp)
	for k, v := range pkgs {
		t.AppendRow(table.Row{k, v, revmp[k]})
	}
	t.Render()
}

func GetPacmanPackages() map[string]string {
	o, err := system.Call("pacman -Q")
	CheckErr(err)
	o = strings.Trim(o, "\n")
	pkgs := map[string]string{}
	for _, pkg := range strings.Split(o, "\n") {
		spl := strings.Split(pkg, " ")
		pkgs[spl[0]] = spl[1]
	}
	return pkgs
}
