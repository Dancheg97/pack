// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   tmpl.InstallShort,
	Run:     Install,
}

// Cli command installing packages into system.
func Install(cmd *cobra.Command, pkgs []string) {
	err := lock.TryLock()
	CheckErr(err)
	err = pack.Install(pkgs)
	CheckErr(err)
}
