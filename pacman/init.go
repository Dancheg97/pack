// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// Dependecy packages.
const (
	pacman  = `pacman`
	sudo    = `sudo`
	makepkg = `makepkg`
	repoadd = `repo-add`
)

// Global lock for operations with pacman database.
var mu sync.Mutex

func init() {
	checkDependency(pacman)
	checkDependency(sudo)
	checkDependency(makepkg)
	checkDependency(repoadd)
}

func checkDependency(p string) {
	_, err := exec.LookPath(p)
	if err != nil {
		fmt.Printf("unable to find %s in system\n", p)
		os.Exit(1)
	}
}

func SudoCommand(sudo bool, cmd string, args ...string) *exec.Cmd {
	if sudo {
		args := append([]string{cmd}, args...)
		return exec.Command("sudo", args...)
	}
	return exec.Command(cmd, args...)
}

func formOptions[Options any](arr []Options, dv *Options) *Options {
	if len(arr) != 1 {
		return dv
	}
	return &arr[0]
}
