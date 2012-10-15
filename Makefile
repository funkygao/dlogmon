install:clean
	go install kx/dlog
	go install kx/dlogmon

test:install
	go test -v kx/dlog

clean:
	rm -f bin/*
	rm -rf pkg/*

run:install
	./bin/dlogmon -f fixture/lz.121015-104410

loc:
	find . -name '*.go' | xargs wc -l | tail -1
