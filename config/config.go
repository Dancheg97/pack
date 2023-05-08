package config

import (
	"fmt"
	"os"
	"os/user"

	"fmnx.io/core/pack/system"
	"gopkg.in/yaml.v2"
)

const DefaultConfig = `# Remove git repositroy after package installation
remove-git-repo: %t
# Remove .pkg.tar.zst file after installation
remove-built-packages: %t
# Print additional debug information
debug-mode: %t
# Disable colors in output
pack-disable-prettyprint: %t
# Cache dir for repositories
repo-cache-dir: %s/.pack
# Where pack will store built .pkg.tar.zst files
package-cache-dir: /var/cache/pacman/pkg
# Location of mapping file (pack packages and related pacman packages)
map-file: %s/.pack/mapping.json
# Location of lock file
lock-file: /tmp/pack.lock
`

// Default template configuration for pack can be found on top.
type Config struct {
	RemoveGitRepos      bool   `yaml:"remove-git-repo"`
	RemoveBuiltPackages bool   `yaml:"remove-built-packages"`
	DebugMode           bool   `yaml:"debug-mode"`
	DisablePrettyPrint  bool   `yaml:"pack-disable-prettyprint"`
	RepoCacheDir        string `yaml:"repo-cache-dir"`
	PackageCacheDir     string `yaml:"package-cache-dir"`
	MapFile             string `yaml:"map-file"`
	LockFile            string `yaml:"lock-file"`
}

func GetConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	cfg, err := os.Stat(usr.HomeDir + "/.pack/config.yml")
	if err != nil || cfg.IsDir() {
		err = system.WriteFile(usr.HomeDir+"/.pack/config.yml", fmt.Sprintf(
			DefaultConfig,
			GetBoolEnv(`PACK_REMOVE_GIT_REPOS`),
			GetBoolEnv(`PACK_REMOVE_BUILT_PACKAGES`),
			GetBoolEnv(`PACK_DEBUG_MODE`),
			GetBoolEnv(`PACK_DISABLE_PRETTYPRINT`),
			usr.HomeDir,
			usr.HomeDir,
		))
		if err != nil {
			return nil, err
		}
		return &Config{
			RemoveGitRepos:      GetBoolEnv(`PACK_REMOVE_GIT_REPOS`),
			RemoveBuiltPackages: GetBoolEnv(`PACK_REMOVE_BUILT_PACKAGES`),
			DebugMode:           GetBoolEnv(`PACK_DEBUG_MODE`),
			DisablePrettyPrint:  GetBoolEnv(`PACK_DISABLE_PRETTYPRINT`),
			RepoCacheDir:        usr.HomeDir + "/.pack",
			PackageCacheDir:     "/var/cache/pacman/pkg",
			MapFile:             usr.HomeDir + "/.pack/mapping.json",
			LockFile:            "/tmp/pack.lock",
		}, nil
	}
	b, err := os.ReadFile(usr.HomeDir + "/.pack/config.yml")
	if err != nil {
		return nil, err
	}
	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

// Function that is made to reduce amount of dependencies.
func GetBoolEnv(v string) bool {
	return os.Getenv(v) == `true`
}
