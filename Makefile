install:
	go install kx/dlog
	go install kx/dlogmon

test:
	go test -v kx/dlog

clean:
	rm -f bin/*
	rm -rf pkg/*

run:install
	./bin/dlogmon -f fixture/lz.121015-104410

