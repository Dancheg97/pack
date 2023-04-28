package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Flag struct {
	Cmd         *cobra.Command
	Name        string
	ShortName   string
	Env         string
	Value       string
	IsRequired  bool
	Type        string
	Description string
}

func AddFlag(cmd Flag) {
	if cmd.Type == "" {
		cmd.Cmd.PersistentFlags().StringP(cmd.Name, cmd.ShortName, cmd.Value, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		checkErr(err)
	}

	if cmd.Type == "strarr" {
		cmd.Cmd.PersistentFlags().StringArrayP(cmd.Name, cmd.ShortName, nil, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		checkErr(err)
	}

	if cmd.Type == "bool" {
		cmd.Cmd.PersistentFlags().BoolP(cmd.Name, cmd.ShortName, false, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		checkErr(err)
	}

	if cmd.Type == "int" {
		if cmd.Value != "" {
			i, err := strconv.Atoi(cmd.Value)
			if err != nil {
				err = fmt.Errorf("value for flag "+cmd.Name+" should be int: %w", err)
				checkErr(err)
			}
			cmd.Cmd.PersistentFlags().IntP(cmd.Name, cmd.ShortName, i, cmd.Description)
			err = viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
			checkErr(err)
			return
		}
		cmd.Cmd.PersistentFlags().IntP(cmd.Name, cmd.ShortName, 0, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		checkErr(err)
	}

	if cmd.Env != "" {
		err := viper.BindEnv(cmd.Name, cmd.Env)
		checkErr(err)
	}

	if cmd.IsRequired {
		err := cmd.Cmd.MarkFlagRequired(cmd.Name)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
