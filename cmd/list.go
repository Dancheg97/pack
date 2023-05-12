// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
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
		i, err := pack.GetByPacmanName(pkg)
		if err != nil {
			prnt.Custom([]prnt.ColoredMessage{
				{
					Message: pkg + " ",
					Color:   prnt.COLOR_WHITE,
				},
				{
					Message: version,
					Color:   prnt.COLOR_BLUE,
				},
			})
			continue
		}
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: i.PacmanName + " ",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: i.PackName + " ",
				Color:   prnt.COLOR_YELLOW,
			},
			{
				Message: i.DefaultBranch,
				Color:   prnt.COLOR_BLUE,
			},
			{
				Message: "-",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: i.Version,
				Color:   prnt.COLOR_BLUE,
			},
		})
	}
}
