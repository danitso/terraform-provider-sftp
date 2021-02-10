GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: build

build:
	goreleaser build --parallelism 2 --rm-dist --snapshot --timeout 1h

fmt:
	gofmt -s -w $(GOFMT_FILES)

init:
	go get ./...

test:
	go test -v ./...

.PHONY: build fmt init test
