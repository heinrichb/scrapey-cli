# File: Makefile
.PHONY: build run

# The binary target depends on all .go files, go.mod, and go.sum.
# The find command collects all Go files in the project.
build: $(shell find . -name '*.go') go.mod go.sum
	@mkdir -p build
	@if [ ! -f build/.stamp ] || [ -n "$$(find . -name '*.go' -newer build/.stamp)" ] || [ go.mod -nt build/.stamp ] || [ go.sum -nt build/.stamp ]; then \
		echo "Changes detected, rebuilding..."; \
		go mod tidy; \
		go build -o build/scrapeycli ./cmd/scrapeycli && touch build/.stamp; \
	else \
		echo "No changes detected, skipping rebuild."; \
	fi


CONFIG_FLAG =
ifdef CONFIG
	CONFIG_FLAG := --config $(CONFIG)
endif

URL_FLAG =
ifdef URL
	URL_FLAG := --url $(URL)
endif

run: build
	./build/scrapeycli $(CONFIG_FLAG) $(URL_FLAG)

test:
	@go test ./...
	@if [ -d test ] && ls test/*.go > /dev/null 2>&1; then \
	    go test ./test; \
	fi
