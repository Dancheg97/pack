// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"os"
	"strings"
)

var Help = `📦 Simplified version of pacman

operations:
	pack {-S --sync}    [options] [repository/owner/package(s)]
	pack {-P --push}    [options] [repository/owner/package(s)]
	pack {-R --remove}  [options] [package(s)]
	pack {-Q --query}   [options] [package(s)]
	pack {-B --build}   [options]
	pack {-O --open}    [options]

use 'pack {-h --help}' with an operation for available options`

var QueryHelp = `🔎 Query packages

options:
	-i, --info     view package information (-ii for backup files)
	-l, --list     list the files owned by the queried package
	-e, --explicit list packages explicitly installed [filter]
	-t, --unreq    list packages not (optionally) required by any
	    --file     query a package file instead of the database
	    --deps     list packages installed as dependencies [filter]
	    --foreign  list installed packages not found in sync db(s) [filter]
	    --native   list installed packages only found in sync db(s) [filter]
	    --check    check that package files exist (-kk for file properties)

usage:  pack {-Q --query} [options] [package(s)]`

var RemoveHelp = `🚫 Remove packages

options:
	-o, --confirm  ask for confirmation when deleting package   
	-a, --norecurs leave package dependencies in the system (removed by default)
	-w, --nocfgs   leave package configs in the system (removed by default)
	    --cascade  remove packages and all packages that depend on them

usage:  pack {-R --remove} [options] <package(s)>`

var SyncHelp = `🔧 Syncronize packages

options:
	-q, --quick       do not ask for any confirmation (noconfirm)
	-y, --refresh     download fresh package databases from the server (-yy force)
	-u, --sysupgrade  upgrade installed packages (-uu enables downgrade)
	-i, --info        view package information (-ii for extended information)
	-l, --list <repo> view a list of packages in a repo
	-f, --no-timeout  use relaxed timeouts for download
	-r, --reinstall   reinstall up to date packages

usage:  pack {-S --sync} [options] <registry/owner/package(s)>`

var PushHelp = `🚀 Push packages

options:
	--http      push over http (default https)
	--dir <dir> use custom source dir different from (default /var/cache/pacman/pkg)

usage:  pack {-P --push} [options] <registry/owner/package(s)>`

var BuildHelp = `🧰 Build packages

options:
	-q, --quick     do not ask for any confirmation (noconfirm)
	-d, --dir <dir> use custom cache dir (default /var/cache/pacman/pkg)
	-r, --reinstall reinstall up to date packages
	-e, --install   install package after successfull build
	-z, --rmdeps    remove installed dependencies after a successful build
	-g, --garbage   don't clean workspace before and after build

usage:  pack {-B --build} [options] <git repository(s)/current directory>`

var OpenHelp = `🌐 Open repository

options:
	-n, --name <domain> registry name, should match the domain name
	-d, --dir  <dir>    directory to expose, by default /var/cache/pacman/pkg
	-m, --mirr <link>   create mirror for existing pacman package repository
	-p, --port <port>   port to launch server on
	-c, --cert <file>   certificate file for TLS
	-k, --key  <file>   key file for TLS

usage:  pack {-B --build} [options] <git repository(s)/current directory>`

func init() {
	b, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println("unable to read pacman configuration")
		os.Exit(1)
	}
	if !strings.Contains(string(b), "\nColor\n") {
		Help = strings.Replace(Help, "📦 ", "", 1)
		QueryHelp = strings.Replace(QueryHelp, "🔎 ", "", 1)
		RemoveHelp = strings.Replace(RemoveHelp, "🚫 ", "", 1)
		SyncHelp = strings.Replace(SyncHelp, "🔧 ", "", 1)
		PushHelp = strings.Replace(PushHelp, "🚀 ", "", 1)
		BuildHelp = strings.Replace(BuildHelp, "🧰 ", "", 1)
		OpenHelp = strings.Replace(OpenHelp, "🌐 ", "", 1)
	}
}
