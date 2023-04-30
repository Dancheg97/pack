package cmd

import (
	"fmt"
	"log"
	"os"

	"fmnx.io/dev/pack/config"
	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fmnx-gen",
		Short: "📦 git-based pacman-compatible package manager",
	}
	flags = []Flag{}
	cfg   *config.Config
	lf    *lockfile.Lockfile
)

func init() {
	conf, err := config.GetConfig()
	CheckErr(err)
	cfg = conf
	lock, err := lockfile.New(conf.LockFile)
	CheckErr(err)
	lf = &lock
	err = lf.TryLock()
	CheckErr(err)
}

func Execute() {
	for _, flag := range flags {
		AddFlag(flag)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Printf("%+v", err)
		lf.Unlock()
		os.Exit(1)
	}
}
