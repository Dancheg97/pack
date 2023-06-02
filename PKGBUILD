# Maintainer: Danila Fominykh <dancheg97@fmnx.su>

pkgname=pack
pkgver=0.1
pkgrel=1
pkgdesc="Speed up your package related operations."
arch=('x86_64')
url="https://fmnx.su/core/pack"
license=('GPL')
depends=(
  'pacman'
  'git'
  'wget'
)
makedepends=('go')

build() {
  cd ..
  go build -o src/p .
}

package() {
  install -Dm755 p $pkgdir/usr/bin/pack
}
