
self:
	sudo rm -r ~/.pack
	sudo rm -f pack
	sudo rm /tmp/pack.lock
	go run . get fmnx.io/dev/pack  
