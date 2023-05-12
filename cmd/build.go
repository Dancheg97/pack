// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/prnt"
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
	dir := system.Pwd()
	prnt.Blue("Preparing package: ", dir)
	out, err := system.Call("makepkg -sfi --noconfirm")
	if err != nil {
		prnt.Red("Unable to execute: ", "makepkg")
		fmt.Println(out)
		os.Exit(1)
	}
	i := GetInstallLink()
	SavePackageInfo(i)
	prnt.Green("Package prepared and installed: ", i.FullName)
}

// Save information about installed package.
func SavePackageInfo(i RepositoryInfo) {
	dir := GetCurrDir()
	branch, err := git.DefaultBranch(dir)
	CheckErr(err)
	version, err := git.LastCommitDir(dir, branch)
	CheckErr(err)
	pack.Update(pack.Package{
		PacmanName: i.ShortName,
		PackName:   i.FullName,
		Version:    version,
		Branch:     branch,
	})
}
