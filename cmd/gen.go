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

var template = `run-deps:
  - vlc
  - wget
  - git
build-deps:
  - flutter
  - clang
  - cmake
build-script:
  - flutter build linux
pack:
  build/linux/x64/release/bundle: /usr/share/pkg
  pkg: /usr/bin/pkg
  pkg.desktop: /usr/share/applications/pkg.desktop
  logo.png: /usr/share/icons/hicolor/512x512/apps/pkg.png
`

func Gen(cmd *cobra.Command, args []string) {
	core.WriteFile("pack.yml", template)
}
