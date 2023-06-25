// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"fmnx.su/core/pack/msgs"
)

// Parameters for util.
type UtilParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Generate flutter pacakge template and exit.
	Flutter bool
	// Generate go CLI pacakge template and exit.
	Gocli bool
	// Export public armored string key to Stdout and exit.
	ExportKey bool
}

func utildefault() *UtilParameters {
	return &UtilParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

func Util(args []string, prms ...UtilParameters) error {
	p := formOptions(prms, utildefault)

	if p.ExportKey {
		return armored(p.Stdout)
	}

	if p.Flutter {
		return fluttertemplate()
	}

	if p.Gocli {
		return goclitemplate()
	}

	return nil
}

// Return armored public key string from GnuPG.
func armored(o io.Writer) error {
	cmd := exec.Command("gpg", "--armor", "--export")
	cmd.Stdout = o
	return call(cmd)
}

// Function generates project template for flutter desktop application based on
// current directory name and identity in GnuPG.
func fluttertemplate() error {
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

	d := fmt.Sprintf(msgs.Desktop, n, n, n, n, n)
	derr := os.WriteFile(n+".desktop", []byte(d), 0600)

	s := fmt.Sprintf(msgs.ShFile, n, n)
	serr := os.WriteFile(n+".sh", []byte(s), 0600)

	p := fmt.Sprintf(msgs.PKGBUILDflutter, ident, n, n, n, n, n, n, n, n)
	perr := os.WriteFile(`PKGBUILD`, []byte(p), 0600)

	return errors.Join(derr, serr, perr)
}

// Function generates project template for go cli utility based on
// current directory name and identity in GnuPG.
func goclitemplate() error {
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

	p := fmt.Sprintf(msgs.PKGBUILDgocli, ident, n, n)
	return os.WriteFile(`PKGBUILD`, []byte(p), 0600)
}
