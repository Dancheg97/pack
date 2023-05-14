// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmt"
	"os"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/tmpl"
	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "pack",
	Short:        tmpl.RootShort,
	Long:         tmpl.RootLong,
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
	rootCmd.SetUsageTemplate(tmpl.Cobra)
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

// Herlper function to exit on unexpected errors.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
