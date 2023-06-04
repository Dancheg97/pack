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
	"os/exec"
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
		mirr = viper.GetStringSlice("mirr")
	)

	go PkgDirDaemon(name)
	go FsMirrDaemon(mirr)

	fmt.Printf("Launching database %s on port %s...\n", name, port)

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

	if cert != "" && key != "" {
		CheckErr(s.ListenAndServeTLS(cert, key))
	}
	CheckErr(s.ListenAndServe())
}

// This function is launching watcher for pacman cache directory, and constatly
// adding new packages to database.
func PkgDirDaemon(name string) {
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move)
	CheckErr(w.Add(pacmancache))
	go w.Start(time.Second) //nolint:errcheck

	for event := range w.Event {
		file := event.FileInfo.Name()
		if strings.HasSuffix(file, pkgext) {
			err := pacman.RepoAdd(
				path.Join(pacmancache, name+dbext),
				path.Join(pacmancache, file),
			)
			if err != nil {
				fmt.Println("error: unable to add package to database")
			}
		}
	}
}

// This function start mirror watcher, which loads packages from remote file
// server to pacman cache directory every 24 hours.
func FsMirrDaemon(links []string) {
	for {
		for _, link := range links {
			err := exec.Command( //nolint:gosec
				"sudo", "wget", "-nd", "-np", "-P",
				pacmancache, "--recursive", link,
			).Run()
			if err != nil {
				fmt.Println("[MIRR] - Failed to load: " + link)
			}
		}
		time.Sleep(time.Hour * 24)
	}
}

// Handler that can be used to upload user packages.
func PushHandler(w http.ResponseWriter, r *http.Request) {
	file := r.Header.Get(file)
	if !strings.HasSuffix(file, pkgext) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(path.Join(pacmancache, file)); err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	sign := r.Header.Get(sign)
	if sign == "" {
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
