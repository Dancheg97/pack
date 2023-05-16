// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"os"
	"strings"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/git"
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/system"
	"fmnx.su/core/pack/tmpl"
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
	if len(pkgs) == 1 && pkgs[0] == `gen` {
		dir := system.Pwd()
		url, err := git.Url(dir)
		CheckErr(err)
		err = pacman.Generate(dir, url)
		CheckErr(err)
		prnt.Green("Generated file: ", "PKGBUILD")
		return
	}
	if len(pkgs) == 0 {
		dir := system.Pwd()
		BuildDirectory(dir, ``)
		return
	}
	for _, pkg := range pkgs {
		i := pack.GetPackInfo(pkg)
		err := git.Clone(i.GitUrl, i.Directory)
		CheckErr(err)
		BuildDirectory(i.Directory, i.Version)
		err = system.MvExt(i.Directory, config.PackageCacheDir, ".pkg.tar.zst")
		CheckErr(err)
	}
	prnt.Blue("Build complete, results in: ", config.PackageCacheDir)
}

// Build package in specified directory. Assumes this directory has cloned git
// repository with PKGBUILD in it.
func BuildDirectory(dir string, version string) string {
	pkgname := ValidateBuildDir(dir)
	if version == `` {
		branch, err := git.DefaultBranch(dir)
		CheckErr(err)
		version, err = git.LastCommitDir(dir, branch)
		CheckErr(err)
	}
	err := git.Checkout(dir, version)
	CheckErr(err)
	ResolveDependencies(dir, pkgname)
	err = pacman.Build(dir, pkgname)
	CheckErr(err)
	return version
}

// Validate directory to be valid pack package - git repo name matching package
// name defined in PKGBUILD.
func ValidateBuildDir(dir string) string {
	url, err := git.Url(dir)
	CheckErr(err)
	pkgname, err := pacman.PkgbuildParam(dir, "pkgname")
	CheckErr(err)
	splt := strings.Split(url, "/")
	if pkgname != splt[len(splt)-1] {
		prnt.Red("package name is not matching git link, can't build: ", dir)
		os.Exit(1)
	}
	return pkgname
}

// Resolve pack dependencies for package in provided directory.
func ResolveDependencies(dir string, pkg string) {
	deps, err := pacman.GetDeps(dir)
	CheckErr(err)
	groups := GroupPackages(deps)
	nf := pack.GetUninstalled(groups.PackPackages)
	if len(nf) > 0 {
		prnt.Blue("Resolving pack deps for "+pkg+": ", strings.Join(nf, " "))
		Install(nil, nf)
	}
	err = pack.SwapDeps(dir+"/PKGBUILD", groups.PackPackages)
	CheckErr(err)
	deps, err = pacman.GetDeps(dir)
	CheckErr(err)
	prnt.Blue("Resolving pacman deps for "+pkg+": ", strings.Join(deps, " "))
	err = pacman.Install(deps)
	CheckErr(err)
}
