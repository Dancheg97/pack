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
type Server struct {
	Dir  string
	Port string
	Repo string
	db   *leveldb.DB
}

// This function runs a server on a specified directory. This directory will be
// exposed to public.
func (s *Server) Serve() error {
	opts := pacman.RepoAddDefault
	opts.Dir = s.Dir

	db, err := leveldb.OpenFile(s.Dir+"/users.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()
	s.db = db

	err = s.initPkgs()
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(s.Dir))
	http.Handle("/", fs)

	log.Print(":: Listening on " + s.Port + "...")
	return http.ListenAndServe(":"+s.Port, nil)
}

// Initialize packages.
func (s *Server) initPkgs() error {
	rootFileInfo, err := os.ReadDir(s.Dir)
	if err != nil {
		return err
	}
	for _, ufi := range rootFileInfo {
		if ufi.IsDir() {
			userDir := s.Dir + "/" + ufi.Name()
			userFileInfo, err := os.ReadDir(userDir)
			if err != nil {
				return err
			}
			for _, userFile := range userFileInfo {
				if strings.HasSuffix(userFile.Name(), ".pkg.tar.zst") {
					err = pacman.RepoAdd(
						userDir+"/"+userFile.Name(),
						userDir+"/"+s.Repo+"."+userFile.Name()+".db.tar.gz",
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
				s.Dir+"/"+ufi.Name(),
				s.Dir+"/"+s.Repo+".db.tar.gz",
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
