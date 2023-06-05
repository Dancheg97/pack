// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Options for building packages.
type MakepkgOptions struct {
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Directory where process will be executed.
	Dir string
	// Specify a key to use for gpg signing instead of the default. [--key <key>]
	GpgKey string
	// Use an alternate build script (not 'PKGBUILD'). [-p <file>]
	File string
	// Ignore incomplete arch field in PKGBUILD. [--ignorearch]
	IgnoreEach bool
	// Clean up work files after build. [--clean]
	Clean bool
	// Remove $srcdir/ dir before building the package. [--cleanbuild]
	CleanBuild bool
	// Skip all dependency checks. [--nodeps]
	NoDeps bool
	// Do not extract source files (use existing $srcdir/ dir). [--noextract]
	NoExtract bool
	// Overwrite existing package. [--force]
	Force bool
	// Generate integrity checks for source files. [--geninteg]
	Geinteg bool
	// Install package after successful build. [--install]
	Install bool
	// Log package build process. [--log]
	Log bool
	// Disable colorized output messages. [--nocolor]
	NoColor bool
	// Download and extract files only. [--nobuild]
	NpBuild bool
	// Remove installed dependencies after a successful build. [--rmdeps]
	RmDeps bool
	// Repackage contents of the package without rebuilding. [--repackage]
	Repackage bool
	// Install missing dependencies with pacman. [--syncdeps]
	SyncDeps bool
	// Use an alternate config file (not '/etc/makepkg.conf'). [--config <file>]
	Config string
	// Do not update VCS sources. [--holdver]
	HoldVer bool
	// Do not create package archive. [--noarchive]
	NoArchive bool
	// Do not run the check() function in the PKGBUILD. [--nocheck]
	NoCheck bool
	// Do not run the prepare() function in the PKGBUILD. [--noprepare]
	NoPrepare bool
	// Do not create a signature for the package. [--nosign]
	NoSign bool
	// Sign the resulting package with gpg. [--sign]
	Sign bool
	// Do not verify checksums of the source files. [--skipchecksums]
	SkipCheckSums bool
	// Do not perform any verification checks on source files. [--skipinteg]
	SkipIntegrityChecks bool
	// Do not verify source files with PGP signatures. [--skippgpcheck]
	SkipPgpCheck bool
	// Do not reinstall up to date packages. [--needed]
	Needed bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Do not show a progress bar when downloading files. [--noprogressbar]
	NoProgressBar bool
	// Install packages as non-explicitly installed. [--asdeps]
	AsDeps bool
}

func MakepkgDefault() *MakepkgOptions {
	return &MakepkgOptions{
		Clean:      true,
		Force:      true,
		Needed:     true,
		Sign:       true,
		CleanBuild: true,
		Install:    true,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Stdin:      os.Stdin,
	}
}

// This command will build a package in directory provided in options.
// Function is safe for concurrent usage. Can be called from multiple
// goruotines, when options Install or SyncDeps are false.
func Makepkg(opts ...MakepkgOptions) error {
	o := formOptions(opts, MakepkgDefault)

	var args []string
	if o.IgnoreEach {
		args = append(args, "--ignorearch")
	}
	if o.Clean {
		args = append(args, "--clean")
	}
	if o.CleanBuild {
		args = append(args, "--cleanbuild")
	}
	if o.NoDeps {
		args = append(args, "--nodeps")
	}
	if o.NoExtract {
		args = append(args, "--noextract")
	}
	if o.Force {
		args = append(args, "--force")
	}
	if o.Geinteg {
		args = append(args, "--geninteg")
	}
	if o.Log {
		args = append(args, "--log")
	}
	if o.NoColor {
		args = append(args, "--nocolor")
	}
	if o.NpBuild {
		args = append(args, "--nobuild")
	}
	if o.RmDeps {
		args = append(args, "--rmdeps")
	}
	if o.Repackage {
		args = append(args, "--repackage")
	}
	if o.HoldVer {
		args = append(args, "--holdver")
	}
	if o.NoArchive {
		args = append(args, "--noarchive")
	}
	if o.NoCheck {
		args = append(args, "--nocheck")
	}
	if o.NoPrepare {
		args = append(args, "--noprepare")
	}
	if o.NoSign {
		args = append(args, "--nosign")
	}
	if o.Sign {
		args = append(args, "--sign")
	}
	if o.SkipCheckSums {
		args = append(args, "--skipchecksums")
	}
	if o.SkipIntegrityChecks {
		args = append(args, "--skipinteg")
	}
	if o.SkipPgpCheck {
		args = append(args, "--skippgpcheck")
	}
	if o.Needed {
		args = append(args, "--needed")
	}
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if o.AsDeps {
		args = append(args, "--asdeps")
	}
	if o.File != `` {
		args = append(args, "-p")
		args = append(args, o.File)
	}
	if o.Config != "" {
		args = append(args, "--config")
		args = append(args, o.Config)
	}
	if o.GpgKey != "" {
		args = append(args, "--key")
		args = append(args, o.GpgKey)
	}
	if o.Install {
		args = append(args, "--install")
		mu.Lock()
		defer mu.Unlock()
	}
	if o.SyncDeps {
		args = append(args, "--syncdeps")
		if mu.TryLock() {
			defer mu.Unlock()
		}
	}
	args = append(args, o.AdditionalParams...)

	cmd := exec.Command(makepkg, args...)
	cmd.Dir = o.Dir
	cmd.Stdin = o.Stdin
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr

	return cmd.Run()
}

// Get parameters from a shell file (might be usefull to resolve dependencies
// before package build/installation process).
func GetShellParams(file string, arg string) ([]string, error) {
	const tmpl = "source %s; for i in ${%s[@]}; do echo $i;done"
	command := fmt.Sprintf(tmpl, file, arg)

	var b bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		fmt.Println(b.String())
		return nil, err
	}
	return strings.Split(strings.Trim(b.String(), "\n"), "\n"), nil
}

// Validate, that packager defined in /etc/makepkg.conf matches signer
// authority in GnuPG.
func ValidatePackager() error {
	gnukey := `gpg --with-colons -k | awk -F: '$1=="uid" {print $10; exit}'`
	cmd := exec.Command("bash", "-c", gnukey)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return errors.New(b.String())
	}
	f, err := os.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return err
	}
	splt := strings.Split(string(f), "\nPACKAGER=\"")
	if len(splt) != 2 {
		return errors.New(
			"Packager is not defined in /etc/makepkg.conf. " +
				"Add PACKAGER variable matching your GnuPG authority " +
				"in /etc/makepkg.conf\n" +
				"Example: PACKAGER=\"John Doe <john@doe.com>\"\n",
		)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	keySigner := strings.ReplaceAll(b.String(), "\n", "")
	if confPackager != keySigner {
		return fmt.Errorf(
			"Gnu key signer should match makepkg packager: %s / %s",
			keySigner, confPackager,
		)
	}
	return nil
}

// Validates all file signatures in provided directory.
func ValideSignature(dirs string) error {
	sigloc := path.Join(dirs, "*.sig")
	command := "gpg --keyserver-options auto-key-retrieve --verify " + sigloc
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Puts all packages and signatures from provided dir to pacakge cache.
func CacheBuiltPackage(src string, dst string) error {
	command := "sudo mv " + path.Join(src, "*.pkg.tar.zst*") + " " + dst
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
