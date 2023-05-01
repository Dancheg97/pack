package cmd

import (
	"fmt"
	"os"

	"fmnx.io/dev/pack/config"
	"fmnx.io/dev/pack/core"
	"github.com/fatih/color"
	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "pack",
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
		lf.Unlock()
		os.Exit(1)
	}
}

func CheckErr(err error) {
	if err != nil {
		RedPrint("Error occured: ", fmt.Sprintf("%+v", err))
		lf.Unlock()
		os.Exit(1)
	}
}

func ExecuteCheck(script string) {
	out, err := core.SystemCall(script)
	if err != nil {
		RedPrint("Command did not succed: ", script)
		fmt.Println("System output: ", out)
		RedPrint("Error occured: ", fmt.Sprintf("%+v", err))
		lf.Unlock()
		os.Exit(1)
	}
}

func RedPrint(white string, red string) {
	fmt.Printf(white + " " + color.RedString(red) + "\n")
}

func BluePrint(white string, blue string) {
	fmt.Printf(white + " " + color.BlueString(blue) + "\n")
}

func GreenPrint(white string, green string) {
	fmt.Printf(white + " " + color.GreenString(green) + "\n")
}

func YellowPrint(white string, yellow string) {
	fmt.Printf(white + " " + color.YellowString(yellow) + "\n")
}
