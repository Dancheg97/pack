package cmd

import (
	"fmt"
	"strings"

	"fmnx.io/dev/pack/core"
	"github.com/spf13/cobra"
)

const (
	CacheDir = `~/.pack-cache`
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "📥 insatll new packages",
	Run:   Get,
}

func Get(cmd *cobra.Command, pkgs []string) {
	err := core.SystemCallf("mkdir -p %s", CacheDir)
	CheckErr(err)

	for _, pkg := range pkgs {
		pkgLink := "https://" + strings.Split(pkg, "@")[0]
		pkgsplit := strings.Split(pkgLink, "/")
		pkgName := pkgsplit[len(pkgsplit)-1]
		err := core.SystemCallf("git clone %s %s/%s", pkgLink, CacheDir, pkgName)
		if strings.Contains(err.Error(), "already exists and is not an empty") {
			fmt.Println("pulling changes")
			core.SystemCallf("git -C %s/%s pull ", CacheDir, pkgName)
		}
	}
}
