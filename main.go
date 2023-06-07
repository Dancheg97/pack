// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package main

import (
	"fmt"
	"os"
	"strings"

	"fmnx.su/core/pack/pack"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/tmpl"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Help bool `long:"help" short:"h"`

		// Root options.
		Query  bool `short:"Q" long:"query"`
		Remove bool `short:"R" long:"remove"`
		Sync   bool `short:"S" long:"sync"`
		Push   bool `short:"P" long:"push"`
		Open   bool `short:"O" long:"open"`
		Build  bool `short:"B" long:"build"`

		// Shared options.
		Info []bool `short:"i" long:"info"`
		List []bool `short:"l" long:"list"`
		Dir  string `short:"d" long:"dir"`

		// Query options.
		Explicit bool   `short:"e" long:"explicit"`
		Unreq    bool   `short:"t" long:"unreq"`
		File     string `long:"file"`
		Foreign  bool   `long:"foreign"`
		Deps     bool   `long:"deps"`
		Native   bool   `long:"native"`
		Check    []bool `long:"check"`

		// Remove options.
		Confirm     bool `short:"o" long:"confirm"`
		Norecursive bool `short:"a" long:"norecursive"`
		Nocfgs      bool `short:"w" long:"nocfgs"`
		Cascade     bool `long:"cascade"`

		// Sync options.
		Refresh   bool   `short:"y" long:"refresh"`
		Upgrade   []bool `short:"u" long:"sysupgrade"`
		Reinstall bool   `short:"r" long:"reinstall"`
		Quick     bool   `short:"q" long:"quick"`

		// Open options.
		Name string   `short:"n" long:"name"`
		Mirr []string `short:"m" long:"mirr"`
		Port string   `short:"p" long:"port"`
		Cert string   `short:"c" long:"cert"`
		Key  string   `short:"k" long:"key"`

		// Push options.
		HTTP bool `long:"http"`
	}

	_, err := flags.NewParser(&opts, flags.None).Parse()
	CheckErr(err)

	switch {
	case opts.Query && opts.Help:
		fmt.Println(tmpl.QueryHelp)
		return

	case opts.Query:
		CheckErr(pacman.Query(args(), pacman.QueryOptions{
			Explicit:         false,
			Deps:             false,
			Native:           opts.Native,
			Foreign:          false,
			Unrequired:       opts.Unreq,
			Groups:           false,
			Info:             opts.Info,
			Check:            []bool{},
			List:             opts.List,
			File:             opts.File,
			Stdout:           os.Stdout,
			Stderr:           os.Stderr,
			Stdin:            os.Stdin,
			AdditionalParams: []string{},
		}))
		return

	case opts.Remove && opts.Help:
		fmt.Println(tmpl.RemoveHelp)
		return

	case opts.Remove:
		CheckErr(pacman.RemoveList(args(), pacman.RemoveOptions{
			Sudo:        true,
			NoConfirm:   !opts.Confirm,
			Recursive:   !opts.Norecursive,
			WithConfigs: !opts.Nocfgs,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		}))
		return

	case opts.Sync && opts.Help:
		fmt.Println(tmpl.SyncHelp)
		return

	case opts.Sync:
		CheckErr(pacman.SyncList(args(), pacman.SyncOptions{
			Sudo:      true,
			Needed:    !opts.Reinstall,
			NoConfirm: opts.Quick,
			Refresh:   opts.Refresh,
			Upgrade:   opts.Upgrade,
			List:      opts.List,
			Stdout:    os.Stdout,
			Stderr:    os.Stderr,
			Stdin:     os.Stdin,
		}))
		return

	case opts.Push && opts.Help:
		fmt.Println(tmpl.PushHelp)
		return

	case opts.Push:
		CheckErr(pack.Push(args()))
		return

	case opts.Build && opts.Help:
		fmt.Println(tmpl.BuildHelp)
		return

	case opts.Help:
		fmt.Println(tmpl.Help)
		return

	default:
		fmt.Println("Please, specify at least one root flag (pack -h)")
		os.Exit(1)
		return
	}
}

// Herlper function to exit on unexpected errors.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func args() []string {
	var filtered []string
	for _, v := range os.Args {
		if !strings.HasPrefix(v, "/") &&
			!strings.HasPrefix(v, "-") &&
			!strings.HasPrefix(v, ".") {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
