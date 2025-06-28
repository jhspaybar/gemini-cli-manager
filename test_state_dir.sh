#!/bin/bash
# Test script for state directory functionality

set -e

echo "Testing Gemini CLI Manager with custom state directories"
echo "======================================================="
echo ""

# Test directory 1
TEST_DIR_1="/tmp/gemini-test-1"
echo "1. Creating first test setup in: $TEST_DIR_1"
mkdir -p "$TEST_DIR_1"
echo "   Running: ./gemini-cli-manager --state-dir $TEST_DIR_1"
echo "   (Press 'q' to quit when the app opens)"
echo ""
read -p "Press Enter to continue..."
./gemini-cli-manager --state-dir "$TEST_DIR_1"

echo ""
echo "2. Checking created directories:"
echo "   Contents of $TEST_DIR_1:"
ls -la "$TEST_DIR_1"
echo ""

# Test directory 2
TEST_DIR_2="/tmp/gemini-test-2"
echo "3. Creating second independent setup in: $TEST_DIR_2"
mkdir -p "$TEST_DIR_2"
echo "   Running: ./gemini-cli-manager --state-dir $TEST_DIR_2"
echo "   (Press 'q' to quit when the app opens)"
echo ""
read -p "Press Enter to continue..."
./gemini-cli-manager --state-dir "$TEST_DIR_2"

echo ""
echo "4. Checking both directories are independent:"
echo "   Contents of $TEST_DIR_1:"
ls -la "$TEST_DIR_1"
echo ""
echo "   Contents of $TEST_DIR_2:"
ls -la "$TEST_DIR_2"
echo ""

echo "5. Testing with home directory shortcut:"
echo "   Running: ./gemini-cli-manager --state-dir ~/gemini-test"
echo "   (Press 'q' to quit when the app opens)"
echo ""
read -p "Press Enter to continue..."
./gemini-cli-manager --state-dir ~/gemini-test

echo ""
echo "   Created in home directory:"
ls -la ~/gemini-test
echo ""

echo "Test complete! Each state directory maintains independent:"
echo "- Extensions"
echo "- Profiles"
echo "- Settings"
echo ""
echo "Clean up test directories:"
echo "  rm -rf $TEST_DIR_1 $TEST_DIR_2 ~/gemini-test"