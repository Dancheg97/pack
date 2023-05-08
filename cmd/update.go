package cmd

// import (
// 	"github.com/spf13/cobra"
// )

// func init() {
// 	rootCmd.AddCommand(updateCmd)
// }

// var updateCmd = &cobra.Command{
// 	Use:     "update",
// 	Aliases: []string{"upd", "u"},
// 	Short:   "🗳️  update packages",
// 	Long: `🗳️  update packages

// You can specify packages with versions, that you need them to update to, or
// provide provide just links to get latest version from default branch.

// If you don't specify any arguements, all packages will be updated.

// Examples:
// pack update
// pack update fmnx.io/core/aist@v0.21
// pack update git.xmpl.sh/pkg
// `,
// 	Run: Update,
// }

// var Updating = false

// func Update(cmd *cobra.Command, pkgs []string) {
// 	Updating = true
// 	if len(pkgs) == 0 {
// 		print.Blue("Starting pacman update: ", "pacman -Syu")
// 		ExecuteCheck("sudo pacman --noconfirm -Syu")
// 		for pkg := range ReadMapping() {
// 			pkgs = append(pkgs, pkg)
// 		}
// 	}
// 	Get(cmd, pkgs)
// }
