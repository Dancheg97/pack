pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

serve:
	mkdir -p tmp/public
	sudo cp /var/cache/pacman/pkg/nano* tmp/public
	go run . serve --serve-auto-tls --serve-db-path tmp --serve-dir tmp/public
