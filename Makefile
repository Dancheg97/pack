
self:
	sudo rm -rf ~/.pack
	sudo rm -f pack
	sudo rm -f /tmp/pack.lock
	go run . get fmnx.io/dev/pack  
