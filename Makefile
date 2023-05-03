
docker:
	docker build -t fmnx.io/core/pack:latest .
	docker run --rm -it fmnx.io/core/pack --help

clean:
	sudo rm -rf ~/.pack
	sudo rm -f pack
	sudo rm -f /tmp/pack.lock

delpacman:
	sudo pacman -R install flutter-fmnx-package

delpack:
	pack remove fmnx.io/dev/install fmnx.io/dancheg97/flutter-fmnx-package

test:
	make clean
	pack get fmnx.io/dancheg97/flutter-fmnx-package fmnx.io/dev/install
	pack list
	pack remove fmnx.io/dancheg97/flutter-fmnx-package
	pack list
	pack get fmnx.io/dancheg97/flutter-fmnx-package
	pack remove fmnx.io/dev/install
	pack list
