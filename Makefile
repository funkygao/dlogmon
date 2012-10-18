PKGS = kx/dlogmon kx/db kx/dlog kx/progress kx/sb kx/stream kx/trace
SRC = src
BIN = bin
PKG = pkg
VAR = var

install:mkvar
	go install ${PKGS}
	@strip bin/dlogmon 2> /dev/null

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dlogmon_linux kx/dlogmon

test:install
	go test ${PKGS}

fmt:
	gofmt -s -tabs=false -tabwidth=4 -w=true ${SRC}

clean:
	rm -rf ${BIN} ${PKG} ${VAR}

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
	@mkdir -p ${VAR}

