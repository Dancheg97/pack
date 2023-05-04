package cmd

import (
	"strings"

	"fmnx.io/dev/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   "📦 prepare package in current directory",
	Long: `📦 prepare package in current directory

This script will read .pack.yml, generate PKGBUILD and prepare .pkg.tar.zst
package. You can use it to test .pack.yml, to get PKGBUILD template for project
or validate installation for pack.

To double check installation, you can run it inside pack docker.

`,
	Run: Package,
}

func Package(cmd *cobra.Command, pkgs []string) {
	pack := ReadPackYml()
	for _, script := range pack.BuildScripts {
		system.Debug = true
		ExecuteCheck(script)
	}
	info := GetCurrentRepositoryInformation()
	GeneratePkgbuild(info, pack)
	GreenPrint("file generated: ", "PKGBUILD")
	ExecuteCheck("makepkg")
	GreenPrint("package build: ", "success")
}

func GetCurrentRepositoryInformation() PkgInfo {
	out, err := system.Call("git rev-parse --abbrev-ref HEAD")
	CheckErr(err)
	link := GetInstallLink()
	splitted := strings.Split(link, "/")
	return PkgInfo{
		FullName:  link,
		ShortName: splitted[len(splitted)-1],
		HttpsLink: "https://" + link,
		Version:   strings.Trim(out, "\n"),
	}
}
