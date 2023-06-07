// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"fmnx.su/core/pack/tmpl"
	"github.com/fatih/color"
)

// Parameters that can be used to build packages.
type BuildParameters struct {
	// Do not ask for any confirmation on build/installation.
	Quick bool
	// Directory where resulting package and signature will be moved.
	Dir string
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

func Build(prms ...BuildParameters) error {
	_ = formOptions(prms, builddefault)

	return CheckGnupg()
}

// // Build package with pack.
// func Build(args []string) error {
// 	// CheckErr(CheckGnupg())
// 	// CheckErr(pacman.ValidatePackager())
// 	// CheckErr(pacman.Makepkg())
// 	// CheckErr(exec.Command(
// 	// 	"bash", "-c",
// 	// 	"sudo mv *.pkg.tar.zst* /var/cache/pacman/pkg",
// 	// ).Run())
// 	return nil
// }

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func CheckGnupg() error {
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
	fmt.Printf(
		"%s unable to find GnuPG private keys, required for package signing\n",
		color.RedString("error:"),
	)
	fmt.Println(tmpl.Gnupgerr)
	return errors.New("GnuPG private keys are missing")
}

// // Validate, that packager defined in /etc/makepkg.conf matches signer
// // authority in GnuPG.
// func ValidatePackager() error {
// 	keySigner, err := GetGnupgIdentity()
// 	if err != nil {
// 		return err
// 	}
// 	f, err := os.ReadFile("/etc/makepkg.conf")
// 	if err != nil {
// 		return err
// 	}
// 	splt := strings.Split(string(f), "\nPACKAGER=\"")
// 	if len(splt) != 2 {
// 		return errors.New(
// 			"packager is not defined in /etc/makepkg.conf. " +
// 				"Add PACKAGER variable matching your GnuPG authority " +
// 				"in /etc/makepkg.conf\n" +
// 				"Example: PACKAGER=\"John Doe <john@doe.com>\"",
// 		)
// 	}
// 	confPackager := strings.Split(splt[1], "\"\n")[0]
// 	if confPackager != keySigner {
// 		return fmt.Errorf(
// 			"gnu key signer should match makepkg packager: %s / %s",
// 			keySigner, confPackager,
// 		)
// 	}
// 	return nil
// }

func GetGnupgIdentity() (string, error) {
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
