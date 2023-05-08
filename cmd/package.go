package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   "📦 prepare .pkg.tar.zst and install it",
	Long: `📦 prepare .pkg.tar.zst and install it

This script will read prepare .pkg.tar.zst package. You can use it to test 
PKGBUILD template for project or validate installation for pack.

To double check installation, you can test it inside pack docker container:
docker run --rm -it ${pwd}:/src -v /src fmnx.io/core/pack package
`,
	Run: Package,
}

type PackMap map[string]string

func Package(cmd *cobra.Command, pkgs []string) {
	BluePrint("Preparing package: ", "makepkg -sfri --noconfirm")
	out, err := system.Call("makepkg -sfri --noconfirm")
	if err != nil {
		RedPrint("Unable to execute:", "makepkg")
		fmt.Println(out)
		os.Exit(1)
	}
	i := GetInstallLink()
	AddToPackMapping(i)
	GreenPrint("Package prepared and installed: ", i.FullName)
}

func AddToPackMapping(i PackageInfo) {
	mp := ReadMapping()
	mp[i.FullName] = i.ShortName
	WriteMapping(mp)
}

func ReadMapping() PackMap {
	_, err := os.Stat(cfg.MapFile)
	if err != nil {
		system.AppendToFile(cfg.MapFile, "{}")
		return PackMap{}
	}
	b, err := os.ReadFile(cfg.MapFile)
	CheckErr(err)
	var mapping PackMap
	err = json.Unmarshal(b, &mapping)
	CheckErr(err)
	return mapping
}

func WriteMapping(m PackMap) {
	if len(m) == 0 {
		system.WriteFile(cfg.MapFile, "{}")
		return
	}
	jsonData, err := json.Marshal(&m)
	CheckErr(err)
	system.WriteFile(cfg.MapFile, string(jsonData))
}
