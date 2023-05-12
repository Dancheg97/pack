// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pacman

// This package acts as library wrapper over pacman and makepkg.

import (
	"errors"
	"os"
	"strings"

	"fmnx.io/core/pack/system"
)

// This command will build package in provided directory. Also installing
// missing packages with pacman.
func Build(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	o, err := system.Call("makepkg -sf --noconfirm")
	if err != nil {
		return errors.New("pacman unable to build:\n" + o)
	}
	return nil
}

// Update listed pacman packages with sync command.
func Update(pkgs []string) error {
	joined := strings.Join(pkgs, " ")
	o, err := system.Call("sudo pacman --noconfirm -Sy " + joined)
	if err != nil {
		return errors.New("pacman unable to update:\n" + o)
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
		return nil, errors.New("pacman unable to get outdated:\n" + o)
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
