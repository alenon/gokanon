#!/usr/bin/env bash
# Test script for gokanon zsh completion

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

# Check if zsh is installed
if ! command -v zsh &> /dev/null; then
    log_error "zsh is not installed. Please install zsh to run this test."
    exit 1
fi

log_info "Testing gokanon zsh completion..."

# Build gokanon if needed
cd "$PROJECT_ROOT"
if [ ! -f "./bin/gokanon" ]; then
    log_info "Building gokanon..."
    make build
fi

# Add gokanon to PATH for testing
export PATH="$PROJECT_ROOT/bin:$PATH"

# Save completion output to file to avoid subshell issues
COMPLETION_FILE="/tmp/gokanon_completion_test.zsh"
gokanon completion zsh > "$COMPLETION_FILE" 2>&1

# Test 1: Check if completion command produces output
log_info "Test 1: Check if 'gokanon completion zsh' produces output"
if [ -s "$COMPLETION_FILE" ]; then
    log_success "gokanon completion zsh produces output"
else
    log_error "gokanon completion zsh produced no output"
fi

# Test 2: Check if completion script contains the function definition
log_info "Test 2: Check if completion script contains _gokanon function"
if grep -q "_gokanon()" "$COMPLETION_FILE"; then
    log_success "Completion script contains _gokanon function"
else
    log_error "Completion script does not contain _gokanon function"
fi

# Test 3: Check if completion script registers the function with compdef
log_info "Test 3: Check if completion script registers _gokanon with compdef"
if grep -q 'compdef _gokanon gokanon' "$COMPLETION_FILE"; then
    log_success "Completion script registers function with compdef"
else
    log_error "Completion script does not register function with compdef"
fi

# Test 4: Check if completion script contains compdef directive
log_info "Test 4: Check if completion script has #compdef directive"
if grep -q "#compdef gokanon" "$COMPLETION_FILE"; then
    log_success "Completion script has #compdef directive"
else
    log_error "Completion script missing #compdef directive"
fi

# Test 5: Test that function loads in zsh
log_info "Test 5: Test that completion function loads in zsh"
ZSH_TEST_SCRIPT='autoload -U compinit; compinit -u 2>/dev/null; source <(gokanon completion zsh) 2>/dev/null; if typeset -f _gokanon > /dev/null 2>&1; then echo "FUNCTION_LOADED"; else echo "FUNCTION_NOT_LOADED"; exit 1; fi'

if ZSH_RESULT=$(timeout 5 zsh -c "$ZSH_TEST_SCRIPT" 2>&1) && echo "$ZSH_RESULT" | grep -q "FUNCTION_LOADED"; then
    log_success "Zsh completion function loads successfully"
else
    log_error "Zsh completion function failed to load: $ZSH_RESULT"
fi

# Test 6: Test that main commands are included in completions
log_info "Test 6: Check if completion script includes main commands"
EXPECTED_COMMANDS=("run" "list" "compare" "export" "stats" "baseline" "completion")
ALL_FOUND=true
for cmd in "${EXPECTED_COMMANDS[@]}"; do
    if ! grep -q "'$cmd:" "$COMPLETION_FILE"; then
        log_error "Completion script missing command: $cmd"
        ALL_FOUND=false
    fi
done
if [ "$ALL_FOUND" = true ]; then
    log_success "All main commands present in completion script"
fi

# Test 7: Test that completion script has required command handlers
log_info "Test 7: Verify completion script contains command definitions"
REQUIRED_CMDS=("run" "compare" "baseline" "completion")
MISSING_CMDS=0
for cmd in "${REQUIRED_CMDS[@]}"; do
    # Look for case handlers for these commands
    if ! grep -q "^[[:space:]]*$cmd)" "$COMPLETION_FILE"; then
        log_error "Completion script missing handler for command: $cmd"
        ((MISSING_CMDS++))
    fi
done
if [ $MISSING_CMDS -eq 0 ]; then
    log_success "All required command handlers present"
fi

# Test 8: Verify no syntax errors in zsh
log_info "Test 8: Check for zsh syntax errors"
if zsh -n "$COMPLETION_FILE" 2>&1; then
    log_success "No zsh syntax errors detected"
else
    log_error "Zsh syntax errors found in completion script"
fi

# Test 9: Verify sourcing doesn't produce _arguments errors
log_info "Test 9: Verify sourcing doesn't produce _arguments errors"
ZSH_SOURCE_ERRORS=$(zsh -c 'autoload -U compinit; compinit -u 2>/dev/null; source <(gokanon completion zsh)' 2>&1)
if echo "$ZSH_SOURCE_ERRORS" | grep -q "_arguments:comparguments"; then
    log_error "Sourcing produces _arguments error: $ZSH_SOURCE_ERRORS"
else
    log_success "Sourcing completes without _arguments errors"
fi

# Cleanup
rm -f "$COMPLETION_FILE"

# Summary
echo ""
echo "================================"
echo "Test Summary"
echo "================================"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    echo ""
    echo "Please fix the failing tests before committing."
    exit 1
else
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    echo ""
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi
