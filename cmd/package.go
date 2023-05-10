// Copyright 2023 FMNX Linux team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pkgCmd)
}

var pkgCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg", "p"},
	Short:   "📦 prepare and install package",
	Long: `📦 prepare .pkg.tar.zst in current directory and install it

This script will read prepare .pkg.tar.zst package. You can use it to test 
PKGBUILD template for project or validate installation for pack.

To double check installation, you can test it inside pack docker container:
docker run --rm -it fmnx.io/core/pack i example.com/package
`,
	Run: Package,
}

// Cli command preparing package in current directory.
func Package(cmd *cobra.Command, pkgs []string) {
	print.Blue("Preparing package: ", "makepkg -sfri --noconfirm")
	out, err := system.Call("makepkg -sfri --noconfirm")
	if err != nil {
		print.Red("Unable to execute: ", "makepkg")
		fmt.Println(out)
		os.Exit(1)
	}
	i := GetInstallLink()
	AddToPackMapping(i)
	print.Green("Package prepared and installed: ", i.FullName)
}

// Function writes package to pack mapping file.
func AddToPackMapping(i RepositoryInfo) {
	mp := ReadMapping()
	mp[i.FullName] = i.ShortName
	WriteMapping(mp)
}

type PackMap map[string]string

// Function reads pack mapping file. Packages are mapped from pack to pacman.
func ReadMapping() PackMap {
	_, err := os.Stat(config.MapFile)
	if err != nil {
		system.AppendToFile(config.MapFile, "{}")
		return PackMap{}
	}
	b, err := os.ReadFile(config.MapFile)
	CheckErr(err)
	var mapping PackMap
	err = json.Unmarshal(b, &mapping)
	CheckErr(err)
	return mapping
}

// Function writes pack mapping file.
func WriteMapping(m PackMap) {
	if len(m) == 0 {
		system.WriteFile(config.MapFile, "{}")
		return
	}
	jsonData, err := json.Marshal(&m)
	CheckErr(err)
	system.WriteFile(config.MapFile, string(jsonData))
}
