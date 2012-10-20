PKGS = kx/dlogmon kx/db kx/dlog kx/log kx/progress kx/sb kx/size kx/stream kx/trace
SRC = src
BIN = bin
PKG = pkg
VAR = var

install:mkvar
	go install ${PKGS}
	@strip bin/dlogmon 2> /dev/null

linux:
	@echo 'cd /usr/local/go/src; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ./make.bash --no-clean'
	@echo 'cd ~/github/dlogmon;  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dlogmon_linux kx/dlogmon'

test:install
	go test ${PKGS}

bench:
	go test -test.bench=".*" -test.benchtime 5 kx/dlog

fmt:
	gofmt -s -tabs=false -tabwidth=4 -w=true ${SRC}

clean:
	rm -rf ${BIN} ${PKG} ${VAR}

run:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -tick 300 -cpuprofile var/cpu.prof -memprofile var/mem.prof

prof:run
	go tool pprof ./bin/dlogmon var/cpu.prof

trace:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -t -d

mr:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -mapper ./mr/amfMapper.py

loc:
	@echo `find src/kx -name '*.go' | xargs wc -l | tail -1` lines
	@echo `find src/kx -name '*.go' | wc -l | tail -1` files

mkvar:
	@mkdir -p ${VAR}

help:
	@echo 'make [install | test | bench | fmt | clean | run | prof | trace | mr | loc]'

