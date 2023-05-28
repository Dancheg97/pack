// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"net/http"
	"time"

	"fmnx.su/core/pack/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-addr",
		Desc:    "🌐 server adress",
		Default: ":8080",
		Env:     "PACK_SERVE_ADDR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-repo",
		Desc:    "📋 name of repository, should match the domain",
		Default: "localhost:8080",
		Env:     "PACK_SERVE_PORT",
	})
	AddStringListFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "serve-users",
		Short: "u",
		Desc:  "😀 initial users, format: name|pwd",
		Env:   "PACK_SERVE_USERS",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-dir",
		Desc: "📂 directory with packages, publically exposed",
		Env:  "PACK_SERVE_DIR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-cert",
		Desc: "📃 certificate file for TLS server",
		Env:  "PACK_SERVE_CERT",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-key",
		Desc: "🔑 key file for TLS server",
		Env:  "PACK_SERVE_KEY",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-db-path",
		Desc: "💾 path to local leveldb with user info",
		Env:  "PACK_SERVE_DB_PATH",
	})
	AddBoolFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-auto-tls",
		Desc: "🔒 automatically generate certs in db dir (depends on openssl)",
		Env:  "PACK_SERVE_AUTO_CERT",
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
		Server: http.Server{
			Addr:         viper.GetString("serve-addr"),
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		},
		ServeDir: viper.GetString("serve-dir"),
		RepoName: viper.GetString("serve-repo"),
		Cert:     viper.GetString("serve-cert"),
		Key:      viper.GetString("serve-key"),
		DbPath:   viper.GetString("serve-db-path"),
		Autocert: viper.GetBool("serve-auto-tls"),
		Users:    viper.GetStringSlice("serve-users"),
	}
	CheckErr(s.Serve())
}
