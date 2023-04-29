package cfg

import (
	"fmt"
	"os"
	"os/user"

	"fmnx.io/dev/pack/core"
	"gopkg.in/yaml.v2"
)

const DefaultConfig = `# Remove git repositroy after package installation
remove-git-repo: false
# Remove build dependencies after installation
remove-build-deps: false
# Remove .pkg.tar.zst file after installation
remove-built-packages: false
# Cache dir for repositories
repo-cache-dir: %s/.pack
# Where pack will store built .pkg.tar.zst files
package-cache-dir: /var/cache/pacman/pkg
`

type Config struct {
	RemoveGitRepos      bool   `yaml:"remove-git-repo"`
	RemoveBuildDeps     bool   `yaml:"remove-build-deps"`
	RemoveBuiltPackages bool   `yaml:"remove-built-packages"`
	RepoCacheDir        string `yaml:"repo-cache-dir"`
	PackageCacheDir     string `yaml:"package-cache-dir"`
}

func GetConfig() *Config {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("unable to get current user")
		os.Exit(1)
	}
	cfg, err := os.Stat(usr.HomeDir + "/.pack/pack.yml")
	if err != nil || cfg.IsDir() {
		contents := fmt.Sprintf(DefaultConfig, usr.HomeDir)
		err = core.WriteFile(usr.HomeDir+"/.pack/pack.yml", contents)
		if err != nil {
			fmt.Println("Unable to wrute default configuration ", err)
			os.Exit(1)
		}
		return &Config{
			RemoveGitRepos:      false,
			RemoveBuildDeps:     false,
			RemoveBuiltPackages: false,
			RepoCacheDir:        usr.HomeDir + "/.pack",
			PackageCacheDir:     "/var/cache/pacman/pkg",
		}
	}
	b, err := os.ReadFile(usr.HomeDir + "/.pack/pack.yml")
	if err != nil {
		fmt.Println("unable to read config ~/.pack/pack.yml")
		os.Exit(1)
	}
	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		fmt.Println("unable unmarshall yaml in cfg: ~/.pack/pack.yml")
		os.Exit(1)
	}
	return &conf
}
