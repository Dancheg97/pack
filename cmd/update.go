package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upd", "u"},
	Short:   "🗳️  update packages",
	Long: `🗳️  update packages

You can specify packages with versions, that you need them to update to, or 
provide provide just links to get latest version from default branch.

If you don't specify any arguements, the whole system will be updated.

Examples:
pack update
pack update fmnx.io/dev/aist@v0.21
pack update git.xmpl.sh/pkg
`,
	Run: Update,
}

func Update(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) == 0 {
		for pkg := range ReadMapping() {
			pkgs = append(pkgs, pkg)
		}
	}
	Get(cmd, pkgs)
}
