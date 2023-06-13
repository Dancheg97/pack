// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package local

import (
	"io"
	"os"
	"path"
)

// Local directory with files to store the contents.
type DirStorage struct {
	Dir string
}

func (d *DirStorage) Get(key string) ([]byte, error) {
	return os.ReadFile(path.Join(d.Dir, key))
}

func (d *DirStorage) Save(key string, content io.Reader) error {
	path := path.Join(d.Dir, key)
	os.RemoveAll(path)
	f, err := os.Open(path)
	if err != nil {
		f, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	_, err = f.ReadFrom(content)
	return err
}
