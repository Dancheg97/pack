// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"os"

	"fmnx.su/core/pack/pacman"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"u"},
	Short:   "💡 upgrade system",
	Long: `💡 upgrade system
	
This command will update all packages in they system, installing latest 
versions from connected registries.`,
	Run: Upgrade,
}

func Upgrade(cmd *cobra.Command, args []string) {
	err := pacman.SyncList(nil, pacman.SyncOptions{
		Sudo:    true,
		Needed:  true,
		Refresh: true,
		Upgrade: true,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	})
	CheckErr(err)
}
