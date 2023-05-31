// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmnx.su/core/pack/pack"
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

This command will create template, containing files required for arch desktop
applications:

- PKGBUILD
- app.desktop
- app.sh`,
	Run: Create,
}

// Cli command installing packages into system.
func Create(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, "")
	}
	err := pack.Create(args[0])
	CheckErr(err)
}
