
docker:
	sudo rm -rf pkg
	sudo rm -rf src
	sudo rm -f packbin
	sudo rm -f *.pkg.tar.zst
	docker build -t fmnx.io/core/pack:latest .
	docker run --rm -it fmnx.io/core/pack --help
	docker run --rm -it fmnx.io/core/pack d pack

clean:
	sudo rm -rf pkg
	sudo rm -rf src
	sudo rm -rf ~/.pack
	sudo rm -f packbin
	sudo rm -f /tmp/pack.lock

test:
	docker run --rm -it -e PACK_DEBUG_MODE=true fmnx.io/core/pack i fmnx.io/pkg/gnome-browser-connector fmnx.io/pkg/gnome-shell-extension-dash-to-dock fmnx.io/pkg/zsh-theme-powerlevel10k-bin-git fmnx.io/pkg/zsh-syntax-highlighting-git fmnx.io/pkg/zsh-autosuggestions fmnx.io/pkg/adw-gtk3 fmnx.io/pkg/flutter fmnx.io/pkg/onlyoffice-bin fmnx.io/pkg/visual-studio-code-bin fmnx.io/pkg/neovim-git fmnx.io/pkg/vlang fmnx.io/pkg/adw-gtk-theme fmnx.io/pkg/papirus-icon-theme
