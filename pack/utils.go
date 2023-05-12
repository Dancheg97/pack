// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pack

// This library contains operations related only to pack json database.
// All functions are safe for concurrent usage.

import (
	"os"
	"strings"
)

// Swap dependencies in PKGBUILD file for proper installation with pacman.
func SwapDeps(pkgbuild string, deps []string) error {
	b, err := os.ReadFile(pkgbuild)
	if err != nil {
		return err
	}
	var rez = string(b)
	for _, dep := range deps {
		dashsplt := strings.Split(dep, "/")
		shortname := dashsplt[len(dashsplt)-1]
		rez = strings.ReplaceAll(rez, dep, shortname)
	}
	return os.WriteFile(pkgbuild, []byte(rez), 0o600)
}
