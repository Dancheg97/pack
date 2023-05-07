package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/dev/pack/system"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"descr", "d"},
	Short:   "🪪  describe package",
	Long: `🪪  view information about package

This tool provides information about package retrieved from pacman or pack.

Example:
pack describe fmnx.io/dev/ainst`,
	Run: Info,
}

func Info(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) != 1 {
		RedPrint(
			"Please, specify single arguement, provided: ",
			color.RedString(strings.Join(pkgs, " ")),
		)

		os.Exit(1)
	}
	pkg := pkgs[0]
	packageMapping := ReadMapping()
	pacmanpkg, ok := packageMapping[pkg]
	if !ok {
		pacmanpkg = pkg
	}
	info, err := system.Call("pacman -Qi " + pacmanpkg)
	if err != nil {
		RedPrint("Error: ", strings.ReplaceAll(info, "error: ", ""))
		os.Exit(1)
	}
	fmt.Println(info)
}
