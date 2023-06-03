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
	AddBoolFlag(&FlagParameters{
		Cmd:   installCmd,
		Name:  "keep",
		Short: "k",
		Desc:  "do not remove database after package installation",
	})
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "🪛 install packages",
	Long: `🪛 install packages

This command is split into 2 pa...`,
	Run: Install,
}

func Install(cmd *cobra.Command, args []string) {
	err := pacman.SyncList(args)
	CheckErr(err)
}
