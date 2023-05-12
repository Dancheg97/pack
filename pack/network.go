// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package pack

import (
	"context"
	"sync"

	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/pacman"
	"fmnx.io/core/pack/prnt"
	"golang.org/x/sync/errgroup"
)

// Get list of pack outdated packages.
func Outdated() []pacman.OutdatedPackage {
	pkgs := List()
	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var rez []pacman.OutdatedPackage
	for _, info := range pkgs {
		sinfo := info
		g.Go(func() error {
			link := "https://" + sinfo.PackName
			last, err := git.LastCommitUrl(link, sinfo.Branch)
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
