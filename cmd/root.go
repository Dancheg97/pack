package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
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

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{})
	rootCmd.SetUsageTemplate(descrTmpl)
	lock, err := lockfile.New(config.LockFile)
	CheckErr(err)
	err = lock.TryLock()
	CheckErr(err)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func CheckErr(err error) {
	if err != nil {
		print.Red("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func ExecuteCheck(script string) {
	out, err := system.Call(script)
	if err != nil {
		print.Red("Command did not succed: ", script)
		fmt.Println("System output: ", out)
		print.Red("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func Chdir(dir string) {
	CheckErr(os.Chdir(dir))
}

const descrTmpl = `{{if gt (len .Aliases) 0}}Aliases:
{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{end}}
`
