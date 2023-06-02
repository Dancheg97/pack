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
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "🛠️ build package",
	Long: `🛠️ build package
	
This command will build package in current directory and store the resulting
package and signature in /var/cache/pacman/pkg.`,
	Run: Build,
}

func Build(cmd *cobra.Command, args []string) {
	err := pacman.Makepkg()
	CheckErr(err)
	
}
