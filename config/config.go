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
	RemoveBuildDeps bool
	SyncPackages    bool
	RemoveGitRepos  bool
	CachePackages   bool
	VerboseOutput   bool
	PrettyPrint     bool
	RepoCacheDir    string
	PackageCacheDir string
	LogFile         string
	MapFile         string
	LockFile        string
)

var (
	homedir string
	cfgfile string
	cfg     config
)

// Configuration variables.
type config struct {
	RemoveBuildDeps bool   `yaml:"remove-build-deps"`
	SyncPackages    bool   `yaml:"sync-packages"`
	RemoveGitRepos  bool   `yaml:"remove-git-repos"`
	CachePackages   bool   `yaml:"cache-packages"`
	VerboseOutput   bool   `yaml:"verbose-output"`
	PrettyPrint     bool   `yaml:"pretty-print"`
	RepoCacheDir    string `yaml:"repo-cache-dir"`
	PackageCacheDir string `yaml:"package-cache-dir"`
	LogFile         string `yaml:"log-file"`
	MapFile         string `yaml:"map-file"`
	LockFile        string `yaml:"lock-file"`
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
	cfg.RemoveBuildDeps = false
	cfg.SyncPackages = false
	cfg.RemoveGitRepos = false
	cfg.CachePackages = true
	cfg.VerboseOutput = false
	cfg.PrettyPrint = true
	cfg.RepoCacheDir = homedir + "/.pack"
	cfg.PackageCacheDir = "/var/cache/pacman/pkg"
	cfg.LogFile = "/tmp/pack.log"
	cfg.MapFile = homedir + "/.pack/mapping.json"
	cfg.LockFile = "/tmp/pack.lock"

	RemoveBuildDeps = cfg.RemoveBuildDeps
	SyncPackages = cfg.SyncPackages
	RemoveGitRepos = cfg.RemoveGitRepos
	CachePackages = cfg.CachePackages
	VerboseOutput = cfg.VerboseOutput
	PrettyPrint = cfg.PrettyPrint
	RepoCacheDir = cfg.RepoCacheDir
	PackageCacheDir = cfg.PackageCacheDir
	LogFile = cfg.LogFile
	MapFile = cfg.MapFile
	LockFile = cfg.LockFile
}

// Save configuration with all new variables.
func Save() {
	b, err := yaml.Marshal(&config{
		RemoveBuildDeps: RemoveBuildDeps,
		SyncPackages:    SyncPackages,
		RemoveGitRepos:  RemoveGitRepos,
		CachePackages:   CachePackages,
		VerboseOutput:   VerboseOutput,
		PrettyPrint:     PrettyPrint,
		RepoCacheDir:    RepoCacheDir,
		PackageCacheDir: PackageCacheDir,
		LogFile:         LogFile,
		MapFile:         MapFile,
		LockFile:        LockFile,
	})
	checkErr(err)
	err = os.WriteFile(cfgfile, b, 0o600)
	checkErr(err)
}
