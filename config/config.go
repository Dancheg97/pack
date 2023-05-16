// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package config

// Project configuration.

import (
	"fmt"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

// This variables are automatically initialized in init().
var (
	RemoveGitRepos      bool
	RemoveBuiltPackages bool
	DebugMode           bool
	DisablePrettyPrint  bool
	CacheDir            string
	PackageCacheDir     string
	MapFile             string
	LockFile            string

	homedir string
	cfgfile string
	cfg     config
)

// Configuration variables.
type config struct {
	RemoveGitRepos      bool   `yaml:"remove-git-repo"`
	RemoveBuiltPackages bool   `yaml:"remove-built-packages"`
	DebugMode           bool   `yaml:"debug-mode"`
	DisablePrettyPrint  bool   `yaml:"pack-disable-prettyprint"`
	PackCacheDir        string `yaml:"pack-cache-dir"`
	PackageCacheDir     string `yaml:"package-cache-dir"`
	MapFile             string `yaml:"map-file"`
	LockFile            string `yaml:"lock-file"`
}

// Initialize runtime configuration variables.
func init() {
	usr, err := user.Current()
	checkErr(err)
	homedir = usr.HomeDir
	cfgfile = homedir + "/.pack/config.yml"

	_, err = os.Stat(cfgfile)
	if err != nil {
		SetDefaults()
		Save()
		return
	}

	b, err := os.ReadFile(cfgfile)
	checkErr(err)
	err = yaml.Unmarshal(b, &cfg)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("configuration error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

// SetDefaults configuration to default values and set save config file.
func SetDefaults() {
	cfg.RemoveGitRepos = false
	cfg.RemoveBuiltPackages = false
	cfg.DebugMode = false
	cfg.DisablePrettyPrint = false
	cfg.PackCacheDir = homedir + "/.pack"
	cfg.PackageCacheDir = "/var/cache/pacman/pkg"
	cfg.MapFile = homedir + "/.pack/mapping.json"
	cfg.LockFile = "/tmp/pack.lock"

	RemoveGitRepos = cfg.RemoveGitRepos
	RemoveBuiltPackages = cfg.RemoveBuiltPackages
	DebugMode = cfg.DebugMode
	DisablePrettyPrint = cfg.DisablePrettyPrint
	CacheDir = cfg.PackCacheDir
	PackageCacheDir = cfg.PackageCacheDir
	MapFile = cfg.MapFile
	LockFile = cfg.LockFile
}

// Save configuration with all new variables.
func Save() {
	b, err := yaml.Marshal(&config{
		RemoveGitRepos:      RemoveGitRepos,
		RemoveBuiltPackages: RemoveBuiltPackages,
		DebugMode:           DebugMode,
		DisablePrettyPrint:  DisablePrettyPrint,
		PackCacheDir:        PackageCacheDir,
		PackageCacheDir:     PackageCacheDir,
		MapFile:             MapFile,
		LockFile:            LockFile,
	})
	checkErr(err)
	err = os.WriteFile(cfgfile, b, 0o600)
	checkErr(err)
}
