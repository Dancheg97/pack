// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"fmt"
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
	// Generate pacakge template and exit.
	Template bool
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

	if p.Template {
		return template()
	}

	tmpl.Amsg(p.Stdout, "Building package")

	tmpl.Smsg(p.Stdout, "Running GnuPG check", 1, 4)
	err := checkGnupg()
	if err != nil {
		return err
	}

	tmpl.Smsg(p.Stdout, "Validating packager identity", 2, 3)
	err = validatePackager()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	tmpl.Smsg(p.Stdout, "Calling makepkg", 3, 3)
	err = pacman.Makepkg(pacman.MakepkgOptions{
		Sign:       true,
		Stdout:     p.Stdout,
		Stderr:     p.Stdout,
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
		return errors.Join(err, errors.New(b.String()))
	}

	tmpl.Amsg(p.Stdout, "Moving package to cache")
	b.Reset()
	cmd := exec.Command("bash", "-c", "sudo mv *.pkg.tar.zst* "+p.Dir)
	cmd.Stderr = &b
	err = cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New(b.String()))
	}
	return nil
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
		o := b.String()
		return ``, errors.New("unable to get gnupg identity: " + o)
	}
	return strings.ReplaceAll(b.String(), "\n", ""), nil
}

func template() error {
	ident, err := gnuPGIdentity()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	splt := strings.Split(dir, "/")
	n := splt[len(splt)-1]

	d := fmt.Sprintf(tmpl.Desktop, n, n, n, n, n)
	derr := os.WriteFile(n+".desktop", []byte(d), 0600)

	s := fmt.Sprintf(tmpl.ShFile, n, n)
	serr := os.WriteFile(n+".sh", []byte(s), 0600)

	p := fmt.Sprintf(tmpl.PKGBUILD, ident, n, n, n, n, n, n, n, n)
	perr := os.WriteFile(`PKGBUILD`, []byte(p), 0600)

	return errors.Join(derr, serr, perr)
}
