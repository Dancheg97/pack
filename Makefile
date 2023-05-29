pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

serve:
	# go run . serve --serve-autocert -u 'user::password' -m https://de.arch.mirror.kescher.at/core/os/x86_64/
	go run . serve --serve-autocert -b https://aur.archlinux.org/nvm

clean:
	sudo rm -rf public
	sudo rm -rf users
