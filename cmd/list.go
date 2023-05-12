// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"strings"

	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   tmpl.ListShort,
	Run:     List,
}

// Cli command listing installed packages and version.
func List(cmd *cobra.Command, args []string) {
	pkgs := GetPacmanPackages()
	for pkg, version := range pkgs {
		i, err := packdb.Get(pkg, packdb.PACMAN)
		if err != nil {
			print.Custom([]print.ColoredMessage{
				{
					Message: pkg + " ",
					Color:   print.WHITE,
				},
				{
					Message: version,
					Color:   print.BLUE,
				},
			})
			continue
		}
		print.Custom([]print.ColoredMessage{
			{
				Message: i.PacmanName + " ",
				Color:   print.WHITE,
			},
			{
				Message: i.PackName + " ",
				Color:   print.YELLOW,
			},
			{
				Message: i.Branch,
				Color:   print.BLUE,
			},
			{
				Message: "-",
				Color:   print.WHITE,
			},
			{
				Message: i.Version,
				Color:   print.BLUE,
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
