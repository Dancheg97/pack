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
	Use:   "remove",
	Short: "🚫 remove package from system",
	Run:   Remove,
}

func Remove(cmd *cobra.Command, pkgs []string) {
	mp := ReadMapping()
	for _, pkg := range pkgs {
		core.SystemCall("sudo pacman -R " + mp[pkg])
		delete(mp, pkg)
	}
}

func WriteMapping(m PackMap) {
	yamlData, err := yaml.Marshal(&m)
	CheckErr(err)
	core.WriteFile(cfg.MapFile, string(yamlData))
}
