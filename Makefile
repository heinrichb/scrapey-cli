# File: Makefile

.PHONY: build run test install

install:
	go mod tidy
	@if ! git diff --exit-code go.sum; then \
		echo "go.sum is not up-to-date. Please commit the changes to go.sum."; \
		exit 1; \
	fi

build: install $(shell find . -name '*.go') go.mod go.sum
	@mkdir -p build
	@if [ ! -f build/.stamp ] || [ -n "$$(find . -name '*.go' -newer build/.stamp)" ] || [ go.mod -nt build/.stamp ] || [ go.sum -nt build/.stamp ]; then \
		echo "Changes detected, rebuilding..."; \
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
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "gotestsum not found, installing..."; \
		go install gotest.tools/gotestsum@latest; \
	fi; \
	gotestsum --format short-verbose ./... -- -coverprofile=./coverage/coverage.txt; \
	if [ -d test ] && ls test/*.go > /dev/null 2>&1; then \
	    gotestsum --format short-verbose ./test -- -coverprofile=./coverage/coverage.txt; \
	fi; \
	# Print a coverage summary using go tool cover.
	go tool cover -func=./coverage/coverage.txt
