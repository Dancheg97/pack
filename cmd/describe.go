// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"descr", "d"},
	Short:   tmpl.DescribeShort,
	Long:    tmpl.DescribeLong,
	Run:     Describe,
}

// Cli command giving package description.
func Describe(cmd *cobra.Command, pkgs []string) {
	groups := SplitPackages(pkgs)
	for _, pkg := range groups.PackPackages {
		i, err := packdb.Get(pkg, packdb.PACK)
		if err != nil {
			print.Red("unable to find pack package: ", pkg)
			os.Exit(1)
		}
		groups.PacmanPackages = append(groups.PacmanPackages, i.PacmanName)
	}
	for _, pkg := range groups.PacmanPackages {
		d, err := pacman.Describe(pkg)
		if err != nil {
			print.Red("unable to find pacman package: ", pkg)
			os.Exit(1)
		}
		fd := packdb.DescribeAppend(d)
		fmt.Println(fd)
	}
}
