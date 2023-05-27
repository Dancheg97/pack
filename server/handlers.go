// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"context"
	"net/http"
)

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

// Functions with additional adjustments that are easier to write.
type RPCFunc func(RPCtx, any) (any, error)

// Extended context to
type RPCtx struct {
	context.Context
}

// Wrapper function to transform rpc-like function to http.Handler format.
func RpcTransform(f RPCFunc) http.HandlerFunc {
	return nil
}
