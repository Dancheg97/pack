package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"descr", "d"},
	Short:   "🪪  describe packages",
	Long: `🪪  view information about packages

This tool provides information about package retrieved from pacman or pack.

Example:
pack describe fmnx.io/core/ainst`,
	Run: Describe,
}

// Cli command giving package description.
func Describe(cmd *cobra.Command, pkgs []string) {
	for _, pkg := range pkgs {
		packageMapping := ReadMapping()
		pacmanpkg, ok := packageMapping[pkg]
		if !ok {
			pacmanpkg = pkg
		}
		info, err := system.Call("pacman -Qi " + pacmanpkg)
		if err != nil {
			print.Red("Error: ", strings.ReplaceAll(info, "error: ", ""))
			os.Exit(1)
		}
		fmt.Println(info)
	}
}
