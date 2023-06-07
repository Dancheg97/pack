// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/tmpl"
)

// Parameters that can be used to build packages.
type BuildParameters struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Directory where resulting package and signature will be moved.
	Dir string
	// Do not ask for any confirmation on build/installation.
	Quick bool
	// Syncronize/reinstall package after build.
	Syncbuild bool
	// Remove dependencies after successful build.
	Rmdeps bool
	// Do not clean workspace before and after build.
	Garbage bool
}

func builddefault() *BuildParameters {
	return &BuildParameters{
		Dir:       "/var/cache/pacman/pkg",
		Syncbuild: true,
	}
}

// Build package in current directory with provided arguements
func Build(prms ...BuildParameters) error {
	p := formOptions(prms, builddefault)

	err := checkGnupg()
	if err != nil {
		return err
	}

	err = validatePackager()
	if err != nil {
		return err
	}

	err = pacman.Makepkg(pacman.MakepkgOptions{
		Sign:       true,
		Stdout:     p.Stdout,
		Stderr:     p.Stderr,
		Stdin:      p.Stdin,
		Clean:      !p.Garbage,
		CleanBuild: !p.Garbage,
		Force:      !p.Garbage,
		Install:    p.Syncbuild,
		RmDeps:     p.Rmdeps,
		SyncDeps:   p.Syncbuild,
		Needed:     !p.Syncbuild,
		NoConfirm:  p.Quick,
	})
	if err != nil {
		return err
	}

	return exec.Command("bash", "-c", "sudo mv *.pkg.tar.zst* "+p.Dir).Run() //nolint
}

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func checkGnupg() error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	gpgdir, err := os.ReadDir(path.Join(hd, ".gnupg"))
	if err != nil {
		return err
	}
	for _, de := range gpgdir {
		if strings.Contains(de.Name(), "private-keys") {
			return nil
		}
	}
	return errors.New(tmpl.ErrGnuPGprivkeyNotFound)
}

// Validate, that packager defined in /etc/makepkg.conf matches signer
// authority in GnuPG.
func validatePackager() error {
	keySigner, err := gnuPGIdentity()
	if err != nil {
		return err
	}
	f, err := os.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return err
	}
	splt := strings.Split(string(f), "\nPACKAGER=\"")
	if len(splt) != 2 {
		return errors.New(tmpl.ErrNoPackager)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	if confPackager != keySigner {
		return errors.New(tmpl.ErrSignerMissmatch)
	}
	return nil
}

func gnuPGIdentity() (string, error) {
	gnukey := `gpg --with-colons -k | awk -F: '$1=="uid" {print $10; exit}'`
	cmd := exec.Command("bash", "-c", gnukey)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return ``, errors.New("unable to get gnupg identity: " + b.String())
	}
	return strings.ReplaceAll(b.String(), "\n", ""), nil
}
