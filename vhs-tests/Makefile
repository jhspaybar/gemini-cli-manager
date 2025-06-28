.PHONY: all test-all test-navigation test-extensions test-profiles test-state-dir clean setup help

# Default target
all: help

# Setup test environment
setup:
	@echo "Setting up VHS test environment..."
	@mkdir -p output
	@mkdir -p temp
	@./helpers/setup-test-env.sh

# Run all tests
test-all: setup test-navigation test-extensions test-profiles test-state-dir
	@echo "✅ All VHS tests completed!"
	@echo "View results in output/ directory"

# Navigation tests
test-navigation: setup
	@echo "Running navigation tests..."
	@for tape in tests/navigation/*.tape; do \
		echo "  Running $$tape..."; \
		vhs $$tape || exit 1; \
	done

# Extension tests  
test-extensions: setup
	@echo "Running extension tests..."
	@for tape in tests/extensions/*.tape; do \
		echo "  Running $$tape..."; \
		vhs $$tape || exit 1; \
	done

# Profile tests
test-profiles: setup
	@echo "Running profile tests..."
	@for tape in tests/profiles/*.tape; do \
		echo "  Running $$tape..."; \
		vhs $$tape || exit 1; \
	done

# State directory tests
test-state-dir: setup
	@echo "Running state directory tests..."
	@for tape in tests/state-dir/*.tape; do \
		echo "  Running $$tape..."; \
		vhs $$tape || exit 1; \
	done

# Run a specific test
test: setup
	@if [ -z "$(TAPE)" ]; then \
		echo "Usage: make test TAPE=tests/navigation/sidebar.tape"; \
		exit 1; \
	fi
	@echo "Running $(TAPE)..."
	@vhs $(TAPE)

# Clean generated files
clean:
	@echo "Cleaning up..."
	@rm -rf output/*
	@rm -rf temp/*
	@find . -name "*.gif" -type f -delete
	@find . -name "*.png" -type f -delete
	@echo "✅ Cleaned"

# Help
help:
	@echo "VHS Visual Testing Suite"
	@echo ""
	@echo "Available targets:"
	@echo "  make setup           - Set up test environment"
	@echo "  make test-all        - Run all tests"
	@echo "  make test-navigation - Run navigation tests"
	@echo "  make test-extensions - Run extension tests"
	@echo "  make test-profiles   - Run profile tests"
	@echo "  make test-state-dir  - Run state directory tests"
	@echo "  make test TAPE=...   - Run specific tape file"
	@echo "  make clean           - Clean generated files"
	@echo "  make help            - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make test-all"
	@echo "  make test TAPE=tests/navigation/sidebar.tape"