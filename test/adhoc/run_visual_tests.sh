#!/bin/bash
# Run visual tests for UI components

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the root directory
ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
VISUAL_TEST_CMD="$ROOT_DIR/cmd/visual-tests/main.go"

# Check if visual test command exists
if [ ! -f "$VISUAL_TEST_CMD" ]; then
    echo "Error: Visual test command not found at $VISUAL_TEST_CMD"
    exit 1
fi

# Function to show usage
show_usage() {
    echo "Usage: ./run_visual_tests.sh [command]"
    echo ""
    echo "Commands:"
    echo "  list       List all available tests"
    echo "  all        Run all tests"
    echo "  <name>     Run specific test (e.g., tabs, emoji-width)"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./run_visual_tests.sh list"
    echo "  ./run_visual_tests.sh tabs"
    echo "  ./run_visual_tests.sh all"
}

# Parse arguments
if [ $# -eq 0 ]; then
    show_usage
    exit 0
fi

case "$1" in
    help|-h|--help)
        show_usage
        ;;
    list)
        echo -e "${YELLOW}Available visual tests:${NC}"
        go run "$VISUAL_TEST_CMD" --list
        ;;
    all)
        echo -e "${GREEN}Running all visual tests...${NC}"
        go run "$VISUAL_TEST_CMD" all
        ;;
    *)
        echo -e "${GREEN}Running visual test: $1${NC}"
        go run "$VISUAL_TEST_CMD" "$1"
        ;;
esac