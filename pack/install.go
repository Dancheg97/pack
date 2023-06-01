// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"fmt"
	"strings"

	"fmnx.su/core/pack/pacman"
)

// Parameters for new packages that will be installed.
type InstallParameters struct {
	// List of packages for installation.
	Packages []string
	// Set optional+trust all flag to new database.
	TrustAll bool
	// Use http connection instead of https.
	HTTP bool
}

// This function can be used to install additional packages into the system.
func Install(p *InstallParameters) error {
	confDbs, err := pacman.GetConfigDatabases()
	if err != nil {
		return err
	}
	pkgDbs, err := PackagesDatabases(p.Packages)
	if err != nil {
		return err
	}
	missing := CheckDatabases(pkgDbs, confDbs)
	for _, missdb := range missing {
		err = pacman.AddConfigDatabase(&pacman.RepositoryParameters{
			Database: missdb,
			TrustAll: p.TrustAll,
			HTTP:     p.HTTP,
		})
		if err != nil {
			return err
		}
	}
	return pacman.SyncList(FormatPkgs(p.Packages))
}

// Get list of pacman databases, that have to be inluded in configuration to
// install packages.
func PackagesDatabases(pkgs []string) ([]string, error) {
	var dbs []string
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 1:
			continue
		case 2:
			dbs = append(dbs, splt[0])
		default:
			return nil, fmt.Errorf("bad package format: %s", pkg)
		}
	}
	return dbs, nil
}

// Get databases that are not currently on a list. This function will return
// databases that are missing from current list.
func CheckDatabases(check []string, current []string) []string {
	var missing []string
	for _, chk := range check {
		found := false
		for _, curr := range current {
			if chk == curr {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, chk)
		}
	}
	return missing
}

// Format packages to appropriate format before installation.
func FormatPkgs(pkgs []string) []string {
	var rez []string
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		if len(splt) == 3 {
			rez = append(rez, splt[0]+"."+splt[1]+"/"+splt[2])
			continue
		}
		rez = append(rez, pkg)
	}
	return rez
}
