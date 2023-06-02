// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:     "gitsync",
	Aliases: []string{"gitsync"},
	Short:   "🧿 rebuild git repositories",
	Run:     Watch,
}

func Watch(cmd *cobra.Command, args []string) {

}
