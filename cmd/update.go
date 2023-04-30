package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "🔁 update packages",
	Run:   Update,
}

func Update(cmd *cobra.Command, pkgs []string) {
	// mp := ReadMapping()
	// if len(pkgs) == 0 {

	// }
}
