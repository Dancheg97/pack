// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "📨 push packages",
	Long: `📨 push packages

`,
	Run: Mirror,
}

func Push(cmd *cobra.Command, args []string) {

}

type Package struct {
	Registry      string
	PackageName   string
	PackageFile   string
	SignatureFile string
}

// This function will find the latest version of package in cache direcotry and
// then push it to registry specified in package name.
func CheckPackage(pkg string) (*Package, error) {
	splt := strings.Split(pkg, "/")
	if len(splt) != 2 {
		msg := "error: package should contain registry and name: "
		return nil, errors.New(msg + pkg)
	}
	return nil, nil
	// exec.Command("bash", "-c", "ls "+pacmancache+"")
}
