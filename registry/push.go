// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package registry

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
)

// Push packages to registry.
func (p *Registry) Push(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("file")
	email := r.Header.Get("email")
	sign := r.Header.Get("sign")
	owner := r.Header.Get("owner")

	p.Amsg("Arch package request revieced: " + filename)

	p.Smsg("Decoding signature", 1, 12)
	sigdata, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	pgpsig := crypto.NewPGPSignature(sigdata)

	p.Smsg("Retrieving related GPG keys", 2, 12)
	armoredKeys, err := p.KeyReader.ReadKey(owner, email)
	if err != nil {
		p.end(w, http.StatusBadRequest, err)
		return
	}

	p.Smsg("Validating GPG identities", 3, 12)
	var matchedKeyring *crypto.KeyRing
	for _, armor := range armoredKeys {
		pgpkey, err := crypto.NewKeyFromArmored(armor)
		if err != nil {
			p.end(w, http.StatusBadRequest, err)
			return
		}
		keyring, err := crypto.NewKeyRing(pgpkey)
		if err != nil {
			p.end(w, http.StatusBadRequest, err)
			return
		}
		for _, idnt := range keyring.GetIdentities() {
			if idnt.Email == email {
				matchedKeyring = keyring
				break
			}
		}
		if matchedKeyring != nil {
			break
		}
	}
	if matchedKeyring == nil {
		msg := "GPG key related to " + email + " not found"
		p.end(w, http.StatusBadRequest, errors.New(msg))
		return
	}

	p.Smsg("Reading request body", 4, 12)
	pkgdata, err := io.ReadAll(r.Body)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()
	pgpmes := crypto.NewPlainMessage(pkgdata)

	p.Smsg("Validating package signature", 5, 12)
	err = matchedKeyring.VerifyDetached(pgpmes, pgpsig, crypto.GetUnixTime())
	if err != nil {
		p.end(w, http.StatusUnauthorized, err)
		return
	}

	p.Smsg("Preparing temporary directory", 6, 12)
	tmpdir := path.Join(p.TmpDir, uuid.New().String())
	err = os.MkdirAll(tmpdir, 0600)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	p.Smsg("Loading old database if exists", 7, 12)
	dbname := strings.Join([]string{owner, p.Dbname, "db"}, ".")
	db, err := p.FileStorage.Get(dbname)
	if err == nil {
		err = os.WriteFile(path.Join(tmpdir, dbname), db, 0600)
		if err != nil {
			p.end(w, http.StatusInternalServerError, err)
			return
		}
	}

	p.Smsg("Loading old database archive if exists", 8, 12)
	dbarchivename := strings.Join([]string{owner, p.Dbname, "db.tar.gz"}, ".")
	dbarchive, err := p.FileStorage.Get(dbarchivename)
	if err == nil {
		err = os.WriteFile(path.Join(tmpdir, dbarchivename), dbarchive, 0600)
		if err != nil {
			p.end(w, http.StatusInternalServerError, err)
			return
		}
	}

	p.Smsg("Adding package to database", 9, 12)
	pkgfilepath := path.Join(tmpdir, filename)
	err = os.WriteFile(pkgfilepath, pkgdata, 0600)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	var rabuf bytes.Buffer
	repoaddcmd := exec.Command("repo-add", dbname, pkgfilepath)
	repoaddcmd.Stderr = &rabuf
	if err := repoaddcmd.Run(); err != nil {
		p.end(w, http.StatusInternalServerError, errors.New(rabuf.String()))
		return
	}

	p.Smsg("Saving new package and signature files", 10, 12)
	var filebuf bytes.Buffer
	_, err = filebuf.Write(pkgdata)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	savefile := strings.Join([]string{owner, filename}, ".")
	err = p.FileStorage.Save(savefile, &filebuf)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	var sigbuf bytes.Buffer
	_, err = sigbuf.Write(sigdata)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	err = p.FileStorage.Save(savefile+".sig", &sigbuf)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	p.Smsg("Saving new database and archive files", 11, 12)
	f, err := os.Open(path.Join(tmpdir, dbname))
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	err = p.FileStorage.Save(dbname, f)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	f, err = os.Open(path.Join(tmpdir, dbarchivename))
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	err = p.FileStorage.Save(dbarchivename, f)
	if err != nil {
		p.end(w, http.StatusInternalServerError, err)
		return
	}

	p.Smsg("Saving package metadata", 12, 12)
	if p.MetadataSaver != nil {
		var mdcmdbuf bytes.Buffer
		mdcmd := exec.Command("pacman", "-Qpi", pkgfilepath)
		mdcmd.Stdout = &mdcmdbuf
		if err = mdcmd.Run(); err != nil {
			p.end(w, http.StatusInternalServerError, err)
			return
		}
		out := mdcmdbuf.String()
		var metadata = struct {
			Name           string
			Version        string
			Description    string
			Architecture   string
			URL            string
			Licenses       string
			Groups         string
			Provides       string
			DependsOn      string
			OptionalDeps   string
			ConflictsWith  string
			Replaces       string
			CompressedSize string
			InstalledSize  string
			Packager       string
			BuildDate      string
			InstallScript  string
			ValidatedBy    string
			Signatures     string
		}{
			Name:           parseArchMdField(out, "Name            : "),
			Version:        parseArchMdField(out, "Version         : "),
			Description:    parseArchMdField(out, "Description     : "),
			Architecture:   parseArchMdField(out, "Architecture    : "),
			URL:            parseArchMdField(out, "URL             : "),
			Licenses:       parseArchMdField(out, "Licenses        : "),
			Groups:         parseArchMdField(out, "Groups          : "),
			Provides:       parseArchMdField(out, "Provides        : "),
			DependsOn:      parseArchMdField(out, "Depends On      : "),
			OptionalDeps:   parseArchMdField(out, "Optional Deps   : "),
			ConflictsWith:  parseArchMdField(out, "Conflicts With  : "),
			Replaces:       parseArchMdField(out, "Replaces        : "),
			CompressedSize: parseArchMdField(out, "Compressed Size : "),
			InstalledSize:  parseArchMdField(out, "Installed Size  : "),
			Packager:       parseArchMdField(out, "Packager        : "),
			BuildDate:      parseArchMdField(out, "Build Date      : "),
			InstallScript:  parseArchMdField(out, "Install Script  : "),
			ValidatedBy:    parseArchMdField(out, "Validated By    : "),
			Signatures:     parseArchMdField(out, "Signatures      : "),
		}
		err = p.MetadataSaver.UpdateMetadata(filename, owner, metadata)
		if err != nil {
			p.end(w, http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
