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
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "⬇ install packages",
	Run:     Install,
}

// Cli command installing packages into system.
func Install(cmd *cobra.Command, pkgs []string) {
	err := lock.TryLock()
	CheckErr(err)
	err = pack.Install(pkgs)
	CheckErr(err)
}
