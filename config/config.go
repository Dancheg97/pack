// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package config

// Project runtime configuration.

import (
	"fmt"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

// This variables will be automatically initialized in init().
var (
	RemoveGitRepos      bool
	RemoveBuiltPackages bool
	DebugMode           bool
	DisablePrettyPrint  bool
	CacheDir            string
	PackageCacheDir     string
	MapFile             string
	LockFile            string
)

// Default template configuration for pack can be found on top.
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

func init() {
	usr, err := user.Current()
	checkErr(err)
	err = os.MkdirAll(usr.HomeDir+"/.pack", 0755)
	checkErr(err)
	cfg, err := os.Stat(usr.HomeDir + "/.pack/config.yml")
	if err != nil || cfg.IsDir() {
		cfgString := fmt.Sprintf(
			defaultConfig,
			getBoolEnv(`PACK_REMOVE_GIT_REPOS`),
			getBoolEnv(`PACK_REMOVE_BUILT_PACKAGES`),
			getBoolEnv(`PACK_DEBUG_MODE`),
			getBoolEnv(`PACK_DISABLE_PRETTYPRINT`),
			usr.HomeDir,
			usr.HomeDir,
		)
		err = os.WriteFile(
			usr.HomeDir+"/.pack/config.yml",
			[]byte(cfgString), 0o600,
		)
		checkErr(err)
		RemoveGitRepos = getBoolEnv(`PACK_REMOVE_GIT_REPOS`)
		RemoveBuiltPackages = getBoolEnv(`PACK_REMOVE_BUILT_PACKAGES`)
		DebugMode = getBoolEnv(`PACK_DEBUG_MODE`)
		DisablePrettyPrint = getBoolEnv(`PACK_DISABLE_PRETTYPRINT`)
		CacheDir = usr.HomeDir + "/.pack"
		PackageCacheDir = "/var/cache/pacman/pkg"
		MapFile = usr.HomeDir + "/.pack/mapping.json"
		LockFile = "/tmp/pack.lock"
	}
	b, err := os.ReadFile(usr.HomeDir + "/.pack/config.yml")
	checkErr(err)
	var conf config
	err = yaml.Unmarshal(b, &conf)
	checkErr(err)
	RemoveGitRepos = conf.RemoveGitRepos
	RemoveBuiltPackages = conf.RemoveBuiltPackages
	DebugMode = conf.DebugMode
	DisablePrettyPrint = conf.DisablePrettyPrint
	CacheDir = conf.PackCacheDir
	PackageCacheDir = conf.PackageCacheDir
	MapFile = conf.MapFile
	LockFile = conf.LockFile
}

// Function that is made to reduce amount of dependencies.
func getBoolEnv(v string) bool {
	return os.Getenv(v) == `true`
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

const defaultConfig = `# Remove git repositroy after package installation
remove-git-repo: %t

# Don't cache .pkg.tar.zst file after installation
remove-built-packages: %t

# Print every system call execution
debug-mode: %t

# Disable colors and emojis in output
disable-prettyprint: %t

# Location where pack will store package repositories
pack-cache-dir: %s/.pack

# Location of mapping file (pack packages and related pacman packages)
map-file: %s/.pack/mapping.json

# Location to put .pkg.tar.zst packages after installation
package-cache-dir: /var/cache/pacman/pkg

# Location of lock file
lock-file: /tmp/pack.lock
`
