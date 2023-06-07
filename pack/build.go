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
	"strings"
)

// Build package with pack.
func Build(args []string) error {
	// CheckErr(CheckGnupg())
	// CheckErr(pacman.ValidatePackager())
	// CheckErr(pacman.Makepkg())
	// CheckErr(exec.Command(
	// 	"bash", "-c",
	// 	"sudo mv *.pkg.tar.zst* /var/cache/pacman/pkg",
	// ).Run())
	return nil
}

// const gnupgerr = `GPG key is not found in user directory ~/.gnupg
// It is required for package signing, run:

// 1) Install gnupg:
// pack i gnupg

// 2) Generate a key:
// gpg --gen-key

// 3) Get KEY-ID, paste it to next command:
// gpg -k

// 4) Send it to key server:
// gpg --send-keys KEY-ID`

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func CheckGnupg() error {
	// hd, err := os.UserHomeDir()
	// CheckErr(err)
	// _, err = os.Stat(path.Join(hd, ".gnupg"))
	// if err != nil {
	// 	fmt.Println(gnupgerr)
	// }
	return nil
}

// Validate, that packager defined in /etc/makepkg.conf matches signer
// authority in GnuPG.
func ValidatePackager() error {
	keySigner, err := GetGnupgIdentity()
	if err != nil {
		return err
	}
	f, err := os.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return err
	}
	splt := strings.Split(string(f), "\nPACKAGER=\"")
	if len(splt) != 2 {
		return errors.New(
			"packager is not defined in /etc/makepkg.conf. " +
				"Add PACKAGER variable matching your GnuPG authority " +
				"in /etc/makepkg.conf\n" +
				"Example: PACKAGER=\"John Doe <john@doe.com>\"",
		)
	}
	confPackager := strings.Split(splt[1], "\"\n")[0]
	if confPackager != keySigner {
		return fmt.Errorf(
			"gnu key signer should match makepkg packager: %s / %s",
			keySigner, confPackager,
		)
	}
	return nil
}

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
