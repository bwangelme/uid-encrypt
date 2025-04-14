.PHONY: build test bench

build:
	go build ./cmd/uid

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./crypto
