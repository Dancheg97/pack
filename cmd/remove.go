// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file corresponding a single command, including root cmd.

import (
	"os"
	"strings"

	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/system"
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
	groups := pack.Split(pkgs)
	restPacmanPkgs := GetPacmanPackagesFromPackNames(groups.PackPackages)
	groups.PacmanPackages = append(groups.PacmanPackages, restPacmanPkgs...)
	RemovePacmanPackages(groups.PacmanPackages)
	for _, pkg := range groups.PacmanPackages {
		pack.Remove(pkg, pack.PACMAN)
	}
}

// Try to remove all packages at once.
func RemovePacmanPackages(pkgs []string) {
	pkgsStr := strings.Join(pkgs, " ")
	o, err := system.Callf("sudo pacman --noconfirm -R %s", pkgsStr)
	if err != nil {
		PrintNotFoundPackages(o)
		os.Exit(1)
	}
	prnt.Yellow("Packages removed: ", pkgsStr)
}

// Get pacman packages from parsed removal command.
func PrintNotFoundPackages(o string) {
	o = strings.ReplaceAll(o, "\n", " ")
	o = strings.ReplaceAll(o, `error: target not found: `, "")
	prnt.Red("Packages not found: ", o)
}

// Get pacman packages related to pack names.
func GetPacmanPackagesFromPackNames(pkgs []string) []string {
	var out []string
	for _, pkg := range pkgs {
		pkgInfo, err := pack.Get(pkg, pack.PACK)
		if err != nil {
			prnt.Red("Unable to find package: ", pkg)
			os.Exit(1)
		}
		out = append(out, pkgInfo.PacmanName)
	}
	return out
}
