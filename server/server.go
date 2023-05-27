// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"fmnx.su/core/pack/pacman"
	"github.com/syndtr/goleveldb/leveldb"
)

// Server that will provide access to packages.
// You can add custom endpoints to mux, they will be added to server.
type Server struct {
	http.Server
	Mux      *http.ServeMux
	db       *leveldb.DB
	ServeDir string
	RepoName string
	Cert     string
	Key      string
}

// This function runs a server on a specified directory. This directory will be
// exposed to public.
func (s *Server) Serve() error {
	err := s.initDatabase()
	if err != nil {
		return err
	}

	err = s.initPkgs()
	if err != nil {
		return err
	}

	return s.runServer()
}

// Initialize server database.
func (s *Server) initDatabase() error {
	db, err := leveldb.OpenFile(s.ServeDir+"/users.db", nil)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

// Initialize packages.
func (s *Server) initPkgs() error {
	rootFileInfo, err := os.ReadDir(s.ServeDir)
	if err != nil {
		return err
	}
	for _, ufi := range rootFileInfo {
		if ufi.IsDir() {
			userDir := s.ServeDir + "/" + ufi.Name()
			userFileInfo, err := os.ReadDir(userDir)
			if err != nil {
				return err
			}
			for _, userFile := range userFileInfo {
				if strings.HasSuffix(userFile.Name(), ".pkg.tar.zst") {
					err = pacman.RepoAdd(
						userDir+"/"+userFile.Name(),
						userDir+"/"+s.RepoName+"."+userFile.Name()+".db.tar.gz",
					)
					if err != nil {
						return err
					}
				}
			}
			continue
		}
		if strings.HasSuffix(ufi.Name(), ".pkg.tar.zst") {
			err = pacman.RepoAdd(
				s.ServeDir+"/"+ufi.Name(),
				s.ServeDir+"/"+s.RepoName+".db.tar.gz",
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Initialize server for packages.
func (s *Server) runServer() error {
	if s.Mux == nil {
		s.Mux = http.DefaultServeMux
	}

	fs := http.FileServer(http.Dir(s.ServeDir))
	s.Mux.Handle("/pacman/", http.StripPrefix("/pacman/", fs))

	s.Server.Handler = s.Mux

	log.Print(":: Listening on " + s.Addr + "...")
	if s.Cert != `` && s.Key != `` {
		return s.Server.ListenAndServeTLS(s.Cert, s.Key)
	}
	return s.Server.ListenAndServe()
}
