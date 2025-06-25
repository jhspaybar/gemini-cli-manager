.PHONY: build test verify fmt clean run

# Default target
all: verify

# Quick build check
build:
	@echo "🏗️  Building..."
	@/usr/local/go/bin/go build ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	@/usr/local/go/bin/go test ./... -v

# Run short tests (faster)
test-short:
	@echo "🧪 Running short tests..."
	@/usr/local/go/bin/go test ./... -short

# Format code
fmt:
	@echo "📝 Formatting code..."
	@/usr/local/go/bin/go fmt ./...

# Full verification
verify: fmt
	@echo "🔨 Running full verification..."
	@./scripts/verify.sh

# Build and run
run: build
	@./gemini-cli-manager

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f gemini-cli-manager
	@/usr/local/go/bin/go clean ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@/usr/local/go/bin/go mod download
	@/usr/local/go/bin/go mod tidy