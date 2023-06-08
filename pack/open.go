// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"fmnx.su/core/pack/server"
	"fmnx.su/core/pack/tmpl"
)

// Parameters to run pack registry.
type OpenParameters struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Endpoint that will be mounted for provided directory. This endpoint
	// can be used as connection in pacman.conf file. Also creates push subpath
	// /push, which can accept packages from end user
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
	p := formOptions(prms, opendefault)

	d := server.LocalDirDb{
		Dir:    p.Dir,
		DbName: p.Name,
	}

	k := server.LocalKeyring{
		File: p.Ring,
	}

	s := server.Pusher{
		Stdout:          p.Stdout,
		Stderr:          p.Stderr,
		GPGVireivicator: &k,
		DbFormer:        &d,
	}

	fs := http.FileServer(http.Dir(p.Dir))
	http.Handle(p.Endpoint, http.StripPrefix(p.Endpoint, fs))
	http.HandleFunc(p.Endpoint+"/push", s.Push)

	startmes := fmt.Sprintf("%s %s%s\n", tmpl.Dots, tmpl.Launching, p.Name)
	p.Stdout.Write([]byte(startmes))

	if p.Cert != "" && p.Key != "" {
		return http.ListenAndServeTLS(
			":"+p.Port, p.Cert, p.Key,
			http.DefaultServeMux,
		)
	}
	return http.ListenAndServe(":"+p.Port, http.DefaultServeMux)
}
