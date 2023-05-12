// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   tmpl.PackageShort,
	Long:    tmpl.PackageLong,
	Run:     Package,
}

// Cli command preparing package in current directory.
func Package(cmd *cobra.Command, pkgs []string) {
	dir := GetCurrDir()
	print.Blue("Preparing package: ", dir)
	out, err := system.Call("makepkg -sfi --noconfirm")
	if err != nil {
		print.Red("Unable to execute: ", "makepkg")
		fmt.Println(out)
		os.Exit(1)
	}
	i := GetInstallLink()
	SavePackageInfo(i)
	print.Green("Package prepared and installed: ", i.FullName)
}

// Save information about installed package.
func SavePackageInfo(i RepositoryInfo) {
	dir := GetCurrDir()
	branch, err := git.DefaultBranch(dir)
	CheckErr(err)
	version, err := git.LastCommitDir(dir, branch)
	CheckErr(err)
	packdb.Update(packdb.Package{
		PacmanName: i.ShortName,
		PackName:   i.FullName,
		Version:    version,
		Branch:     branch,
	})
}

// Get directory for current process.
func GetCurrDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to find curr dir")
		os.Exit(1)
	}
	return dir
}
