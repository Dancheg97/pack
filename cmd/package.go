package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   "📦 prepare package in current directory",
	Long: `📦 prepare package in current directory

This script will read .pack.yml, generate PKGBUILD and prepare .pkg.tar.zst
package. You can use it to test .pack.yml, to get PKGBUILD template for project
or validate installation for pack.

To double check installation, you can run it inside pack docker.

`,
	Run: Package,
}

func Package(cmd *cobra.Command, pkgs []string) {
}
