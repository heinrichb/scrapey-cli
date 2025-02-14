# File: Makefile
.PHONY: build run test

build:
	go mod tidy
	go build -o build/scrapeycli ./cmd/scrapeycli

CONFIG_FLAG =
ifdef CONFIG
	CONFIG_FLAG := --config $(CONFIG)
endif

URL_FLAG =
ifdef URL
	URL_FLAG := --url $(URL)
endif

run:
	./build/scrapeycli $(CONFIG_FLAG) $(URL_FLAG)

test:
	@go test ./...
	@if [ -d test ] && ls test/*.go > /dev/null 2>&1; then \
	    go test ./test; \
	fi
