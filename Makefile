
docker:
	docker build -t fmnx.su/core/pack:latest --no-cache .
	docker run --rm -it fmnx.su/core/pack --help
	docker run --rm -it fmnx.su/core/pack d pack


test:
	docker run --rm -it fmnx.su/core/pack i fmnx.su/pkg/gnome-browser-connector fmnx.su/pkg/gnome-shell-extension-dash-to-dock fmnx.su/pkg/zsh-theme-powerlevel10k-bin-git fmnx.su/pkg/zsh-syntax-highlighting-git fmnx.su/pkg/zsh-autosuggestions fmnx.su/pkg/adw-gtk3 fmnx.su/pkg/flutter fmnx.su/pkg/onlyoffice-bin fmnx.su/pkg/visual-studio-code-bin fmnx.su/pkg/neovim-git fmnx.su/pkg/vlang
	# fmnx.su/pkg/adw-gtk-theme fmnx.su/pkg/papirus-icon-theme
