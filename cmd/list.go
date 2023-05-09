package cmd

import (
	"strings"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
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

// Cli command listing installed packages and their status.
func List(cmd *cobra.Command, args []string) {
	// Display outdated last.
	// Display outdated last version.
	pkgs := GetPacmanPackages()
	reversePackMapping := ReverseMapping(ReadMapping())
	for pkg, version := range pkgs {
		print.Custom([]print.ColoredMessage{
			{
				Message: pkg + " ",
				Color:   print.WHITE,
			},
			{
				Message: version + " ",
				Color:   print.BLUE,
			},
			{
				Message: reversePackMapping[pkg],
				Color:   print.YELLOW,
			},
		})
	}
}

// Get all installed packages from pacman.
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
