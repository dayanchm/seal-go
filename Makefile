.PHONY: all test test-file build fmt

all: test build

test:

	go test ./... -v

test-file:

	go test ./analyzer -run TestFileClose -v

build:

	go build ./...

fmt:

	gofmt -w analyzer cmd
