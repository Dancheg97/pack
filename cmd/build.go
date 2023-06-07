// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"fmnx.su/core/pack/pacman"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "🧰 build package",
	Long: `🧰 build package
	
This command will build, install and sign package in current directory, then
store the resulting package and signature in pacman cache directory.`,
	Run: Build,
}

// Build package with pack.
func Build(cmd *cobra.Command, args []string) {
	CheckErr(CheckGnupg())
	CheckErr(pacman.ValidatePackager())
	CheckErr(pacman.Makepkg())
	CheckErr(exec.Command(
		"bash", "-c",
		"sudo mv *.pkg.tar.zst* /var/cache/pacman/pkg",
	).Run())
}

const gnupgerr = `GPG key is not found in user directory ~/.gnupg
It is required for package signing, run:

1) Install gnupg:
pack i gnupg

2) Generate a key:
gpg --gen-key

3) Get KEY-ID, paste it to next command:
gpg -k

4) Send it to key server:
gpg --send-keys KEY-ID`

// Ensure, that user have created gnupg keys for package signing before package
// is built and cached.
func CheckGnupg() error {
	hd, err := os.UserHomeDir()
	CheckErr(err)
	_, err = os.Stat(path.Join(hd, ".gnupg"))
	if err != nil {
		fmt.Println(gnupgerr)
	}
	return err
}
