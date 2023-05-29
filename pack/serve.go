// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"fmnx.su/core/pack/pacman"
)

// Server that will provide access to packages.
// You can add custom endpoints to mux, they will be added to server.
type Server struct {
	http.Server
	Mux *http.ServeMux

	// Simple database to store user credentials.
	Db Database

	// Publicly exposed directory with packages and pacman repository files.
	ServeDir string

	// Repository name: should match the domain name.
	RepoName string

	// TLS certificate file path, don't store it in serve-dir.
	Cert string
	// TLS key file path, don't store it in serve-dir.
	Key string

	// If true, certificate will be generated automatically by openssl.
	Autocert bool

	// Additional handlers that will be registered under /pacman/ path.
	Handlers []Handler

	// Mirror links which will be pulled by 24h timeout with wget.
	PullMirr []string

	// Links to git repositories, which should contain valid PKGBUILD files,
	// (for example AUR links), they will be built and provided as packages.
	BuildGit []string
}

// Additional handlers that can be added to server.
type Handler struct {
	http.HandlerFunc
	Path string
}

// By default, just a simple file which stores information about users, but you
// can use any database for your specific use case to integrate this system in
// your personal workflow.
type Database interface {
	List() ([]string, error)
	Validate(name string, password string) bool
	Update(name string, password string) error
}

// This function runs a server on a specified directory. This directory will be
// exposed to public.
func (s *Server) Serve() error {
	var startFuncs []func() error = []func() error{
		s.prepareDirectories,
		s.pullMirrors,
		s.prepareRepositories,
		s.prepareCertificates,
		s.initRoutes,
		s.runServer,
	}
	for _, startFunc := range startFuncs {
		err := startFunc()
		if err != nil {
			return err
		}
	}
	return nil
}

// Used to create all directories, that are required for operating server.
// Fucntion will ensure, that root directory and nested user directories exist,
// otherwise it will create them.
func (s *Server) prepareDirectories() error {
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

// Pull all packages and signatures from specified mirrors to root folder,
// packages then will be added to root repository.
func (s *Server) pullMirrors() error {
	for _, mirr := range s.PullMirr {
		go LaunchMirrorDaemon(&MirrorParams{
			Dir:    s.ServeDir,
			Link:   mirr,
			Dbname: s.Addr,
			Dur:    time.Hour * 36,
		})
	}
	return nil
}

type MirrorParams struct {
	Dir    string
	Link   string
	Dbname string
	Dur    time.Duration
}

// Start mirroring specified repository.
func LaunchMirrorDaemon(p *MirrorParams) {
	err := os.MkdirAll(p.Dir, 0755)
	if err != nil {
		fmt.Println("[MIRROR] Failed to create mirr dir: ", p.Dir)
		return
	}
	for {
		err := exec.Command( //nolint:gosec
			"wget", "-nd", "-np", "-P",
			p.Dir, "--recursive", p.Link,
		).Run()
		if err != nil {
			fmt.Println("[MIRROR] Error in mirror occured: ", p.Link)
		}

		err = prepareDirRepo(p.Dir, p.Dbname)
		if err != nil {
			fmt.Println("[MIRROR] Failed to create mirror DB", p.Dbname)
		}

		time.Sleep(p.Dur)
	}
}

// Function is used to initialize database and all nested user databases with
// pacman packages.
func (s *Server) prepareRepositories() error {
	err := prepareDirRepo(s.ServeDir, s.Addr)
	if err != nil {
		return err
	}
	users, err := s.Db.List()
	if err != nil {
		return err
	}
	for _, user := range users {
		err = prepareDirRepo(path.Join(s.ServeDir, user), s.Addr+"."+user)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function that is used to create database for all pacman packages in
// directory.
func prepareDirRepo(dir string, db string) error {
	rootFileInfo, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range rootFileInfo {
		filename := fi.Name()
		if strings.HasSuffix(filename, ".pkg.tar.zst") {
			err = pacman.RepoAdd(
				path.Join(dir, filename),
				path.Join(dir, db+".db.tar.gz"),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Initialize server secutrity, if autosert is provided - generating
// certs.
func (s *Server) prepareCertificates() error {
	if !s.Autocert {
		return nil
	}
	fmt.Println(":: Generating certificates...")
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	s.Key = path.Join(d, "key.pem")
	s.Cert = path.Join(d, "cert.pem")
	return exec.Command( //nolint:gosec
		"openssl", "req", "-x509", "-newkey", "rsa:4096",
		"-keyout", s.Key, "-out", s.Cert,
		"-sha256", "-days", "3650", "-nodes", "-subj",
		"/C=XX/ST=_/L=_/O=_/OU=_/CN=_",
	).Run()
}

// Initialize default handlers for server.
func (s *Server) initRoutes() error {
	if s.Mux == nil {
		s.Mux = http.DefaultServeMux
	}
	fs := http.FileServer(http.Dir(s.ServeDir))
	s.Mux.Handle("/pacman/", http.StripPrefix("/pacman/", fs))

	s.Handlers = append(s.Handlers, Handler{
		HandlerFunc: s.push,
		Path:        "/push",
	})

	for _, h := range s.Handlers {
		s.Mux.Handle("/pacman"+h.Path, http.StripPrefix("/pacman"+h.Path, h))
	}
	s.Server.Handler = s.Mux
	return nil
}

// Handler that can be used to upload user packages.
func (s *Server) push(w http.ResponseWriter, r *http.Request) {
	u := r.Header.Get("user")
	p := r.Header.Get("password")
	if !s.Db.Validate(u, p) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file := r.Header.Get("file")
	f, err := os.Create(fmt.Sprintf("%s/%s/%s", s.ServeDir, u, file))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = prepareDirRepo(path.Join(s.ServeDir, u), s.Addr+"."+u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Initialize server for packages.
func (s *Server) runServer() error {
	fmt.Print(":: Listening on " + s.Addr + "...")
	if s.Cert != `` && s.Key != `` {
		return s.Server.ListenAndServeTLS(s.Cert, s.Key)
	}
	return s.Server.ListenAndServe()
}
