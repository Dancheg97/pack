// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "📦 create new package",
	Long: `📦 create new package

This command will generate PKGBUILD file and create initialize git submodule 
for your package. After you set up build scripts, you can use pack build to
build your package.`,
	Run: Create,
}

// Cli command installing packages into system.
func Create(cmd *cobra.Command, args []string) {

}
