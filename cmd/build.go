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
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "🛠️ build package",
	Long: `🛠️ build package
	
This command will build package in current directory and store the resulting
package and signature in /var/cache/pacman/pkg.
`,
	Run: Build,
}

// Build package with pack.
func Build(cmd *cobra.Command, args []string) {
	CheckGnupg()
	CheckErr(pacman.Makepkg())
	ValideSignature()
	CacheBuiltPackage()
}

const gnupgerr = `Before installation you should create GPG key.
It is required for package signing, run:

pack i gnupg
gpg --gen-key

Complete documentation can be found here:
https://wiki.archlinux.org/title/DeveloperWiki:Signing_Packages`

func CheckGnupg() {
	hd, err := os.UserHomeDir()
	CheckErr(err)
	_, err = os.Stat(path.Join(hd, ".gnupg"))
	if err != nil {
		fmt.Println(gnupgerr)
		os.Exit(1)
	}
}

func ValideSignature() {
	command := "gpg --keyserver-options auto-key-retrieve --verify *.sig"
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	CheckErr(cmd.Run())
}

func CacheBuiltPackage() {
	command := "sudo mv *.pkg.tar.zst* /var/cache/pacman/pkg"
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	CheckErr(cmd.Run())
}
