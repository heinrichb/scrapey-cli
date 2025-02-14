# File: Makefile

.PHONY: build run test install tree

# The install target runs go mod tidy and verifies that go.sum is up-to-date.
install:
	go mod tidy
	@if ! git diff --exit-code go.sum; then \
		echo "go.sum is not up-to-date. Please commit the changes to go.sum."; \
		exit 1; \
	fi

# The build target depends on source files, go.mod, and go.sum.
# It rebuilds the binary only when changes are detected.
build: install $(shell find . -name '*.go') go.mod go.sum
	@mkdir -p build
	@if [ ! -f build/.stamp ] || [ -n "$$(find . -name '*.go' -newer build/.stamp)" ] || [ go.mod -nt build/.stamp ] || [ go.sum -nt build/.stamp ]; then \
		echo "Changes detected, rebuilding..."; \
		go build -o build/scrapeycli ./cmd/scrapeycli && touch build/.stamp; \
	else \
		echo "No changes detected, skipping rebuild."; \
	fi

# Set optional command-line flags for the run target.
CONFIG_FLAG =
ifdef CONFIG
	CONFIG_FLAG := --config $(CONFIG)
endif

URL_FLAG =
ifdef URL
	URL_FLAG := --url $(URL)
endif

# The run target executes the built binary with any provided configuration or URL flags.
run: build
	./build/scrapeycli $(CONFIG_FLAG) $(URL_FLAG)

# The test target ensures that gotestsum is installed and uses it to run tests with coverage reporting.
test:
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "gotestsum not found, installing..."; \
		go install gotest.tools/gotestsum@latest; \
	fi; \
	gotestsum --format short-verbose ./... -- -coverprofile=./coverage/coverage.txt; \
	if [ -d test ] && ls test/*.go > /dev/null 2>&1; then \
	    gotestsum --format short-verbose ./test -- -coverprofile=./coverage/coverage.txt; \
	fi; \
	go tool cover -func=./coverage/coverage.txt

# The tree target displays the project directory structure.
# If the 'tree' command is not installed, it attempts to install it based on the OS.
tree:
	@if ! command -v tree >/dev/null 2>&1; then \
		echo "tree command not found. Attempting to install..."; \
		OS=$$(uname); \
		if [ "$$OS" = "Linux" ]; then \
			sudo apt-get update && sudo apt-get install -y tree; \
		elif [ "$$OS" = "Darwin" ]; then \
			brew install tree; \
		else \
			echo "Automatic installation not supported on $$OS. Please install tree manually."; \
			exit 1; \
		fi; \
	else \
		echo "tree command found, skipping installation."; \
	fi; \
	tree -n -I "vendor|.git"
