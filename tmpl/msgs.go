// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Color bool

func init() {
	b, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println("unable to read pacman configuration")
		os.Exit(1)
	}
	Color = strings.Contains(string(b), "\nColor\n")
	if !Color {
		color.NoColor = true
	}
	if Color {
		Help = strings.Join([]string{"📦", Help}, " ")
		QueryHelp = strings.Join([]string{"🔎", QueryHelp}, " ")
		RemoveHelp = strings.Join([]string{"📍", RemoveHelp}, " ")
		SyncHelp = strings.Join([]string{"⚡", SyncHelp}, " ")
		PushHelp = strings.Join([]string{"🚀", PushHelp}, " ")
		BuildHelp = strings.Join([]string{"🔐", BuildHelp}, " ")
		OpenHelp = strings.Join([]string{"🌐", OpenHelp}, " ")
	}
}

var Help = `Simplified version of pacman

operations:
	pack {-S --sync}   [options] [registry/(owner)/package(s)]
	pack {-P --push}   [options] [registry/(owner)/package(s)]
	pack {-R --remove} [options] [package(s)]
	pack {-Q --query}  [options] [package(s)]
	pack {-B --build}  [options]
	pack {-O --open}   [options]

use 'pack {-h --help}' with an operation for available options`

var SyncHelp = `Syncronize packages

options:
	-q, --quick       Do not ask for any confirmation (noconfirm)
	-y, --refresh     Download fresh package databases from the server (-yy force)
	-u, --upgrade     Upgrade installed packages (-uu enables downgrade)
	-i, --info        View package information (-ii for extended information)
	-l, --list <repo> View a list of packages in a repo
	-j, --notimeout   Use relaxed timeouts for download
	-f, --force       Reinstall up to date targets
	-k, --keepcfg     Do not save new registries in pacman.conf

usage:  pack {-S --sync} [options] <registry/owner/package(s)>`

var PushHelp = `Push packages

options:
	-d, --dir <dir> Use custom source dir with packages (default /var/cache/pacman/pkg)
	-f, --force     Owerwrite package with same version if exists
	    --protocol  Protocol that will be used for client (default https)
	    --endpoint  Use custom API endpoint (default /api/pack)

usage:  pack {-P --push} [options] <registry/owner/package(s)>`

var RemoveHelp = `Remove packages

options:
	-o, --confirm  Ask for confirmation when deleting package
	-a, --norecurs Leave package dependencies in the system (removed by default)
	-w, --nocfgs   Leave package configs in the system (removed by default)
	    --cascade  Remove packages and all packages that depend on them

usage:  pack {-R --remove} [options] <package(s)>`

var QueryHelp = `Query packages

options:
	-i, --info     View package information (-ii for backup files)
	-l, --list     List the files owned by the queried package
	    --explicit List packages explicitly installed [filter]
	    --unreq    List packages not (optionally) required by any
	    --file     Query a package file instead of the database
	    --deps     List packages installed as dependencies [filter]
	    --foreign  List installed packages not found in sync db(s) [filter]
	    --native   List installed packages only found in sync db(s) [filter]
	    --check    Check that package files exist (-kk for file properties)
	    --groups   View all members of a package group

usage:  pack {-Q --query} [options] [package(s)]`

var BuildHelp = `Build package in current directory

options:
	-q, --quick     Do not ask for any confirmation (noconfirm)
	-d, --dir <dir> Use custom dir to store result (default /var/cache/pacman/pkg)
	-s, --syncbuild Syncronize dependencies and build target
	-r, --rmdeps    Remove installed dependencies after a successful build
	-g, --garbage   Do not clean workspace before and after build

usage:  pack {-B --build} [options]`

var OpenHelp = `Launch package registry

options:
	-d, --dir      Exposed directory (default /var/cache/pacman/pkg)
	-p, --port     Port to launch server (default 80)
	-n, --name     Name of registry database (default localhost)
	    --cert     Certificate file for TLS server
	    --key      Key file for TLS server
	    --gpgdir   Directory containing files with GnuPG keys (sign validation)
	    --endpoint Use custom API endpoint (default /api/pack)

usage:  pack {-O --open} [options]`

var Version = `
$$$$$$\   $$$$$$\   $$$$$$$\  $$ |  $$\             Pack - package manager.
$$  __$$\  \____$$\ $$  ____| $$ | $$  |          Copyright (C) 2023 FMNX team
$$ /  $$ | $$$$$$$ |$$ /      $$$$$$  / 
$$ |  $$ |$$  __$$ |$$ |      $$  _$$\   This program may be freely redistributed under
$$$$$$$  |\$$$$$$$ |\$$$$$$$\ $$ | \$$\   the terms of the GNU General Public License.
$$  ____/  \_______| \_______|\__|  \__|      Web page: https://fmnx.su/core/pack
$$ |
$$ |                                                    Version: 0.4.5
\__|`

var Dots = color.New(color.FgWhite, color.Bold, color.FgHiBlue).Sprintf("::")
var DbAdded = color.New(color.FgWhite, color.Bold).Sprintf(" Database added: ")
var SigAccess = color.New(color.FgWhite, color.Bold).Sprintf(" Reading package signature: ")
var Accepted = color.New(color.Bold, color.FgHiGreen).Sprintf("accepted")
var Launching = color.New(color.Bold, color.Bold).Sprintf("Launching pack server: ")