// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"fmnx.su/core/pack/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Cmd:     serveCmd,
		Name:    "name",
		Short:   "n",
		Desc:    "database name, should match the domain",
		Default: "localhost:4572",
	})
	AddStringFlag(&FlagParameters{
		Cmd:     serveCmd,
		Name:    "dir",
		Short:   "d",
		Desc:    "exposed directory where packages will be stored",
		Default: pacmancache,
	})
	AddStringFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "key",
		Short: "k",
		Desc:  "key file path for TLS connection",
	})
	AddStringFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "cert",
		Short: "c",
		Desc:  "certificate file path for TLS connection",
	})
	AddStringListFlag(&FlagParameters{
		Cmd:   serveCmd,
		Name:  "mirror",
		Short: "m",
		Desc:  "pull mirror used to load packages every 24 hours",
	})
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "🌐 run package registry",
	Long: `🌐 run package registry

This command will expose your pacman cache directory, create database and
provide access to your packages for other users. Command should be runned in
as root, because it is modifying contents of pacman cache directory.

Also this command will create endpoint, which allows users to upload signed 
packages to your database. Signatures will be validated with gnupg.

By default pack client uses HTTPS, so you should provide valid certificates or
set them up in reverse-proxy (nginx/traefik).`,
	Run: Serve,
}

func Serve(cmd *cobra.Command, args []string) {
	var (
		name = viper.GetString("name")
		port = viper.GetString("port")
		cert = viper.GetString("cert")
		key  = viper.GetString("key")
		dir  = viper.GetString("dir")
		mirr = viper.GetStringSlice("mirror")
	)

	go func() {
		err := server.PkgDirDaemon(server.PkgDirParams{
			DbName:     name,
			WatchDir:   dir,
			MkDirMode:  fs.ModePerm,
			InfoLogger: log.Default(),
			ErrLogger:  log.Default(),
		})
		CheckErr(err)
	}()

	go func() {
		for _, link := range mirr {
			err := server.MirrFsDaemon(server.MirrFsParams{
				Link:        link,
				Dir:         dir,
				Dur:         time.Hour * 24,
				ErrorLogger: log.Default(),
				InfoLogger:  log.Default(),
			})
			CheckErr(err)
		}
	}()

	fmt.Printf("Launching registry %s on port %s...\n", name, port)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(pacmancache))
	pushHandler := server.PushHandler{
		CacheDir:   dir,
		TmpDir:     "/tmp",
		ErrLogger:  log.Default(),
		InfoLogger: log.Default(),
	}

	mux.Handle(fsendpoint, http.StripPrefix(fsendpoint, fs))
	mux.HandleFunc(pushendpoint, pushHandler.Push)

	s := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}

	if cert != "" && key != "" {
		CheckErr(s.ListenAndServeTLS(cert, key))
	}
	CheckErr(s.ListenAndServe())
}
