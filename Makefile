.PHONY: build test check clean run help check-warnings

# Default target
all: check build

# Build the project
build:
	@echo "Building project..."
	@cargo build

# Build release version
release:
	@echo "Building release version..."
	@cargo build --release

# Run tests
test:
	@echo "Running tests..."
	@cargo test

# Check code (includes format, lint, and warnings)
check: check-warnings
	@echo "Checking code format..."
	@cargo fmt -- --check
	@echo "Running clippy..."
	@cargo clippy -- -D warnings

# Check for warnings
check-warnings:
	@echo "Checking for warnings..."
	@./check-warnings.sh

# Format code
fmt:
	@echo "Formatting code..."
	@cargo fmt

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@cargo clean

# Run the application
run:
	@cargo run

# Run with debug logging
run-debug:
	@RUST_LOG=debug cargo run

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the project"
	@echo "  make release       - Build release version"
	@echo "  make test          - Run tests"
	@echo "  make check         - Check code (format, clippy, warnings)"
	@echo "  make check-warnings - Check for compilation warnings"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make run           - Run the application"
	@echo "  make run-debug     - Run with debug logging"
	@echo "  make help          - Show this help message"