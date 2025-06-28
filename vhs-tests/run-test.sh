#!/bin/bash
# Simple test runner for VHS tests

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "🎬 VHS Test Runner"
echo "=================="
echo ""

# Check if VHS is installed
if ! command -v vhs &> /dev/null; then
    echo -e "${RED}❌ VHS is not installed!${NC}"
    echo "Install with: brew install vhs"
    exit 1
fi

# Run setup
echo "Setting up environment..."
./helpers/setup-test-env.sh
echo ""

# Handle arguments
if [ $# -eq 0 ]; then
    echo "Usage: ./run-test.sh <test-name|all>"
    echo ""
    echo "Examples:"
    echo "  ./run-test.sh navigation/sidebar"
    echo "  ./run-test.sh state-dir/multiple-setups"
    echo "  ./run-test.sh all"
    echo ""
    echo "Available tests:"
    find tests -name "*.tape" -type f | sed 's|tests/||' | sed 's|.tape||' | sed 's|^|  |'
    exit 0
fi

TEST_NAME=$1

# Function to run a single test
run_test() {
    local tape_file=$1
    local test_name=$(echo $tape_file | sed 's|tests/||' | sed 's|.tape||')
    local test_md="${tape_file%.tape}.test.md"
    
    echo -e "${YELLOW}Running test: $test_name${NC}"
    echo "Tape file: $tape_file"
    
    if [ ! -f "$tape_file" ]; then
        echo -e "${RED}❌ Tape file not found: $tape_file${NC}"
        return 1
    fi
    
    # Run VHS
    if vhs "$tape_file"; then
        echo -e "${GREEN}✅ Test recording completed${NC}"
        
        # Check if test spec exists
        if [ -f "$test_md" ]; then
            echo "📋 Test specification: $test_md"
            echo "   Review the GIF against the specification"
        else
            echo -e "${YELLOW}⚠️  No test specification found at: $test_md${NC}"
        fi
        
        # List generated files
        echo "📁 Generated files:"
        find output -name "*${test_name//\//-}*" -type f 2>/dev/null | sed 's|^|   |' || echo "   None found"
        
        return 0
    else
        echo -e "${RED}❌ Test failed${NC}"
        return 1
    fi
}

# Run tests
if [ "$TEST_NAME" = "all" ]; then
    echo "Running all tests..."
    echo ""
    
    failed=0
    total=0
    
    for tape in tests/*/*.tape; do
        ((total++))
        if ! run_test "$tape"; then
            ((failed++))
        fi
        echo ""
    done
    
    echo "Summary: $((total-failed))/$total tests passed"
    
    if [ $failed -gt 0 ]; then
        echo -e "${RED}❌ Some tests failed${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ All tests passed${NC}"
    fi
else
    # Run single test
    tape_file="tests/${TEST_NAME}.tape"
    if run_test "$tape_file"; then
        echo ""
        echo -e "${GREEN}✅ Test completed successfully${NC}"
        echo "View the GIF at: output/$(basename ${TEST_NAME//\//-}).gif"
    else
        exit 1
    fi
fi