.PHONY: build test lint test-lab

build:
	go build -o bin/spar ./cmd/spar

test:
	go test ./...

lint:
	golangci-lint run

test-lab:
	go test -tags lab ./...
