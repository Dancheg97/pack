// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/git"
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/system"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upd", "u"},
	Short:   tmpl.UpdateShort,
	Long:    tmpl.UpdateLong,
	Run:     Update,
}

// Cli command performing package update.
func Update(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) == 0 {
		FullPacmanUpdate()
		FullPackUpdate()
		return
	}
	groups := GroupPackages(pkgs)
	nfPacman := pacman.GetUninstalled(groups.PacmanPackages)
	nfPack := pack.GetUninstalled(groups.PackPackages)
	if len(nfPack) > 0 || len(nfPacman) > 0 {
		nfPack = append(nfPack, nfPacman...)
		prnt.Red("Some packages are not installed", strings.Join(nfPack, " "))
		os.Exit(1)
	}
	err := pacman.Update(groups.PacmanPackages)
	CheckErr(err)
	Install(nil, groups.PackPackages)
}

// Perform full pacman update.
func FullPacmanUpdate() {
	o, err := system.Call("sudo pacman --noconfirm -Syu")
	if err != nil {
		prnt.Red("Unable to update pacman packages: ", o)
		os.Exit(1)
	}
	prnt.Green("Pacman update: ", "done")
}

// Perform full pack update.
func FullPackUpdate() {
	outdatedpkgs := GetPackOutdated()
	var pkgs []string
	for _, pkg := range outdatedpkgs {
		pkgs = append(pkgs, fmt.Sprintf("%s@%s", pkg.Name, pkg.NewVersion))
	}
	Install(nil, pkgs)
	prnt.Green("Pack update: ", "done")
}

// Get list of pack outdated packages.
func GetPackOutdated() []pacman.OutdatedPackage {
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
			err = git.Pull(config.PkgCacheDir + "/" + sinfo.PacmanName)
			if err != nil {
				return err
			}
			return nil
		})
	}
	CheckErr(g.Wait())
	return rez
}
