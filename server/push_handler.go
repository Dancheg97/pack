// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Structure, that allows to create handler for incoming arch packages
// push requests.
type PushHandler struct {
	// Direcotry, where push handler will store the resulting packages.
	CacheDir string

	// Public key source, get public key related to specific user. If not
	// provided default gnupg signature verification scheme will be used.
	PubkeySource

	// Subdir source, if enabled packages would be created in subdirectories
	// instead of base directory, allowing
	// SubdirSource

	ErrLogger  Logger
	InfoLogger Logger
}

// Handler that can be used to upload user packages.
func (p *PushHandler) Push(w http.ResponseWriter, r *http.Request) {
	ep := weberrwriter{respWriter: w, logger: p.ErrLogger}

	file := r.Header.Get("file")
	if !strings.HasSuffix(file, pkgext) {
		ep.write(http.StatusBadRequest, "unable to get file name from header")
		return
	}

	if _, err := os.Stat(path.Join(p.CacheDir, file)); err == nil {
		ep.write(http.StatusConflict, "package exists")
		return
	}

	sign := r.Header.Get("sign")
	if sign == "" {
		ep.write(http.StatusConflict, "unable to get signature from header")
		return
	}

	email := r.Header.Get("email")
	if sign == "" {
		ep.write(http.StatusConflict, "unable to get email from header")
		return
	}

	pkgdata, err := io.ReadAll(r.Body)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to read body")
		return
	}
	pkgmes := crypto.NewPlainMessage(pkgdata)

	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to decode sign base64")
		return
	}
	signature := crypto.NewPGPSignature(sigdata)

	keys, err := p.PubkeySource.Get(email)
	if err != nil {
		ep.write(http.StatusUnauthorized, "no GPG keys for email: "+email)
		return
	}

	// Verification with all public keys associated to specific email adress.
	var verified bool
	var trace []string
	for _, key := range keys {
		if verified {
			break
		}
		pk, err := crypto.NewKeyFromArmored(key)
		if err != nil {
			trace = append(trace, "unable to get key from armored")
			continue
		}
		kr, err := crypto.NewKeyRing(pk)
		if err != nil {
			trace = append(trace, "unable to get keyring from key")
			continue
		}
		for _, ident := range kr.GetIdentities() {
			if ident.Email != email {
				continue
			}
			err = kr.VerifyDetached(pkgmes, signature, crypto.GetUnixTime())
			if err != nil {
				trace = append(trace, "verification with key failed")
				break
			}
			verified = true
			break
		}
	}
	if !verified {
		ep.write(http.StatusUnauthorized, strings.Join(trace, " "))
		return
	}

	err = os.WriteFile(path.Join(p.CacheDir, file), pkgdata, os.ModePerm)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to write package file")
		return
	}

	err = os.WriteFile(path.Join(p.CacheDir, file+".sig"), sigdata, os.ModePerm)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to write signature file")
		return
	}

	p.InfoLogger.Printf("[PUSH] - package accepted: " + file)
	w.WriteHeader(http.StatusOK)
}

// Structure that will log errors, form response bodies and send http codes.
type weberrwriter struct {
	respWriter http.ResponseWriter
	logger     Logger
}

func (e *weberrwriter) write(status int, msg string) {
	e.logger.Printf(msg)
	e.respWriter.WriteHeader(status)
	e.respWriter.Write([]byte(msg)) //nolint:errcheck
}
