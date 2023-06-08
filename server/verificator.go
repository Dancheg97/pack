// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"errors"
	"io"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Local keyring can serve as GPG verificator for local package server. It
// requires path to specific keyring file.
type LocalKeyring struct {
	// Path to file containing GnuPG keyring.
	File string
}

// Parameters for package verification.
type VerificationParameters struct {
	Email     string
	Owner     string
	PkgReader io.Reader
	Signature []byte
}

func (l *LocalKeyring) Verify(p VerificationParameters) error {
	f, err := os.Open(l.File)
	if err != nil {
		return err
	}

	pgpkey, err := crypto.NewKeyFromArmoredReader(f)
	if err != nil {
		return err
	}

	pgpsig := crypto.NewPGPSignature(p.Signature)

	msg, err := io.ReadAll(p.PkgReader)
	if err != nil {
		return err
	}
	pgpmes := crypto.NewPlainMessage(msg)

	keyring, err := crypto.NewKeyRing(pgpkey)
	if err != nil {
		return err
	}

	var found bool
	for _, idnt := range keyring.GetIdentities() {
		if idnt.Email == p.Email {
			found = true
		}
	}
	if !found {
		return errors.New("unable to find email in keyring identities")
	}

	return keyring.VerifyDetached(pgpmes, pgpsig, crypto.GetUnixTime())
}
