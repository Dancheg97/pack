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

// Optional parameters for pacman remove command.
type RemoveOptions struct {
	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Remove with all unnecessary packages. [--recursive]
	Recursive bool
	// Remove with all explicitly installed deps. [-ss]
	ForceRecursive bool
	// Remove configuration files aswell. [--nosave]
	WithConfigs bool
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
}

var RemoveDefault = RemoveOptions{
	Sudo:        true,
	NoConfirm:   true,
	Recursive:   true,
	WithConfigs: true,
	Stdout:      os.Stdout,
	Stderr:      os.Stderr,
	Stdin:       os.Stdin,
}

// Remove packages from system.
func Remove(pkgs string, opts ...RemoveOptions) error {
	return RemoveList(strings.Split(pkgs, " "), opts...)
}

// Remove packages from system.
func RemoveList(pkgs []string, opts ...RemoveOptions) error {
	o := formOptions(opts, &RemoveDefault)

	var args []string
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.Recursive {
		args = append(args, "--recursive")
	}
	if o.ForceRecursive {
		args = append(args, "-ss")
	}
	if o.WithConfigs {
		args = append(args, "--nosave")
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, pkgs...)

	cmd := SudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	mu.Lock()
	defer mu.Unlock()
	return cmd.Run()
}
