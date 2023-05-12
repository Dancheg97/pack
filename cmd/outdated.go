// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file corresponding a single command, including root cmd.

import (
	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(outdatedCmd)
}

var outdatedCmd = &cobra.Command{
	Use:     "outdated",
	Aliases: []string{"out", "o"},
	Short:   tmpl.OutdatedShort,
	Long:    tmpl.OutdatedLong,
	Run:     Outdated,
}

// Cli command listing installed packages and their status.
func Outdated(cmd *cobra.Command, args []string) {
	pacmanOutdated, err := pacman.Outdated()
	CheckErr(err)
	packoutdated := pack.Outdated()
	allOutdated := append(pacmanOutdated, packoutdated...)
	for _, info := range allOutdated {
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: info.Name + " ",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: info.CurrentVersion + " ",
				Color:   prnt.COLOR_YELLOW,
			},
			{
				Message: info.NewVersion,
				Color:   prnt.COLOR_BLUE,
			},
		})
	}
}
