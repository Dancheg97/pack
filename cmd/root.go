package cmd

import (
	"fmt"
	"os"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/system"
	"github.com/fatih/color"
	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
)

const descrTmpl = `{{if gt (len .Aliases) 0}}Aliases:
{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{end}}
`

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

var (
	cfg *config.Config
	lf  *lockfile.Lockfile
)

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{})
	rootCmd.SetUsageTemplate(descrTmpl)

	conf, err := config.GetConfig()
	CheckErr(err)
	system.Debug = conf.DebugMode
	cfg = conf
	lock, err := lockfile.New(conf.LockFile)
	CheckErr(err)
	lf = &lock
	err = lf.TryLock()
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
		RedPrint("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func ExecuteCheck(script string) {
	out, err := system.Call(script)
	if err != nil {
		RedPrint("Command did not succed: ", script)
		fmt.Println("System output: ", out)
		RedPrint("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func Chdir(dir string) {
	CheckErr(os.Chdir(dir))
}

func RedPrint(white string, red string) {
	if cfg.DisablePrettyPrint {
		fmt.Printf(white + red + "\n")
		return
	}
	fmt.Printf(white + color.RedString(red) + "\n")
}

func BluePrint(white string, blue string) {
	if cfg.DisablePrettyPrint {
		fmt.Printf(white + blue + "\n")
		return
	}
	fmt.Printf(white + color.BlueString(blue) + "\n")
}

func GreenPrint(white string, green string) {
	if cfg.DisablePrettyPrint {
		fmt.Printf(white + green + "\n")
		return
	}
	fmt.Printf(white + color.GreenString(green) + "\n")
}

func YellowPrint(white string, yellow string) {
	if cfg.DisablePrettyPrint {
		fmt.Printf(white + yellow + "\n")
		return
	}
	fmt.Printf(white + color.YellowString(yellow) + "\n")
}
