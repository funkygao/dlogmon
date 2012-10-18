install:
	mkdir -p var
	go install kx/dlog
	go install kx/dlogmon
	go install kx/stream
	go install kx/trace
	go install kx/progress
	@strip bin/dlogmon 2> /dev/null

test:install
	go test -v kx/dlog

fmt:
	gofmt -s -spaces=true -tabindent=false -tabwidth=4

clean:
	rm -rf bin/
	rm -rf pkg/

run:install
	./bin/dlogmon -f test/fixture/lz.121015-104410

T:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -t -d

mr:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -mapper ./mr/amfMapper.py

loc:
	@find . -name '*.go' | xargs wc -l | tail -1

help:
	@echo 'make [install | test | fmt | clean | run | mr | loc]'
