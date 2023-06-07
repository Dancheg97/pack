// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"d"},
	Short:   "📃 describe packages",
	Run:     Describe,
}

func Describe(cmd *cobra.Command, args []string) {
	for _, pkg := range args {
		cmd := exec.Command("pacman", "-Qi", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		CheckErr(cmd.Run())
	}
}
