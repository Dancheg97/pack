package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "",
	Run:   Gen,
}

func Gen(cmd *cobra.Command, args []string) {
}
