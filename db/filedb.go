// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package db

import (
	"errors"
	"os"
	"path"
	"strings"
)

// Simple file with user creds.
type FileDB struct {
	Path string
}

func GetFileDb(p string) (*FileDB, error) {
	return &FileDB{Path: path.Join(p, "users")}, nil
}

func (d *FileDB) Fill(users []string) error {
	for _, u := range users {
		splt := strings.Split(u, ":")
		err := d.Update(splt[0], splt[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *FileDB) List() ([]string, error) {
	b, err := os.ReadFile(d.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var users []string
	for _, u := range strings.Split(string(b), "\n") {
		users = append(users, strings.Split(u, ":")[0])
	}
	return users, nil
}

func (d *FileDB) Validate(name string, password string) bool {
	b, err := os.ReadFile(d.Path)
	if err != nil {
		return false
	}
	for _, u := range strings.Split(string(b), "\n") {
		splt := strings.Split(u, ":")
		if name == splt[0] && password == splt[1] {
			return true
		}
	}
	return false
}

func (d *FileDB) Update(name string, password string) error {
	usr := []byte(name + ":" + password)
	b, err := os.ReadFile(d.Path)
	if err != nil {
		return os.WriteFile(d.Path, usr, 0600)
	}
	splt := strings.Split(string(b), "\n")
	for i, v := range splt {
		usplt := strings.Split(v, ":")
		if usplt[0] == name {
			splt[i] = name + ":" + password
			return os.WriteFile(d.Path, []byte(strings.Join(splt, "\n")), 0600)
		}
	}
	splt = append(splt, string(usr))
	return os.WriteFile(d.Path, []byte(strings.Join(splt, "\n")), 0600)
}
