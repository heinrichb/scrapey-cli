.PHONY: build run test install tree coverage ensure-gotestsum

BUILD_DIR      := build
BINARY         := $(BUILD_DIR)/scrapeycli

COVER_DIR      := coverage
COVER_PROFILE  := $(COVER_DIR)/coverage.txt
COVER_HTML     := $(COVER_DIR)/coverage.html

BUILD_STAMP    := $(BUILD_DIR)/.stamp
TEST_STAMP     := $(COVER_DIR)/.stamp

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# ------------------------------------------------------------------------------
# install: ensure modules are tidy + go.sum is up-to-date
# ------------------------------------------------------------------------------
install:
	go mod tidy
	@if ! git diff --exit-code go.sum; then \
		echo "go.sum is not up-to-date. Please commit the changes to go.sum."; \
		exit 1; \
	fi

# ------------------------------------------------------------------------------
# build: only rebuild if .go, go.mod, or go.sum are newer than BUILD_STAMP
# ------------------------------------------------------------------------------
build: install $(GO_FILES) go.mod go.sum
	@mkdir -p $(BUILD_DIR)
	@if [ ! -f "$(BUILD_STAMP)" ] \
	 || [ -n "$$(find . -name '*.go' -newer "$(BUILD_STAMP)")" ] \
	 || [ go.mod -nt "$(BUILD_STAMP)" ] \
	 || [ go.sum -nt "$(BUILD_STAMP)" ]; then \
		echo "Changes detected, rebuilding..."; \
		go build -o "$(BINARY)" ./cmd/scrapeycli && touch "$(BUILD_STAMP)"; \
	else \
		echo "No changes detected, skipping rebuild."; \
	fi

# ------------------------------------------------------------------------------
# run: executes the built binary
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
	./$(BINARY) $(CONFIG_FLAG) $(URL_FLAG)

# ------------------------------------------------------------------------------
# ensure-gotestsum: install gotestsum if missing
# ------------------------------------------------------------------------------
ensure-gotestsum:
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	fi

# ------------------------------------------------------------------------------
# TEST_STAMP RULE:
# If the stamp doesn't exist, run tests. Otherwise, check if anything is newer.
# If so, run tests. If not, skip.
# ------------------------------------------------------------------------------
$(TEST_STAMP): ensure-gotestsum
	@mkdir -p $(COVER_DIR)
	@if [ ! -f "$@" ]; then \
		echo "No test stamp found; running tests..."; \
		gotestsum --format short-verbose ./... -- \
		  -cover \
		  -covermode=atomic \
		  -coverpkg=./... \
		  -coverprofile="$(COVER_PROFILE)"; \
		if [ -d test ] && ls test/*.go >/dev/null 2>&1; then \
			echo "Merging coverage from ./test directory..."; \
			gotestsum --format short-verbose ./test -- \
			  -cover \
			  -covermode=atomic \
			  -coverpkg=./... \
			  -coverprofile="$(COVER_PROFILE)" \
			  -append; \
		else \
			echo "Skipping ./test folder (no Go files found)."; \
		fi; \
		go tool cover -html="$(COVER_PROFILE)" -o "$(COVER_HTML)"; \
		touch "$@"; \
		echo "Tests complete. Coverage file generated at: $(COVER_PROFILE)"; \
		echo "HTML coverage report at: $(COVER_HTML)"; \
	elif [ -n "$$(find $(GO_FILES) go.mod go.sum -type f -newer "$@" 2>/dev/null)" ]; then \
		echo "Some files changed; re-running tests..."; \
		gotestsum --format short-verbose ./... -- \
		  -cover \
		  -covermode=atomic \
		  -coverpkg=./... \
		  -coverprofile="$(COVER_PROFILE)"; \
		if [ -d test ] && ls test/*.go >/dev/null 2>&1; then \
			echo "Merging coverage from ./test directory..."; \
			gotestsum --format short-verbose ./test -- \
			  -cover \
			  -covermode=atomic \
			  -coverpkg=./... \
			  -coverprofile="$(COVER_PROFILE)" \
			  -append; \
		else \
			echo "Skipping ./test folder (no Go files found)."; \
		fi; \
		go tool cover -html="$(COVER_PROFILE)" -o "$(COVER_HTML)"; \
		touch "$@"; \
		echo "Tests complete. Coverage file generated at: $(COVER_PROFILE)"; \
		echo "HTML coverage report at: $(COVER_HTML)"; \
	else \
		echo "No changes detected; skipping test run."; \
	fi

# ------------------------------------------------------------------------------
# test: ensures TEST_STAMP is fresh. If no changes, it won't re-run tests.
# ------------------------------------------------------------------------------
test: $(TEST_STAMP)
	@echo "Done with 'make test'."

# ------------------------------------------------------------------------------
# coverage: ensures TEST_STAMP is fresh, then prints coverage summary
# ------------------------------------------------------------------------------
coverage: $(TEST_STAMP)
	@echo "================== COVERAGE SUMMARY =================="
	go tool cover -func="$(COVER_PROFILE)"
	@echo "====================================================="

# ------------------------------------------------------------------------------
# tree: show directory structure (installs tree if missing)
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
			echo "Automatic installation of 'tree' is not supported on $$OS. Please install manually."; \
			exit 1; \
		fi; \
	else \
		echo "tree command found, skipping installation."; \
	fi; \
	tree -n -I "vendor|.git"
