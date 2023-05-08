package cmd

import (
	"fmt"
	"strings"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/system"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "📄 list packages installed with pack",
	Run:     List,
}

func List(cmd *cobra.Command, args []string) {
	pkgs := GetPacmanPackages()
	mp := ReadMapping()
	revmp := ReverseMapping(mp)
	for k, v := range pkgs {
		if config.DisablePrettyPrint {
			fmt.Println(k, v, revmp[k])
			continue
		}
		fmt.Println(k, color.BlueString(v), color.YellowString(revmp[k]))
	}
}

func GetPacmanPackages() map[string]string {
	o, err := system.Call("pacman -Q")
	CheckErr(err)
	o = strings.Trim(o, "\n")
	pkgs := map[string]string{}
	for _, pkg := range strings.Split(o, "\n") {
		spl := strings.Split(pkg, " ")
		pkgs[spl[0]] = spl[1]
	}
	return pkgs
}
