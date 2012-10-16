install:
	go install kx/dlog
	go install kx/dlogmon
	go install kx/stream
	@strip bin/dlogmon 2> /dev/null

test:install
	go test -v kx/dlog

clean:
	rm -rf bin/
	rm -rf pkg/

run:install
	./bin/dlogmon -f fixture/lz.121015-104410

loc:
	@find . -name '*.go' | xargs wc -l | tail -1

help:
	@echo 'make [install | test | clean | run | loc]'
