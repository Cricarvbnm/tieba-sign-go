pkgname=tieba-sign
pkgver=1.0.0
pkgrel=1

arch=(any)

_repo_dir="tieba-sign-go-$pkgver-systemd"

makedepends=(go)
build() {
  cd "$_repo_dir" || exit 1
  go build
}

package() {
  install -Dt "$pkgdir/usr/bin" "$_repo_dir/tieba-sign"
  install -Dt "$pkgdir/usr/lib/systemd/user" "$_repo_dir"/tieba-sign.{service,timer}
}

source=(https://github.com/Cricarvbnm/tieba-sign-go/archive/refs/tags/v1.0.0-systemd.zip)
sha256sums=('0c00e454f70ea339a1faa82c6e31574ca8b6f63045789a44cfc87dcf78017843')
