// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"fmnx.su/core/pack/pacman"
	"github.com/google/uuid"
	"github.com/radovskyb/watcher"
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
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "🌐 run package registry",
	Long: `🌐 run package registry

This command will expose your pacman cache directory, create database and
provide access to your packages for other users.

Also this command will create endpoint, which allows users to upload signed 
packages to your database. Signatures will be validated with gnupg.`,
	Run: Serve,
}

func Serve(cmd *cobra.Command, args []string) {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(pacmancache))
	mux.Handle("/pack/", http.StripPrefix("/pack/", fs))
	mux.HandleFunc("/pack/push", PushHandler)
	s := http.Server{
		Addr:         ":" + viper.GetString("port"),
		Handler:      mux,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}
	go RunDbDaemon(path.Join(pacmancache, viper.GetString("name")+dbext))
	CheckErr(s.ListenAndServe())
}

// Run package database daemon, that will add new packages to database.
func RunDbDaemon(dbpath string) {
	w := watcher.New()
	err := w.Add(pacmancache)
	CheckErr(err)
	for event := range w.Event {
		file := event.FileInfo.Name()
		if strings.HasSuffix(file, pkgext) {
			err = pacman.RepoAdd(dbpath, path.Join(pacmancache, file))
			if err != nil {
				fmt.Println("error: unable to add package to database")
			}
		}
	}
}

// Handler that can be used to upload user packages.
func PushHandler(w http.ResponseWriter, r *http.Request) {
	file := r.Header.Get("file")
	if !strings.HasSuffix(file, pkgext) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sign := r.Header.Get("sign")
	if file == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpdir := path.Join("/tmp", uuid.New().String())
	err := os.MkdirAll(tmpdir, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpdir)

	f, err := os.Create(path.Join(tmpdir, file))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err = os.Create(path.Join(tmpdir, file+".sig"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.ReadFrom(base64.NewDecoder(
		base64.RawStdEncoding,
		strings.NewReader(sign)),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ValideSignature(tmpdir)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = CacheBuiltPackage(tmpdir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
