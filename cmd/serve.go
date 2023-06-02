// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "port",
		Short:   "p",
		Desc:    "port to run on",
		Default: "4572",
	})
	AddStringFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "cert",
		Short: "c",
		Desc:  "path to certificate file",
	})
	AddStringFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "key",
		Short: "k",
		Desc:  "path to key file",
	})
	AddStringFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "name",
		Short: "n",
		Desc:  "database name, should match the domain",
	})
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "🌐 run package registry",
	Long: `🌐 run package registry

This command will expose your /var/cache/pacman/pkg directory, create database
and provide access to your packages for other users.

`,
	Run: Serve,
}

func Serve(cmd *cobra.Command, args []string) {

}
