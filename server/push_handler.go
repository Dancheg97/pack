// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"encoding/base64"
	"net/http"
	"os"
	"path"
	"strings"

	"fmnx.su/core/pack/pacman"
	"github.com/google/uuid"
)

// Structure, that allows to create handler for incoming arch packages
// push requests.
type PushHandler struct {
	// Direcotry, where push handler will store the resulting packages.
	CacheDir string
	// Directory, where handler will temporarily store file to check signature.
	TmpDir string

	// Public key source, get public key related to specific user. If not
	// provided default gnupg signature verification scheme will be used.
	// PubkeySource

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

	tmpdir := path.Join(p.TmpDir, "pack-"+uuid.New().String())
	err := os.MkdirAll(tmpdir, os.ModePerm)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to create cache directory")
		return
	}
	defer os.RemoveAll(tmpdir)

	f, err := os.Create(path.Join(tmpdir, file))
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to create file")
		return
	}

	if _, err = f.ReadFrom(r.Body); err != nil {
		ep.write(http.StatusInternalServerError, "unable read file body")
		return
	}

	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to decode sign base64")
		return
	}
	err = os.WriteFile(path.Join(tmpdir, file+".sig"), sigdata, os.ModePerm)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to write sign file")
		return
	}

	err = pacman.ValideSignature(tmpdir)
	if err != nil {
		ep.write(http.StatusInternalServerError, "signature is not validate with gnupg")
		return
	}

	err = pacman.CacheBuiltPackage(tmpdir, p.CacheDir)
	if err != nil {
		ep.write(http.StatusInternalServerError, "unable to move package to cache")
		w.WriteHeader(http.StatusInternalServerError)
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
