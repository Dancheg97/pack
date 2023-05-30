pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

serve:
	go run . serve --serve-autocert -u 'user::password' \
	 -m https://de.arch.mirror.kescher.at/core/os/x86_64/ \
	 -g https://aur.archlinux.org/yay

clean:
	sudo rm -rf public
	sudo rm -rf users

docker-compose:
	sudo rm -rf docker
	docker build -t fmnx.su/core/pack .
	docker compose up
