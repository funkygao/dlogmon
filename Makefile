PKGS = kx/dlogmon \
	   kx/db \
	   kx/dlog \
	   kx/logger \
	   kx/progress \
	   kx/sb \
	   kx/size \
	   kx/stream \
	   kx/trace \
	   kx/cache \
	   kx/util \
	   kx/stats
SRC = src
BIN = bin
PKG = pkg
VAR = var

install:mkvar
	go install ${PKGS}

linux:
	@echo 'cd /usr/local/go/src; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ./make.bash --no-clean'
	@echo 'cd ~/github/dlogmon;  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dlogmon_linux kx/dlogmon'

up:
	go get -u github.com/bmizerany/assert
	go get -u github.com/kless/goconfig/config
	go get -u github.com/mattn/go-sqlite3

rb:clean install

dep:
	@find src/kx -name '*.go' | xargs grep -e 'github.com' -e 'code.google.com'

test:install
	@go test ${PKGS}

bench:
	go test -test.bench=".*" -test.benchtime 5 kx/dlog

fmt:
	@gofmt -s -tabs=false -tabwidth=4 -w=true ${SRC}

clean:
	rm -rf ${BIN} ${PKG} ${VAR}

run:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -tick 500 -progress -cpuprofile var/cpu.prof -memprofile var/mem.prof

real:install
	./bin/dlogmon -tick 2000 -progress -d

prof:run
	@go tool pprof ./bin/dlogmon var/cpu.prof

trace:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -t -d

mr:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -d -mapper ./mr/amfMapper.py -progress

loc:
	@echo `find src/kx -name '*.go' | xargs wc -l | tail -1` lines
	@echo `find src/kx -name '*.go' | wc -l | tail -1` files

mkvar:
	@mkdir -p ${VAR}

help:
	@echo 'make [install | test | bench | fmt | clean | run | prof | trace | mr | loc]'

