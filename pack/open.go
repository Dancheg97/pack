// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"io"
	"os"
)

// Parameters to run pack registry.
type OpenParameters struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Endpoint that will be mounted for provided directory (default /api/pack).
	Endpoint string
	// Directory with packages that will be opened. Make sure you have access
	// to read and write (default /var/cache/pacman/pkg).
	Dir string
	// Name of domain registry, database name will be automatically assigned.
	Name string
	// Port to run registry on.
	Port string
	// Path to certificate file for TLS.
	Cert string
	// Path to key file for TLS.
	Key string
	// Path to custom keyring, by default pacman keyring. (pacman-key -e)
	Ring string
}

func opendefault() *OpenParameters {
	return &OpenParameters{
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Stdin:    os.Stdin,
		Endpoint: "/api/pack",
		Dir:      "/var/cache/pacman/pkg",
		Ring:     "/usr/share/pacman/keyrings/archlinux.gpg",
		Name:     "localhost",
		Port:     "8080",
	}
}

func Open(prms ...OpenParameters) error {
	_ = formOptions(prms, opendefault)
	return nil
}

// func a() {
// 	var (
// 		name = viper.GetString("name")
// 		port = viper.GetString("port")
// 		cert = viper.GetString("cert")
// 		key  = viper.GetString("key")
// 		dir  = viper.GetString("dir")
// 		mirr = viper.GetStringSlice("mirror")
// 	)

// 	go func() {
// 		d := server.PkgDirDaemon{
// 			DbName:   name,
// 			WatchDir: dir,
// 		}
// 		CheckErr(d.Run())
// 	}()

// 	go func() {
// 		for _, link := range mirr {
// 			d := server.FsMirrorDaemon{
// 				Link: link,
// 				Dir:  dir,
// 				Dur:  time.Hour * 24,
// 			}
// 			CheckErr(d.Run())
// 		}
// 	}()

// 	fmt.Printf("Launching registry %s on port %s...\n", name, port)

// 	mux := http.NewServeMux()
// 	fs := http.FileServer(http.Dir(pacmancache))
// 	pushHandler := server.PushHandler{
// 		CacheDir:   dir,
// 		ErrLogger:  log.Default(),
// 		InfoLogger: log.Default(),
// 	}

// 	mux.Handle(fsendpoint, http.StripPrefix(fsendpoint, fs))
// 	mux.HandleFunc(pushendpoint, pushHandler.Push)

// 	s := http.Server{
// 		Addr:         ":" + port,
// 		Handler:      mux,
// 		ReadTimeout:  time.Minute * 15,
// 		WriteTimeout: time.Minute * 15,
// 	}

// 	if cert != "" && key != "" {
// 		CheckErr(s.ListenAndServeTLS(cert, key))
// 	}
// 	CheckErr(s.ListenAndServe())
// }
