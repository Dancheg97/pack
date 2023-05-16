// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
)

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c", "cfg"},
	Short:   tmpl.ConfigShort,
	Long:    tmpl.ConfigLong,
	Run:     Install,
}

// View and change config
func Config(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		//
	}

	if len(args) == 1 && args[0] == "reset" {
		config.SetDefaults()
		config.Save()

		return
	}

}
