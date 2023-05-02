package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/dev/pack/core"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "📋 generate .pack.yml and update README.md",
	Long: `📋 generate .pack.yml and update README.md

This command will generate .pack.yml template and add some lines to README.md
to provide information about installation with pack.`,
	Run: Gen,
}

const (
	packYmlTemplate = `# Dependencies, that are required for project at runtime.
run-deps:
  - vlc
  - wget
  - git
# Dependencies, that are required to build project.
build-deps:
  - flutter
  - clang
  - cmake
  # Scripts, that would be executed in root directory to get build files.
scripts:
  - flutter build linux
# File mapping for resulting build files and directories from project root
# to resulting file system.
# Each file or folder will be installed as it is mapped in this file.
mapping:
  pkg: /usr/bin/pkg
  pkg.desktop: /usr/share/applications/pkg.desktop
  logo.png: /usr/share/icons/hicolor/512x512/apps/pkg.png
  build/linux/x64/release/bundle: /usr/share/pkg
`
	readmeTemplate = `

---

### 📦 Install package with [pack](https://fmnx.io/dev/pack):

%s
pack get %s
%s
`
)

func Gen(cmd *cobra.Command, args []string) {
	core.WriteFile(".pack.yml", packYmlTemplate)
	insatllMd := fmt.Sprintf(readmeTemplate, "```", GetInstallLink(), "```")
	core.AppendToFile("README.md", insatllMd)
	fmt.Printf(
		"Generated file: %s \nUpdated files: %s, %s\n",
		color.RedString(".pack.yml"),
		color.CyanString("README.md"),
		color.HiYellowString(".gitignore"),
	)
}

func GetInstallLink() string {
	gitconf, err := os.ReadFile(`.git/config`)
	CheckErr(err)
	for _, line := range strings.Split(string(gitconf), "\n") {
		if strings.Contains(line, "url = ") {
			line = strings.Split(line, "url = ")[1]
			line = strings.ReplaceAll(line, "https://", "")
			line = strings.ReplaceAll(line, "git@", "")
			return strings.ReplaceAll(line, ".git", "")
		}
	}
	RedPrint("Error occured: ", "unable to find ref in git config")
	lf.Unlock()
	os.Exit(1)
	return ""
}
