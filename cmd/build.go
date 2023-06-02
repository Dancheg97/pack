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

const gnupgerr = `Before installation you should create GPG key.
It is required for package signing, run:

pack i gnupg
gpg --gen-key

Complete documentation can be found here:
https://wiki.archlinux.org/title/DeveloperWiki:Signing_Packages`

func Build(cmd *cobra.Command, args []string) {
	hd, err := os.UserHomeDir()
	CheckErr(err)
	_, err = os.Stat(path.Join(hd, ".gnupgx"))
	if err != nil {
		fmt.Println(gnupgerr)
		os.Exit(1)
	}

	err = pacman.Makepkg()
	CheckErr(err)
	mv := "sudo mv *.pkg.tar.zst* /var/cache/pacman/pkg"
	c := exec.Command("bash", "-c", mv)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	CheckErr(c.Run())
}
