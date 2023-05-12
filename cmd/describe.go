// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"os"
	"strings"

	"fmnx.io/core/pack/packdb"
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
	Aliases: []string{"descr", "d"},
	Short:   tmpl.DescribeShort,
	Long:    tmpl.DescribeLong,
	Run:     Describe,
}

// Cli command giving package description.
func Describe(cmd *cobra.Command, pkgs []string) {
	groups := SplitPackages(pkgs)
	var notfound []string
	for _, pkg := range groups.PackPackages {
		i, err := packdb.Get(pkg, packdb.PACK)
		if err != nil {
			notfound = append(notfound, pkg)
			continue
		}
		groups.PacmanPackages = append(groups.PacmanPackages, i.PacmanName)
	}
	var desclist []packdb.Description
	for _, pkg := range groups.PacmanPackages {
		d, err := pacman.Describe(pkg)
		if err != nil {
			notfound = append(notfound, pkg)
			continue
		}
		fd := packdb.DescribeAppend(d)
		desclist = append(desclist, fd)
	}
	if len(notfound) > 0 {
		prnt.Red("unable to find packages: ", strings.Join(notfound, " "))
		os.Exit(1)
	}

}
