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
		Cmd:   pushCmd,
		Name:  "http",
		Short: "p",
		Desc:  "📍 use http instead of https",
	})
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "⬆ push package",
	Long: `⬆ push packages

Provide repository name and packages you want to push.

This command will search for built package in cache, then it will try to upload
this package to compatible pack repo.

You can provide multiple packge names, all of them will be installed.`,
	Run: Push,
}

func Push(cmd *cobra.Command, args []string) {
	err := pack.Push(&pack.PushParameters{
		Registry: args[0],
		Packages: args[1:],
		HTTP:     viper.GetBool("http"),
	})
	CheckErr(err)
}
