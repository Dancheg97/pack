# Maintainer: Danila Fominykh <dancheg97@fmnx.su>

pkgname=pack
pkgver='0.6.2'
pkgrel=1
pkgdesc="Decentralized package manager based on pacman."
arch=('x86_64')
url="https://fmnx.su/core/pack"
license=('GPL')
depends=(
  'sudo'
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
