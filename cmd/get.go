package cmd

import (
	"fmt"
	"strings"

	"fmnx.io/dev/pack/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type PackageInfo struct {
	Name    string
	Link    string
	Version string
	Owner   string
}

type PackYml struct {
	RunDeps     []string          `yaml:"run-deps"`
	BuildDeps   []string          `yaml:"build-deps"`
	BuildScript []string          `yaml:"build-script"`
	PackMap     map[string]string `yaml:"pack-map"`
}

const (
	CacheDir = `~/.pack-cache`
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "📥 insatll new packages",
	Run:   Get,
}

func Get(cmd *cobra.Command, pkgs []string) {
	err := core.SystemCallf("mkdir -p %s", CacheDir)
	CheckErr(err)

	for _, pkg := range pkgs {
		info := EjectInfo(pkg)
		PrepareRepo(info)
		SwitchToVersion(info)
		packyml := ReadPackYml(info)
		allDeps := append(packyml.RunDeps, packyml.BuildDeps...)
		pacmanPkgs, packPkgs := SplitDependencies(allDeps)
		ResolvePacmanDeps(pacmanPkgs)
		Get(cmd, packPkgs)
		
	}
}

func EjectInfo(pkg string) PackageInfo {
	link := "https://" + strings.Split(pkg, "@")[0]
	split := strings.Split(link, "/")
	name := split[len(split)-1]
	owner := strings.Join(split[0:len(split)-1], "/")
	version := ""
	if len(strings.Split(pkg, "@")) == 1 {
		branch, err := GetDefaultBranch(pkg)
		CheckErr(err)
		version = branch
	} else {
		version = strings.Split(pkg, "@")[1]
	}
	return PackageInfo{
		Name:    name,
		Link:    link,
		Version: version,
		Owner:   owner,
	}
}

func GetDefaultBranch(pkg string) (string, error) {
	pkgLink := "https://" + strings.Split(pkg, "@")[0]
	out, err := core.SystemCallOutf("git remote show %s | sed -n '/HEAD branch/s/.*: //p'", pkgLink)
	CheckErr(err)
	return strings.Trim(out, "\n"), nil
}

func PrepareRepo(i PackageInfo) {
	err := core.SystemCallf("git clone %s %s/%s", i.Link, CacheDir, i.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			CheckErr(err)
		}
		fmt.Println("pulling changes")
		err = core.SystemCallf("git -C %s/%s pull ", CacheDir, i.Name)
		CheckErr(err)
	}
}

func SwitchToVersion(i PackageInfo) {
	err := core.SystemCallf("git -C %s/%s checkout %s", CacheDir, i.Name, i.Version)
	CheckErr(err)
}

func ReadPackYml(i PackageInfo) PackYml {
	content, err := core.SystemCallOutf("cat %s/%s/pack.yml", CacheDir, i.Name)
	CheckErr(err)
	var packyml PackYml
	err = yaml.Unmarshal([]byte(content), &packyml)
	CheckErr(err)
	return packyml
}

func SplitDependencies(deps []string) ([]string, []string) {
	var pacmandeps []string
	var packdeps []string
	for _, dep := range deps {
		if strings.Contains(dep, ".") {
			packdeps = append(packdeps, dep)
			continue
		}
		pacmandeps = append(pacmandeps, dep)
	}
	return pacmandeps, packdeps
}

func ResolvePacmanDeps(deps []string) {
	err := core.SystemCall("sudo pacman --noconfirm -Sy " + strings.Join(deps, " "))
	CheckErr(err)
}

// func ResolveDependecies(i PackageInfo) error {

// }

// func GeneratePkgbuild(i PackageInfo) {

// }

// func InstallPackage() {

// }
