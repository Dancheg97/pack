# PKGBUILD generated by pack.
# More info at: https://fmnx.io/core/pack

pkgname="pack"
pkgdesc="📦 decentralized package manager based on git and pacman"
pkgver="1"
pkgrel="1"
arch=('any')
url="https://fmnx.io/core/pack"
depends=(
  "git"
  "pacman"
)
makedepends=(
  "go"
)

build() {
  cd ..
  go build -buildvcs=false -o packbin .
}

package() {
  cd ..
  install -Dm755 packbin $pkgdir/usr/bin/pack
}
