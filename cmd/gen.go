package cmd

import (
	"fmt"
	"log"
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
	color.Cyan("README.md")
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")

	// Or just add them to New()
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("This prints bold cyan %s\n", "too!.")

	// Mix up foreground and background colors, create new mixes!
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")

	whiteBackground := red.Add(color.BgWhite)
	whiteBackground.Println("Red text with white background.")
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
	log.Println("unable to find ref in git config")
	lf.Unlock()
	os.Exit(1)
	return ""
}
