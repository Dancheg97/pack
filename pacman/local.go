// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pacman

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
)

// Check directory to exist and contain PKGBUILD.
func Check(dir string) error {
	_, err := os.Stat(dir + "/PKGBUILD")
	return err
}

// Geberate PKGBUILD file template. Provide pacman package name and git url.
func Generate(dir string, name string, url string) error {
	content := fmt.Sprintf(tmpl.PKGBUILD, name, url)
	err := os.WriteFile(dir+"/PKGBUILD", []byte(content), 0o600)
	if err != nil {
		return errors.New("PKGBUILD generation failed:\n" + err.Error())
	}
	return nil
}

// Some pacman package description fields.
type Description struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Url         string `json:"url"`
}

// Get package description from pacman and parse it.
func Describe(pkg string) (Description, error) {
	o, err := system.Callf("pacman -Qi %s", pkg)
	if err != nil {
		return Description{}, errors.New("package not found: " + pkg)
	}
	return Description{
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
