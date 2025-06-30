#!/bin/bash
# Test runner script for Gemini CLI Manager

set -e

echo "🧪 Running Gemini CLI Manager Tests"
echo "===================================="

# Run all tests with nice output
echo ""
echo "📋 Running unit tests..."
cargo test --test unit -- --nocapture

echo ""
echo "🔗 Running integration tests..."
cargo test --test integration -- --nocapture

echo ""
echo "📸 Running snapshot tests..."
cargo insta test

echo ""
echo "✅ All tests completed!"
echo ""

# Optional: generate coverage report
if command -v cargo-tarpaulin &> /dev/null; then
    echo "📊 Generating coverage report..."
    cargo tarpaulin --out Html --output-dir target/coverage
    echo "Coverage report generated at: target/coverage/tarpaulin-report.html"
fi