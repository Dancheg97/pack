package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Run:   Update,
}

func Update(cmd *cobra.Command, args []string) {
	
}
