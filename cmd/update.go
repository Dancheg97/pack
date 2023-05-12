// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upd", "u"},
	Short:   tmpl.UpdateShort,
	Long:    tmpl.UpdateLong,
	Run:     Update,
}

// Cli command performing package update.
func Update(cmd *cobra.Command, pkgs []string) {
	Updating = true
	if len(pkgs) == 0 {
		FullPacmanUpdate()
		FullPackUpdate()
		return
	}
	groups := SplitPackages(pkgs)
	VerifyPacmanPackages(groups.PacmanPackages)
	VerifyPackPackages(groups.PackPackages)
	err := pacman.Update(groups.PacmanPackages)
	CheckErr(err)
	Install(nil, groups.PackPackages)
}

// Perform full pacman update.
func FullPacmanUpdate() {
	o, err := system.Call("sudo pacman --noconfirm -Syu")
	if err != nil {
		prnt.Red("Unable to update pacman packages: ", o)
		os.Exit(1)
	}
	prnt.Green("Pacman update: ", "done")
}

var Updating bool

// Perform full pack update.
func FullPackUpdate() {
	outdatedpkgs := GetPackOutdated()
	var pkgs []string
	for _, pkg := range outdatedpkgs {
		pkgs = append(pkgs, fmt.Sprintf("%s@%s", pkg.Name, pkg.NewVersion))
	}
	Install(nil, pkgs)
	prnt.Green("Pack update: ", "done")
}

// Verify pacman packages exist in system.
func VerifyPacmanPackages(pkgs []string) {
	o, err := system.Callf("pacman -Q %s", strings.Join(pkgs, " "))
	var nfpkgs []string
	if err != nil {
		for _, line := range strings.Split(strings.Trim(o, "\n"), "\n") {
			if strings.Contains(line, "was not found") {
				line = strings.ReplaceAll(line, "error: package '", "")
				line = strings.ReplaceAll(line, "' was not found", "")
				nfpkgs = append(nfpkgs, line)
			}
		}
		prnt.Red("Unable to find: ", strings.Join(nfpkgs, " "))
		os.Exit(1)
	}
}

// Verify pack packages are installed in system.
func VerifyPackPackages(pkgs []string) {
	var nfpkgs []string
	for _, pkg := range pkgs {
		_, err := pack.Get(pkg, pack.PACK)
		if err != nil {
			nfpkgs = append(nfpkgs, pkg)
		}
	}
	if len(nfpkgs) > 0 {
		prnt.Red("Unable to find: ", strings.Join(nfpkgs, " "))
		os.Exit(1)
	}
}
