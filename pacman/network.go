// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pacman

import (
	"errors"
	"os"
	"strings"

	"fmnx.io/core/pack/system"
)

// This command will build package in provided directory. Also installing
// missing packages with pacman.
func Build(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	o, err := system.Call("makepkg -sf --noconfirm")
	if err != nil {
		return errors.New("pacman unable to build:\n" + o)
	}
	return nil
}

// Update listed pacman packages with sync command.
func Update(pkgs []string) error {
	joined := strings.Join(pkgs, " ")
	o, err := system.Call("sudo pacman --noconfirm -Sy " + joined)
	if err != nil {
		return errors.New("pacman unable to update:\n" + o)
	}
	return nil
}
