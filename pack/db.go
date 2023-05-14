// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pack

// This library contains operations related only to pack json database.
// All functions are safe for concurrent usage.

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
)

// Data about pack package stored in pack database.
type Package struct {
	PackName      string `json:"pack-name"`
	PacmanName    string `json:"pacman-name"`
	Version       string `json:"version"`
	DefaultBranch string `json:"default-branch"`
}

var (
	ErrNotFound      = errors.New("package is not stored in pack")
	ErrAlreadyExists = errors.New("package already exists")
	mu               sync.Mutex
	packages         []Package
)

func init() {
	b, err := os.ReadFile(config.MapFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &packages)
	if err != nil {
		prnt.Red("Unable to parse package mapping file: ", config.MapFile)
		os.Exit(1)
	}
	pacmanList := pacman.List()
	for i, p := range packages {
		if _, ok := pacmanList[p.PacmanName]; !ok {
			packages = append(packages[:i], packages[i+1:]...)
		}
	}
}

func savePackages() {
	b, err := json.Marshal(packages)
	if err != nil {
		prnt.Red("Unable to parse packages: ", config.MapFile)
		os.Exit(1)
	}
	err = os.WriteFile(config.MapFile, b, 0o600)
	if err != nil {
		prnt.Red("Unable to save package mapping file: ", config.MapFile)
		os.Exit(1)
	}
}

// Get list of pack packages.
func List() []Package {
	return packages
}

// Update package in database, if package does not exist, it will be added. If
// it exists, new information will be appended.
func Update(pkg Package) {
	mu.Lock()
	defer mu.Unlock()
	for i, p := range packages {
		if pkg.PackName == p.PackName {
			packages[i] = pkg
			savePackages()
			return
		}
	}
	packages = append(packages, pkg)
	savePackages()
}

// Get package by pack package name.
func Get(name string) (*Package, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, p := range packages {
		if p.PackName == name {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

// Get package by pack package name.
func GetByPacmanName(name string) (*Package, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, p := range packages {
		if p.PacmanName == name {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

// Remove package from database. If package does not exist in database no
// action will be done.
func Remove(pkgs []string) {
	mu.Lock()
	defer mu.Unlock()
	for _, target := range pkgs {
		for i, check := range packages {
			if check.PackName == target {
				packages = append(packages[:i], packages[i+1:]...)
			}
		}
	}
	savePackages()
}

// Get uninstalled pack packages.
func GetUninstalled(pkgs []string) []string {
	var uninstalled []string
	for _, pkg := range pkgs {
		_, err := Get(pkg)
		if err != nil {
			uninstalled = append(uninstalled, pkg)
		}
	}
	return uninstalled
}
