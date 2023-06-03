// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
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
	Run: Push,
}

func Push(cmd *cobra.Command, args []string) {
	var pkgs []*Package
	for _, pkg := range args {
		p, err := FormPackage(pkg)
		CheckErr(err)
		pkgs = append(pkgs, p)
	}
	fmt.Println(pkgs[0])
}

type Package struct {
	Registry string
	PkgName  string
	PkgFile  string
	SigFile  string
}

// This function will find the latest version of package in cache direcotry and
// then push it to registry specified in package name provided in argiement.
func FormPackage(pkg string) (*Package, error) {
	splt := strings.Split(pkg, "/")
	if len(splt) != 2 {
		msg := "error: package should contain registry and name: "
		return nil, errors.New(msg + pkg)
	}
	registry := splt[0]
	pkgname := splt[1]
	des, err := os.ReadDir(pacmancache)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		if !strings.HasSuffix(de.Name(), pkgext) {
			continue
		}
		pkgsplt := strings.Split(de.Name(), "-")
		if len(pkgsplt) < 4 {
			return nil, errors.New("invalid package in cache: " + de.Name())
		}
		namecheck := strings.Join(pkgsplt[:len(pkgsplt)-3], "-")
		if pkgname == namecheck {
			return &Package{
				Registry: registry,
				PkgName:  pkgname,
				PkgFile:  path.Join(pacmancache, de.Name()),
				SigFile:  path.Join(pacmancache, de.Name()+".sig"),
			}, err
		}
	}
	return nil, errors.New("package not found in cache: " + pkgname)
}
