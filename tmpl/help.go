// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

const Help = `📦 Simplified version of pacman

operations:
	pack {-S --sync}    [options] [registry/owner/package(s)]
	pack {-R --remove}  [options] [package(s)]
	pack {-Q --query}   [options] [package(s)]
	pack {-P --push}    [options] [package(s)]
	pack {-B --build}   [options] [repository(s)]
	pack {-O --open}    [options] [registry(s)]

use 'pack {-h --help}' with an operation for available options`

const QueryHelp = `🔎 Query packages

options:
	-d, --deps     list packages installed as dependencies [filter]
	-e, --explicit list packages explicitly installed [filter]
	-i, --info     view package information (-ii for backup files)
	-k, --check    check that package files exist (-kk for file properties)
	-l, --list     list the files owned by the queried package
	-m, --foreign  list installed packages not found in sync db(s) [filter]
	-n, --native   list installed packages only found in sync db(s) [filter]
	-p, --file     query a package file instead of the database
	-t, --unreq    list packages not (optionally) required by any

usage:  pack {-Q --query} [options] [package(s)]`

const RemoveHelp = `🚫 Remove packages

options:
	-o, --confirm  ask for confirmation when deleting package   
	-c, --cascade  remove packages and all packages that depend on them
	-a, --norecurs leave package dependencies in the system (removed by default)
	-w, --nocfgs   leave package configs in the system (removed by default)

usage:  pack {-R --remove} [options] <package(s)>`

const SyncHelp = `🔧 Syncronize packages

options:
	-y, --refresh     download fresh package databases from the server (-yy force)
	-u, --sysupgrade  upgrade installed packages (-uu enables downgrade)
	-g, --garbage     leave old packages in cache directory
	-i, --info        view package information (-ii for extended information)
	-l, --list <repo> view a list of packages in a repo
	-f, --no-timeout  use relaxed timeouts for download
	-r, --reinstall   reinstall up to date packages
	-q, --quick       do not ask for any confirmation (noconfirm)

usage:  pack {-S --sync} [options] <registry/owner/package(s)>`
