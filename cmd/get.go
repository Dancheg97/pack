package cmd

import (
	"fmt"
	"os"
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
	RunDeps      []string          `yaml:"run-deps"`
	BuildDeps    []string          `yaml:"build-deps"`
	BuildScripts []string          `yaml:"build-scripts"`
	PackMap      map[string]string `yaml:"pack-map"`
}

const (
	cacheDir = `~/.pack-cache`
	pkgbuild = `pkgname="%s"
pkgver="%s"
pkgrel="1"
arch=('i686' 'pentium4' 'x86_64' 'arm' 'armv7h' 'armv6h' 'aarch64' 'riscv64')
url="%s"
%s
%s
package() {
  cd ..
  %s
}`
	writeFileCmd = `tee -a %s << END
%s
END
`
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "📥 install new packages",
	Run:   Get,
}

func Get(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) != 0 {
		err := core.SystemCallf("mkdir -p %s", cacheDir)
		CheckErr(err)
	}
	for _, pkg := range pkgs {
		info := EjectInfo(pkg)
		PrepareRepo(info)
		packyml := ReadPackYml(info)
		allDeps := append(packyml.RunDeps, packyml.BuildDeps...)
		pacmanPkgs, packPkgs := SplitDependencies(allDeps)
		ResolvePacmanDeps(pacmanPkgs)
		Get(cmd, packPkgs)
		BuildPackage(info, packyml)
		GeneratePkgbuild(info, packyml)
	}
}

func EjectInfo(pkg string) PackageInfo {
	link := "https://" + strings.Split(pkg, "@")[0]
	split := strings.Split(link, "/")
	name := split[len(split)-1]
	owner := strings.Join(split[0:len(split)-1], "/")
	version := ""
	if len(strings.Split(pkg, "@")) == 1 {
		branch := GetDefaultBranch(pkg)
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

func GetDefaultBranch(pkg string) string {
	pkgLink := "https://" + strings.Split(pkg, "@")[0]
	out, err := core.SystemCallOutf("git remote show %s | sed -n '/HEAD branch/s/.*: //p'", pkgLink)
	CheckErr(err)
	return strings.Trim(out, "\n")
}

func PrepareRepo(i PackageInfo) {
	CheckErr(os.Chdir(cacheDir))
	err := core.SystemCallf("git clone %s", i.Link, cacheDir, i.Name)
	CheckErr(os.Chdir(cacheDir + "/" + i.Name))
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			CheckErr(err)
		}
		fmt.Println("pulling changes")
		err = core.SystemCallf("git -C %s/%s pull ", cacheDir, i.Name)
		CheckErr(err)
	}
	err = core.SystemCallf("git checkout %s", cacheDir, i.Name, i.Version)
	CheckErr(err)
}

func ReadPackYml(i PackageInfo) PackYml {
	content, err := core.SystemCallOutf("cat %s/%s/pack.yml", cacheDir, i.Name)
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

func ResolvePacmanDeps(pkgs []string) {
	for _, pkg := range pkgs {
		_, err := core.SystemCallOut("pacman -Q " + pkg)
		if err != nil {
			err := core.SystemCall("sudo pacman --noconfirm -Sy " + pkg)
			CheckErr(err)
		}
	}
}

func BuildPackage(i PackageInfo, y PackYml) {
	CheckErr(os.Chdir(cacheDir + "/" + i.Name))
	for _, script := range y.BuildScripts {
		CheckErr(core.SystemCall(script))
	}
}

func GeneratePkgbuild(i PackageInfo, y PackYml) {
	CheckErr(os.Chdir(cacheDir + "/" + i.Name))
	deps := "depends=(\n  \"" + strings.Join(y.RunDeps, "\"\n  \"") + "\"\n)\n"
	if len(y.RunDeps) == 0 {
		deps = ""
	}
	makedeps := "makedepends=(\n  \"" + strings.Join(y.BuildDeps, "\"\n  \"") + "\"\n)\n"
	if len(makedeps) == 0 {
		makedeps = ""
	}
	var installScripts []string
	for src, dst := range y.PackMap {
		installScripts = append(installScripts, FormatInstallSrc(src, dst))
	}
	install := strings.Join(installScripts, "\n  ")
	pkgb := fmt.Sprintf(pkgbuild, i.Name, i.Version, i.Link, deps, makedeps, install)
	err := core.SystemCallf(writeFileCmd, "PKGBUILD", pkgb)
	CheckErr(err)
}

func FormatInstallSrc(src string, dst string) string {
	filetype, err := core.SystemCallOutf("stat -c %%F %s", src)
	CheckErr(err)
	if filetype == "directory" {
		return fmt.Sprintf(`cd %s && find . -type f -exec install -Dm755 {} "${pkgdir}%s/{}" \; && cd $srcdir && cd ..`, src, dst)
	}
	return fmt.Sprintf(`install -Dm755 %s "${pkgdir}%s"`, src, dst)
}

// func InstallPackage() {

// }
