// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/pack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AddBoolFlag(&FlagParameters{
		Cmd:  installCmd,
		Name: "install-http",
		Desc: "📍 use http instead of https",
	})
	AddBoolFlag(&FlagParameters{
		Cmd:  installCmd,
		Name: "install-trust-all",
		Desc: "📌 set optioanl and trust-all mode to new database",
	})
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
	err = pack.Install(&pack.InstallParameters{
		Packages: pkgs,
		TrustAll: viper.GetBool("install-trust-all"),
		HTTP:     viper.GetBool("install-http"),
	})
	CheckErr(err)
}
