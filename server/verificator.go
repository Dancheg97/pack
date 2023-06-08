// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
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

func (l *LocalKeyring) Verify(owner string, email string, pkg io.Reader, sign []byte) error {
	f, err := os.Open(l.File)
	if err != nil {
		return err
	}

	pgpkey, err := crypto.NewKeyFromArmoredReader(f)
	if err != nil {
		return err
	}

	pgpsig := crypto.NewPGPSignature(sign)

	msg, err := io.ReadAll(pkg)
	if err != nil {
		return err
	}
	pgpmes := crypto.NewPlainMessage(msg)

	keyring, err := crypto.NewKeyRing(pgpkey)
	if err != nil {
		return err
	}

	return keyring.VerifyDetached(pgpmes, pgpsig, crypto.GetUnixTime())
}
