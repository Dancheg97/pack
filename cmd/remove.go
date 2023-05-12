// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"os"

	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   tmpl.RemoveShort,
	Long:    tmpl.RemoveLong,
	Run:     Remove,
}

// Cli command removing packages from system.
func Remove(cmd *cobra.Command, pkgs []string) {
	groups := GroupPackages(pkgs)
	restPacmanPkgs := GetPacmanPackagesFromPackNames(groups.PackPackages)
	groups.PacmanPackages = append(groups.PacmanPackages, restPacmanPkgs...)
	err := pacman.Remove(groups.PacmanPackages)
	CheckErr(err)
	pack.Remove(groups.PackPackages)
}

// Get pacman packages related to pack names.
func GetPacmanPackagesFromPackNames(pkgs []string) []string {
	var out []string
	for _, pkg := range pkgs {
		pkgInfo, err := pack.Get(pkg)
		if err != nil {
			prnt.Red("Unable to find package: ", pkg)
			os.Exit(1)
		}
		out = append(out, pkgInfo.PacmanName)
	}
	return out
}
