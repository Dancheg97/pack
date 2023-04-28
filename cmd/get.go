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
		PrepareRepo(pkg)
		SwitchToVersion(pkg)
	}
}

func PrepareRepo(pkg string) {
	pkgLink := "https://" + strings.Split(pkg, "@")[0]
	pkgSplit := strings.Split(pkgLink, "/")
	pkgName := pkgSplit[len(pkgSplit)-1]
	err := core.SystemCallf("git clone %s %s/%s", pkgLink, CacheDir, pkgName)
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			CheckErr(err)
		}
		fmt.Println("pulling changes")
		err = core.SystemCallf("git -C %s/%s pull ", CacheDir, pkgName)
		CheckErr(err)
	}
}

func SwitchToVersion(pkg string) {
	pkgRepo := strings.Split(pkg, "@")[0]
	pkgSplit := strings.Split(pkgRepo, "/")
	pkgName := pkgSplit[len(pkgSplit)-1]
	if len(strings.Split(pkg, "@")) == 1 {
		branch, err := GetDefaultBranch(pkg)
		CheckErr(err)
		err = core.SystemCallf("git -C %s/%s checkout %s", CacheDir, pkgName, branch)
		CheckErr(err)
		return
	}
	pkgVer := strings.Split(pkg, "@")[1]
	err := core.SystemCallf("git -C %s/%s checkout %s", CacheDir, pkgName, pkgVer)
	CheckErr(err)
}

func GetDefaultBranch(pkg string) (string, error) {
	pkgLink := "https://" + strings.Split(pkg, "@")[0]
	out, err := core.SystemCallOutf("git remote show %s | sed -n '/HEAD branch/s/.*: //p'", pkgLink)
	CheckErr(err)
	return strings.Trim(out, "\n"), nil
}
