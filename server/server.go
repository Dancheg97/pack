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
	Mux         *http.ServeMux
	Db          db.Database
	ServeDir    string
	RepoName    string
	Cert        string
	Key         string
	LevelDbPath string
	Autocert    bool
	Users       []string
	Handlers    []Handler
}

// Additional handlers that can be added to server.
type Handler struct {
	http.HandlerFunc
	Path string
}

// This function runs a server on a specified directory. This directory will be
// exposed to public.
func (s *Server) Serve() error {
	err := s.initDir()
	if err != nil {
		return err
	}

	err = s.initDatabase()
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

// If directory is not provided by user, set up current process dir.
func (s *Server) initDir() error {
	if s.ServeDir == `` {
		d, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.MkdirAll("public", 0644)
		if err != nil {
			return err
		}
		s.ServeDir = path.Join(d, "public")
	}
	return nil
}

// Initialize server database.
func (s *Server) initDatabase() error {
	ldb, err := db.GetLevelDB(s.LevelDbPath)
	if err != nil {
		return err
	}
	for _, u := range s.Users {
		splt := strings.Split(u, "|")
		err = ldb.Update(db.User{
			Name:     splt[0],
			Password: splt[1],
		})
		if err != nil {
			return err
		}
	}
	s.Db = ldb
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
	users, err := s.Db.List()
	if err != nil {
		return err
	}
	for _, u := range users {
		err = s.initPkgs(path.Join(s.ServeDir, u.Name), "."+u.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Initialize default handlers for server.
func (s *Server) initDefaultHandlers() error {
	if s.Mux == nil {
		s.Mux = http.DefaultServeMux
	}
	fs := http.FileServer(http.Dir(s.ServeDir))
	s.Mux.Handle("/pacman/", http.StripPrefix("/pacman/", fs))

	s.Handlers = append(s.Handlers, Handler{
		HandlerFunc: s.pushHandler,
		Path:        "/push",
	})

	for _, h := range s.Handlers {
		s.Mux.Handle("/pacman"+h.Path, http.StripPrefix("/pacman"+h.Path, h))
	}
	s.Server.Handler = s.Mux
	return nil
}

// Generate certificates for secure connection with server.
func (s *Server) generateCerts() error {
	fmt.Println(":: Generating certificates...")
	cert := path.Join(s.LevelDbPath, "key.pem")
	key := path.Join(s.LevelDbPath, "cert.pem")
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
