// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"fmnx.su/core/pack/pacman"
	"github.com/google/uuid"
)

// Structure, that allows to create handler for incoming push requests.
type PushHandler struct {
	// Direcotry, where push handler will store the resulting packages.
	CacheDir string
}

// Handler that can be used to upload user packages.
func (p *PushHandler) Push(w http.ResponseWriter, r *http.Request) {
	file := r.Header.Get("file")
	if !strings.HasSuffix(file, pkgext) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(path.Join(p.CacheDir, file)); err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	sign := r.Header.Get("sign")
	if sign == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpdir := path.Join("/tmp", "pack-"+uuid.New().String())
	err := os.MkdirAll(tmpdir, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpdir)

	f, err := os.Create(path.Join(tmpdir, file))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(path.Join(tmpdir, file+".sig"), sigdata, 0600)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = pacman.ValideSignature(tmpdir)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = pacman.CacheBuiltPackage(tmpdir, p.CacheDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[PUSH] - package accepted: " + file)
	w.WriteHeader(http.StatusOK)
}
