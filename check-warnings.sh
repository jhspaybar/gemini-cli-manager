#!/bin/bash
# Check for warnings in the build

echo "Checking for warnings..."

# Run cargo check with warnings as errors
if RUSTFLAGS="-D warnings" cargo check --all-targets 2>&1; then
    echo "✅ No warnings found!"
    exit 0
else
    echo "❌ Warnings found! Please fix them before proceeding."
    echo ""
    echo "Run 'cargo build' to see the warnings."
    exit 1
fi