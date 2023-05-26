// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

import (
	"io"
	"os"
	"strings"
)

// Options to apply when searching for some package.
type UpgradeOptions struct {
	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Do not reinstall up to date packages. [--needed]
	Needed bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Do not show a progress bar when downloading files. [--noprogressbar]
	NoProgressBar bool
	// Do not execute the install scriptlet if one exists. [--noscriptlet]
	NoScriptlet bool
	// Install packages as non-explicitly installed. [--asdeps]
	AsDeps bool
	// Install packages as explictly installed. [--asexplict]
	AsExplict bool
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
}

var UpgradeDefault = UpgradeOptions{
	Sudo:      true,
	Needed:    true,
	NoConfirm: true,
	Stdout:    os.Stdout,
	Stderr:    os.Stderr,
	Stdin:     os.Stdin,
}

// Install packages from files.
func Upgrade(files string, opts ...UpgradeOptions) error {
	return UpgradeList(strings.Split(files, " "), opts...)
}

// Install packages from files.
func UpgradeList(files []string, opts ...UpgradeOptions) error {
	o := formOptions(opts, &UpgradeDefault)

	args := []string{"-U"}
	if o.Needed {
		args = append(args, "--needed")
	}
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if o.NoScriptlet {
		args = append(args, "--noscriptlet")
	}
	if o.AsDeps {
		args = append(args, "--asdeps")
	}
	if o.AsExplict {
		args = append(args, "--asexplict")
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, files...)

	cmd := SudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	mu.Lock()
	defer mu.Unlock()
	return cmd.Run()
}
