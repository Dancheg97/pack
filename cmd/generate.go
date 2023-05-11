// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"fmt"
	"strings"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   tmpl.GenerateShort,
	Long:    tmpl.GenerateLong,
	Run:     Generate,
}

// Cli command modifying files in workdir for pack compatability.
func Generate(cmd *cobra.Command, args []string) {
	info := GetInstallLink()
	WritePackageBuild(info)
	ModifyReadmeFile(info.Link)
	print.Green("Updated files: ", "PKGBUILD README.md")
}

// Some basic info about repository.
type RepositoryInfo struct {
	ShortName string
	FullName  string
	Link      string
}

// Get information about repository located in current directory.
func GetInstallLink() RepositoryInfo {
	link, err := system.Callf("git config --get remote.origin.url")
	CheckErr(err)
	link = strings.Trim(link, "\n")
	link = strings.ReplaceAll(link, "https://", "")
	link = strings.ReplaceAll(link, "git@", "")
	link = strings.ReplaceAll(link, ":", "/")
	link = strings.ReplaceAll(link, ".git", "")
	splt := strings.Split(link, "/")
	return RepositoryInfo{
		ShortName: splt[len(splt)-1],
		FullName:  link,
		Link:      "https://" + link,
	}
}

// Owerwrite current PKGBUILD file.
func WritePackageBuild(i RepositoryInfo) {
	tmpl := fmt.Sprintf(tmpl.PKGBUILD, i.ShortName, i.Link)
	err := system.WriteFile("PKGBUILD", tmpl)
	CheckErr(err)
}

// Append some lines to README.md file.
func ModifyReadmeFile(link string) {
	insatllMd := fmt.Sprintf(tmpl.READMEmd, "```", link, "```")
	system.AppendToFile("README.md", insatllMd)
}
