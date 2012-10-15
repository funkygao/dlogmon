install:
	go install kx/dlog
	go install kx/mon

test:
	go test kx/dlog

clean:
	rm -f bin/*
	rm -rf pkg/*

