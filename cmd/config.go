// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"strconv"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/tmpl"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c", "cfg"},
	Short:   tmpl.ConfigShort,
	Long:    tmpl.ConfigLong,
	Run:     Config,
}

// View and change config
func Config(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		ShowConfigDescription()
		return
	}
	if len(args) == 1 && args[0] == "reset" {
		config.SetDefaults()
		config.Save()
		return
	}
	if len(args) != 2 {
		fmt.Println("bad input for config, run 'pack config -h'")
		return
	}
	switch args[0] {
	case `needed`:
		config.Needed = ParseBool(args[1])
	case `rm-deps`:
		config.RmDeps = ParseBool(args[1])
	case `rm-repos`:
		config.RmRepos = ParseBool(args[1])
	case `cache-pkgs`:
		config.CachePkgs = ParseBool(args[1])
	case `verbose`:
		config.Verbose = ParseBool(args[1])
	case `pretty-print`:
		config.PrettyPrint = ParseBool(args[1])
	case `repo-cache-dir`:
		config.RepoCacheDir = args[1]
	case `pkg-cache-dir`:
		config.PkgCacheDir = args[1]
	case `log-file`:
		config.LogFile = args[1]
	case `map-file`:
		config.MapFile = args[1]
	case `lock-file`:
		config.LockFile = args[1]
	default:
		fmt.Println("unable to find config arguement: ", args[0])
		return
	}
	config.Save()
}

// Show configuration variables of configuration and describe them.
func ShowConfigDescription() {
	fmt.Printf(
		tmpl.PrettyConfig,
		color.CyanString(fmt.Sprint(config.Needed)),
		color.CyanString(fmt.Sprint(config.RmDeps)),
		color.CyanString(fmt.Sprint(config.RmRepos)),
		color.CyanString(fmt.Sprint(config.CachePkgs)),
		color.CyanString(fmt.Sprint(config.Verbose)),
		color.CyanString(fmt.Sprint(config.PrettyPrint)),
		color.CyanString(config.RepoCacheDir),
		color.CyanString(config.PkgCacheDir),
		color.CyanString(config.LogFile),
		color.CyanString(config.MapFile),
		color.CyanString(config.LockFile),
	)
}

// Parse boolean variable from string for config.
func ParseBool(s string) bool {
	o, err := strconv.ParseBool(s)
	CheckErr(err)
	return o
}
