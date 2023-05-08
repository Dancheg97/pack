package cmd

import (
	"os"
	"strings"

	"fmnx.io/dev/pack/system"
	"github.com/spf13/cobra"
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

var (
	// Dependencies that could be resolved with pacman (packages that you can
	// load from pacman servers).
	ResolvableByPacman = map[string]struct{}{}

	PacmanDependencies = map[string]struct{}{}
	PackDependencies   = map[string]struct{}{}
	InstalledPkgs      []string
)

func Get(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) == 0 {
		return
	}
	CheckCacheDirExist()
	BluePrint("Installing packages: ", strings.Join(pkgs, " "))
	FillDependeciesResolvableByPacman()
	splittedPkgs := SplitPackages(pkgs)
	CheckUnreachablePacmanPackages(splittedPkgs.PacmanPackages)

}

// Prepare cache directories for package repositories.
func CheckCacheDirExist() {
	err := system.PrepareDir(cfg.RepoCacheDir)
	CheckErr(err)
	err = system.PrepareDir(cfg.PackageCacheDir)
	CheckErr(err)
}

// Fill struct that shows which packages could be resolved with pacman.
func FillDependeciesResolvableByPacman() {
	o, err := system.Call("pacman -Ssq")
	CheckErr(err)
	for _, pkg := range strings.Split(o, "\n") {
		ResolvableByPacman[pkg] = struct{}{}
	}
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
	var unreachable []string
	for _, pkg := range pkgs {
		if _, ok := ResolvableByPacman[pkg]; !ok {
			unreachable = append(unreachable, pkg)
		}
	}
	if len(unreachable) != 0 {
		pkgs := strings.Join(unreachable, " ")
		RedPrint("Unable to resolve those pacman packages: ", pkgs)
		os.Exit(1)
	}
}

// Check if some pack packages could not be installed.
func CheckUnreachablePackPackages(pkgs []string) {

}

// Validate pack package to be reachable via network.
func CheckPackPackage(pkg string) {

}
