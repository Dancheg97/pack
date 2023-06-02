// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/pacman"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "🚫 remove packages",
	Run:     Remove,
}

func Remove(cmd *cobra.Command, args []string) {
	err := pacman.RemoveList(args)
	CheckErr(err)
}
