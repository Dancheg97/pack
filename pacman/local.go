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

	"fmnx.io/core/pack/prnt"
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
		return errors.New("PKGBUILD generation failed")
	}
	return nil
}

// Some pacman package description fields.
type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Url         string `json:"url"`
	BuildDate   string `json:"build-date"`
	DependsOn   string `json:"depends-on"`
	RequiredBy  string `json:"required-by"`
}

// Get package description from pacman and parse it.
func Describe(pkg string) (*Package, error) {
	o, err := system.Callf("pacman -Qi %s", pkg)
	if err != nil {
		return nil, errors.New("package not found: " + pkg)
	}
	return &Package{
		Name:        parseDescField(o, "Name            : "),
		Version:     parseDescField(o, "Version         : "),
		Description: parseDescField(o, "Description     : "),
		Size:        parseDescField(o, "Installed Size  : "),
		Url:         parseDescField(o, "URL             : "),
		BuildDate:   parseDescField(o, "Build Date      : "),
		DependsOn:   parseDescField(o, "Depends On      : "),
		RequiredBy:  parseDescField(o, "Required By     : "),
	}, nil
}

func parseDescField(o string, field string) string {
	splt1 := strings.Split(o, field)
	return strings.Split(splt1[1], "\n")[0]
}

// Install all .pkg.tar.zst files in provided directory.
func Install(dir string) error {
	_, err := system.Call("sudo pacman -U " + dir + " *.pkg.tar.zst")
	if err != nil {
		return errors.New("pacman unable to install in dir " + dir)
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

// Get parameter from PKGBUILD proving path to direcotry.
func PkgbuildParam(file string, param string) (string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return ``, err
	}
	splitted := strings.Split(string(f), "\n"+param+"=")
	if len(splitted) < 2 {
		return ``, nil
	}
	value := strings.Split(splitted[1], "\n")[0]
	value = strings.ReplaceAll(value, "'", "")
	value = strings.ReplaceAll(value, "\"", "")
	return value, nil
}

// Eject list parameters from shell file, typically PKGBUILD.
func PkgbuildParamList(file string, param string) ([]string, error) {
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

// Function that takes a list of packages and returns those that are not
// installed in the end system.
func GetUninstalled(pkgs []string) []string {
	l := List()
	var out []string
	for _, pkg := range pkgs {
		if _, ok := l[pkg]; !ok {
			out = append(out, pkg)
		}
	}
	return out
}

// Get dependecies from PKGBUILD file.
func GetDeps(pkgbuild string) ([]string, error) {
	deps, err := PkgbuildParamList(pkgbuild, "depends")
	if err != nil {
		return nil, err
	}
	makedeps, err := PkgbuildParamList(pkgbuild, "makedepends")
	if err != nil {
		return nil, err
	}
	return append(deps, makedeps...), nil
}

// Try to remove all packages at once.
func Remove(pkgs []string) error {
	pkgsStr := strings.Join(pkgs, " ")
	_, err := system.Callf("sudo pacman --noconfirm -R %s", pkgsStr)
	if err != nil {
		return errors.New("pacman unable to remove")
	}
	prnt.Yellow("Packages removed: ", pkgsStr)
	return nil
}

// Get pacman packages from parsed removal command.
func PrintNotFoundPackages(o string) {
	o = strings.ReplaceAll(o, "\n", " ")
	o = strings.ReplaceAll(o, `error: target not found: `, "")
	prnt.Red("Packages not found: ", o)
}
