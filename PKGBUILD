# Maintainer: Danila Fominykh <dancheg97@fmnx.su>

pkgname=pack
pkgver=0.1
pkgrel=1
pkgdesc="Another way of arch package distribution"
arch=('x86_64')
url="https://fmnx.su/core/pack"
license=('GPL')
depends=(
  'pacman'
  'gnupg'
  'git'
)
makedepends=('go')

build() {
  cd ..
  go build -o src/p .
}

package() {
  install -Dm755 p $pkgdir/usr/bin/pack
}
