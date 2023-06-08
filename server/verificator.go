// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
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

// Parameters for package verification.
type VerificationParameters struct {
	Email     string
	Owner     string
	PkgReader io.Reader
	Signature []byte
}

func (l *LocalGpgDir) Verify(p VerificationParameters) ([]byte, error) {
	f, err := os.Open(path.Join(l.GpgDir, p.Owner, p.Email+".gpg"))
	if err != nil {
		return nil, err
	}

	pgpkey, err := crypto.NewKeyFromArmoredReader(f)
	if err != nil {
		return nil, err
	}

	pgpsig := crypto.NewPGPSignature(p.Signature)

	msg, err := io.ReadAll(p.PkgReader)
	if err != nil {
		return nil, err
	}
	pgpmes := crypto.NewPlainMessage(msg)

	keyring, err := crypto.NewKeyRing(pgpkey)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, idnt := range keyring.GetIdentities() {
		if idnt.Email == p.Email {
			found = true
		}
	}
	if !found {
		return nil, errors.New("unable to find email in keyring identities")
	}

	return msg, keyring.VerifyDetached(pgpmes, pgpsig, crypto.GetUnixTime())
}
