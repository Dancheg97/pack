// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmnx.su/core/pack/git"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/system"
	"fmnx.su/core/pack/tmpl"
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
	prnt.Green("Generated file: ", "PKGBUILD")
}
