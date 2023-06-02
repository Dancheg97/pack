// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "📤 push packages",
	Run:     Mirror,
}

func Push(cmd *cobra.Command, args []string) {

}

// This function will find the latest version of package in cache direcotry and
// then push it to registry specified in package name.
func PushPackage(pkg string) error {
	return nil
}
