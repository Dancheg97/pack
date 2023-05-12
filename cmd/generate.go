// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"strings"

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
	splt := strings.Split(url, "/")
	err = pacman.Generate(dir, splt[len(splt)-1], url)
	CheckErr(err)
	ModifyReadmeFile(url)
	prnt.Green("Updated files: ", "PKGBUILD README.md")
}

// Append some lines to README.md file.
func ModifyReadmeFile(link string) {
	insatllMd := fmt.Sprintf(tmpl.READMEmd, "```", link, "```")
	system.AppendToFile("README.md", insatllMd)
}
