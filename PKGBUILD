# Maintainer: Danila Fominykh <dancheg97@fmnx.su>

pkgname=pack
pkgver='0.4.5'
pkgrel=1
pkgdesc="Simplified version of pacman written in go."
arch=('x86_64')
url="https://fmnx.su/core/pack"
license=('GPL')
depends=(
  'pacman'
  'gnupg'
)
makedepends=('go')

build() {
  cd ..
  go build -o src/p .
}

package() {
  install -Dm755 p $pkgdir/usr/bin/pack
}
