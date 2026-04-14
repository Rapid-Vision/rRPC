.PHONY: all test staticcheck lint

all: test staticcheck lint

test:
	go test ./...

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run 
