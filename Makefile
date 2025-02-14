# File: Makefile
# Purpose:
#   - "install": ensure modules are tidy if go.mod or go.sum changed.
#   - "build": recompile binary if Go source or mod files changed.
#   - "test": run coverage if code changed.
#   - "coverage": display coverage summary from coverage.out
#   - "run": executes the compiled binary
#   - All skip with "No changes detected, skipping X." if nothing changed.

.PHONY: install build run test coverage ensure-gotestsum tree

# Directories for build artifacts and stamp files
BUILD_DIR       := build
STAMPS_DIR      := $(BUILD_DIR)/.stamps

# Stamp file paths for each step
INSTALL_STAMP   := $(STAMPS_DIR)/install.stamp
BUILD_STAMP     := $(STAMPS_DIR)/build.stamp
TEST_STAMP      := $(STAMPS_DIR)/test.stamp

BINARY          := $(BUILD_DIR)/scrapeycli

# Coverage output
COVER_DIR       := coverage
COVER_PROFILE   := $(COVER_DIR)/coverage.txt
COVER_HTML      := $(COVER_DIR)/coverage.html

# All Go source files (including _test.go files)
GO_FILES        := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Reusable messages
SKIP_MSG        := No changes detected, skipping
CHANGE_MSG      := Changes detected

# ------------------------------------------------------------------------------
# install: Ensure go.mod/go.sum are tidy.
# This target always runs its recipe (via a phony dependency) and prints either
# "Changes detected, running 'go mod tidy' for install..." or "No changes detected, skipping install."
# ------------------------------------------------------------------------------
install: $(INSTALL_STAMP)

$(INSTALL_STAMP): go.mod go.sum
	@mkdir -p $(STAMPS_DIR)
	@TARGET=install; \
	if [ ! -f "$@" ] || [ go.mod -nt "$@" ] || [ go.sum -nt "$@" ]; then \
		echo "$(CHANGE_MSG), running 'go mod tidy' for $$TARGET..."; \
		go mod tidy; \
		if ! git diff --exit-code go.sum; then \
			echo "go.sum updated. Please commit the changes."; \
			exit 1; \
		fi; \
		touch "$@"; \
		echo "Done with 'make $$TARGET'."; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

# ------------------------------------------------------------------------------
# build: Rebuild binary if any .go file or go.mod/go.sum changed.
# ------------------------------------------------------------------------------
build: $(BUILD_STAMP)

$(BUILD_STAMP): $(INSTALL_STAMP) $(GO_FILES) go.mod go.sum
	@mkdir -p $(BUILD_DIR) $(STAMPS_DIR)
	@TARGET=build; \
	if [ ! -f "$@" ] || [ -n "$$(find . -name '*.go' -newer "$@")" ] || [ go.mod -nt "$@" ] || [ go.sum -nt "$@" ]; then \
		echo "$(CHANGE_MSG), rebuilding $$TARGET..."; \
		go build -o "$(BINARY)" ./cmd/scrapeycli; \
		touch "$@"; \
		echo "Done with 'make $$TARGET'."; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

# ------------------------------------------------------------------------------
# run: Execute the compiled binary.
# ------------------------------------------------------------------------------
CONFIG_FLAG =
ifdef CONFIG
	CONFIG_FLAG := --config $(CONFIG)
endif
URL_FLAG =
ifdef URL
	URL_FLAG := --url $(URL)
endif
run: build
	@./$(BINARY) $(CONFIG_FLAG) $(URL_FLAG)

# ------------------------------------------------------------------------------
# ensure-gotestsum: Installs gotestsum if missing.
# ------------------------------------------------------------------------------
ensure-gotestsum:
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	fi

# ------------------------------------------------------------------------------
# test: Run tests (and update coverage) if any Go files (including _test.go) have changed.
# ------------------------------------------------------------------------------
test: $(TEST_STAMP)

$(TEST_STAMP): ensure-gotestsum
	@mkdir -p $(COVER_DIR) $(STAMPS_DIR)
	@TARGET=test; \
	if [ ! -f "$@" ]; then \
		echo "No $$TARGET stamp found; running $$TARGET..."; \
		$(MAKE) --no-print-directory do-coverage-run; \
	elif [ -n "$$(find $(GO_FILES) go.mod go.sum -type f -newer "$@" 2>/dev/null)" ]; then \
		echo "Some files changed; re-running $$TARGET..."; \
		$(MAKE) --no-print-directory do-coverage-run; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

.PHONY: do-coverage-run
do-coverage-run:
	gotestsum --format short-verbose ./... -- -cover -covermode=atomic -coverpkg=./... -coverprofile="$(COVER_PROFILE)"
	@if [ -d test ] && ls test/*.go >/dev/null 2>&1; then \
		echo "Merging coverage from ./test directory..."; \
		gotestsum --format short-verbose ./test -- -cover -covermode=atomic -coverpkg=./... -coverprofile="$(COVER_PROFILE)" -append; \
	else \
		echo "Skipping ./test folder (no Go files found)."; \
	fi
	go tool cover -html="$(COVER_PROFILE)" -o "$(COVER_HTML)"
	touch "$(TEST_STAMP)"
	echo "Done with 'make test'."

# ------------------------------------------------------------------------------
# coverage: Display a colorized coverage summary.
# ------------------------------------------------------------------------------
coverage: test
	@echo "================== COVERAGE SUMMARY =================="
	@go tool cover -func="$(COVER_PROFILE)" | go run ./scripts/coverage_formatter.go
	@echo "====================================================="

# ------------------------------------------------------------------------------
# tree: Display directory structure (installs tree if missing)
# ------------------------------------------------------------------------------
tree:
	@if ! command -v tree >/dev/null 2>&1; then \
		echo "tree command not found. Attempting to install..."; \
		OS=$$(uname); \
		if [ "$$OS" = "Linux" ]; then \
			sudo apt-get update && sudo apt-get install -y tree; \
		elif [ "$$OS" = "Darwin" ]; then \
			brew install tree; \
		else \
			echo "Automatic installation not supported on $$OS. Please install manually."; \
			exit 1; \
		fi; \
	else \
		echo "tree command found, skipping installation."; \
	fi; \
	tree -n -I "vendor|.git"
