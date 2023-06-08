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

	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/tmpl"
)

// Syncronize packages with pack.
type SyncParameters struct {
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Download fresh package databases from the server (-yy force)
	Refresh []bool
	// Upgrade installed packages (-uu enables downgrade)
	Upgrade []bool
	// View package information (-ii for extended information)
	Info []bool
	// View a list of packages in a repo
	List []bool
	// Don't ask for any confirmation (--noconfirm)
	Quick bool
	// Use relaxed timeouts for download
	Notimeout bool
	// Reinstall up to date targets
	Force bool
	// Do not save new registries in pacman.conf
	Keepcfg bool
}

func syncdefault() *SyncParameters {
	return &SyncParameters{
		Quick:   true,
		Refresh: []bool{true},
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}
}

// Syncronize provided packages with provided parameters.
func Sync(args []string, prms ...SyncParameters) error {
	p := formOptions(prms, syncdefault)

	pkgs, fmtpkgs, err := formatpkgs(args)
	if err != nil {
		return err
	}

	initial, err := prepareconf(pkgs, p.Stdout)
	if err != nil {
		return err
	}

	err = pacman.SyncList(fmtpkgs, pacman.SyncOptions{
		Sudo:      true,
		Needed:    p.Force,
		NoConfirm: p.Quick,
		Refresh:   p.Refresh,
		Upgrade:   p.Upgrade,
		List:      p.List,
		Stdout:    p.Stdout,
		Stderr:    p.Stderr,
		Stdin:     p.Stdin,
	})
	if err != nil || p.Keepcfg {
		rollbackconf(*initial)
		return err
	}
	return nil
}

// Pakcage with owner and registry for further pack operations.
type registrypkg struct {
	Registry string
	Owner    string
	Name     string
}

// Format packages to pack compatible formats for operations with registries.
func formatpkgs(pkgs []string) ([]registrypkg, []string, error) {
	var rez []registrypkg
	var fmtpkgs []string
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		switch len(splt) {
		case 1:
			rez = append(rez, registrypkg{
				Name: splt[0],
			})
			fmtpkgs = append(fmtpkgs, pkg)
		case 2:
			rez = append(rez, registrypkg{
				Registry: splt[0],
				Name:     splt[1],
			})
			fmtpkgs = append(fmtpkgs, pkg)
		case 3:
			rez = append(rez, registrypkg{
				Registry: splt[0],
				Owner:    splt[1],
				Name:     splt[2],
			})
			fmtpkgs = append(fmtpkgs, splt[0]+"/"+splt[2])
		default:
			return nil, nil, errors.New(tmpl.Err + " broken package: " + pkg)
		}
	}
	return rez, fmtpkgs, nil
}

// Add missing registries to pacman configuration file and return file before
// modifications.
func prepareconf(pkgs []registrypkg, ow io.Writer) (*string, error) {
	f, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		return nil, err
	}
	conf := string(f)

	for _, pkg := range pkgs {
		switch {
		case pkg.Registry != "" && pkg.Owner != "":
			if !checkexistowner(conf, pkg.Registry, pkg.Owner) {
				err = addconfdb(pkg, ow)
				if err != nil {
					return nil, err
				}
			}
		case pkg.Registry != "":
			if !checkexistsroot(conf, pkg.Registry) {
				err = addconfdb(pkg, ow)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &conf, nil
}

func checkexistowner(conf string, registry string, owner string) bool {
	return strings.Contains(conf, "\n["+registry+"."+owner+"]\n")
}

func checkexistsroot(conf string, registry string) bool {
	return strings.Contains(conf, "\n["+registry+"]\n")
}

func addconfdb(pkg registrypkg, ow io.Writer) error {
	var t string
	if pkg.Owner == "" {
		t = fmt.Sprintf(confroot, pkg.Registry, pkg.Registry)
	} else {
		t = fmt.Sprintf(confuser, pkg.Registry, pkg.Owner, pkg.Registry)
	}
	command := "cat <<EOF >> /etc/pacman.conf" + t + "EOF"
	err := exec.Command("sudo", "bash", "-c", command).Run()
	if err != nil {
		return errors.New(tmpl.Err + " unable to add to pacman.conf: " + t)
	}
	ow.Write([]byte(tmpl.Dots + tmpl.DbAdded + pkg.Registry + "\n")) //nolint
	return nil
}

func rollbackconf(s string) {
	exec.Command( //nolint
		"sudo", "bash", "-c",
		"cat <<EOF > /etc/pacman.conf\n"+s+"EOF",
	).Run() //nolint
}

const confroot = `
[%s]
Server = https://%s/api/pack
`

const confuser = `
[%s.%s]
Server = https://%s/api/pack
`
