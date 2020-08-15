.PHONY: clean lint lint-fix test

export GO111MODULE=on

default: lint test

test:
	go test -v -cover ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix