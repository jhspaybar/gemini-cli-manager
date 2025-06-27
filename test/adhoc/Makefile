# Makefile for visual tests
.PHONY: help list all test clean

# Default target
help:
	@echo "Visual Test Commands:"
	@echo "  make list      - List all available visual tests"
	@echo "  make all       - Run all visual tests"
	@echo "  make test      - Run visual tests as Go tests"
	@echo "  make tabs      - Run tabs visual test"
	@echo "  make emoji     - Run emoji width test"
	@echo "  make clean     - Clean backup files"
	@echo ""
	@echo "Run specific test: make <test-name>"
	@echo "Example: make tabs-dynamic"

# List available tests
list:
	@go run ../../cmd/visual-tests/main.go --list

# Run all tests
all:
	@go run ../../cmd/visual-tests/main.go all

# Run as Go tests
test:
	@echo "Running visual tests (set VISUAL_TESTS=true to see output):"
	@VISUAL_TESTS=true go test -v ./. -run TestVisual

# Individual test targets
tabs:
	@go run ../../cmd/visual-tests/main.go tabs

tabs-dynamic:
	@go run ../../cmd/visual-tests/main.go tabs-dynamic

tabs-overflow:
	@go run ../../cmd/visual-tests/main.go tabs-overflow

tabs-switch:
	@go run ../../cmd/visual-tests/main.go tabs-switch

emoji: emoji-width
emoji-width:
	@go run ../../cmd/visual-tests/main.go emoji-width

gear: gear-spacing
gear-spacing:
	@go run ../../cmd/visual-tests/main.go gear-spacing

# Clean backup files
clean:
	@echo "Cleaning backup files..."
	@rm -f *.bak
	@echo "Done!"