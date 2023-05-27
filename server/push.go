// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"fmt"
	"net/http"
)

func (s *Server) pushHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("req recieved")
	u := r.Header.Get("user")
	p := r.Header.Get("password")
	fmt.Println(u, p)
}
