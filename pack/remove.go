// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"io"
	"os"
	"strings"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

type RemoveParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Ask for confirmation when deleting package.
	Confirm bool
	// Leave package dependencies in the system (removed by default).
	Norecursive bool
	// Leave package configs in the system (removed by default).
	Nocfgs bool
	// Remove packages and all packages that depend on them.
	Cascade bool
}

func removeDefault() *RemoveParameters {
	return &RemoveParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

func Remove(args []string, prms ...RemoveParameters) error {
	p := formOptions(prms, removeDefault)

	local, remote := splitRemoved(args)

	if len(local) > 0 {
		err := pacman.RemoveList(local, pacman.RemoveParameters{
			Sudo:        true,
			NoConfirm:   !p.Confirm,
			Recursive:   !p.Norecursive,
			WithConfigs: !p.Nocfgs,
			Cascade:     p.Cascade,
			Stdout:      p.Stdout,
			Stderr:      p.Stderr,
			Stdin:       p.Stdin,
		})
		if err != nil {
			return err
		}
	}

	if len(remote) > 0 {
		msgs.Amsg(p.Stdout, "Removing remote packages")
		for i, v := range remote {
			msgs.Smsg(p.Stdout, "Removing "+v, i, len(remote))
			err := rmRemote(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Splits packages that will be removed locally and on remote.
func splitRemoved(pkgs []string) ([]string, []string) {
	var local []string
	var remote []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "/") {
			remote = append(remote, pkg)
			continue
		}
		local = append(local, pkg)
	}
	return local, remote
}

// Function that will be used to remove remote package.
func rmRemote(pkg string) error {
	return nil
}
