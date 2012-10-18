PKGS = kx/dlogmon kx/dlog kx/stream kx/sb kx/stream kx/progress kx/db
SRC = src

install:mkvar
	go install ${PKGS}
	@strip bin/dlogmon 2> /dev/null

test:install
	go test ${PKGS}

fmt:
	gofmt -s -tabs=false -tabwidth=4 -w=true ${SRC}

clean:
	rm -rf bin/ pkg/ var/

run:install
	./bin/dlogmon -f test/fixture/lz.121015-104410

T:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -t -d

mr:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -mapper ./mr/amfMapper.py

loc:
	@find src/kx -name '*.go' | xargs wc -l | tail -1

help:
	@echo 'make [install | test | fmt | clean | run | mr | loc]'

mkvar:
	@mkdir -p var

