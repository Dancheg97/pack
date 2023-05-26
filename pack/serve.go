// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pack

import (
	"log"
	"net/http"

	"fmnx.su/core/pack/pacman"
)

// Parameters that used in pack server.
type ServeParameters struct {
	Dir  string
	Port string
	Repo string
}

// This command starts a server in specified directory. This command will
// automatically create pack instance with 2 main areas:
// 1 - user space zone (for packages that are separated in main database, but
// they can be checked from endpoints)
// 2 - root zone for packages in main database
func Serve(p ServeParameters) error {
	opts := pacman.RepoAddDefault
	opts.Dir = p.Dir
	err := pacman.RepoAdd(p.Dir+"/*.pkg.tar.zst", p.Repo+".db.tar.gz")
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(p.Dir))
	http.Handle("/pacman/", fs)

	log.Print(":: Listening on " + p.Port + "...")
	return http.ListenAndServe(":"+p.Port, nil)
}
