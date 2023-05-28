// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package db

import (
	"os"
	"path"
	"strings"
)

// Simple file with user creds.
type FileDB struct {
	Path string
}

func GetFileDb(p string) (Database, error) {
	return &FileDB{Path: path.Join(p, "users")}, nil
}

func (d *FileDB) List() ([]string, error) {
	b, err := os.ReadFile(d.Path)
	if err != nil {
		return nil, err
	}
	var users []string
	for _, u := range strings.Split(string(b), "\n") {
		users = append(users, strings.Split(u, " ")[0])
	}
	return users, nil
}

func (d *FileDB) Validate(name string, password string) bool {
	b, err := os.ReadFile(d.Path)
	if err != nil {
		return false
	}
	for _, u := range strings.Split(string(b), "\n") {
		splt := strings.Split(u, " ")
		if name == splt[0] && password == splt[1] {
			return true
		}
	}
	return false
}

func (d *FileDB) Add(name string, password string) error {
	newUser := []byte(name + " " + password + "\n")
	b, err := os.ReadFile(d.Path)
	if err != nil {
		return os.WriteFile(d.Path, newUser, 0600)
	}
	return os.WriteFile(d.Path, append(b, newUser...), 0600)
}
