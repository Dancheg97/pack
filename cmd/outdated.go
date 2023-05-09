package cmd

import (
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(outdatedCmd)
}

var outdatedCmd = &cobra.Command{
	Use:     "outdated",
	Aliases: []string{"out", "o"},
	Short:   "📌 show outdated packages",
	Run:     Outdated,
}

// Cli command listing installed packages and their status.
func Outdated(cmd *cobra.Command, args []string) {
	outdatedPkgs := GetOutdated()
	fmt.Println(outdatedPkgs)
}

type OutdatedPackage struct {
	Name        string
	CurrVersion string
	NewVersion  string
}

// Get outdated packages and their versions.
func GetOutdated() []OutdatedPackage {
	o, err := system.Call("sudo pacman -Syup")
	if err != nil {
		print.Red("Unable to connect to pacman servers: ", "network error")
		os.Exit(1)
	}
	splt := strings.Split(o, "is up to date\n")
	pkgsLinksString := splt[len(splt)-1]
	var out []OutdatedPackage
	for _, link := range strings.Split(pkgsLinksString, "\n") {
		splt = strings.Split(link, "/")
		file := splt[len(splt)-1]
		splt = strings.Split(file, "-")
		curr := GetCurrentVersion(splt[0])
		out = append(out, OutdatedPackage{
			Name:        splt[0],
			CurrVersion: curr,
			NewVersion:  splt[1],
		})
	}
	return out
}

// Get current package version.
func GetCurrentVersion(pkg string) string {
	o, err := system.Callf("pacman -Q %s", pkg)
	fmt.Println(o)
	CheckErr(err)
	return strings.Trim(strings.Split(o, " ")[1], "\n")
}
