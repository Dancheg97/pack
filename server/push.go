// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import "net/http"

func init() {
	Handlers = append(Handlers, Handler{
		HandlerFunc: pushHandler,
		Path:        "push",
	})
}

func pushHandler(w http.ResponseWriter, r *http.Request) {

}
