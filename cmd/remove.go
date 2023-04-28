package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "",
	Run:   Remove,
}

func Remove(cmd *cobra.Command, args []string) {

}
