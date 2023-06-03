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
	var (
		name = viper.GetString("name")
		port = viper.GetString("port")
	)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(pacmancache))
	mux.Handle(fsendpoint, http.StripPrefix(fsendpoint, fs))
	mux.HandleFunc(pushendpoint, PushHandler)
	s := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}
	go RunDbDaemon(path.Join(pacmancache, name+dbext))
	fmt.Printf("Launching database %s on port %s...\n", name, port)
	CheckErr(s.ListenAndServe())
}

// Run package database daemon, that will add new packages to database.
func RunDbDaemon(dbname string) {
	w := watcher.New()
	err := w.Add(pacmancache)
	CheckErr(err)
	for event := range w.Event {
		file := event.FileInfo.Name()
		if strings.HasSuffix(file, pkgext) {
			err = pacman.RepoAdd(dbname, path.Join(pacmancache, file))
			if err != nil {
				fmt.Println("error: unable to add package to database")
			}
		}
	}
}

// Handler that can be used to upload user packages.
func PushHandler(w http.ResponseWriter, r *http.Request) {
	file := r.Header.Get(file)
	if !strings.HasSuffix(file, pkgext) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sign := r.Header.Get(sign)
	if file == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpdir := path.Join("/tmp", "pack-"+uuid.New().String())
	err := os.MkdirAll(tmpdir, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpdir)

	filepath := path.Join(tmpdir, file)
	if _, err := os.Stat(filepath); err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	f, err := os.Create(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(path.Join(tmpdir, file+".sig"), sigdata, 0600)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ValideSignature(tmpdir)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = CacheBuiltPackage(tmpdir + "/")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[PUSH] - package accepted: " + file)
	w.WriteHeader(http.StatusOK)
}
