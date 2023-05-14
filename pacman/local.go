// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

// This package acts as library wrapper over pacman and makepkg.
// Package is safe for concurrent usage.

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/system"
	"fmnx.su/core/pack/tmpl"
)

// Global pacman mutext for safe installing and building packages.
var mu sync.Mutex

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

// InstallDir all .pkg.tar.zst files in provided directory.
func InstallDir(dir string) error {
	mu.Lock()
	defer mu.Unlock()
	cmd := "sudo pacman --noconfirm --needed -U " + dir + "/*.pkg.tar.zst"
	o, err := system.Call(cmd)
	if err != nil {
		return fmt.Errorf("pacman cant install %s\n%s", dir, o)
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
func PkgbuildParam(dir string, param string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	err := os.Chdir(dir)
	if err != nil {
		return ``, err
	}
	o, err := system.Callf("source PKGBUILD && echo $%s", param)
	if err != nil {
		return ``, err
	}
	return strings.Trim(o, "\n"), nil
}

// Eject list parameters from shell file, typically PKGBUILD.
func PkgbuildParamList(dir string, param string) ([]string, error) {
	mu.Lock()
	defer mu.Unlock()
	err := os.Chdir(dir)
	if err != nil {
		return nil, err
	}
	err = system.WriteFile(dir+"/l.sh", fmt.Sprintf(tmpl.PkgbuildList, param))
	if err != nil {
		return nil, err
	}
	o, err := system.Call("sh l.sh")
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Trim(o, "\n"), "\n"), nil
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
func GetDeps(dir string) ([]string, error) {
	deps, err := PkgbuildParamList(dir, "depends")
	if err != nil {
		return nil, err
	}
	makedeps, err := PkgbuildParamList(dir, "makedepends")
	if err != nil {
		return nil, err
	}
	return append(deps, makedeps...), nil
}

// Try to remove all packages at once.
func Remove(pkgs []string) error {
	mu.Lock()
	defer mu.Unlock()
	pkgsStr := strings.Join(pkgs, " ")
	o, err := system.Callf("sudo pacman --noconfirm -R %s", pkgsStr)
	if err != nil {
		return errors.New(o)
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

// Check which packages from list could not be installed with pacman. Takes a
// list of packages and returning those, that could not be installed.
func GetUnreachable(pkgs []string) []string {
	o, err := system.Call("pacman -Ssq")
	if err != nil {
		prnt.Red("Unexpected error on system call: ", "pacman -Ssq")
		fmt.Println(o)
		os.Exit(1)
	}
	available := map[string]struct{}{}
	for _, pkg := range strings.Split(o, "\n") {
		available[pkg] = struct{}{}
	}
	var unavailable []string
	for _, pkg := range pkgs {
		_, ok := available[pkg]
		if !ok {
			unavailable = append(unavailable, pkg)
		}
	}
	return unavailable
}
