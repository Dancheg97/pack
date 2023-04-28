package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fmnx-gen",
	Short: "🧰 Cli tool for something awesome.",
	Long:  "long example",
}

var flags = []Flag{
	{
		Cmd:         rootCmd,
		Name:        "flg",
		ShortName:   "f",
		Env:         "FLG",
		Value:       "value",
		Description: "📄 cool description",
	},
}

func Execute() {
	for _, flag := range flags {
		AddFlag(flag)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
