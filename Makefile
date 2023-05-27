pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

serve:
	mkdir -p tmp
	sudo cp /var/cache/pacman/pkg/nano* tmp
	go run . serve --serve-dir tmp

servetls:
	mkdir -p tmp
	sudo cp /var/cache/pacman/pkg/nano* tmp
	openssl req -x509 -newkey rsa:4096 -keyout tmp/key.pem -out tmp/cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"
	go run . serve --serve-dir tmp --serve-key tmp/key.pem --serve-cert tmp/cert.pem
