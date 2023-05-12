// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"context"
	"sync"

	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/pack"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
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
	pacmanOutdated, err := pacman.Outdated()
	CheckErr(err)
	packoutdated := GetPackOutdated()
	allOutdated := append(pacmanOutdated, packoutdated...)
	for _, info := range allOutdated {
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: info.Name + " ",
				Color:   prnt.WHITE,
			},
			{
				Message: info.CurrentVersion + " ",
				Color:   prnt.YELLOW,
			},
			{
				Message: info.NewVersion,
				Color:   prnt.BLUE,
			},
		})
	}
}

// Get pack outdated packages.
func GetPackOutdated() []pacman.OutdatedPackage {
	pkgs := pack.List()
	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var rez []pacman.OutdatedPackage
	for _, info := range pkgs {
		sinfo := info
		g.Go(func() error {
			link := "https://" + sinfo.PackName
			last, err := git.LastCommitUrl(link, sinfo.Branch)
			if err != nil {
				prnt.Yellow("Unable to get versoin for: ", link)
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
