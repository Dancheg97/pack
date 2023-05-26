// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
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
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   tmpl.ServeShort,
	Long:    tmpl.ServeLong,
	Run:     Serve,
}

// Cli command installing packages into system.
func Serve(cmd *cobra.Command, pkgs []string) {
	err := pack.Serve(pack.ServeParameters{
		Dir:  "",
		Port: "8080",
		Repo: "localhost:8080",
	})
	CheckErr(err)
}
