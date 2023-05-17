// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmt"

	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"desc", "d"},
	Short:   tmpl.DescribeShort,
	Long:    tmpl.DescribeLong,
	Run:     Describe,
}

// Cli command giving package description.
func Describe(cmd *cobra.Command, pkgs []string) {
	for _, pkg := range pkgs {
		pm, pk := DescribePackage(pkg)
		if pm == nil {
			prnt.Yellow("--------------------------------\nNot found: ", pkg)
			continue
		}
		PrintDescription(pm, pk)
	}
}

// Get pacman description, and append pack information if package has it.
func DescribePackage(pkg string) (*pacman.Package, *pack.Package) {
	pm, err := pacman.Describe(pkg)
	if err != nil {
		return nil, nil
	}
	pk, err := pack.GetByPacmanName(pkg)
	if err != nil {
		pk = &pack.Package{
			PackName:      "None",
			Version:       "None",
			DefaultBranch: "None",
		}
	}
	return pm, pk
}

// Print package description.
func PrintDescription(pm *pacman.Package, pk *pack.Package) {
	fmt.Printf(
		tmpl.PrettyDesc,
		pm.Name,
		pm.Version,
		pm.Description,
		pm.Size,
		pm.Url,
		pm.BuildDate,
		pk.PackName,
		pk.Version,
		pk.DefaultBranch,
		pm.DependsOn,
		pm.RequiredBy,
	)
}
