// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pacman

// This package acts as library wrapper over pacman and makepkg.

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
)

// ValidateDir directory to exist and contain PKGBUILD.
func ValidateDir(dir string) error {
	_, err := os.Stat(dir + "/PKGBUILD")
	return err
}

// Geberate PKGBUILD file template. Provide pacman package name and git url.
func Generate(dir string, url string) error {
	splt := strings.Split(url, "/")
	content := fmt.Sprintf(tmpl.PKGBUILD, splt[len(splt)-1], url)
	err := os.WriteFile(dir+"/PKGBUILD", []byte(content), 0o600)
	if err != nil {
		return errors.New("PKGBUILD generation failed:\n" + err.Error())
	}
	return nil
}

// Some pacman package description fields.
type PkgInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Url         string `json:"url"`
}

// Get package description from pacman and parse it.
func Describe(pkg string) (PkgInfo, error) {
	o, err := system.Callf("pacman -Qi %s", pkg)
	if err != nil {
		return PkgInfo{}, errors.New("package not found: " + pkg)
	}
	return PkgInfo{
		Name:        parseDescField(o, "Name            : "),
		Version:     parseDescField(o, "Version         : "),
		Description: parseDescField(o, "Description     : "),
		Size:        parseDescField(o, "Installed Size  : "),
		Url:         parseDescField(o, "URL             : "),
	}, nil
}

func parseDescField(o string, field string) string {
	splt1 := strings.Split(o, field)
	return strings.Split(splt1[1], "\n")[0]
}

// Install all .pkg.tar.zst files in provided directory.
func Install(dir string) error {
	out, err := system.Call("sudo pacman -U " + dir + " *.pkg.tar.zst")
	if err != nil {
		return errors.New("pacman unable to install:\n" + out)
	}
	return nil
}

var listcache map[string]string

// Get all installed packages from pacman.
func List() map[string]string {
	if listcache != nil {
		return listcache
	}
	o, err := system.Call("pacman -Q")
	if err != nil {
		fmt.Println(o)
		os.Exit(1)
	}
	o = strings.Trim(o, "\n")
	pkgs := map[string]string{}
	for _, pkg := range strings.Split(o, "\n") {
		spl := strings.Split(pkg, " ")
		pkgs[spl[0]] = spl[1]
	}
	listcache = pkgs
	return pkgs
}

// Get current version of specified package.
func Version(pkg string) string {
	o, err := system.Callf("pacman -Q %s", pkg)
	if err != nil {
		return ``
	}
	verAndRel := strings.Split(o, " ")[1]
	return strings.Trim(strings.Split(verAndRel, "-")[0], "\n")
}

// Eject list parameters from shell file, typically PKGBUILD.
func PkgbuildParams(file string, param string) ([]string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(string(f), "\n"+param+"=(")
	if len(splitted) < 2 {
		return nil, nil
	}
	splitted = strings.Split(splitted[1], ")")
	dirtyParams := splitted[0]
	var cleanParams []string
	for _, param := range splitParams(dirtyParams) {
		cleanParams = append(cleanParams, cleanParameter(param))
	}
	return cleanParams, nil
}

func splitParams(params string) []string {
	// TODO rework add quotas check
	params = strings.ReplaceAll(params, "\n", " ")
	for strings.Contains(params, "  ") {
		params = strings.ReplaceAll(params, "  ", " ")
	}
	return strings.Split(strings.Trim(params, " "), " ")
}

func cleanParameter(param string) string {
	param = strings.ReplaceAll(param, "'", "")
	param = strings.ReplaceAll(param, "\"", "")
	return param
}
