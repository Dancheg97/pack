package cmd

import (
	"fmnx.io/dev/pack/system"
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
pack rm fmnx.io/dev/ainst`,
	Run: Remove,
}

func Remove(cmd *cobra.Command, pkgs []string) {
	mp := ReadMapping()
	revmp := ReverseMapping(mp)
	for _, pkg := range pkgs {
		pacmanpkg, ok := mp[pkg]
		if !ok {
			_, err := system.Call("pacman -Q " + pkg)
			if err != nil {
				continue
			}
			delete(mp, revmp[pkg])
			ExecuteCheck("sudo pacman --noconfirm -R " + pkg)
			continue
		}
		_, err := system.Call("sudo pacman --noconfirm -R " + pacmanpkg)
		if err != nil {
			YellowPrint("Pack package was not found in pacman: ", pkg)
		}
		delete(mp, pkg)
	}
	WriteMapping(mp)
}

func ReverseMapping(in map[string]string) map[string]string {
	r := map[string]string{}
	for k, v := range in {
		r[v] = k
	}
	return r
}
