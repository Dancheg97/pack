// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package main

import (
	"fmt"
	"os"
	"strings"

	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/tmpl"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Help bool `long:"help" short:"h"`

		// Root flags.
		Query  bool `short:"Q" long:"query"`
		Remove bool `short:"R" long:"remove"`
		Sync   bool `short:"S" long:"sync"`
		Push   bool `short:"P" long:"push"`
		Open   bool `short:"O" long:"open"`
		Build  bool `short:"B" long:"build"`

		// Shared flags.
		Info []bool `short:"i" long:"info"`
		List []bool `short:"l" long:"list"`

		// Query options.
		Deps     bool   `short:"d" long:"deps"`
		Explicit bool   `short:"e" long:"explicit"`
		Foreign  bool   `short:"m" long:"foreign"`
		Native   bool   `short:"n" long:"native"`
		Unreq    bool   `short:"t" long:"unreq"`
		File     string `short:"p" long:"file"`
		Check    []bool `short:"k" long:"check"`

		// Remove options.
		Confirm     bool `short:"o" long:"confirm"`
		Cascade     bool `short:"c" long:"cascade"`
		Norecursive bool `short:"a" long:"norecursive"`
		Nocfgs      bool `short:"w" long:"nocfgs"`

		// Sync options.
		Garbage   bool   `short:"g" long:"garbage"`
		Refresh   bool   `short:"y" long:"refresh"`
		Reinstall bool   `short:"r" long:"reinstall"`
		Quick     bool   `short:"q" long:"quick"`
		Upgrade   []bool `short:"u" long:"sysupgrade"`
	}

	_, err := flags.NewParser(&opts, flags.None).Parse()
	CheckErr(err)

	switch {
	case opts.Query && opts.Help:
		fmt.Println(tmpl.QueryHelp)
		return

	case opts.Query:
		CheckErr(pacman.Query(args(), pacman.QueryOptions{
			Explicit:   opts.Explicit,
			Deps:       opts.Deps,
			Native:     opts.Native,
			Foreign:    opts.Foreign,
			Unrequired: opts.Unreq,
			Info:       opts.Info,
			Check:      opts.Check,
			List:       opts.List,
			File:       opts.File,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Stdin:      os.Stdin,
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
			CleanAll:  !opts.Garbage,
			Stdout:    os.Stdout,
			Stderr:    os.Stderr,
			Stdin:     os.Stdin,
		}))
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
