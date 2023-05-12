// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"fmnx.io/core/pack/config"
	"fmnx.io/core/pack/git"
	"fmnx.io/core/pack/packdb"
	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"fmnx.io/core/pack/tmpl"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Example: "pack install fmnx.io/core/ainst fmnx.io/core/keks@main",
	Aliases: []string{"i"},
	Short:   tmpl.InstallShort,
	Long:    tmpl.InstallLong,
	Run:     Install,
}

// Cli command installing packages into system.
func Install(cmd *cobra.Command, upkgs []string) {
	pkgs := SplitPackages(upkgs)
	CheckUnreachablePacmanPackages(pkgs.PacmanPackages)
	CheckUnreachablePackPackages(pkgs.PackPackages)
	InstallPacmanPackages(pkgs.PacmanPackages)
	InstallPackPackages(pkgs.PackPackages)
}

type PackageGroups struct {
	PacmanPackages []string
	PackPackages   []string
}

// Split packages into pacman and pack to resolve dependencies differently.
func SplitPackages(pkgs []string) PackageGroups {
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
	deps := GetDependeciesResolvableByPacman()
	var unreachable []string
	for _, pkg := range pkgs {
		if _, ok := deps[pkg]; !ok {
			unreachable = append(unreachable, pkg)
		}
	}
	if len(unreachable) != 0 {
		pkgs := strings.Join(unreachable, " ")
		print.Red("Unable to resolve those pacman packages: ", pkgs)
		os.Exit(1)
	}
}

// Fill struct that shows which packages could be resolved with pacman
// (packages that you can load from pacman servers).
func GetDependeciesResolvableByPacman() map[string]struct{} {
	o, err := system.Call("pacman -Ssq")
	CheckErr(err)
	deps := map[string]struct{}{}
	for _, pkg := range strings.Split(o, "\n") {
		deps[pkg] = struct{}{}
	}
	return deps
}

// Check if some pack packages could not be installed.
func CheckUnreachablePackPackages(pkgs []string) {
	g, _ := errgroup.WithContext(context.Background())
	var unreachable []string
	for _, pkg := range pkgs {
		spkg := pkg
		g.Go(func() error {
			_, err := packdb.Get(spkg, packdb.PACK)
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
		print.Red("Some pack packages are unreachable: ", out)
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

// Install pacman packages.
func InstallPacmanPackages(pkgs []string) {
	uninstalled := CleanAlreadyInstalled(pkgs)
	if len(uninstalled) == 0 {
		return
	}
	joined := strings.Join(uninstalled, " ")
	o, err := system.Callf("sudo pacman --noconfirm -S %s", joined)
	if err != nil {
		print.Red("Unable to get pacman packages: ", joined)
		fmt.Println(o)
		os.Exit(1)
	}
	print.Green("Pacman packages installed: ", joined)
}

// Removes pacman packages that are already installed in the system.
func CleanAlreadyInstalled(pkgs []string) []string {
	var uninstalledPkgs []string
	for _, pkg := range pkgs {
		_, err := system.Callf("pacman -Q %s", pkg)
		if err != nil {
			uninstalledPkgs = append(uninstalledPkgs, pkg)
		}
	}
	return uninstalledPkgs
}

// Checks if packages are not installed and installing them.
func InstallPackPackages(pkgs []string) {
	for _, pkg := range pkgs {
		_, err := packdb.Get(pkg, packdb.PACK)
		if err == nil && !Updating {
			continue
		}
		InstallPackPackage(EjectInfoFromPackLink(pkg))
	}
	if len(pkgs) > 0 {
		pkglist := strings.Join(pkgs, " ")
		if !Updating {
			print.Green("Installed: ", pkglist)
		}
	}
}

// Install pack package.
func InstallPackPackage(i PackInfo) {
	err := git.Clean(i.Directory)
	CheckErr(err)
	err = git.Pull(i.Directory)
	CheckErr(err)
	branch, err := git.DefaultBranch(i.Directory)
	CheckErr(err)
	if i.Version == `` {
		i.Version, err = git.LastCommitUrl(i.Url, branch)
		CheckErr(err)
	}
	err = git.Checkout(i.Directory, i.Version)
	CheckErr(err)
	packDeps := EjectPackDependencies(i.Pkgbuild)
	Install(nil, packDeps)
	SwapPackDependencies(i.Pkgbuild, packDeps)
	InstallPackageWithMakepkg(i)
	packdb.Update(packdb.Package{
		PacmanName: i.PacmanName,
		PackName:   i.PackName,
		Version:    i.Version,
		Branch:     branch,
	})
	CachePackage(i.Directory)
	git.Clean(i.Directory)
}

// Get dependencies and make dependencies related to pack from PKGBUILD file.
func EjectPackDependencies(pkgbuild string) []string {
	deps, err := system.GetShellList(pkgbuild, "depends")
	CheckErr(err)
	makedeps, err := system.GetShellList(pkgbuild, "makedepends")
	CheckErr(err)
	alldeps := append(deps, makedeps...)
	groups := SplitPackages(alldeps)
	return groups.PackPackages
}

// Temporarily swap pack dependencies in PKGBUILD to pacman package name for
// installation process.
func SwapPackDependencies(pkgbuild string, deps []string) {
	b, err := os.ReadFile(pkgbuild)
	CheckErr(err)
	var rez = string(b)
	for _, dep := range deps {
		dashsplt := strings.Split(dep, "/")
		shortname := dashsplt[len(dashsplt)-1]
		rez = strings.ReplaceAll(rez, dep, shortname)
	}
	err = os.WriteFile(pkgbuild, []byte(rez), 0o600)
	CheckErr(err)
}

// Install package with makepkg.
func InstallPackageWithMakepkg(i PackInfo) {
	CheckErr(os.Chdir(i.Directory))
	print.Yellow("Building package: ", i.PackName)
	out, err := system.Call("makepkg -sfi --noconfirm")
	if err != nil {
		print.Red("Unable to build and install package: ", i.PackName)
		fmt.Println(out)
		os.Exit(1)
	}
}

// Move prepared .pkg.tar.zst package into pacman cache.
func CachePackage(dir string) {
	if !config.RemoveBuiltPackages {
		const command = "sudo mv %s/*.pkg.tar.zst %s"
		_, err := system.Callf(command, dir, config.PackageCacheDir)
		CheckErr(err)
	}
}
