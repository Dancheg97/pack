// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"context"
	"os"
	"strings"
	"sync"

	"fmnx.io/core/pack/database"
	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(outdatedCmd)
}

var outdatedCmd = &cobra.Command{
	Use:     "outdated",
	Aliases: []string{"out", "o"},
	Short:   tmpl.OutdatedShort,
	Long:    tmpl.OutdatedLong,
	Run:     Outdated,
}

// Cli command listing installed packages and their status.
func Outdated(cmd *cobra.Command, args []string) {
	pacmanOutdated := GetPacmanOutdated()
	packoutdated := GetPackOutdated()
	allOutdated := append(pacmanOutdated, packoutdated...)
	for _, info := range allOutdated {
		print.Custom([]print.ColoredMessage{
			{
				Message: info.Name + " ",
				Color:   print.WHITE,
			},
			{
				Message: info.CurrVersion + " ",
				Color:   print.YELLOW,
			},
			{
				Message: info.NewVersion,
				Color:   print.BLUE,
			},
		})
	}
}

type OutdatedPackageInfo struct {
	Name        string
	CurrVersion string
	NewVersion  string
}

// Get outdated packages and their versions.
func GetPacmanOutdated() []OutdatedPackageInfo {
	links := GetUpdateLinks()
	var out []OutdatedPackageInfo
	for _, link := range links {
		linkSplit := strings.Split(link, "/")
		file := linkSplit[len(linkSplit)-1]
		fileSplit := strings.Split(file, "-")
		packageName := strings.Join(fileSplit[0:len(fileSplit)-3], "-")
		newVersion := fileSplit[len(fileSplit)-3]
		currVersion := GetCurrentVersion(packageName)
		out = append(out, OutdatedPackageInfo{
			Name:        packageName,
			CurrVersion: currVersion,
			NewVersion:  newVersion,
		})
	}
	return out
}

// Get found update links for pacman packges.
func GetUpdateLinks() []string {
	o, err := system.Call("sudo pacman -Syup")
	if err != nil {
		print.Red("Unable to connect to pacman servers: ", "network error")
		os.Exit(1)
	}
	if !strings.Contains(o, "https://") {
		return nil
	}
	splt := strings.Split(o, "downloading...\n")
	pkgsLinksString := strings.Trim(splt[len(splt)-1], "\n")
	return strings.Split(pkgsLinksString, "\n")
}

// Get current package version.
func GetCurrentVersion(pkg string) string {
	o, err := system.Callf("pacman -Q %s", pkg)
	CheckErr(err)
	verAndRel := strings.Split(o, " ")[1]
	return strings.Trim(strings.Split(verAndRel, "-")[0], "\n")
}

// Get pack outdated packages.
func GetPackOutdated() []OutdatedPackageInfo {
	pkgs := database.List()
	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var rez []OutdatedPackageInfo
	for _, info := range pkgs {
		sinfo := info
		g.Go(func() error {
			link := "https://" + sinfo.PackName
			last, err := git.LastCommitUrl(link, sinfo.Branch)
			if err != nil {
				print.Yellow("Unable to get versoin for: ", link)
				return nil
			}
			if sinfo.Version == last {
				return nil
			}
			mu.Lock()
			rez = append(rez, OutdatedPackageInfo{
				Name:        sinfo.PackName,
				CurrVersion: sinfo.Version,
				NewVersion:  last,
			})
			mu.Unlock()
			return nil
		})
	}
	g.Wait()
	return rez
}
