package cmd

import (
	"context"
	"os"
	"strings"

	"fmnx.io/dev/pack/system"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:     "install",
	Example: "pack install fmnx.io/dev/ainst fmnx.io/dev/keks@main",
	Aliases: []string{"i", "insall", "get", "g"},
	Short:   "📥 install new packages",
	Long: `📥 install new packages

You can mix pacman and pack packages, provoding names and git links. If you
need to specify version, you can provide it after @ symbol.

Examples:
pack install fmnx.io/dev/aist@v0.21
pack install fmnx.io/dev/ainst github.com/exm/pkg@v1.23 nano`,
	Run: Get,
}

func Get(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) == 0 {
		return
	}
	CheckCacheDirExist()
	BluePrint("Installing packages: ", strings.Join(pkgs, " "))
	splittedPkgs := SplitPackages(pkgs)
	CheckUnreachablePacmanPackages(splittedPkgs.PacmanPackages)
	CheckUnreachablePackPackages(splittedPkgs.PackPackages)
}

// Prepare cache directories for package repositories.
func CheckCacheDirExist() {
	err := system.PrepareDir(cfg.RepoCacheDir)
	CheckErr(err)
	err = system.PrepareDir(cfg.PackageCacheDir)
	CheckErr(err)
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
		RedPrint("Unable to resolve those pacman packages: ", pkgs)
		os.Exit(1)
	}
}

// Fill struct that shows which packages could be resolved with pacman (packages
// that you can load from pacman servers).
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
		syncpkg := pkg
		g.Go(func() error {
			err := CheckPackPackage(syncpkg)
			if err != nil {
				unreachable = append(unreachable, syncpkg)
			}
			return err
		})
	}
	err := g.Wait()
	if err != nil {
		out := strings.Join(unreachable, " ")
		RedPrint("Some pack packages are unreachable: ", out)
		os.Exit(1)
	}
}

// Validate pack package to be reachable via network.
func CheckPackPackage(pkg string) error {
	splt := strings.Split(pkg, "/")
	dir := cfg.RepoCacheDir + "/" + splt[len(splt)-1]
	out, err := system.Callf("git clone https://%s %s", pkg, dir)
	if err != nil {
		if !strings.Contains(out, "already exists and is not an empty dir") {
			RedPrint("Unable to reach git for: ", pkg)
			return err
		}
	}
	_, err = os.Stat(dir + "/PKGBUILD")
	if err != nil {
		RedPrint("Unable to find PKGBUILD for: ", pkg)
	}
	return err
}
