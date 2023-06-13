// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package local

import (
	"os"
	"path"
)

// Local directory, which stores GnuPG keys with emails in GPG key name.
type LocalKeyDir struct {
	Dir string
}

func (l *LocalKeyDir) ReadKey(owner, email string) ([]string, error) {
	keypath := path.Join(l.Dir, owner, email+".gpg")
	data, err := os.ReadFile(keypath)
	if err != nil {
		return nil, err
	}
	return []string{string(data)}, nil
}
