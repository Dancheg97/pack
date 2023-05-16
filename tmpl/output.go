// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package tmpl

// This package contains string templates for utility command output.
// Output might be modified based on configuration - prettyprint.

import (
	"strings"

	"fmnx.su/core/pack/config"
)

var RootShort = "📦 decentralized package manager based on git and pacman"
var RootLong = `📦 decentralized package manager based on git and pacman

Configuration file: ~/.pack/config.yml
Official web page: https://fmnx.su/core/pack.

Usage:
pack [command] <package(s)>`

var DescribeShort = `📝 describe packages`
var DescribeLong = `📝 view information about packages

This tool provides information about package retrieved from pacman or pack.

Example:
pack d fmnx.su/core/ainst`

var GenerateShort = "📋 generate PKGBUILD"
var GenerateLong = `📋 generate PKGBUILD and update README.md with installation helper

This command will generate .pack.yml template and add some lines to README.md
to provide information about installation with pack.`

var InstallShort = "📥 install packages"
var InstallLong = `📥 install packages

You can provide links to repositories, regular package names, network links or
pathes to packages in local file system. This command this split packages into
groups and install. If some are already installed pack will skip them.

Example:
pack i github.com/exm/pkg@v1.23 nano https://example.com/package.pkg.tar.zst`

var ListShort = "📄 show installed packages"

var OutdatedShort = "📌 show outdated packages"
var OutdatedLong = `📌 show outdated packages

This command will make a call to pacman servers and collect information about
all remote repos for packages installed with pack. Then it will print a list
of packages that require update displaying current and new available version.`

var BuildShort = "🛠️  build packages"
var BuildLong = `🛠️  build packages

If no arguements provided, this script will build package in current directory.
If there are some arguements, pack will treat them as pack links, clone 
repositories, build packages inside and store them in cache, which is defined
in '~/.pack/config.yml'. You can also specify versions of packages that you
want to build after @, similarly with install.

To double check build process, you can test it inside pack docker container:
docker run --rm -it fmnx.su/core/pack b example.com/package`

var RemoveShort = "❌ remove packages"
var RemoveLong = `❌ remove packages

Use this command to remove packages from system. You can specify both pacman 
packages and pack links.

Example:
pack rm fmnx.su/core/ainst`

var UpdateShort = "🗳️  update packages"
var UpdateLong = `🗳️  update packages

You can specify packages with versions, that you need them to update to, or
provide provide just links to get latest version from default branch.

If you don't specify any arguements, all packages will be updated.

Examples:
pack update
pack update fmnx.su/core/aist@v0.21`

func init() {
	if config.DisablePrettyPrint {
		DescribeShort = strings.ReplaceAll(DescribeShort, `📝 `, ``)
		DescribeLong = strings.ReplaceAll(DescribeLong, `📝 `, ``)
		GenerateShort = strings.ReplaceAll(GenerateShort, `📋 `, ``)
		GenerateLong = strings.ReplaceAll(GenerateLong, `📋 `, ``)
		InstallShort = strings.ReplaceAll(InstallShort, `📥 `, ``)
		InstallLong = strings.ReplaceAll(InstallLong, `📥 `, ``)
		ListShort = strings.ReplaceAll(ListShort, `📄 `, ``)
		OutdatedShort = strings.ReplaceAll(OutdatedShort, `📌 `, ``)
		OutdatedLong = strings.ReplaceAll(OutdatedLong, `📌 `, ``)
		BuildShort = strings.ReplaceAll(BuildShort, `🛠️ `, ``)
		BuildLong = strings.ReplaceAll(BuildLong, `🛠️ `, ``)
		RemoveShort = strings.ReplaceAll(RemoveShort, `❌ `, ``)
		RemoveLong = strings.ReplaceAll(RemoveLong, `❌ `, ``)
		RootShort = strings.ReplaceAll(RootShort, `📦 `, ``)
		RootLong = strings.ReplaceAll(RootLong, `📦 `, ``)
		UpdateShort = strings.ReplaceAll(UpdateShort, `🗳️  `, ``)
		UpdateLong = strings.ReplaceAll(UpdateLong, `🗳️  `, ``)
	}
}
