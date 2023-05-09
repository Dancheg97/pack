// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information could be found on official web page: https://fmnx.io/
// Email for contacts: help@fmnx.io
package cmd

import (
	"os"
	"strings"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   "❌ remove packages",
	Long: `❌ remove packages

Use this command to remove packages from system. You can specify both pacman 
packages and pack links.

Example:
pack rm fmnx.io/core/ainst`,
	Run: Remove,
}

// Cli command removing packages from system.
func Remove(cmd *cobra.Command, pkgs []string) {
	mp := ReadMapping()
	groups := SplitPackages(pkgs)
	rmtargets := GetRemoveTargetPackages(groups, mp)
	TryRemove(rmtargets)
	WriteNewMapping(mp, rmtargets)
}

// Try to remove all packages at once.
func TryRemove(pkgs []string) {
	pkgsStr := strings.Join(pkgs, " ")
	o, err := system.Callf("sudo pacman --noconfirm -R %s", pkgsStr)
	if err != nil {
		PrintNotFoundPackages(o)
		os.Exit(1)
	}
	print.Yellow("Packages removed: ", pkgsStr)
}

// Get pacman packages from parsed removal command.
func PrintNotFoundPackages(o string) {
	o = strings.ReplaceAll(o, "\n", " ")
	o = strings.ReplaceAll(o, `error: target not found: `, "")
	print.Red("Packages not found: ", o)
}

// Get all packages that would be removed in pacman format.
func GetRemoveTargetPackages(groups PackageGroups, mp PackMap) []string {
	for _, pkg := range groups.PackPackages {
		groups.PacmanPackages = append(groups.PacmanPackages, mp[pkg])
	}
	return groups.PacmanPackages
}

// This function will form new package mapping and write it.
func WriteNewMapping(mp PackMap, pkgs []string) {
	revmp := ReverseMapping(mp)
	for _, pkg := range pkgs {
		delete(revmp, pkg)
	}
	WriteMapping(ReverseMapping(revmp))
}

// Returns mapping from pacman package to pack package.
func ReverseMapping(in map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range in {
		r[v] = k
	}
	return r
}
