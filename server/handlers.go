// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"fmt"
	"net/http"
	"os"
)

// Initialize default handlers for server.
func (s *Server) initDefaultHandlers() error {
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

// Handler that can be used to upload packages.
func (s *Server) push(w http.ResponseWriter, r *http.Request) {
	u := r.Header.Get("login")
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
}
