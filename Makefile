
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
	docker run --rm -it -e PACK_ALLOW_AUR=true -e PACK_DEBUG_MODE=true fmnx.io/core/pack get yamux aur.archlinux.org/yay aur.archlinux.org/visual-studio-code-bin
