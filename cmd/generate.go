// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   tmpl.GenerateShort,
	Long:    tmpl.GenerateLong,
	Run:     Generate,
}

// Cli command modifying files in workdir for pack compatability.
func Generate(cmd *cobra.Command, args []string) {
	dir := system.Pwd()
	url, err := git.Url(dir)
	CheckErr(err)
	err = pacman.Generate(dir, url)
	CheckErr(err)
	prnt.Green("Updated files: ", "PKGBUILD README.md")
}
