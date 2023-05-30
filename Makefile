pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

serve:
	go run . serve --serve-autocert -u 'user::password' \
	 -m https://de.arch.mirror.kescher.at/core/os/x86_64/

clean:
	sudo rm -rf public
	sudo rm -rf users
	sudo rm -rf cert.pem
	sudo rm -rf key.pem

docker-compose:
	sudo rm -rf docker
	docker build -t fmnx.su/core/pack .
	docker compose up
