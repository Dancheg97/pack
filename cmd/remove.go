package cmd

import (
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

func Remove(cmd *cobra.Command, pkgs []string) {
	mp := ReadMapping()
	// revmp := ReverseMapping(mp)
	groups := SplitPackages(pkgs)
	rmtargets := GetRemoveTargetPackages(groups, mp)
	TryRemove(rmtargets)
	// // TODO make so that for each target result is printed.
	// for _, pkg := range pkgs {
	// 	pacmanpkg, ok := mp[pkg]
	// 	if !ok {
	// 		_, err := system.Call("pacman -Q " + pkg)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		delete(mp, revmp[pkg])
	// 		_, err = system.Call("sudo pacman --noconfirm -R " + pkg)
	// 		CheckErr(err)
	// 		continue
	// 	}
	// 	_, err := system.Call("sudo pacman --noconfirm -R " + pacmanpkg)
	// 	if err != nil {
	// 		print.Yellow("Pack package was not found in pacman: ", pkg)
	// 	}
	// 	delete(mp, pkg)
	// }
	// WriteMapping(mp)
}

// Try to remove all packages at once.
func TryRemove(pkgs []string) {
	o, err := system.Callf("sudo pacman -R %s", strings.Join(pkgs, " "))
	if err != nil {
		PrintNotFoundPackages(o)
		return
	}
	print.Green("Packages removed: ", strings.Join(pkgs, " "))
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

// Returns mapping from pacman package to pack package.
func ReverseMapping(in map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range in {
		r[v] = k
	}
	return r
}
