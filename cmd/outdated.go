package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(outdatedCmd)
}

var outdatedCmd = &cobra.Command{
	Use:     "outdated",
	Aliases: []string{"out", "o"},
	Short:   "📌 show outdated packages",
	Long: `📌 show outdated packages

This command will make a call to pacman servers and collect information about
all remote repos for packages installed with pack. Then it will print a list
of packages that require update displaying current and new available version.
`,
	Run: Outdated,
}

// Cli command listing installed packages and their status.
func Outdated(cmd *cobra.Command, args []string) {
	pacmanOutdated := GetPacmanOutdated()
	packoutdated := GetPackOutdated()
	allOutdated := append(pacmanOutdated, packoutdated...)
	for i, op := range allOutdated {
		print.Custom([]print.ColoredMessage{
			{
				Message: fmt.Sprintf("%d - %s ", i+1, op.Name),
				Color:   print.WHITE,
			},
			{
				Message: op.CurrVersion + " ",
				Color:   print.YELLOW,
			},
			{
				Message: op.NewVersion,
				Color:   print.BLUE,
			},
		})
	}
}

type OutdatedPackage struct {
	Name        string
	CurrVersion string
	NewVersion  string
}

// Get outdated packages and their versions.
func GetPacmanOutdated() []OutdatedPackage {
	links := GetUpdateLinks()
	var out []OutdatedPackage
	for _, link := range links {
		linkSplit := strings.Split(link, "/")
		file := linkSplit[len(linkSplit)-1]
		fileSplit := strings.Split(file, "-")
		ver := GetCurrentVersion(fileSplit[0])
		out = append(out, OutdatedPackage{
			Name:        fileSplit[0],
			CurrVersion: ver,
			NewVersion:  fileSplit[1],
		})
	}
	return out
}

// Get found update links for pacman packges.
func GetUpdateLinks() []string {
	o, err := system.Call("sudo pacman -Syup")
	if err != nil {
		print.Red("Unable to connect to pacman servers: ", "network error")
		os.Exit(1)
	}
	if !strings.Contains(o, "https://") {
		return nil
	}
	splt := strings.Split(o, "downloading...\n")
	pkgsLinksString := strings.Trim(splt[len(splt)-1], "\n")
	return strings.Split(pkgsLinksString, "\n")
}

// Get current package version.
func GetCurrentVersion(pkg string) string {
	o, err := system.Callf("pacman -Q %s", pkg)
	CheckErr(err)
	verAndRel := strings.Split(o, " ")[1]
	return strings.Trim(strings.Split(verAndRel, "-")[0], "\n")
}

// Get pack outdated packages.
func GetPackOutdated() []OutdatedPackage {
	mp := ReadMapping()
	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var rez []OutdatedPackage
	for pack, pacman := range mp {
		spack, spacman := pack, pacman
		var curr, last, branch string
		g.Go(func() error {
			branch, curr = GetPackVerInfo(spacman)
			last = GetRemoteVersionForBranch("https://"+spack, branch)
			if curr == last {
				return nil
			}
			mu.Lock()
			rez = append(rez, OutdatedPackage{
				Name:        spack,
				CurrVersion: curr,
				NewVersion:  last,
			})
			mu.Unlock()
			return nil
		})
	}
	g.Wait()
	return rez
}

// Get branch and git commit hash of pack package. (branch, hash)
func GetPackVerInfo(pkg string) (string, string) {
	o, err := system.Callf("pacman -Qi %s", pkg)
	CheckErr(err)
	const versionField = "Version         : "
	rawver := strings.Split(strings.Split(o, versionField)[1], "\n")[0]
	rawverSplit := strings.Split(rawver, ".")
	branch := rawverSplit[0]
	hash := strings.Split(rawverSplit[1], "-")[0]
	return branch, hash
}

// Get remote version for specific branch of git repository.
func GetRemoteVersionForBranch(link string, branch string) string {
	o, err := system.Callf("git ls-remote -h %s", link)
	if err != nil {
		return "unable to connect"
	}
	refs := strings.Split(strings.Trim(o, "\n"), "\n")
	for _, ref := range refs {
		if strings.HasSuffix(ref, branch) {
			return strings.Split(ref, "	")[0]
		}
	}
	return "unable to find branch in remote repo"
}
