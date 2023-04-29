package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/dev/pack/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "📋 generate pack.yml template",
	Run:   Gen,
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
build-scripts:
  - flutter build linux
# File mapping for resulting build files and directories from project root
# to resulting file system.
# Each file or folder will be installed as it is mapped in this file.
pack-map:
  pkg: /usr/bin/pkg
  pkg.desktop: /usr/share/applications/pkg.desktop
  logo.png: /usr/share/icons/hicolor/512x512/apps/pkg.png
  build/linux/x64/release/bundle: /usr/share/pkg
`
	gitignoreTemplate = `
pkg/**
src/**
**.pkg.tar.zst
PKGBUILD
`
	readmeTemplate = `

---

### 📦 Install package with pack:

%s
pack get %s
%s
`
)

func Gen(cmd *cobra.Command, args []string) {
	core.WriteFile("pack.yml", packYmlTemplate)
	core.AppendToFile(".gitignore", gitignoreTemplate)
	insatllMd := fmt.Sprintf(readmeTemplate, "```", GetInstallLink(), "```")
	core.AppendToFile("README.md", insatllMd)
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
	fmt.Println("unable to find ref in git config")
	os.Exit(1)
	return ""
}
