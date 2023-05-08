package cmd

import (
	"context"
	"fmt"
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

func Get(cmd *cobra.Command, upkgs []string) {
	PrepareForInstallation(upkgs)
	pkgs := SplitPackages(upkgs)
	CheckUnreachablePacmanPackages(pkgs.PacmanPackages)
	CheckUnreachablePackPackages(pkgs.PackPackages)
	InstallPacmanPackages(pkgs.PacmanPackages)
}

// Exis if there is no target packages, prepare cache directories.
func PrepareForInstallation(pkgs []string) {
	if len(pkgs) == 0 {
		os.Exit(0)
	}
	CheckCacheDirExist()
	BluePrint("Installing packages: ", strings.Join(pkgs, " "))
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
	i := EjectInfoFromPackLink(pkg)
	out, err := system.Callf("git clone %s %s", i.HttpsLink, i.Directory)
	if err != nil {
		if !strings.Contains(out, "already exists and is not an empty dir") {
			RedPrint("Unable to reach git for: ", pkg)
			return err
		}
	}
	_, err = os.Stat(i.Directory + "/PKGBUILD")
	if err != nil {
		RedPrint("Unable to find PKGBUILD for: ", pkg)
	}
	return err
}

// Info formed from pack link about all information related to that package.
type PackInfo struct {
	ShortName string
	FullName  string
	Directory string
	Version   string
	HttpsLink string
}

// Eject pack information for provided pack link.
func EjectInfoFromPackLink(pkg string) PackInfo {
	rez := PackInfo{}
	versplt := strings.Split(pkg, "@")
	rez.FullName = versplt[0]
	rez.HttpsLink = "https://" + versplt[0]
	if len(versplt) > 1 {
		rez.Version = versplt[1]
	}
	dashsplt := strings.Split(pkg, "/")
	rez.ShortName = dashsplt[len(dashsplt)-1]
	rez.Directory = cfg.RepoCacheDir + "/" + rez.ShortName
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
		RedPrint("Unable to install pacman packages: ", joined)
		fmt.Println(o)
		os.Exit(1)
	}
	GreenPrint("Pacman packages installed: ", joined)
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

// Install pack package.
func InstallPackPackage(i PackInfo) {
	SetPackageVersion(i)
	// Eject pack dependencies
	// Resolve pack dependencies
	// Swap pack dependencies
	// Install packge with makepkg
	// Put package to pacman cache
	// Clean git or remove git untracked
}

func SetPackageVersion(i PackInfo) {
	if i.Version == `` {
		i.Version = GetDefaultGitBranch(i.Directory)
	}
	o, err := system.Callf("git -C %s checkout %s", i.Directory, i.Version)
	if err != nil {
		if !strings.HasPrefix(o, "Already on ") {
			RedPrint("Unable to set pack version for: ", i.FullName)
			fmt.Println(o)
			os.Exit(1)
		}
	}
}

// Returns default branch for git repository located in some directory.
func GetDefaultGitBranch(dir string) string {
	origin, err := system.Callf("git -C %s remote show", dir)
	CheckErr(err)
	origin = strings.Trim(origin, "\n")
	remoteInfo, err := system.Callf("git -C %s remote show %s", dir, origin)
	CheckErr(err)
	return EjectDefaultGitBranchFromRemoteInfo(remoteInfo)
}

// Function that parses information about git remote and returns it's default
// branch.
func EjectDefaultGitBranchFromRemoteInfo(rawInfo string) string {
	rawInfo = strings.Split(rawInfo, "HEAD branch: ")[1]
	return strings.Split(rawInfo, "\n")[0]
}

//
