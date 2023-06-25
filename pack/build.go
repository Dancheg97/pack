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

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

// Parameters that can be used to build packages.
type BuildParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

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
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
		Dir:       "/var/cache/pacman/pkg",
		Syncbuild: true,
		Rmdeps:    true,
	}
}

// Build package in current directory with provided arguements
func Build(prms ...BuildParameters) error {
	p := formOptions(prms, builddefault)

	msgs.Amsg(p.Stdout, "Building package")

	msgs.Smsg(p.Stdout, "Running GnuPG check", 1, 2)
	err := checkGnupg()
	if err != nil {
		return err
	}

	msgs.Smsg(p.Stdout, "Validating packager identity", 2, 2)
	err = validatePackager()
	if err != nil {
		return err
	}

	msgs.Amsg(p.Stdout, "Building package with makepkg")
	err = pacman.Makepkg(pacman.MakepkgParameters{
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
		return errors.Join(err)
	}

	msgs.Amsg(p.Stdout, "Moving package to cache")
	cmd := exec.Command("bash", "-c", "sudo mv *.pkg.tar.zst* "+p.Dir)
	return call(cmd)
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
		return errors.New(msgs.ErrGnuPGprivkeyNotFound)
	}
	for _, de := range gpgdir {
		if strings.Contains(de.Name(), "private-keys") {
			return nil
		}
	}
	return errors.New(msgs.ErrGnuPGprivkeyNotFound)
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
		return errors.New(msgs.ErrNoPackager)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	if confPackager != keySigner {
		return errors.New(msgs.ErrSignerMissmatch)
	}
	return nil
}

// Returns name and email from GnuPG. Error, if did not succeed.
func gnuPGIdentity() (string, error) {
	gnukey := `gpg --with-colons -k | awk -F: '$1=="uid" {print $10; exit}'`
	cmd := exec.Command("bash", "-c", gnukey)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		o := b.String()
		return ``, errors.New("unable to get gnupg identity: " + o)
	}
	return strings.ReplaceAll(b.String(), "\n", ""), nil
}
