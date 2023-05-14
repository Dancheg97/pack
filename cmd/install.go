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

	"fmnx.su/core/pack/git"
	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/prnt"
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

// Checks if packages are not installed and installing them.
func InstallPackPackages(pkgs []string) {
	g, _ := errgroup.WithContext(context.Background())
	for _, pkg := range pkgs {
		spkg := pkg
		g.Go(func() error {
			InstallPackPackage(spkg)
			return nil
		})
	}
	g.Wait()
}

// Install pack package.
func InstallPackPackage(pkg string) {
	i := pack.GetPackInfo(pkg)
	_, err := pack.Get(pkg)
	if err == nil {
		return
	}
	err = git.Clone(i.GitUrl, i.Directory)
	CheckErr(err)
	prnt.Blue("Cloned: ", i.GitUrl)
	branch, err := git.DefaultBranch(i.Directory)
	CheckErr(err)
	i.Version = BuildDirectory(i.Directory, i.Version, true)
	pack.Update(pack.Package{
		PacmanName:    i.PacmanName,
		PackName:      i.PackName,
		Version:       i.Version,
		DefaultBranch: branch,
	})
}
