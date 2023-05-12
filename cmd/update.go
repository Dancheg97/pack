// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/print"
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
	UpdatePacmanPackages(groups.PacmanPackages)
	Install(nil, groups.PackPackages)
}

// Perform full pacman update.
func FullPacmanUpdate() {
	o, err := system.Call("sudo pacman --noconfirm -Syu")
	if err != nil {
		print.Red("Unable to update pacman packages: ", o)
		os.Exit(1)
	}
	print.Green("Pacman update: ", "done")
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
	print.Green("Pack update: ", "done")
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
		print.Red("Unable to find: ", strings.Join(nfpkgs, " "))
		os.Exit(1)
	}
}

// Verify pack packages are installed in system.
func VerifyPackPackages(pkgs []string) {
	var nfpkgs []string
	for _, pkg := range pkgs {
		_, err := packdb.Get(pkg, packdb.PACK)
		if err != nil {
			nfpkgs = append(nfpkgs, pkg)
		}
	}
	if len(nfpkgs) > 0 {
		print.Red("Unable to find: ", strings.Join(nfpkgs, " "))
		os.Exit(1)
	}
}

// Update pacman packages.
func UpdatePacmanPackages(pkgs []string) {
	if len(pkgs) == 0 {
		return
	}
	joined := strings.Join(pkgs, " ")
	o, err := system.Callf("sudo pacman --noconfirm -S %s", joined)
	if err != nil {
		print.Red("Unable to update packages: %s", err.Error())
		fmt.Println(o)
		os.Exit(1)
	}
}
