PKGS = kx/dlogmon \
	   kx/db \
	   kx/dlog \
	   kx/logger \
	   kx/mr \
	   kx/progress \
	   kx/sb \
	   kx/size \
	   kx/stream \
	   kx/trace \
	   kx/cache \
	   kx/util \
	   kx/stats \
	   kx/netapi
SRC = src
BIN = bin
PKG = pkg
VAR = var

install:mkvar
	@go install ${PKGS}

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

benchmark:
	go test -v -test.bench=".*" github.com/feyeleanor/gospeed

fmt:
	@gofmt -s -tabs=false -tabwidth=4 -w=true ${SRC}

clean:
	rm -rf ${BIN} ${PKG} ${VAR}

kxi:install
	@./bin/dlogmon -f test/fixture/lz.121015-104410 -d -k kxi -tick 500 -progress -cpuprofile var/cpu.prof -memprofile var/mem.prof

amf:install
	@./bin/dlogmon -f test/fixture/lz.121015-104410 -tick 500 -progress -cpuprofile var/cpu.prof -memprofile var/mem.prof

file:install
	@./bin/dlogmon -f README.rst -filemode -d -tick 500 -progress -k file -progress

noop:install
	@./bin/dlogmon -k noop -f test/fixture/lz.121015-104410 -d -tick 500 -progress -cpuprofile var/cpu.prof -memprofile var/mem.prof

real:install
	@./bin/dlogmon -tick 30000 -progress -n 50

nreal:install
	@./bin/dlogmon -k noop -tick 30000 -progress -n 50

prof:
	@go tool pprof ./bin/dlogmon var/cpu.prof

mprof:
	@go tool pprof ./bin/dlogmon var/mem.prof

trace:install
	./bin/dlogmon -f test/fixture/lz.121015-104410 -t -d

mr:install
	@./bin/dlogmon -f test/fixture/lz.121015-104410 -mapper ./contrib/kxiMapper.py -progress=true -k kxi

todo:
	@find src/kx -name '*.go' | xargs grep -n -1 --color TODO

loc:
	@echo `find src/kx -name '*.go' | xargs wc -l | tail -1` lines
	@echo `find src/kx -name '*.go' | wc -l | tail -1` files

mkvar:
	@mkdir -p ${VAR}

help:
	@echo 'make [install | test | bench | fmt | clean | run | prof | trace | mr | loc]'

