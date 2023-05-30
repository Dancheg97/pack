// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"net/http"
	"os"
	"time"

	"fmnx.su/core/pack/db"
	"fmnx.su/core/pack/pack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-addr",
		Desc:    "🌐 server adress",
		Default: "localhost:8080",
		Env:     "PACK_SERVE_ADDR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "serve-repo",
		Desc:    "📋 name of repository, should match the domain",
		Default: "localhost:8080",
		Env:     "PACK_SERVE_REPO",
	})
	AddStringListFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "serve-users",
		Short: "u",
		Desc:  "😀 initial users, format - user::password",
		Env:   "PACK_SERVE_USERS",
	})
	AddStringListFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "serve-pull-mirr",
		Short: "m",
		Desc:  "🪞 list of pull mirrors, just provide links (depends on wget)",
		Env:   "PACK_SERVE_PULL_MIRR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-log-file",
		Desc: "📝 file for server logs",
		Env:  "PACK_SERVE_LOG_FILE",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-work-dir",
		Desc: "🗃️ directory with private files required for server",
		Env:  "PACK_SERVE_WORK_DIR",
	})
	AddStringFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-public-dir",
		Desc: "📂 directory with packages, publicly exposed",
		Env:  "PACK_SERVE_PUBLIC_DIR",
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
		Cmd:     serveCmd,
		Name:    "serve-db-file",
		Desc:    "💾 path to local file with user info",
		Default: "users",
		Env:     "PACK_SERVE_DB_FILE",
	})
	AddBoolFlag(&FlagParameters{
		Cmd:  serveCmd,
		Name: "serve-autocert",
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
	db, err := db.GetFileDb(viper.GetString("serve-db-file"))
	CheckErr(err)
	err = db.Fill(viper.GetStringSlice("serve-users"))
	CheckErr(err)
	if viper.GetString("serve-work-dir") == "" {
		d, err := os.Getwd()
		CheckErr(err)
		viper.Set("serve-work-dir", d)
	}
	server := pack.Server{
		Server: http.Server{
			Addr:         viper.GetString("serve-addr"),
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		},
		WorkDir:  viper.GetString("serve-work-dir"),
		ServeDir: viper.GetString("serve-public-dir"),
		RepoName: viper.GetString("serve-repo"),
		Cert:     viper.GetString("serve-cert"),
		Key:      viper.GetString("serve-key"),
		Autocert: viper.GetBool("serve-autocert"),
		PullMirr: viper.GetStringSlice("serve-pull-mirr"),
		Db:       db,
	}
	err = server.Serve()
	CheckErr(err)
}
