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
	Needed       bool
	RmDeps       bool
	RmRepos      bool
	CachePkgs    bool
	Verbose      bool
	PrettyPrint  bool
	RepoCacheDir string
	PkgCacheDir  string
	LogFile      string
	MapFile      string
	LockFile     string
)

var (
	homedir string
	cfgfile string
	cfg     config
)

// Configuration variables.
type config struct {
	Needed       bool   `yaml:"needed"`
	RmDeps       bool   `yaml:"remove-build-deps"`
	RmRepos      bool   `yaml:"remove-git-repos"`
	CachePkgs    bool   `yaml:"cache-packages"`
	Verbose      bool   `yaml:"verbose-output"`
	PrettyPrint  bool   `yaml:"pretty-print"`
	RepoCacheDir string `yaml:"repo-cache-dir"`
	PkgCacheDir  string `yaml:"package-cache-dir"`
	LogFile      string `yaml:"log-file"`
	MapFile      string `yaml:"map-file"`
	LockFile     string `yaml:"lock-file"`
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
	cfg.RmDeps = false
	cfg.Needed = false
	cfg.RmRepos = false
	cfg.CachePkgs = true
	cfg.Verbose = false
	cfg.PrettyPrint = true
	cfg.RepoCacheDir = homedir + "/.pack"
	cfg.PkgCacheDir = "/var/cache/pacman/pkg"
	cfg.LogFile = "/tmp/pack.log"
	cfg.MapFile = homedir + "/.pack/mapping.json"
	cfg.LockFile = "/tmp/pack.lock"

	RmDeps = cfg.RmDeps
	Needed = cfg.Needed
	RmRepos = cfg.RmRepos
	CachePkgs = cfg.CachePkgs
	Verbose = cfg.Verbose
	PrettyPrint = cfg.PrettyPrint
	RepoCacheDir = cfg.RepoCacheDir
	PkgCacheDir = cfg.PkgCacheDir
	LogFile = cfg.LogFile
	MapFile = cfg.MapFile
	LockFile = cfg.LockFile
}

// Save configuration with all new variables.
func Save() {
	b, err := yaml.Marshal(&config{
		RmDeps:       RmDeps,
		Needed:       Needed,
		RmRepos:      RmRepos,
		CachePkgs:    CachePkgs,
		Verbose:      Verbose,
		PrettyPrint:  PrettyPrint,
		RepoCacheDir: RepoCacheDir,
		PkgCacheDir:  PkgCacheDir,
		LogFile:      LogFile,
		MapFile:      MapFile,
		LockFile:     LockFile,
	})
	checkErr(err)
	err = os.WriteFile(cfgfile, b, 0o600)
	checkErr(err)
}
