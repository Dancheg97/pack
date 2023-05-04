
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
	docker run --rm -it -e PACK_ALLOW_AUR=true -e PACK_DEBUG_MODE=true fmnx.io/core/pack get cmake clang qemu-desktop edk2-ovmf archiso archinstall aur.archlinux.org/zsh-theme-powerlevel10k-bin-git aur.archlinux.org/zsh-autosuggestions aur.archlinux.org/zsh-syntax-highlighting-git xdg-user-dirs-gtk aur.archlinux.org/adw-gtk3 aur.archlinux.org/papirus-icon-theme aur.archlinux.org/adw-gtk-theme aur.archlinux.org/gnome-browser-connector aur.archlinux.org/gnome-shell-extension-dash-to-dock aur.archlinux.org/onlyoffice-bin aur.archlinux.org/visual-studio-code-bin aur.archlinux.org/flutter aur.archlinux.org/buf-bin aur.archlinux.org/golangci-lint-bin aur.archlinux.org/protoc-gen-go-grpc aur.archlinux.org/bluez-utils aur.archlinux.org/gnome-shell-extension-dash-to-dock aur.archlinux.org/gnome-shell-extension-gtile aur.archlinux.org/neovim-git aur.archlinux.org/vlang
