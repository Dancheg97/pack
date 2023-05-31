// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/pack"
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

This command will build package in current directory and move it to pacman cache 
directory (/var/cache/pacman/pkg), after build you will be able to upload this 
package with 'push' command to pack registry.`,
	Run: Build,
}

func Build(cmd *cobra.Command, args []string) {
	err := lock.TryLock()
	CheckErr(err)
	err = pack.Build()
	CheckErr(err)
}
