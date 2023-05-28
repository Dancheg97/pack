// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"fmnx.su/core/pack/db"
	"fmnx.su/core/pack/pacman"
)

// Server that will provide access to packages.
// You can add custom endpoints to mux, they will be added to server.
type Server struct {
	http.Server
	Mux      *http.ServeMux
	Db       db.Database
	ServeDir string
	RepoName string
	Cert     string
	Key      string
	DbPath   string
	Autocert bool
	Users    []string
	Handlers []Handler
}

// Additional handlers that can be added to server.
type Handler struct {
	http.HandlerFunc
	Path string
}

// This function runs a server on a specified directory. This directory will be
// exposed to public.
func (s *Server) Serve() error {
	err := s.initDatabase()
	if err != nil {
		return err
	}

	err = s.initDirs()
	if err != nil {
		return err
	}

	err = s.initPkgs(s.ServeDir, "")
	if err != nil {
		return err
	}

	err = s.initDefaultHandlers()
	if err != nil {
		return err
	}

	return s.runServer()
}

// Initialize server database.
func (s *Server) initDatabase() error {
	fdb, err := db.GetFileDb(s.DbPath)
	if err != nil {
		return err
	}
	for _, u := range s.Users {
		splt := strings.Split(u, ":")
		err = fdb.Update(splt[0], splt[1])
		if err != nil {
			return err
		}
	}
	s.Db = fdb
	return nil
}

// If directory is not provided by user, set up current process dir and
// directories for user packages.
func (s *Server) initDirs() error {
	if s.ServeDir == `` {
		d, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.MkdirAll("public", 0755)
		if err != nil {
			return err
		}
		s.ServeDir = path.Join(d, "public")
	}
	users, err := s.Db.List()
	if err != nil {
		return err
	}
	for _, u := range users {
		err = os.MkdirAll(path.Join(s.ServeDir, u), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Initializes packages, will recursively walk throw provided dir and add all
// .pkg.tar.zst packages to in each specified repository. Userprefix is used
// to add use names in nested folders.
func (s *Server) initPkgs(dir string, userprefix string) error {
	rootFileInfo, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range rootFileInfo {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".pkg.tar.zst") {
			err = pacman.RepoAdd(
				path.Join(s.ServeDir, fi.Name()),
				path.Join(s.ServeDir, userprefix+s.RepoName+".db.tar.gz"),
			)
			if err != nil {
				return err
			}
		}
	}
	if userprefix == `` {
		users, err := s.Db.List()
		if err != nil {
			return err
		}
		for _, u := range users {
			err = s.initPkgs(path.Join(s.ServeDir, u), "."+u)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Generate certificates for secure connection with server.
func (s *Server) generateCerts() error {
	fmt.Println(":: Generating certificates...")
	cert := path.Join(s.DbPath, "key.pem")
	key := path.Join(s.DbPath, "cert.pem")
	cmd := exec.Command(
		"openssl", "req", "-x509", "-newkey", "rsa:4096",
		"-keyout", key, "-out", cert,
		"-sha256", "-days", "3650", "-nodes", "-subj",
		"/C=XX/ST=_/L=_/O=_/OU=_/CN=_",
	)
	s.Key = key
	s.Cert = cert
	return cmd.Run()
}

// Initialize server for packages.
func (s *Server) runServer() error {
	if s.Autocert {
		err := s.generateCerts()
		if err != nil {
			return err
		}
	}
	fmt.Print(":: Listening on " + s.Addr + "...")
	if s.Cert != `` && s.Key != `` {
		return s.Server.ListenAndServeTLS(s.Cert, s.Key)
	}
	return s.Server.ListenAndServe()
}
