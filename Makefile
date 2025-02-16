# File: Makefile
# Purpose:
#   - "install": ensure modules are tidy if go.mod or go.sum changed.
#   - "test": run coverage if code changed.
#   - "coverage": display coverage summary from coverage.out
#   - "build": recompile binary if Go source or mod files changed.
#   - "run": executes the compiled binary
#   - "tree": display directory structure
#   - All skip with "No changes detected, skipping X." if nothing changed.

.PHONY: install test coverage build run tree

# Directories for build artifacts and stamp files
BUILD_DIR       := build
STAMPS_DIR      := $(BUILD_DIR)/.stamps

# Stamp file paths for each step
INSTALL_STAMP   := $(STAMPS_DIR)/install.stamp
BUILD_STAMP     := $(STAMPS_DIR)/build.stamp
TEST_STAMP      := $(STAMPS_DIR)/test.stamp

BINARY          := $(BUILD_DIR)/scrapeycli

# Coverage output
COVER_DIR       := ${BUILD_DIR}/coverage
COVER_PROFILE   := $(COVER_DIR)/coverage.txt
COVER_HTML      := $(COVER_DIR)/coverage.html

# All Go source files (including _test.go files)
GO_FILES        := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Reusable messages
SKIP_MSG        := No changes detected, skipping
CHANGE_MSG      := Some files changed; re-running

# ------------------------------------------------------------------------------
# install: ensure go.mod/go.sum are tidy if changed.
# ------------------------------------------------------------------------------
install:
	@mkdir -p $(STAMPS_DIR)
	@TARGET=install; \
	if [ ! -f "$(INSTALL_STAMP)" ] || [ go.mod -nt "$(INSTALL_STAMP)" ] || [ go.sum -nt "$(INSTALL_STAMP)" ]; then \
		echo "$(CHANGE_MSG) $$TARGET..."; \
		go mod tidy; \
		if ! git diff --exit-code go.sum; then \
			echo "go.sum updated. Please commit the changes."; \
			exit 1; \
		fi; \
		touch "$(INSTALL_STAMP)"; \
		echo "Done with installing."; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

# ------------------------------------------------------------------------------
# test: run tests and update coverage if any Go source (including _test.go files) have changed.
# Ensures gotestsum is installed before running tests.
# Depends on install.
# ------------------------------------------------------------------------------
test: install
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	fi
	@mkdir -p $(COVER_DIR) $(STAMPS_DIR)
	@TARGET=test; \
	if [ ! -f "$(TEST_STAMP)" ] || [ -n "$$(find $(GO_FILES) -newer "$(TEST_STAMP)" 2>/dev/null)" ]; then \
		echo "$(CHANGE_MSG) $$TARGET..."; \
		> "$(COVER_PROFILE)"; \
		if gotestsum --format short-verbose ./... && \
		   go test -cover -covermode=atomic -coverpkg=./... -coverprofile="$(COVER_PROFILE)" ./... >/dev/null; then \
			if [ -f "$(COVER_PROFILE)" ]; then \
				grep -v "cmd/scrapeycli/main.go:" "$(COVER_PROFILE)" > "$(COVER_PROFILE).tmp" && mv "$(COVER_PROFILE).tmp" "$(COVER_PROFILE)"; \
				go tool cover -html="$(COVER_PROFILE)" -o "$(COVER_HTML)"; \
				echo "Coverage file generated at: $(COVER_PROFILE)"; \
				echo "HTML coverage report at: $(COVER_HTML)"; \
			else \
				echo "ERROR: Coverage file was not generated!"; \
			fi; \
			touch "$(TEST_STAMP)"; \
		else \
			echo "Tests failed! Skipping stamp update."; \
			exit 1; \
		fi; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

# ------------------------------------------------------------------------------
# coverage: displays a colorized coverage summary from the coverage file.
# Depends on test.
# ------------------------------------------------------------------------------
coverage: test
	@echo "================== COVERAGE SUMMARY =================="
	@go tool cover -func="$(COVER_PROFILE)" | go run ./scripts/coverage_formatter.go
	@echo "====================================================="

# ------------------------------------------------------------------------------
# build: compile binary if Go sources changed.
# Depends on install.
# ------------------------------------------------------------------------------
build: install
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(STAMPS_DIR)
	@TARGET=build; \
	if [ ! -f "$(BUILD_STAMP)" ] || [ -n "$$(find $(GO_FILES) -newer "$(BUILD_STAMP)" 2>/dev/null)" ]; then \
		echo "$(CHANGE_MSG) $$TARGET..."; \
		go build -o $(BINARY) ./cmd/scrapeycli; \
		touch "$(BUILD_STAMP)"; \
		echo "Done with building."; \
	else \
		echo "$(SKIP_MSG) $$TARGET."; \
	fi

# ------------------------------------------------------------------------------
# run: execute the compiled binary.
# Depends on build.
# ------------------------------------------------------------------------------
run: build
	@echo "Running application..."
	@$(BINARY)

# ------------------------------------------------------------------------------
# tree: displays directory structure (installs tree if missing).
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
	fi; \
	tree -n -I "vendor|.git"