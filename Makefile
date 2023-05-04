
docker:
	sudo rm -rf pkg
	sudo rm -rf src
	sudo rm -f pack
	sudo rm -f *.pkg.tar.zst
	docker build -t fmnx.io/core/pack:latest .
	docker run --rm -it fmnx.io/core/pack --help
	docker run --rm -it fmnx.io/core/pack info pack

clean:
	sudo rm -rf pkg
	sudo rm -rf src
	sudo rm -rf ~/.pack
	sudo rm -f pack
	sudo rm -f /tmp/pack.lock

delpacman:
	sudo pacman -R install flutter-fmnx-package

delpack:
	pack remove fmnx.io/dev/install fmnx.io/dancheg97/flutter-fmnx-package

test:
	docker run --rm -it -e PACK_ALLOW_AUR=true -e PACK_DEBUG_MODE=true fmnx.io/core/pack get cmake clang bluez-utils qemu-desktop edk2-ovmf archiso archinstall xdg-user-dirs-gtk fmnx.io/pkg/gnome-browser-connector fmnx.io/pkg/zsh-theme-powerlevel10k-bin-git fmnx.io/pkg/zsh-autosuggestions fmnx.io/pkg/zsh-syntax-highlighting-git fmnx.io/pkg/adw-gtk3 fmnx.io/pkg/papirus-icon-theme fmnx.io/pkg/adw-gtk-theme fmnx.io/pkg/gnome-shell-extension-dash-to-dock fmnx.io/pkg/onlyoffice-bin fmnx.io/pkg/visual-studio-code-bin fmnx.io/pkg/flutter fmnx.io/pkg/buf-bin fmnx.io/pkg/golangci-lint-bin fmnx.io/pkg/protoc-gen-go fmnx.io/pkg/protoc-gen-go-grpc fmnx.io/pkg/gnome-shell-extension-dash-to-dock fmnx.io/pkg/gnome-shell-extension-gtile fmnx.io/pkg/neovim-git fmnx.io/pkg/vlang
