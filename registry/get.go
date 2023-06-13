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
	file := strings.Join([]string{vars["owner"], vars["file"]}, ".")
	data, err := p.FileStorage.Get(file)
	if err != nil {
		p.end(w, http.StatusNotFound, err)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
