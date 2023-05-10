// Copyright 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/database"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   "📦 prepare and install package",
	Long: `📦 prepare .pkg.tar.zst in current directory and install it

This script will read prepare .pkg.tar.zst package. You can use it to test 
PKGBUILD template for project or validate installation for pack.

To double check installation, you can test it inside pack docker container:
docker run --rm -it fmnx.io/core/pack i example.com/package
`,
	Run: Package,
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
	branch := GetDefaultGitBranch(dir)
	version := GetLastCommitHash(dir, branch)
	database.Add(database.Package{
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
