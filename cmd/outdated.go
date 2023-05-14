// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"context"
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
