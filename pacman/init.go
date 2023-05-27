// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

// Dependecy packages.
const (
	pacman  = `pacman`
	makepkg = `makepkg`
	repoadd = `repo-add`
)

// Global lock for operations with pacman database.
var mu sync.Mutex

func init() {
	checkDependency(pacman)
	checkDependency(makepkg)
	checkDependency(repoadd)
}

func checkDependency(p string) {
	_, err := exec.LookPath(p)
	if err != nil {
		log.Printf("unable to find %s in system\n", p)
		os.Exit(1)
	}
}

func formOptions[Options any](arr []Options, dv *Options) *Options {
	if len(arr) != 1 {
		return dv
	}
	return &arr[0]
}
