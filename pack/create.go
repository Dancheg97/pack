// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"fmt"
	"os"
	"strings"
)

func Create(pkgname string) error {
	name := "example"
	title := "Example"
	if pkgname != "" {
		name = pkgname
		title = strings.ToTitle(pkgname)
	}
	detmpl := fmt.Sprintf(desktop, title, title, title, name, name, name)
	err := os.WriteFile("example.desktop", []byte(detmpl), 0600)
	if err != nil {
		return err
	}

	shtmpl := fmt.Sprintf(shell, name)
	err = os.WriteFile("example.sh", []byte(shtmpl), 0600)
	if err != nil {
		return err
	}

	pkgbldtmpl := fmt.Sprintf(
		pkgbuikd, name, name, name, name, name, name, name, name, name, name,
	)
	err = os.WriteFile("PKGBUILD", []byte(pkgbldtmpl), 0600)
	if err != nil {
		return err
	}

	fmt.Println(message)
	return nil
}

const desktop = `[Desktop Entry]
Name=%s
GenericName=%s
Comment=%s
Exec=/usr/bin/%s
WMClass=%s
Icon=/usr/share/%s/icons/icon.png
Type=Application
`

const shell = `#!/usr/bin/env sh
exec /usr/bin/%s
`

const pkgbuikd = `# Maintainer: Real Name <maintainer@fmnx.su>

pkgname=%s
pkgver=1
pkgrel=1
pkgdesc="Very helpfull one-line description"
arch=("x86_64")
url="https://fmnx.su/maintainer/%s"
license=('GPL')
depends=(
  "vlc"
)
makedepends=(
  "ninja"
  "clang"
  "cmake"
)
source=("$pkgname::git+https://fmnx.su/maintainer/%s.git")
md5sums=("SKIP")


build() {
  cd "$pkgname"
  flutter build linux
}

package() {
  cd "$pkgname"
  install -Dm755 %s.sh $pkgdir/usr/bin/%s
  install -Dm755 %s.desktop $pkgdir/usr/share/applications/%s.desktop
  install -Dm755 assets/%s.png $pkgdir/usr/share/icons/hicolor/512x512/apps/%s.png
  cd build/linux/x64/release/bundle && find . -type f -exec install -Dm755 {} $pkgdir/usr/share/%s/{} \; && cd $srcdir/..
}`

const message = `Done!
You can adjust the contents of created templates for your package.

Created files:
- PKGBUILD
- example.sh
- example.desktop`
