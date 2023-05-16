// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"context"
	"io"
	"net/http"
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
	LoadLinkPackages(groups.NetworkLinks)
	CopyFilePackages(groups.LocalFiles)
	CheckUnreachablePacmanPackages(groups.PacmanPackages)
	CheckErr(pacman.Install(groups.PacmanPackages))
	InstallPackPackages(groups.PackPackages)
}

// Packages splitted into groups depending on type of installation.
type PackageGroups struct {
	PacmanPackages []string
	PackPackages   []string
	NetworkLinks   []string
	LocalFiles     []string
}

// Existing package groups.
const (
	PACMAN_PACKAGE = 0
	PACK_GIT_REPO  = 1
	LOCAL_FILE     = 2
	HTTP_FILE_LINK = 3
)

// Split packages into groups.
func GroupPackages(pkgs []string) PackageGroups {
	var groups PackageGroups
	for _, pkg := range pkgs {
		group := GetPackageGroup(pkg)
		switch group {
		case PACMAN_PACKAGE:
			groups.PacmanPackages = append(groups.PacmanPackages, pkg)
		case PACK_GIT_REPO:
			groups.PackPackages = append(groups.PackPackages, pkg)
		case HTTP_FILE_LINK:
			groups.NetworkLinks = append(groups.NetworkLinks, pkg)
		case LOCAL_FILE:
			groups.LocalFiles = append(groups.LocalFiles, pkg)
		}
	}
	return groups
}

// Get package group based on it's name.
func GetPackageGroup(pkg string) int {
	hasSuffix := strings.HasPrefix(pkg, "http")
	hasPrefix := strings.HasSuffix(pkg, ".pkg.tar.zst")
	if hasSuffix && hasPrefix {
		return HTTP_FILE_LINK
	}
	if hasPrefix {
		return LOCAL_FILE
	}
	if strings.Contains(pkg, "/") {
		return PACK_GIT_REPO
	}
	return PACMAN_PACKAGE
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

// Load packages that are defined as network links to pack cache directory.
func LoadLinkPackages(links []string) {
	for _, link := range links {
		r, err := http.Get(link)
		CheckErr(err)
		defer r.Body.Close()
		splt := strings.Split(link, "/")
		pkgName := splt[len(splt)-1]
		f, err := os.Create(config.CacheDir + "/" + pkgName)
		CheckErr(err)
		_, err = io.Copy(f, r.Body)
		prnt.Green("Loaded file: ", pkgName)
		CheckErr(err)
	}
}

// Copy packages that are meant to be installed as files to cache dir.
func CopyFilePackages(pkgs []string) {
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		file := splt[len(splt)-1]
		_, err := system.Callf("sudo cp %s %s/%s", pkg, config.CacheDir, file)
		CheckErr(err)
	}
}

// Install packages in pack cache dir and move them to pacman cache dir.
func InstallFromPackCache() {
	err := pacman.InstallDir(config.CacheDir)
	CheckErr(err)
	pkgs, err := system.LsExt(config.CacheDir, ".pkg.tar.zst")
	CheckErr(err)
	prnt.Green("File packages installed: ", strings.Join(pkgs, " "))
	err = system.MvExt(config.CacheDir, config.PackageCacheDir, ".pkg.tar.zst")
	CheckErr(err)
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
	i.Version = BuildDirectory(i.Directory, i.Version)
	err = pacman.InstallDir(i.Directory)
	CheckErr(err)
	prnt.Green("Installed package: ", i.PackName+"@"+i.Version)
	if !config.RemoveBuiltPackages {
		err = system.MvExt(i.Directory, config.PackageCacheDir, ".pkg.tar.zst")
		CheckErr(err)
	}
	if !config.RemoveGitRepos {
		err = git.Clean(i.Directory)
		CheckErr(err)
	}
	if config.RemoveGitRepos {
		err = os.RemoveAll(i.Directory)
		CheckErr(err)
	}
	pack.Update(pack.Package{
		PacmanName:    i.PacmanName,
		PackName:      i.PackName,
		Version:       i.Version,
		DefaultBranch: branch,
	})
}
