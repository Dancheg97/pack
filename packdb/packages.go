// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package packdb

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
)

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

type Package struct {
	PacmanName string `json:"pacman"`
	PackName   string `json:"pack"`
	Version    string `json:"version"`
	Branch     string `json:"branch"`
}

type NameType int

const (
	PACMAN NameType = iota
	PACK   NameType = iota
)

// Get list of packages installed by pack with metadata. This is readonly
// instance that does not affect database.
func List() []Package {
	return packages
}

// Update package in database, if package does not exist, it will be added.
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

// Get package by pacman or pack package name.
func Get(name string, nametype NameType) (*Package, error) {
	mu.Lock()
	defer mu.Unlock()
	switch nametype {
	case PACK:
		for _, p := range packages {
			if p.PackName == name {
				return &p, nil
			}
		}
	case PACMAN:
		for _, p := range packages {
			if p.PacmanName == name {
				return &p, nil
			}
		}
	}
	return nil, ErrNotFound
}

// Remove package from database. If package does not exist in database no
// action will be done.
func Remove(name string, nametype NameType) {
	mu.Lock()
	defer mu.Unlock()
	switch nametype {
	case PACK:
		for i, p := range packages {
			if p.PackName == name {
				packages = append(packages[:i], packages[i+1:]...)
			}
		}
	case PACMAN:
		for i, p := range packages {
			if p.PacmanName == name {
				packages = append(packages[:i], packages[i+1:]...)
			}
		}
	}
	savePackages()
}

// Pacman description with additional pack fields.
type Description struct {
	pacman.Description
	PackName    string `json:"pack-name"`
	PackVersion string `json:"pack-version"`
	PackBranch  string `json:"pack-branch"`
}

// Add pack fields to pacman package description.
func DescribeAppend(d pacman.Description) Description {
	pkg, _ := Get(d.Name, PACMAN)
	return Description{
		Description: d,
		PackName:    pkg.PackName,
		PackVersion: pkg.Version,
		PackBranch:  pkg.Branch,
	}
}
