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

	SearchSources []SearchSource
)

var (
	homedir string
	cfgfile string
	cfg     config
)

// Configuration variables.
type config struct {
	Needed       bool   `yaml:"needed"`
	RmDeps       bool   `yaml:"rm-deps"`
	RmRepos      bool   `yaml:"rm-repos"`
	CachePkgs    bool   `yaml:"cache-pkgs"`
	Verbose      bool   `yaml:"verbose"`
	PrettyPrint  bool   `yaml:"pretty-print"`
	RepoCacheDir string `yaml:"repo-cache-dir"`
	PkgCacheDir  string `yaml:"pkg-cache-dir"`
	LogFile      string `yaml:"log-file"`
	MapFile      string `yaml:"map-file"`
	LockFile     string `yaml:"lock-file"`

	SearchSources []SearchSource `yaml:"search-sources"`
}

type SearchSource struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Field  string `yaml:"regexp"`
	Prefix string `yaml:"prefix"`
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

	Needed = cfg.Needed
	RmDeps = cfg.RmDeps
	RmRepos = cfg.RmRepos
	CachePkgs = cfg.CachePkgs
	Verbose = cfg.Verbose
	PrettyPrint = cfg.PrettyPrint
	RepoCacheDir = cfg.RepoCacheDir
	PkgCacheDir = cfg.PkgCacheDir
	LogFile = cfg.LogFile
	MapFile = cfg.MapFile
	LockFile = cfg.LockFile
	SearchSources = cfg.SearchSources
}

// SetDefaults configuration to default values and set save config file.
func SetDefaults() {
	Needed = true
	RmDeps = false
	RmRepos = false
	CachePkgs = true
	Verbose = false
	PrettyPrint = true
	RepoCacheDir = homedir + "/.pack"
	PkgCacheDir = "/var/cache/pacman/pkg"
	LogFile = "/tmp/pack.log"
	MapFile = homedir + "/.pack/mapping.json"
	LockFile = "/tmp/pack.lock"
	SearchSources = []SearchSource{
		{
			Name:   "fmnx linux packages",
			Url:    "https://fmnx.su/api/v1/repos/search?q={{package}}&team_id=4",
			Field:  "name",
			Prefix: "fmnx.su/",
		},
		{
			Name:   "arch linux repository",
			Url:    "https://aur.archlinux.org/rpc/?v=5&type=search&by=name&arg={{package}}",
			Field:  "PackageBase",
			Prefix: "aur.archlinux.org/",
		},
	}
}

// Save configuration with all new variables.
func Save() {
	err := os.MkdirAll(homedir+"/.pack", os.ModePerm)
	checkErr(err)
	b, err := yaml.Marshal(&config{
		Needed:        Needed,
		RmDeps:        RmDeps,
		RmRepos:       RmRepos,
		CachePkgs:     CachePkgs,
		Verbose:       Verbose,
		PrettyPrint:   PrettyPrint,
		RepoCacheDir:  RepoCacheDir,
		PkgCacheDir:   PkgCacheDir,
		LogFile:       LogFile,
		MapFile:       MapFile,
		LockFile:      LockFile,
		SearchSources: SearchSources,
	})
	checkErr(err)
	err = os.WriteFile(cfgfile, b, 0o600)
	checkErr(err)
}

// Set output for pack logs.
func SetLogger() {
	err := os.WriteFile(LogFile, []byte{}, 0666)
	checkErr(err)
	f, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	log.Default().SetOutput(f)
}
