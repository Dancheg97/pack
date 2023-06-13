// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package registry

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (p *Registry) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trimmedfile := strings.TrimPrefix(vars["file"], vars["owner"]+".")
	data, err := p.FileStorage.Get(Join(vars["owner"], trimmedfile))
	if err != nil {
		p.end(w, http.StatusNotFound, err)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
