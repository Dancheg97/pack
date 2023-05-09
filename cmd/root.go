// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information could be found on official web page: https://fmnx.io/
// Email for contacts: help@fmnx.io
package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/print"
	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pack",
	Short: "📦 git based arch compatible package manager.",
	Long: `📦 git based arch compatible package manager.

This utility accumulates power of git and pacman to provide decentralized way
of arch package distribution. Pack config: '~/.pack/config.yml'. Find more 
information at https://fmnx.io/core/pack.

Usage:
pack [command] <package(s)>

Example:
pack install fmnx.io/core/ainst`,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	},
}

// Prepare cobra and viper templates.
func init() {
	rootCmd.SetHelpCommand(&cobra.Command{})
	rootCmd.SetUsageTemplate(descrTmpl)
	lock, err := lockfile.New(config.LockFile)
	CheckErr(err)
	err = lock.TryLock()
	CheckErr(err)
}

// Main execution of cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Utility to exit on unexpected errors.
func CheckErr(err error) {
	if err != nil {
		print.Red("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

const descrTmpl = `{{if gt (len .Aliases) 0}}Aliases:
{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{end}}
`
