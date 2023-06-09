// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package registry

import (
	"errors"
	"io"
	"os"
	"path"
)

// This structure can be used to pushes calls incoming to pack. Verify method
// makes lookups to provided directory, checks wether there is a file with
// public key and verifies incoming packages using GPG. Also you can use 1
// level nested folders to get members/groups for some databases.
type LocalGpgDir struct {
	// Path to directory containing GnuPG files related to emails by name.
	// Example /folder/name@email.md.gpg
	GpgDir string
}

func (l *LocalGpgDir) ReadKey(owner, email string) (io.Reader, error) {
	f, err := os.Open(path.Join(l.GpgDir, owner, email+".gpg"))
	if err != nil {
		return nil, errors.Join(errors.New("user don't have rights"), err)
	}
	return f, err
}

// Parameters for package verification.
type VerificationParameters struct {
	Email     string
	Owner     string
	PkgReader io.Reader
	Signature []byte
}
