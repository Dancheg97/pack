// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
)

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"upd", "u"},
	Short:   tmpl.SearchShort,
	Long:    tmpl.SearchLong,
	Run:     Search,
}

// Search for packages in connected search sources.
func Search(cmd *cobra.Command, args []string) {

}
