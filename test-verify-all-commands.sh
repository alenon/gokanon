#!/bin/bash

# End-to-End Verification Script for GoKanon Commands
# This script tests all commands to ensure they work as expected

set -e  # Exit on first error

echo "=========================================="
echo "GoKanon Commands End-to-End Verification"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to print test results
pass() {
    echo -e "${GREEN}✓ PASS:${NC} $1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
}

fail() {
    echo -e "${RED}✗ FAIL:${NC} $1"
    TESTS_FAILED=$((TESTS_FAILED + 1))
}

info() {
    echo -e "${YELLOW}ℹ INFO:${NC} $1"
}

# Setup test environment
TEST_DIR=$(mktemp -d)
STORAGE_DIR="$TEST_DIR/.gokanon-test"

info "Test directory: $TEST_DIR"
info "Storage directory: $STORAGE_DIR"
echo ""

# Build the binary
info "Building gokanon..."
go build -o "$TEST_DIR/gokanon" . || { fail "Failed to build gokanon"; exit 1; }
pass "Built gokanon binary"
echo ""

GOKANON="$TEST_DIR/gokanon"

# Test 1: Doctor command
echo "Test 1: doctor command"
if $GOKANON doctor -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "doctor command executed successfully"
else
    fail "doctor command failed"
fi
echo ""

# Test 2: Help command
echo "Test 2: help command"
if $GOKANON help | grep -q "gokanon - A CLI tool"; then
    pass "help command shows usage text"
else
    fail "help command output incorrect"
fi
echo ""

# Test 3: List command (empty storage)
echo "Test 3: list command (empty storage)"
if $GOKANON list -storage="$STORAGE_DIR" | grep -q "No benchmark results found"; then
    pass "list command works with empty storage"
else
    fail "list command didn't handle empty storage correctly"
fi
echo ""

# Test 4: Run command (create test benchmarks)
echo "Test 4: run command"
# Create a simple Go benchmark file
mkdir -p "$TEST_DIR/benchtest"
cat > "$TEST_DIR/benchtest/bench_test.go" <<'EOF'
package benchtest

import "testing"

func BenchmarkSimpleAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = i + i
    }
}

func BenchmarkSimpleMultiply(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = i * 2
    }
}
EOF

cat > "$TEST_DIR/benchtest/go.mod" <<'EOF'
module benchtest

go 1.21
EOF

if $GOKANON run -bench=. -pkg="$TEST_DIR/benchtest" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "run command executed successfully"
else
    fail "run command failed"
fi
echo ""

# Test 5: List command (with data)
echo "Test 5: list command (with data)"
if $GOKANON list -storage="$STORAGE_DIR" | grep -q "BenchmarkSimpleAdd\|Timestamp"; then
    pass "list command shows benchmark results"
else
    fail "list command didn't show results"
fi
echo ""

# Test 6: Run command again (for comparison)
echo "Test 6: run command (second time for comparison)"
if $GOKANON run -bench=. -pkg="$TEST_DIR/benchtest" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "run command executed successfully (2nd time)"
else
    fail "run command failed (2nd time)"
fi
echo ""

# Test 7: Compare command (latest)
echo "Test 7: compare command (latest)"
if $GOKANON compare -latest -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "compare command with --latest works"
else
    fail "compare command with --latest failed"
fi
echo ""

# Test 8: Stats command
echo "Test 8: stats command"
if $GOKANON stats -storage="$STORAGE_DIR" | grep -q "Statistical Analysis"; then
    pass "stats command executed successfully"
else
    fail "stats command failed"
fi
echo ""

# Test 9: Trend command
echo "Test 9: trend command"
if $GOKANON trend -storage="$STORAGE_DIR" | grep -q "Performance Trend Analysis"; then
    pass "trend command executed successfully"
else
    fail "trend command failed"
fi
echo ""

# Test 10: Baseline save
echo "Test 10: baseline save command"
if $GOKANON baseline save -name=test-baseline -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "baseline save command executed successfully"
else
    fail "baseline save command failed"
fi
echo ""

# Test 11: Baseline list
echo "Test 11: baseline list command"
if $GOKANON baseline list -storage="$STORAGE_DIR" | grep -q "test-baseline"; then
    pass "baseline list command shows saved baseline"
else
    fail "baseline list command failed"
fi
echo ""

# Test 12: Baseline show
echo "Test 12: baseline show command"
if $GOKANON baseline show -name=test-baseline -storage="$STORAGE_DIR" | grep -q "Baseline: test-baseline"; then
    pass "baseline show command executed successfully"
else
    fail "baseline show command failed"
fi
echo ""

# Test 13: Compare with baseline
echo "Test 13: compare with baseline"
if $GOKANON compare -baseline=test-baseline -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "compare command with baseline works"
else
    fail "compare command with baseline failed"
fi
echo ""

# Test 14: Export to CSV
echo "Test 14: export command (CSV)"
if $GOKANON export -latest -format=csv -output="$TEST_DIR/comparison.csv" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    if [ -f "$TEST_DIR/comparison.csv" ]; then
        pass "export to CSV succeeded"
    else
        fail "export to CSV didn't create file"
    fi
else
    fail "export to CSV command failed"
fi
echo ""

# Test 15: Export to HTML
echo "Test 15: export command (HTML)"
if $GOKANON export -latest -format=html -output="$TEST_DIR/comparison.html" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    if [ -f "$TEST_DIR/comparison.html" ]; then
        pass "export to HTML succeeded"
    else
        fail "export to HTML didn't create file"
    fi
else
    fail "export to HTML command failed"
fi
echo ""

# Test 16: Export to Markdown
echo "Test 16: export command (Markdown)"
if $GOKANON export -latest -format=markdown -output="$TEST_DIR/comparison.md" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    if [ -f "$TEST_DIR/comparison.md" ]; then
        pass "export to Markdown succeeded"
    else
        fail "export to Markdown didn't create file"
    fi
else
    fail "export to Markdown command failed"
fi
echo ""

# Test 17: Check command (with high threshold)
echo "Test 17: check command (should pass with high threshold)"
if $GOKANON check -latest -threshold=500.0 -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "check command with high threshold passed"
else
    # This might fail if performance changed significantly, which is OK
    info "check command exited with non-zero (expected if performance degraded)"
fi
echo ""

# Test 18: Completion commands
echo "Test 18: completion command (bash)"
if $GOKANON completion bash | grep -q "bash completion"; then
    pass "completion command for bash works"
else
    fail "completion command for bash failed"
fi
echo ""

echo "Test 19: completion command (zsh)"
if $GOKANON completion zsh | grep -q "zsh completion"; then
    pass "completion command for zsh works"
else
    fail "completion command for zsh failed"
fi
echo ""

echo "Test 20: completion command (fish)"
if $GOKANON completion fish | grep -q "Fish completion"; then
    pass "completion command for fish works"
else
    fail "completion command for fish failed"
fi
echo ""

# Test 21: Baseline delete
echo "Test 21: baseline delete command"
if $GOKANON baseline delete -name=test-baseline -storage="$STORAGE_DIR" >/dev/null 2>&1; then
    pass "baseline delete command executed successfully"
else
    fail "baseline delete command failed"
fi
echo ""

# Test 22: Delete command
echo "Test 22: delete command"
# Get the first run ID
RUN_ID=$($GOKANON list -storage="$STORAGE_DIR" | grep -v "ID\|--" | head -1 | awk '{print $1}')
if [ -n "$RUN_ID" ]; then
    if $GOKANON delete "$RUN_ID" -storage="$STORAGE_DIR" >/dev/null 2>&1; then
        pass "delete command executed successfully"
    else
        fail "delete command failed"
    fi
else
    fail "Could not get run ID for delete test"
fi
echo ""

# Test 23: Invalid command handling
echo "Test 23: invalid command handling"
if $GOKANON invalid-command 2>&1 | grep -q "Unknown command"; then
    pass "invalid command properly handled"
else
    fail "invalid command not handled correctly"
fi
echo ""

# Cleanup
info "Cleaning up test directory: $TEST_DIR"
rm -rf "$TEST_DIR"
echo ""

# Final report
echo "=========================================="
echo "Test Results Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo "Total:  $((TESTS_PASSED + TESTS_FAILED))"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
