// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmt"

	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/tmpl"
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
		fd, err := pacman.Describe(pkg)
		if err != nil {
			prnt.Yellow("--------------------------------\nNot found: ", pkg)
			continue
		}
		sd, err := pack.GetByPacmanName(pkg)
		if err != nil {
			sd = &pack.Package{
				PackName:      "None",
				Version:       "None",
				DefaultBranch: "None",
			}
		}
		PrintDescription(fd, sd)
	}
}

// Print package description.
func PrintDescription(r *pacman.Package, o *pack.Package) {
	fmt.Printf(
		tmpl.PrettyDesc,
		r.Name,
		r.Version,
		r.Description,
		r.Size,
		r.Url,
		r.BuildDate,
		o.PackName,
		o.Version,
		o.DefaultBranch,
		r.DependsOn,
		r.RequiredBy,
	)
}
