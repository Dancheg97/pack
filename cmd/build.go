// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/git"
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
		dir := system.Pwd()
		BuildDirectory(dir, ``)
		return
	}
	for _, pkg := range pkgs {
		i := EjectInfoFromPackLink(pkg)
		err := git.Clone(i.Url, i.Directory)
		CheckErr(err)
		BuildDirectory(i.Directory, i.Version)
		err = system.MvExt(i.Directory, config.PackageCacheDir, ".pkg.tar.zst")
		CheckErr(err)
	}
}

// Build package in specified directory. Assumes this directory has cloned git
// repository with PKGBUILD in it.
func BuildDirectory(dir string, version string) {
	if version == `` {
		branch, err := git.DefaultBranch(dir)
		CheckErr(err)
		version, err = git.LastCommitDir(dir, branch)
		CheckErr(err)
	}
	err := git.Checkout(dir, version)
	CheckErr(err)
	err = pacman.Build(dir)
	CheckErr(err)
}
