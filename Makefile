.PHONY: test

all: generate fmt test

generate:
	go generate ./...
	gofmt -w .

fmt:
	gofmt -w .

test:
	go test -v ./tests/...

test-all:
	go test -v ./...
	go test -bench=. ./...
