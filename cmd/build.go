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
package and signature in pacman cache directory.
`,
	Run: Build,
}

// Build package with pack.
func Build(cmd *cobra.Command, args []string) {
	CheckErr(CheckGnupg())
	CheckErr(pacman.ValidatePackager())
	CheckErr(pacman.Makepkg())
	CheckErr(ValideSignature(""))
	CheckErr(CacheBuiltPackage(""))
}

const gnupgerr = `GPG key is not found in user directory ~/.gnupg
It is required for package signing, run:

pack i gnupg
gpg --gen-key

Complete documentation can be found here:
https://wiki.archlinux.org/title/DeveloperWiki:Signing_Packages`

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

// Validates all file signatures in provided directory.
func ValideSignature(dir string) error {
	sigloc := path.Join(dir, "*.sig")
	command := "gpg --keyserver-options auto-key-retrieve --verify " + sigloc
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Puts all packages and signatures from provided dir to pacakge cachegpg --with-colons -k | awk -F: '$1=="uid" {print $10; exit}'.
func CacheBuiltPackage(dir string) error {
	command := "sudo mv *.pkg.tar.zst* " + pacmancache
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
