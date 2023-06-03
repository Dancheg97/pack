// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	pacmancache  = "/var/cache/pacman/pkg"
	pkgext       = ".pkg.tar.zst"
	dbext        = ".db.tar.gz"
	protocol     = "https://"
	fsendpoint   = "/pack/"
	pushendpoint = "/pack/push"
	file         = "file"
	sign         = "sign"
)

var rootCmd = &cobra.Command{
	Use:          "pack",
	Short:        "📦 Simplified version of pacman with extended functionality",
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
	rootCmd.SetUsageTemplate(CobraTmpl)
}

// Main execution of cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Herlper function to exit on unexpected errors.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Parameters to add boolean persistent flag.
type FlagParameters struct {
	// Cobra cmd to add flag.
	Cmd *cobra.Command
	// Flag name, for example --help.
	Name string
	// Flag shortname, for example -h.
	Short string
	// Flag description.
	Desc string
	// Default value (only for string flags).
	Default string
	// Environment variable to be automatically assigned.
	Env string
}

// Add boolean flag to cobra command.
func AddBoolFlag(p *FlagParameters) {
	p.Cmd.PersistentFlags().BoolP(p.Name, p.Short, false, p.Desc)
	err := viper.BindPFlag(p.Name, p.Cmd.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
	err = viper.BindEnv(p.Name, p.Env)
	CheckErr(err)
}

// Add string flag to cobra command.
func AddStringFlag(p *FlagParameters) {
	p.Cmd.PersistentFlags().StringP(p.Name, p.Short, p.Default, p.Desc)
	err := viper.BindPFlag(p.Name, p.Cmd.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
	err = viper.BindEnv(p.Name, p.Env)
	CheckErr(err)
}

// Add string flag to cobra command.
func AddStringListFlag(p *FlagParameters) {
	p.Cmd.PersistentFlags().StringArrayP(p.Name, p.Short, nil, p.Desc)
	err := viper.BindPFlag(p.Name, p.Cmd.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
	err = viper.BindEnv(p.Name, p.Env)
	CheckErr(err)
}

const CobraTmpl = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}} <pkgs>

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{end}}{{if .HasAvailableSubCommands}}{{end}}
`
