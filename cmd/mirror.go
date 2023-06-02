// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mirrCmd)
}

var mirrCmd = &cobra.Command{
	Use:     "mirror",
	Aliases: []string{"m"},
	Short:   "🪞 launch fs mirror",
	Run:     Mirror,
}

func Mirror(cmd *cobra.Command, args []string) {

}
