// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmnx.su/core/pack/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-dir",
		Desc:    "📂 directory with packages",
		Default: "/var/cache/pacman/pkg",
		Env:     "PACK_SERVE_DIR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-port",
		Desc:    "🌐 exposed port, on which server will run",
		Default: "1997",
		Env:     "PACK_EXPOSED_PORT",
	})
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-repo",
		Desc:    "📋 name of repository, should match the domain",
		Default: "localhost:1997",
		Env:     "PACK_EXPOSED_PORT",
	})
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   `🌐 run pack server`,
	Run:     Serve,
}

// Cli command installing packages into system.
func Serve(cmd *cobra.Command, pkgs []string) {
	s := server.Server{
		Dir:  viper.GetString("serve-dir"),
		Port: viper.GetString("serve-port"),
		Repo: viper.GetString("serve-repo"),
	}
	CheckErr(s.Serve())
}
