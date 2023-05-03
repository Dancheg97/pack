package config

import (
	"fmt"
	"os"
	"os/user"

	"fmnx.io/dev/pack/system"
	"gopkg.in/yaml.v2"
)

const DefaultConfig = `# Remove git repositroy after package installation
remove-git-repo: false
# Remove .pkg.tar.zst file after installation
remove-built-packages: false
# Print additional debug information
debug-mode: false
# Cache dir for repositories
repo-cache-dir: %s/.pack
# Where pack will store built .pkg.tar.zst files
package-cache-dir: /var/cache/pacman/pkg
# Location of mapping file (pack packages and related pacman packages)
map-file: %s/.pack/mapping.json
# Location of lock file
lock-file: /tmp/pack.lock
`

type Config struct {
	RemoveGitRepos      bool   `yaml:"remove-git-repo"`
	RemoveBuiltPackages bool   `yaml:"remove-built-packages"`
	DebugMode           bool   `yaml:"debug-mode"`
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
		contents := fmt.Sprintf(DefaultConfig, usr.HomeDir, usr.HomeDir)
		err = system.WriteFile(usr.HomeDir+"/.pack/config.yml", contents)
		if err != nil {
			return nil, err
		}
		return &Config{
			RemoveGitRepos:      false,
			RemoveBuiltPackages: false,
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
