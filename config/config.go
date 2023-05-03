package config

import (
	"fmt"
	"os"
	"os/user"

	"fmnx.io/dev/pack/system"
	"github.com/caarlos0/env/v8"
	"gopkg.in/yaml.v2"
)

const DefaultConfig = `# Remove git repositroy after package installation
remove-git-repo: %t
# Remove .pkg.tar.zst file after installation
remove-built-packages: %t
# Print additional debug information
debug-mode: %t
# Always OK to AUR packages
allow-aur: %t
# Cache dir for repositories
repo-cache-dir: %s/.pack
# Where pack will store built .pkg.tar.zst files
package-cache-dir: /var/cache/pacman/pkg
# Location of mapping file (pack packages and related pacman packages)
map-file: %s/.pack/mapping.json
# Location of lock file
lock-file: /tmp/pack.lock
`

type EnvVars struct {
	RemoveGitRepos      bool `env:"PACK_REMOVE_GIT_REPOS" envDefault:"false"`
	RemoveBuiltPackages bool `env:"PACK_REMOVE_BUILT_PACKAGES" envDefault:"false"`
	DebugMode           bool `env:"PACK_DEBUG_MODE" envDefault:"false"`
	AllowAUR            bool `env:"PACK_ALLOW_AUR" envDefault:"false"`
}

type Config struct {
	RemoveGitRepos      bool   `yaml:"remove-git-repo"`
	RemoveBuiltPackages bool   `yaml:"remove-built-packages"`
	DebugMode           bool   `yaml:"debug-mode"`
	AllowAUR            bool   `yaml:"allow-aur"`
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
		envs := EnvVars{}
		err = env.Parse(&envs)
		if err != nil {
			return nil, err
		}
		err = system.WriteFile(usr.HomeDir+"/.pack/config.yml", fmt.Sprintf(
			DefaultConfig,
			envs.RemoveGitRepos,
			envs.RemoveBuiltPackages,
			envs.DebugMode,
			envs.AllowAUR,
			usr.HomeDir,
			usr.HomeDir,
		))
		if err != nil {
			return nil, err
		}
		return &Config{
			RemoveGitRepos:      envs.RemoveGitRepos,
			RemoveBuiltPackages: envs.RemoveBuiltPackages,
			DebugMode:           envs.DebugMode,
			AllowAUR:            envs.AllowAUR,
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
