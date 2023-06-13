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

	"fmnx.su/core/pack/registry"
	"fmnx.su/core/pack/registry/local"
	"fmnx.su/core/pack/tmpl"
	"github.com/gorilla/mux"
)

// Parameters to run pack registry.
type OpenParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

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
	// Path to custom directory, which contains files with public GnuPG keys,
	// which will further be used to validate pushed packages.
	GpgDir string
}

func opendefault() *OpenParameters {
	return &OpenParameters{
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Stdin:    os.Stdin,
		Endpoint: "/api/packages/arch",
		Dir:      "/var/cache/pacman/pkg",
		Name:     "localhost",
		Port:     "8080",
	}
}

func Open(prms ...OpenParameters) error {
	p := formOptions(prms, opendefault)

	_, err := os.Stat(p.Dir)
	if err != nil {
		return fmt.Errorf("expose dir not found: %+v", err)
	}

	d := local.DirStorage{
		Dir: p.Dir,
	}

	k := local.LocalKeyDir{
		Dir: p.GpgDir,
	}

	r := registry.Registry{
		Stdout:      p.Stdout,
		Stderr:      p.Stderr,
		TmpDir:      "/tmp",
		Dbname:      p.Name,
		FileStorage: &d,
		KeyReader:   &k,
	}

	router := mux.NewRouter()

	router.HandleFunc(p.Endpoint+"/push", r.Push)
	router.HandleFunc(p.Endpoint+"/{owner}/{file}", r.Get)
	router.HandleFunc(p.Endpoint+"/{file}", r.Get)

	msg := fmt.Sprintf("Starting registry %s on port %s", p.Name, p.Port)
	tmpl.Amsg(p.Stdout, msg)

	if p.Key != "" && p.Cert != "" {
		return http.ListenAndServeTLS(":"+p.Port, p.Cert, p.Key, router)
	}
	return http.ListenAndServe(":"+p.Port, router)
}
