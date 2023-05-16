// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pack

// This library contains operations related only to pack json database.
// All functions are safe for concurrent usage.

import (
	"os"
	"strings"

	"fmnx.su/core/pack/config"
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

// Pack metadata parsed from pack link.
type PackMd struct {
	PackName   string
	PacmanName string
	Version    string
	GitUrl     string
	Directory  string
}

// Form pack installation metadata from pack link.
func GetPackInfo(link string) PackMd {
	rez := PackMd{}
	versplt := strings.Split(link, "@")
	rez.PackName = versplt[0]
	rez.GitUrl = "https://" + versplt[0]
	if len(versplt) > 1 {
		rez.Version = versplt[1]
	}
	dashsplt := strings.Split(rez.PackName, "/")
	rez.PacmanName = dashsplt[len(dashsplt)-1]
	rez.Directory = config.RepoCacheDir + "/" + rez.PacmanName
	return rez
}
