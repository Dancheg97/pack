// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package registry

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"fmnx.su/core/pack/tmpl"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Parameters required to get http.Pusher.
type Pusher struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Interface that will be used to verify incoming packages.
	KeySource KeyReader
	// Interface that will be used to add new packages to database.
	DbFormer DbFormer
}

// An interface, that can check that package is signed by valid email and GnuPG
// key belogns to required keyring/exists in other trusted source for specified
// package owner. Verificator returns bytes of package it have verified.
type KeyReader interface {
	ReadKey(owner, email string) (io.Reader, error)
}

// Interface, that accepts package bytes body, writes signature and forms
// database with new packages.
type DbFormer interface {
	AddPkg(p AddPkgParameters) error
}

// Handler that can be used to upload user packages.
func (p *Pusher) Push(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("file")
	email := r.Header.Get("email")
	sign := r.Header.Get("sign")
	owner := r.Header.Get("owner")
	force := r.Header.Get("force") == "true"

	tmpl.Amsg(p.Stdout, "Request: "+filename)

	tmpl.Smsg(p.Stdout, "Decoding signature", 1, 8)
	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	pgpsig := crypto.NewPGPSignature(sigdata)

	tmpl.Smsg(p.Stdout, "Preparing GPG key", 2, 8)
	keyreader, err := p.KeySource.ReadKey(owner, email)
	if err != nil {
		p.end(w, http.StatusBadRequest, err)
		return
	}

	pgpkey, err := crypto.NewKeyFromArmoredReader(keyreader)
	if err != nil {
		p.end(w, http.StatusBadRequest, err)
		return
	}

	tmpl.Smsg(p.Stdout, "Creating GPG keyring", 3, 8)
	keyring, err := crypto.NewKeyRing(pgpkey)
	if err != nil {
		p.end(w, http.StatusBadRequest, err)
		return
	}

	tmpl.Smsg(p.Stdout, "Validating email in keyring", 4, 8)
	var found bool
	for _, idnt := range keyring.GetIdentities() {
		if idnt.Email == email {
			found = true
		}
	}
	if !found {
		err := errors.New("email not found in keyring")
		p.end(w, http.StatusUnauthorized, err)
		return
	}

	tmpl.Smsg(p.Stdout, "Reading request body", 5, 8)
	pkgdata, err := io.ReadAll(r.Body)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	pgpmes := crypto.NewPlainMessage(pkgdata)

	tmpl.Smsg(p.Stdout, "Validating package signature", 6, 8)
	err = keyring.VerifyDetached(pgpmes, pgpsig, crypto.GetUnixTime())
	if err != nil {
		p.end(w, http.StatusUnauthorized, err)
		return
	}

	tmpl.Smsg(p.Stdout, "Updating database", 7, 8)
	err = p.DbFormer.AddPkg(AddPkgParameters{
		Package:  pkgdata,
		Sign:     sigdata,
		Filename: filename,
		Owner:    owner,
		Force:    force,
	})
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	tmpl.Smsg(p.Stdout, "Accepted "+filename, 8, 8)
	w.WriteHeader(http.StatusOK)
}

// Write header, log error and end request.
func (p *Pusher) end(w http.ResponseWriter, status int, msg error) {
	errmsg := []byte(tmpl.Err + msg.Error())
	p.Stderr.Write(errmsg)
	w.WriteHeader(status)
	w.Write(errmsg)
}
