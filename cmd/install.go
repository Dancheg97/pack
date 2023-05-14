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
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   tmpl.InstallShort,
	Long:    tmpl.InstallLong,
	Run:     Install,
}

// Cli command installing packages into system.
func Install(cmd *cobra.Command, upkgs []string) {
	groups := GroupPackages(upkgs)
	CheckUnreachablePacmanPackages(groups.PacmanPackages)
	CheckUnreachablePackPackages(groups.PackPackages)
	CheckErr(pacman.Install(groups.PacmanPackages))
	InstallPackPackages(groups.PackPackages)
}

type PackageGroups struct {
	PacmanPackages []string
	PackPackages   []string
}

// Split packages into pacman and pack groups.
func GroupPackages(pkgs []string) PackageGroups {
	var pacmanPackages []string
	var packPackages []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "/") {
			packPackages = append(packPackages, pkg)
			continue
		}
		pacmanPackages = append(pacmanPackages, pkg)
	}
	return PackageGroups{
		PacmanPackages: pacmanPackages,
		PackPackages:   packPackages,
	}
}

// Check if some pacman packages could not be installed.
func CheckUnreachablePacmanPackages(pkgs []string) {
	unreachable := pacman.GetUnreachable(pkgs)
	if len(unreachable) > 0 {
		pkgs := strings.Join(unreachable, " ")
		prnt.Red("Unable to resolve those pacman packages: ", pkgs)
		os.Exit(1)
	}
}

// Check if some pack packages could not be installed.
func CheckUnreachablePackPackages(pkgs []string) {
	g, _ := errgroup.WithContext(context.Background())
	var unreachable []string
	for _, pkg := range pkgs {
		spkg := pkg
		g.Go(func() error {
			_, err := pack.Get(spkg)
			if err == nil {
				return nil
			}
			info := EjectInfoFromPackLink(spkg)
			err = git.Clone(info.Url, info.Directory)
			if err != nil {
				unreachable = append(unreachable, spkg)
				return err
			}
			_, err = os.Stat(info.Pkgbuild)
			if err != nil {
				unreachable = append(unreachable, spkg)
			}
			return err
		})
	}
	err := g.Wait()
	if err != nil {
		out := strings.Join(unreachable, " ")
		prnt.Red("Some pack packages are unreachable: ", out)
		os.Exit(1)
	}
}

// Info formed from pack link about all information related to that package.
type PackInfo struct {
	PacmanName string
	PackName   string
	Directory  string
	Version    string
	Pkgbuild   string
	Url        string
}

// Eject pack information for provided pack link.
func EjectInfoFromPackLink(pkg string) PackInfo {
	rez := PackInfo{}
	versplt := strings.Split(pkg, "@")
	rez.PackName = versplt[0]
	rez.Url = "https://" + versplt[0]
	if len(versplt) > 1 {
		rez.Version = versplt[1]
	}
	dashsplt := strings.Split(rez.PackName, "/")
	rez.PacmanName = dashsplt[len(dashsplt)-1]
	rez.Directory = config.RepoCacheDir + "/" + rez.PacmanName
	rez.Pkgbuild = rez.Directory + "/PKGBUILD"
	return rez
}

// Checks if packages are not installed and installing them.
func InstallPackPackages(pkgs []string) {
	for _, pkg := range pkgs {
		_, err := pack.Get(pkg)
		if err == nil && !Updating {
			continue
		}
		InstallPackPackage(EjectInfoFromPackLink(pkg))
	}
	if len(pkgs) > 0 {
		pkglist := strings.Join(pkgs, " ")
		if !Updating {
			prnt.Green("Installed: ", pkglist)
		}
	}
}

// Install pack package.
func InstallPackPackage(i PackInfo) {
	err := git.Clean(i.Directory)
	CheckErr(err)
	branch, err := git.DefaultBranch(i.Directory)
	CheckErr(err)
	err = git.Checkout(i.Directory, branch)
	CheckErr(err)
	err = git.Pull(i.Directory)
	CheckErr(err)
	if i.Version == `` {
		i.Version, err = git.LastCommitUrl(i.Url, branch)
		CheckErr(err)
	}
	err = git.Checkout(i.Directory, i.Version)
	CheckErr(err)
	packDeps, err := pacman.GetDeps(i.Pkgbuild)
	CheckErr(err)
	groups := GroupPackages(packDeps)
	Install(nil, groups.PackPackages)
	err = pack.SwapDeps(i.Pkgbuild, packDeps)
	CheckErr(err)
	prnt.Yellow("Staring build: ", i.PackName)
	err = pacman.Build(i.Directory)
	CheckErr(err)
	err = pacman.InstallDir(i.Directory)
	CheckErr(err)
	pack.Update(pack.Package{
		PacmanName:    i.PacmanName,
		PackName:      i.PackName,
		Version:       i.Version,
		DefaultBranch: branch,
	})
	err = system.MvExt(i.Directory, config.PackageCacheDir, ".pkg.tar.zst")
	CheckErr(err)
	git.Clean(i.Directory)
	prnt.Green("Completed build: ", i.PackName)
}
