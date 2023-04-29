package cmd

import (
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

var packYmlTemplate = `# Dependencies, that are required for project at runtime.
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

func Gen(cmd *cobra.Command, args []string) {
	core.WriteFile("pack.yml", packYmlTemplate)
}
