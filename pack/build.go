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
	// Generate pacakge template and exit.
	Template bool
	// Export public armored string key to Stdout and exit.
	ExportKey bool
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

	if p.Template {
		return template()
	}

	if p.ExportKey {
		return armored(p.Stdout)
	}

	tmpl.Amsg(p.Stdout, "Building package")

	tmpl.Smsg(p.Stdout, "Running GnuPG check", 1, 2)
	err := checkGnupg()
	if err != nil {
		return err
	}

	tmpl.Smsg(p.Stdout, "Validating packager identity", 2, 2)
	err = validatePackager()
	if err != nil {
		return err
	}

	tmpl.Amsg(p.Stdout, "Building package with makepkg")
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

	tmpl.Amsg(p.Stdout, "Moving package to cache")
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

func armored(o io.Writer) error {
	cmd := exec.Command("gpg", "--armor", "--export")
	cmd.Stdout = o
	return call(cmd)
}
