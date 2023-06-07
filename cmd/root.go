// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	pacmancache  = "/var/cache/pacman/pkg"
	pkgext       = ".pkg.tar.zst"
	protocol     = "https://"
	fsendpoint   = "/api/pack/"
	pushendpoint = "/api/pack/push"
	file         = "file"
	sign         = "sign"
	email        = "email"
)

// Prepare cobra and viper templates.
func init() {
	Command.SetHelpCommand(&cobra.Command{})
	Command.SetUsageTemplate(`📦 Simplified version of pacman

operations:
    pack {-h --help}
    pack {-Q --query}   [options] [package(s)]
    pack {-R --remove}  [options] [package(s)]
    pack {-S --sync}    [options] [package(s)]
    pack {-P --push}    [options] [package(s)]
	pack {-U --upgrade} [options] [file(s)]
    pack {-O --open}    [options] [database(s)]
    pack {-B --build}   [options] [repository(s)]

use 'pack {-h --help}' with an operation for available options
`)
}

var Command = &cobra.Command{
	Use:          "pack",
	SilenceUsage: true,
	Run:          Runfunc,
}

func Runfunc(cmd *cobra.Command, args []string) {
	fmt.Println("YO")
}

// Main execution of cobra command.
func init() {
	Flg(&Flgp{Name: "help", Short: "h"})

	Flg(&Flgp{Name: "query", Short: "Q"})
	Flg(&Flgp{Name: "remove", Short: "R"})
	Flg(&Flgp{Name: "sync", Short: "S"})
	Flg(&Flgp{Name: "deptest", Short: "T"})
	Flg(&Flgp{Name: "upgrade", Short: "U"})

	SFlg(&Flgp{Name: "dbpath", Short: "b"})
	Flg(&Flgp{Name: "check", Short: "k"})
	Flg(&Flgp{Name: "quiet", Short: "q"})
	SFlg(&Flgp{Name: "root", Short: "r"})
	Flg(&Flgp{Name: "verbose", Short: "v"})
	SFlg(&Flgp{Name: "arch"})
	Flg(&Flgp{Name: "asdeps"})
	Flg(&Flgp{Name: "asexplicit"})
	SFlg(&Flgp{Name: "cachedir"})
	SFlg(&Flgp{Name: "color"})
	SFlg(&Flgp{Name: "config"})
	Flg(&Flgp{Name: "confirm"})
	Flg(&Flgp{Name: "debug"})
	Flg(&Flgp{Name: "disable-download-timeout"})
	SFlg(&Flgp{Name: "gpgdir"})
	SFlg(&Flgp{Name: "hookdir"})
	SFlg(&Flgp{Name: "logfile"})
	SFlg(&Flgp{Name: "noconfirm"})
	SFlg(&Flgp{Name: "sysroot"})
	SFlg(&Flgp{Name: "list", Short: "l"})
	SFlg(&Flgp{Name: "regex", Short: "x"})
	SFlg(&Flgp{Name: "refresh", Short: "y"})
	SFlg(&Flgp{Name: "machinereadable"})
	Flg(&Flgp{Name: "deps", Short: "d"})
	Flg(&Flgp{Name: "explict", Short: "e"})
	Flg(&Flgp{Name: "groups", Short: "g"})
	Flg(&Flgp{Name: "info", Short: "i"})
	Flg(&Flgp{Name: "foreign", Short: "m"})
	Flg(&Flgp{Name: "native", Short: "n"})
	SFlg(&Flgp{Name: "owns", Short: "o"})
	SFlg(&Flgp{Name: "file"})
	SFlg(&Flgp{Name: "search"})
	SFlg(&Flgp{Name: "unrequired", Short: "t"})
	SFlg(&Flgp{Name: "upgrades", Short: "u"})
	SFlg(&Flgp{Name: "cascade", Short: "c"})
	SFlg(&Flgp{Name: "nodeps"})

	SFlg(&Flgp{Name: "port", Default: "80"})
	SFlg(&Flgp{Name: "name", Default: "localhost"})
	SFlg(&Flgp{Name: "dir", Default: pacmancache})
	SFlg(&Flgp{Name: "key"})
	SFlg(&Flgp{Name: "cert"})
	SLFlf(&Flgp{Name: "mirror"})

}

// Herlper function to exit on unexpected errors.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Parameters to add boolean persistent flag.
type Flgp struct {
	// Flag name, for example --help.
	Name string
	// Flag shortname, for example -h.
	Short string
	// Default value (only for string flags).
	Default string
}

// Add boolean flag to cobra command.
func Flg(p *Flgp) {
	Command.PersistentFlags().BoolP(p.Name, p.Short, false, "")
	err := viper.BindPFlag(p.Name, Command.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
}

// Add string flag to cobra command.
func SFlg(p *Flgp) {
	Command.PersistentFlags().StringP(p.Name, p.Short, p.Default, "")
	err := viper.BindPFlag(p.Name, Command.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
}

// Add string flag to cobra command.
func SLFlf(p *Flgp) {
	Command.PersistentFlags().StringArrayP(p.Name, p.Short, nil, "")
	err := viper.BindPFlag(p.Name, Command.PersistentFlags().Lookup(p.Name))
	CheckErr(err)
}
