// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"os"
	"os/exec"

	"fmnx.su/core/pack/pacman"
)

// This command builds package in current directory and stores it in pacman
// cache directory.
func Build() error {
	err := pacman.Makepkg(pacman.MakepkgOptions{
		Clean:     true,
		Force:     true,
		Log:       true,
		HoldVer:   true,
		Needed:    true,
		NoConfirm: true,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
	})
	if err != nil {
		return err
	}
	mvcmd := "sudo mv *.pkg.tar.zst /var/cache/pacman/pkg"
	cmd := exec.Command("bash", "-c", mvcmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
