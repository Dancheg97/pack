// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/print"
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
	pkgs := pacman.List()
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
