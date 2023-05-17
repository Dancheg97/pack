// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"context"
	"os"
	"strings"
	"sync"

	"fmnx.su/core/pack/git"
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   tmpl.ListShort,
	Long:    tmpl.ListLong,
	Run:     List,
}

// Cli command listing installed packages and version.
func List(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		PrintPackages(true, true)
		return
	}
	if len(args) != 1 {
		prnt.Red("Too many arguemnents for list: ", strings.Join(args, " "))
		os.Exit(1)
	}
	switch args[0] {
	case "outdated":
		ShowOutdated()
		return
	case "pack":
		PrintPackages(true, false)
		return
	case "pacman":
		PrintPackages(false, true)
		return
	}
	prnt.Red("Unknown arguement: ", args[0])
	os.Exit(1)
}

// Show all outdated packages for pack and pacman.
func ShowOutdated() {
	pacmanOutdated, err := pacman.Outdated()
	CheckErr(err)
	packoutdated := PackOutdated()
	allOutdated := append(pacmanOutdated, packoutdated...)
	for _, info := range allOutdated {
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: info.Name + " ",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: info.CurrentVersion + " ",
				Color:   prnt.COLOR_YELLOW,
			},
			{
				Message: info.NewVersion,
				Color:   prnt.COLOR_BLUE,
			},
		})
	}
}

// Get list of pack outdated packages.
func PackOutdated() []pacman.OutdatedPackage {
	pkgs := pack.List()
	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var rez []pacman.OutdatedPackage
	for _, info := range pkgs {
		sinfo := info
		g.Go(func() error {
			link := "https://" + sinfo.PackName
			last, err := git.LastCommitUrl(link, sinfo.DefaultBranch)
			if err != nil {
				mu.Lock()
				prnt.Yellow("Unable to get versoin for: ", link)
				mu.Unlock()
				return nil
			}
			if sinfo.Version == last {
				return nil
			}
			mu.Lock()
			rez = append(rez, pacman.OutdatedPackage{
				Name:           sinfo.PackName,
				CurrentVersion: sinfo.Version,
				NewVersion:     last,
			})
			mu.Unlock()
			return nil
		})
	}
	g.Wait()
	return rez
}

// Function that prints packages, with their key params.
func PrintPackages(showPack bool, showPacman bool) {
	pkgs := pacman.List()
	for pkg, version := range pkgs {
		i, err := pack.GetByPacmanName(pkg)
		if err != nil && showPacman {
			prnt.Custom([]prnt.ColoredMessage{
				{
					Message: pkg + " ",
					Color:   prnt.COLOR_WHITE,
				},
				{
					Message: version,
					Color:   prnt.COLOR_BLUE,
				},
			})
			continue
		}
		if err != nil || !showPack {
			continue
		}
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: i.PacmanName + " ",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: i.PackName + " ",
				Color:   prnt.COLOR_YELLOW,
			},
			{
				Message: i.DefaultBranch,
				Color:   prnt.COLOR_BLUE,
			},
			{
				Message: "-",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: i.Version,
				Color:   prnt.COLOR_BLUE,
			},
		})
	}
}
