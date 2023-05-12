// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"strings"

	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   tmpl.BuildShort,
	Long:    tmpl.BuildLong,
	Run:     Build,
}

// Cli command preparing package in current directory.
func Build(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) == 0 {
		BuildCurrentDirectory()
		return
	}
	// TODO add functions to build remote repos.
}

// Build package in current directory.
func BuildCurrentDirectory() {
	dir := system.Pwd()
	err := pacman.Build(dir)
	CheckErr(err)
	branch, err := git.DefaultBranch(dir)
	CheckErr(err)
	commit, err := git.CurrentCommitDir(dir)
	CheckErr(err)
	name, err := pacman.PkgbuildParam(dir+"/PKGBUILD", "pkgname")
	CheckErr(err)
	url, err := git.Url(dir)
	CheckErr(err)
	pack.Update(pack.Package{
		PackName:      strings.Replace(url, "https://", "", 1),
		PacmanName:    name,
		Version:       commit,
		DefaultBranch: branch,
	})
}
