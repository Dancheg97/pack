// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"u"},
	Short:   "〽️ upgrade system",
	Run:     Upgrade,
}

func Upgrade(cmd *cobra.Command, args []string) {

}
