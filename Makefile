install:
	go install kx/dlog
	go install kx/mon

test:
	go test -v kx/dlog

clean:
	rm -f bin/*
	rm -rf pkg/*

run:install
	./bin/mon -f fixture/lz.121015-104410

