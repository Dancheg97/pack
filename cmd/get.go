package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "📥 insatll new packages",
	Run:   Gen,
}

func Get(cmd *cobra.Command, args []string) {
}
