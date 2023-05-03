package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fmnx.io/dev/pack/input"
	"fmnx.io/dev/pack/system"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:     "get",
	Example: "pack get fmnx.io/dev/ainst fmnx.io/dev/keks@main",
	Aliases: []string{"g"},
	Short:   "📥 get and install new packages",
	Long: `📥 get and install new packages

You can mix pacman and pack packages, provoding names and git links. If you 
need to specify version, you can provide it after @ symbol.

Examples:
pack get fmnx.io/dev/aist@v0.21
pack get git.xmpl.sh/pkg
pack get fmnx.io/dev/ainst github.com/exm/pkg@v1.23 nano`,
	Run: Get,
}

type PkgInfo struct {
	FullName  string
	ShortName string
	HttpsLink string
	Version   string
	IsPacman  bool
}

type PackYml struct {
	RunDeps      []string          `yaml:"run-deps"`
	BuildDeps    []string          `yaml:"build-deps"`
	BuildScripts []string          `yaml:"scripts"`
	PackMap      map[string]string `yaml:"mapping"`
}

type PkgbuildInfo struct {
	Exists bool
	Deps   []string
}

type PackMap map[string]string

var (
	depsTmpl     = "depends=(\n  \"%s\"\n)"
	makedepsTmpl = "makedepends=(\n  \"%s\"\n)"
	pkgbuildTmpl = `# PKGBUILD generated by pack.
# More info at: https://fmnx.io/dev/pack

pkgname="%s"
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
)

func Get(cmd *cobra.Command, pkgs []string) {
	if len(pkgs) != 0 {
		err := system.PrepareDir(cfg.RepoCacheDir)
		CheckErr(err)
		err = system.PrepareDir(cfg.PackageCacheDir)
		CheckErr(err)
	}
	for _, pkg := range pkgs {
		info := FormPkgInfoFromLink(pkg)
		if CheckIfInstalled(info) {
			YellowPrint("Package installed, skipping: ", info.FullName)
			continue
		}
		if info.IsPacman {
			BluePrint("Installing package with pacman: ", info.FullName)
			out, err := system.Call("sudo pacman --noconfirm -Sy " + pkg)
			if err != nil {
				if strings.Contains(out, "target not found") {
					RedPrint("Pacman package not found: ", pkg)
					fmt.Printf("Use pack for aur.archlinux.org/%s? [Y/n]\n", pkg)
					confirmed := input.AskForConfirmation()
					if confirmed || cfg.AllowAUR {
						Get(cmd, []string{"aur.archlinux.org/" + pkg})
						continue
					}
					RedPrint("Unable to install package: ", pkg)
					lf.Unlock()
					os.Exit(1)
				}
				fmt.Println("Pacman output: ", out)
			}
			CheckErr(err)
			GreenPrint("Installed: ", info.FullName+" - OK")
			continue
		}
		PrepareRepo(info)
		pkgbuildInfo := ReadPkgbuildInfo()
		if !pkgbuildInfo.Exists {
			packyml := ReadPackYml()
			Get(cmd, append(packyml.RunDeps, packyml.BuildDeps...))
			BuildPackage(info, packyml)
			GeneratePkgbuild(info, packyml)
		} else {
			Get(cmd, pkgbuildInfo.Deps)
		}
		InstallPackage()
		CachePkgTarZst()
		AddToMapping(info)
		CleanGitDir(info.ShortName)
		GreenPrint("Package installed: ", info.FullName)
	}
}

func FormPkgInfoFromLink(pkg string) PkgInfo {
	if !strings.Contains(pkg, ".") {
		return PkgInfo{
			FullName:  pkg,
			ShortName: pkg,
			IsPacman:  true,
		}
	}
	fullName := strings.Split(pkg, "@")[0]
	httpslink := "https://" + fullName
	split := strings.Split(httpslink, "/")
	shortname := split[len(split)-1]
	version := "latest"
	if len(strings.Split(pkg, "@")) != 1 {
		version = strings.Split(pkg, "@")[1]
	}
	return PkgInfo{
		FullName:  fullName,
		ShortName: shortname,
		HttpsLink: httpslink,
		Version:   version,
		IsPacman:  false,
	}
}

func CheckIfInstalled(i PkgInfo) bool {
	mp := ReadMapping()
	if _, packageExists := mp[i.FullName]; packageExists {
		return true
	}
	_, err := system.Call("pacman -Q " + i.ShortName)
	return err == nil
}

func ReadMapping() PackMap {
	_, err := os.Stat(cfg.MapFile)
	if err != nil {
		system.AppendToFile(cfg.MapFile, "{}")
		return PackMap{}
	}
	b, err := os.ReadFile(cfg.MapFile)
	CheckErr(err)
	var mapping PackMap
	err = json.Unmarshal(b, &mapping)
	CheckErr(err)
	return mapping
}

func PrepareRepo(i PkgInfo) {
	CheckErr(os.Chdir(cfg.RepoCacheDir))
	BluePrint("Cloning repository: ", i.HttpsLink)
	out, err := system.SystemCallf("git clone %s", i.HttpsLink)
	CheckErr(os.Chdir(i.ShortName))
	if strings.Contains(out, "already exists and is not an empty directory") {
		YellowPrint("Repository exists: ", "pulling changes...")
		ExecuteCheck("git pull")
		GreenPrint("Changes pulled: ", "success")
		err = nil
	}
	CheckErr(err)
	if i.Version != `latest` {
		BluePrint("Switching repo to version: ", i.Version)
		ExecuteCheck("git checkout " + i.Version)
	}
}

func ReadPkgbuildInfo() PkgbuildInfo {
	_, err := os.Stat("PKGBUILD")
	if err != nil {
		return PkgbuildInfo{Exists: false}
	}
	deps, err := system.EjectShList("PKGBUILD", "depends")
	CheckErr(err)
	buildeps, err := system.EjectShList("PKGBUILD", "makedepends")
	CheckErr(err)
	YellowPrint("Installing with: ", "PKGBUILD")
	return PkgbuildInfo{
		Exists: true,
		Deps:   append(deps, buildeps...),
	}
}

func ReadPackYml() PackYml {
	b, err := os.ReadFile(".pack.yml")
	CheckErr(err)
	var packyml PackYml
	err = yaml.Unmarshal(b, &packyml)
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
		_, err := system.Call("pacman -Q " + pkg)
		if err != nil {
			BluePrint("Gettings dependecy package: ", pkg)
			ExecuteCheck("sudo pacman --noconfirm -Sy " + pkg)
		}
	}
}

func BuildPackage(i PkgInfo, y PackYml) {
	CheckErr(os.Chdir(cfg.RepoCacheDir + "/" + i.ShortName))
	for _, script := range y.BuildScripts {
		BluePrint("Executing build script: ", script)
		ExecuteCheck(script)
	}
}

func GeneratePkgbuild(i PkgInfo, y PackYml) {
	deps := fmt.Sprintf(depsTmpl, strings.Join(y.RunDeps, "\"\n  \""))
	if len(y.RunDeps) == 0 {
		deps = ""
	}
	makedeps := fmt.Sprintf(makedepsTmpl, strings.Join(y.BuildDeps, "\"\n  \""))
	if len(makedeps) == 0 {
		makedeps = ""
	}
	var installScripts []string
	for src, dst := range y.PackMap {
		installScripts = append(installScripts, FormatInstallSrc(src, dst))
	}
	install := strings.Join(installScripts, "\n  ")
	pkgb := fmt.Sprintf(pkgbuildTmpl, i.ShortName, i.Version, i.HttpsLink, deps, makedeps, install)
	CheckErr(system.WriteFile("PKGBUILD", pkgb))
}

func FormatInstallSrc(src string, dst string) string {
	i, err := os.Stat(src)
	CheckErr(err)
	if i.IsDir() {
		return fmt.Sprintf(`cd %s && find . -type f -exec install -Dm755 {} "${pkgdir}%s/{}" \; && cd $srcdir/..`, src, dst)
	}
	return fmt.Sprintf(`install -Dm755 %s "${pkgdir}%s"`, src, dst)
}

func InstallPackage() {
	BluePrint("Building and installing package: ", "makepkg -sfri")
	ExecuteCheck("makepkg --noconfirm -sfri")
}

func AddToMapping(i PkgInfo) {
	mp := ReadMapping()
	mp[i.FullName] = i.ShortName
	WriteMapping(mp)
}

func CleanGitDir(repo string) {
	if !cfg.RemoveGitRepos {
		ExecuteCheck("git clean -fd")
		ExecuteCheck("git reset --hard")
		return
	}
	ExecuteCheck("sudo rm -rf " + cfg.RepoCacheDir + "/" + repo)
}

func CachePkgTarZst() {
	if !cfg.RemoveBuiltPackages {
		ExecuteCheck("sudo mv *.pkg.tar.zst " + cfg.PackageCacheDir)
	}
}
