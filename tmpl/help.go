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
		Help = strings.Replace(Help, "📦 ", "", 1)
		QueryHelp = strings.Replace(QueryHelp, "🔎 ", "", 1)
		RemoveHelp = strings.Replace(RemoveHelp, "🚫 ", "", 1)
		SyncHelp = strings.Replace(SyncHelp, "🔧 ", "", 1)
		PushHelp = strings.Replace(PushHelp, "🚀 ", "", 1)
		BuildHelp = strings.Replace(BuildHelp, "🧰 ", "", 1)
		OpenHelp = strings.Replace(OpenHelp, "🌐 ", "", 1)
	}
}

var Help = `📦 Simplified version of pacman

operations:
	pack {-S --sync}    [options] [repository/owner/package(s)]
	pack {-P --push}    [options] [repository/owner/package(s)]
	pack {-R --remove}  [options] [package(s)]
	pack {-Q --query}   [options] [package(s)]
	pack {-B --build}   [options]
	pack {-O --open}    [options]

use 'pack {-h --help}' with an operation for available options`

var SyncHelp = `🔧 Syncronize packages, repositroy and owner args are optional

options:
	-q, --quick       do not ask for any confirmation (noconfirm)
	-y, --refresh     download fresh package databases from the server (-yy force)
	-u, --upgrade     upgrade installed packages (-uu enables downgrade)
	-i, --info        view package information (-ii for extended information)
	-l, --list <repo> view a list of packages in a repo
	-j, --notimeout   use relaxed timeouts for download
	-f, --force       reinstall up to date targets
	-k, --keepcfg     do not save new registries in pacman.conf

usage:  pack {-S --sync} [options] <registry/owner/package(s)>`

var PushHelp = `🚀 Push packages

options:
	--dir <dir> use custom source dir with packages (default /var/cache/pacman/pkg)
	--protocol  protocol that will be used for client (default https)
	--endpoint  use custom endpoint for push (default /api/pack/push)

usage:  pack {-P --push} [options] <registry/owner/package(s)>`

var RemoveHelp = `📍 Remove packages

options:
	-o, --confirm  ask for confirmation when deleting package
	-a, --norecurs leave package dependencies in the system (removed by default)
	-w, --nocfgs   leave package configs in the system (removed by default)
	    --cascade  remove packages and all packages that depend on them

usage:  pack {-R --remove} [options] <package(s)>`

var QueryHelp = `🔎 Query packages

options:
	-i, --info     view package information (-ii for backup files)
	-l, --list     list the files owned by the queried package
	    --explicit list packages explicitly installed [filter]
	    --unreq    list packages not (optionally) required by any
	    --file     query a package file instead of the database
	    --deps     list packages installed as dependencies [filter]
	    --foreign  list installed packages not found in sync db(s) [filter]
	    --native   list installed packages only found in sync db(s) [filter]
	    --check    check that package files exist (-kk for file properties)
	    --groups   view all members of a package group

usage:  pack {-Q --query} [options] [package(s)]`

var BuildHelp = `⚙️ Build package in current directory

options:
	-q, --quick     do not ask for any confirmation (noconfirm)
	-d, --dir <dir> use custom dir to store result (default /var/cache/pacman/pkg)
	-s, --syncbuild sync/reinstall target packages
	-r, --rmdeps    remove installed dependencies after a successful build
	-g, --garbage   don't clean workspace before and after build

usage:  pack {-B --build} [options]`

var OpenHelp = `🌐 Open registry - launch web-server

options:
	-d, --dir  <dir>    exposed directory (default /var/cache/pacman/pkg)
	-m, --mirr <link>   create mirror for existing pacman package repository
	-p, --port <port>   port to launch server on
	    --cert <file>   certificate file for TLS
	    --key  <file>   key file for TLS

usage:  pack {-O --open} [options]`
