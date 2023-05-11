// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/database"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
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
	for _, pkg := range groups.PacmanPackages {
		info, err := database.Get(pkg, database.PACMAN)
		if err == nil {
			groups.PackPackages = append(groups.PackPackages, info.PackName)
			continue
		}
		descr := GetPacmanDescription(pkg)
		fmt.Println(descr)
	}
	for _, pkg := range groups.PackPackages {
		info, err := database.Get(pkg, database.PACK)
		if err != nil {
			print.Yellow("Package not found: ", pkg)
			continue
		}
		descr := GetPacmanDescription(info.PacmanName)
		fmt.Println(AppendPackParams(descr, info))
	}
}

// Get pacman package description.
func GetPacmanDescription(pkg string) string {
	info, err := system.Call("pacman -Qi " + pkg)
	if err != nil {
		print.Red("Error: ", strings.ReplaceAll(info, "error: ", ""))
		os.Exit(1)
	}
	return info
}

// Append pack information to package description.
func AppendPackParams(info string, p *database.Package) string {
	const (
		ver    = "Version         : "
		branch = "DefaultBranch   : "
	)
	splt1 := strings.Split(info, ver)
	splt2 := strings.Split(splt1[1], "\n")
	rest := strings.Join(splt2[1:], "\n")
	return splt1[0] + ver + p.Version + "\n" + branch + p.Branch + "\n" + rest
}
