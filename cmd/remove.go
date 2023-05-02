package cmd

import (
	"fmnx.io/dev/pack/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   "🚫 remove packages",
	Long: `🚫 remove packages

Use this command to remove packages from system. You can specify both pacman 
packages and pack links.

Example:
pack rm fmnx.io/dev/ainst`,
	Run: Remove,
}

func Remove(cmd *cobra.Command, pkgs []string) {
	mp := ReadMapping()
	for _, packpkg := range pkgs {
		pacmanpkg, ok := mp[packpkg]
		if !ok {
			_, err := core.SystemCall("pacman -Q " + packpkg)
			if err != nil {
				YellowPrint("Package not found, skipping: ", packpkg)
				continue
			}
			ExecuteCheck("sudo pacman --noconfirm -R " + packpkg)
			continue
		}
		ExecuteCheck("sudo pacman --noconfirm -R " + pacmanpkg)
		delete(mp, packpkg)
		RedPrint("Package removed: ", packpkg)
	}
	WriteMapping(mp)
}

func WriteMapping(m PackMap) {
	if len(m) == 0 {
		core.WriteFile(cfg.MapFile, "")
		return
	}
	yamlData, err := yaml.Marshal(&m)
	CheckErr(err)
	core.WriteFile(cfg.MapFile, string(yamlData))
}
