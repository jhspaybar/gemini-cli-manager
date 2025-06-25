#!/bin/bash
# scripts/verify.sh - Complete build verification

set -e

echo "🔨 Running build verification..."

# Format check
echo "📝 Checking formatting..."
if ! /usr/local/go/bin/go fmt ./... | grep -q .; then
    echo "✅ Format check passed"
else
    echo "❌ Format issues found. Run: go fmt ./..."
    exit 1
fi

# Vet check
echo "🔍 Running go vet..."
/usr/local/go/bin/go vet ./...
echo "✅ Vet check passed"

# Build all packages
echo "🏗️  Building all packages..."
/usr/local/go/bin/go build ./...
echo "✅ Build successful"

# Run tests
echo "🧪 Running tests..."
/usr/local/go/bin/go test ./... -short
echo "✅ Tests passed"

# Build binary
echo "📦 Building binary..."
/usr/local/go/bin/go build -o gemini-cli-manager
echo "✅ Binary built successfully"

echo "✨ All verification checks passed!"