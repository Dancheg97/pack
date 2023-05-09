package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(outdatedCmd)
}

var outdatedCmd = &cobra.Command{
	Use:     "outdated",
	Aliases: []string{"out", "o"},
	Short:   "📌 show outdated packages",
	Run:     Outdated,
}

// Cli command listing installed packages and their status.
func Outdated(cmd *cobra.Command, args []string) {

}

type OutdatedPackage struct {
	Name        string
	CurrVersion string
	NewVersion  string
}

// Get outdated packages and their versions.
