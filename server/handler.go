// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"encoding/base64"
	"io"
	"net/http"
)

// Parameters required to get http.Pusher.
type Pusher struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Interface that will be used to verify incoming packages.
	GPGVireivicator GPGVireivicator
	// Interface that will be used to add new packages to database.
	DbFormer DbFormer
}

// An interface, that can check that package is signed by valid email and GnuPG
// key belogns to required keyring/exists in other trusted source for specified
// package owner. Verificator returns bytes of package it have verified.
type GPGVireivicator interface {
	Verify(p VerificationParameters) ([]byte, error)
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
	force := r.Header.Get("force") == "true"
	owner := r.Header.Get("owner")

	sigbytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		p.end(w, http.StatusInternalServerError, "unable to decode signature")
		return
	}

	pkgbytes, err := p.GPGVireivicator.Verify(VerificationParameters{
		Email:     email,
		Owner:     owner,
		PkgReader: r.Body,
		Signature: sigbytes,
	})
	if err != nil {
		p.end(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = p.DbFormer.AddPkg(AddPkgParameters{
		Package:  pkgbytes,
		Sign:     sigbytes,
		Filename: filename,
		Force:    force,
	})
	if err != nil {
		p.end(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Write header, log error and end request.
func (p *Pusher) end(w http.ResponseWriter, status int, msg string) {
	msgbytes := []byte(msg)
	p.Stderr.Write(msgbytes) //nolint
	w.WriteHeader(status)
	w.Write(msgbytes) //nolint
}

// func (p *PushHandler) Push(w http.ResponseWriter, r *http.Request) {
// 	ep := weberrwriter{respWriter: w, logger: p.ErrLogger}

// 	file := r.Header.Get("file")
// 	if !strings.HasSuffix(file, pkgext) {
// 		ep.write(http.StatusBadRequest, "unable to get file name from header")
// 		return
// 	}

// 	if _, err := os.Stat(path.Join(p.CacheDir, file)); err == nil {
// 		ep.write(http.StatusConflict, "package exists")
// 		return
// 	}

// 	sign := r.Header.Get("sign")
// 	if sign == "" {
// 		ep.write(http.StatusConflict, "unable to get signature from header")
// 		return
// 	}

// 	email := r.Header.Get("email")
// 	if sign == "" {
// 		ep.write(http.StatusConflict, "unable to get email from header")
// 		return
// 	}

// 	pkgdata, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		ep.write(http.StatusInternalServerError, "unable to read body")
// 		return
// 	}
// 	pkgmes := crypto.NewPlainMessage(pkgdata)

// 	sigdata, err := base64.StdEncoding.DecodeString(sign)
// 	if err != nil {
// 		ep.write(http.StatusInternalServerError, "unable to decode sign base64")
// 		return
// 	}
// 	signature := crypto.NewPGPSignature(sigdata)

// 	keys, err := p.PubkeySource.Get(email)
// 	if err != nil {
// 		ep.write(http.StatusUnauthorized, "unable to get GPG keys for: "+email)
// 		return
// 	}

// 	// Verification with all public keys associated to specific email adress.
// 	var verified bool
// 	var trace []string
// 	for _, key := range keys {
// 		if verified {
// 			break
// 		}
// 		pk, err := crypto.NewKeyFromArmored(key)
// 		if err != nil {
// 			trace = append(trace, "unable to get key from armored")
// 			continue
// 		}
// 		kr, err := crypto.NewKeyRing(pk)
// 		if err != nil {
// 			trace = append(trace, "unable to get keyring from key")
// 			continue
// 		}
// 		for _, ident := range kr.GetIdentities() {
// 			if ident.Email != email {
// 				continue
// 			}
// 			err = kr.VerifyDetached(pkgmes, signature, crypto.GetUnixTime())
// 			if err != nil {
// 				trace = append(trace, "verification with key failed")
// 				break
// 			}
// 			verified = true
// 			break
// 		}
// 	}
// 	if !verified {
// 		ep.write(http.StatusUnauthorized, strings.Join(trace, " "))
// 		return
// 	}

// 	err = os.WriteFile(path.Join(p.CacheDir, file), pkgdata, os.ModePerm)
// 	if err != nil {
// 		ep.write(http.StatusInternalServerError, "unable to write package file")
// 		return
// 	}

// 	err = os.WriteFile(path.Join(p.CacheDir, file+".sig"), sigdata, os.ModePerm)
// 	if err != nil {
// 		ep.write(http.StatusInternalServerError, "unable to write signature file")
// 		return
// 	}

// 	p.InfoLogger.Printf("[PUSH] - package accepted: " + file)
// 	w.WriteHeader(http.StatusOK)
// }
