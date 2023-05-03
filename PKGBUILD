# Maintainer: dancheg97 <dancheg97@fmnx.io>

pkgname=pack
pkgver="1"
pkgrel="1"
pkgdesc="📦 git-based pacman-compatible package manager."
arch=("x86_64")
url="https://fmnx.io/dev/pack"
license=("GPL3")
makedepends=(
  "go"
)

package() {
  cd ..
  go build .
  install -Dm755 pack "${pkgdir}/usr/bin/pack"
}
