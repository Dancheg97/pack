// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "⬆ push package",
	Long: `⬆ push packages

This command will search for built package in cache, then it will try to upload
this package to compatible pack repo.

You can provide multiple packge names, all of them will be installed.`,
	Run: Push,
}

func Push(cmd *cobra.Command, args []string) {

}
