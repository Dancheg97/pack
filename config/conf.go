// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package config

// Project configuration.

import (
	"fmt"
	"log"
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
	defer SetLogger()
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

	ReadConfigParams()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("configuration error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

// Set params from config to variables.
func ReadConfigParams() {
	b, err := os.ReadFile(cfgfile)
	checkErr(err)
	err = yaml.Unmarshal(b, &cfg)
	checkErr(err)

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

// SetDefaults configuration to default values and set save config file.
func SetDefaults() {
	RmDeps = false
	Needed = false
	RmRepos = false
	CachePkgs = true
	Verbose = false
	PrettyPrint = true
	RepoCacheDir = homedir + "/.pack"
	PkgCacheDir = "/var/cache/pacman/pkg"
	LogFile = "/tmp/pack.log"
	MapFile = homedir + "/.pack/mapping.json"
	LockFile = "/tmp/pack.lock"
}

// Save configuration with all new variables.
func Save() {
	err := os.MkdirAll(homedir+"/.pack", os.ModePerm)
	checkErr(err)
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

// Set output for pack logs.
func SetLogger() {
	err := os.WriteFile(LogFile, []byte{}, 0666)
	checkErr(err)
	f, err := os.OpenFile(
		LogFile,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	checkErr(err)
	log.Default().SetOutput(f)
}
