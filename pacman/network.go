// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

// This package acts as library wrapper over pacman and makepkg.
// Package is safe for concurrent usage.

import (
	"errors"
	"os"
	"strings"

	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/system"
)

// This command will build package in provided directory. Also installing
// missing packages with pacman. Safe for concurrent usage.
func Build(dir string) error {
	mu.Lock()
	defer mu.Unlock()
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	o, err := system.Call("makepkg -sf --noconfirm")
	if err != nil {
		return errors.New(o)
	}
	return nil
}

// Update listed pacman packages with sync command.
func Update(pkgs []string) error {
	mu.Lock()
	defer mu.Unlock()
	joined := strings.Join(pkgs, " ")
	_, err := system.Call("sudo pacman --noconfirm --needed -Sy " + joined)
	if err != nil {
		return errors.New("pacman unable to update")
	}
	return nil
}

// Base information about outdated pacman package.
type OutdatedPackage struct {
	Name           string
	CurrentVersion string
	NewVersion     string
}

// List outdated pacman packages.
func Outdated() ([]OutdatedPackage, error) {
	links, err := oudatedLinks()
	if err != nil {
		return nil, err
	}
	var op []OutdatedPackage
	for _, link := range links {
		name, newver := parsePackageLink(link)
		currver := Version(name)
		op = append(op, OutdatedPackage{
			Name:           name,
			CurrentVersion: currver,
			NewVersion:     newver,
		})
	}
	return op, nil
}

// Get all links for outdated packages.
func oudatedLinks() ([]string, error) {
	o, err := system.Call("sudo pacman -Syup")
	if err != nil {
		return nil, errors.New("pacman unable to get outdated")
	}
	if !strings.Contains(o, "https://") {
		return nil, nil
	}
	splt := strings.Split(o, "downloading...\n")
	pkgsLinksString := strings.Trim(splt[len(splt)-1], "\n")
	return strings.Split(pkgsLinksString, "\n"), nil
}

// Parse link leading to some outdated package, to get package name and new
// version. Returning package name and new version.
func parsePackageLink(link string) (string, string) {
	linkSplit := strings.Split(link, "/")
	file := linkSplit[len(linkSplit)-1]
	fileSplit := strings.Split(file, "-")
	packageName := strings.Join(fileSplit[0:len(fileSplit)-3], "-")
	newVersion := fileSplit[len(fileSplit)-3]
	return packageName, newVersion
}

func Install(pkgs []string) error {
	mu.Lock()
	defer mu.Unlock()
	uninstalled := GetInstalled(pkgs)
	if len(uninstalled) == 0 {
		return nil
	}
	joined := strings.Join(uninstalled, " ")
	o, err := system.Callf("sudo pacman --needed --noconfirm -S %s", joined)
	if err != nil {
		return errors.New(o)
	}
	prnt.Green("Installed: ", joined)
	return nil
}

// Check which packages are already installed and remove them from list.
func GetInstalled(pkgs []string) []string {
	var uninstalledPkgs []string
	for _, pkg := range pkgs {
		_, err := system.Callf("pacman -Q %s", pkg)
		if err != nil {
			uninstalledPkgs = append(uninstalledPkgs, pkg)
		}
	}
	return uninstalledPkgs
}
