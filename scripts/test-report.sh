#!/bin/bash

# GoKanon Test Report Generator
# Generates a beautiful, easy-to-understand test report with coverage metrics

set -e

# Colors for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Output file (optional)
OUTPUT_FILE="${1:-}"
TEMP_OUTPUT=$(mktemp)

# Function to output to both console and file
output() {
    echo -e "$1" | tee -a "$TEMP_OUTPUT"
}

# Function to create a progress bar
progress_bar() {
    local percentage=$1
    local width=40
    local filled=$((percentage * width / 100))
    local empty=$((width - filled))

    printf "["
    printf "${GREEN}%${filled}s${NC}" | tr ' ' '█'
    printf "%${empty}s" | tr ' ' '░'
    printf "] ${BOLD}%3d%%${NC}" "$percentage"
}

# Header
output ""
output "${BOLD}${CYAN}╔══════════════════════════════════════════════════════════════════╗${NC}"
output "${BOLD}${CYAN}║                   GoKanon Test Report                            ║${NC}"
output "${BOLD}${CYAN}╔══════════════════════════════════════════════════════════════════╗${NC}"
output ""
output "${BOLD}Generated:${NC} $(date '+%Y-%m-%d %H:%M:%S')"
output ""

# Run tests and capture output
echo "Running tests..."
TEST_OUTPUT=$(mktemp)
COVERAGE_FILE=$(mktemp)

if go test -v -coverprofile="$COVERAGE_FILE" ./... 2>&1 | tee "$TEST_OUTPUT"; then
    TEST_STATUS="${GREEN}✓ PASS${NC}"
    TEST_RESULT=0
else
    TEST_STATUS="${RED}✗ FAIL${NC}"
    TEST_RESULT=1
fi

# Parse test results
TOTAL_TESTS=$(grep -E "^(PASS|FAIL)" "$TEST_OUTPUT" | wc -l)
PASSED_TESTS=$(grep -E "^PASS" "$TEST_OUTPUT" | wc -l)
FAILED_TESTS=$(grep -E "^FAIL" "$TEST_OUTPUT" | wc -l)

# Calculate overall coverage
OVERALL_COVERAGE=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "total:" | awk '{print $3}' | sed 's/%//' || echo "0")
OVERALL_COVERAGE_INT=${OVERALL_COVERAGE%.*}

# Test Summary Section
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output "${BOLD}  TEST SUMMARY${NC}"
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output ""
output "  Status:        $TEST_STATUS"
output "  Total Tests:   ${BOLD}$TOTAL_TESTS${NC}"
output "  ${GREEN}Passed:${NC}        ${BOLD}$PASSED_TESTS${NC}"
if [ "$FAILED_TESTS" -gt 0 ]; then
    output "  ${RED}Failed:${NC}        ${BOLD}$FAILED_TESTS${NC}"
fi
output ""

# Coverage Summary Section
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output "${BOLD}  COVERAGE SUMMARY${NC}"
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output ""
output "  Overall Coverage: $(progress_bar "$OVERALL_COVERAGE_INT")"
output ""

# Coverage grading
if [ "$OVERALL_COVERAGE_INT" -ge 80 ]; then
    GRADE="${GREEN}A - Excellent${NC}"
elif [ "$OVERALL_COVERAGE_INT" -ge 70 ]; then
    GRADE="${CYAN}B - Good${NC}"
elif [ "$OVERALL_COVERAGE_INT" -ge 60 ]; then
    GRADE="${YELLOW}C - Fair${NC}"
elif [ "$OVERALL_COVERAGE_INT" -ge 50 ]; then
    GRADE="${YELLOW}D - Needs Improvement${NC}"
else
    GRADE="${RED}F - Poor${NC}"
fi

output "  Coverage Grade:   $GRADE"
output ""

# Package Coverage Details
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output "${BOLD}  PACKAGE COVERAGE${NC}"
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output ""

# Get coverage by package
go test -coverprofile="$COVERAGE_FILE" ./... >/dev/null 2>&1 || true

# Parse and display package coverage
while IFS= read -r line; do
    if [[ $line == ok* ]] || [[ $line == "?"* ]]; then
        PACKAGE=$(echo "$line" | awk '{print $2}')
        COVERAGE=$(echo "$line" | grep -oP 'coverage: \K[0-9.]+' || echo "0")

        # Skip if no coverage info
        if [ "$COVERAGE" == "0" ] && [[ $line == "?"* ]]; then
            continue
        fi

        # Shorten package name for display
        SHORT_PKG=$(echo "$PACKAGE" | sed 's/github.com\/alenon\/gokanon\///' | cut -c1-40)

        # Get integer part of coverage
        COVERAGE_INT=${COVERAGE%.*}

        # Color code based on coverage
        if [ "$COVERAGE_INT" -ge 80 ]; then
            COLOR=$GREEN
            ICON="✓"
        elif [ "$COVERAGE_INT" -ge 60 ]; then
            COLOR=$YELLOW
            ICON="◐"
        elif [ "$COVERAGE_INT" -gt 0 ]; then
            COLOR=$RED
            ICON="✗"
        else
            COLOR=$CYAN
            ICON="○"
            COVERAGE="no tests"
        fi

        printf "  ${COLOR}%-5s${NC} %-42s ${BOLD}%10s${NC}\n" "$ICON" "$SHORT_PKG" "$COVERAGE%" | tee -a "$TEMP_OUTPUT"
    fi
done < <(go test -cover ./... 2>&1)

output ""

# Command Coverage Details (if exists)
if [ -f "$COVERAGE_FILE" ]; then
    COMMANDS_COVERAGE=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "internal/cli/commands/" | grep -E "\.go:" || true)

    if [ -n "$COMMANDS_COVERAGE" ]; then
        output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        output "${BOLD}  COMMAND COVERAGE (Key Functions)${NC}"
        output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        output ""

        # Parse command functions and their coverage
        echo "$COMMANDS_COVERAGE" | grep -E "(Run|List|Compare|Export|Stats|Trend|Check|Flamegraph|Serve|Delete|Baseline|Doctor|Interactive|Completion)\s" | while read -r line; do
            FUNC=$(echo "$line" | awk '{print $2}' | sed 's/.*://')
            COV=$(echo "$line" | awk '{print $3}' | sed 's/%//')
            COV_INT=${COV%.*}

            # Color code
            if [ "$COV_INT" -ge 80 ]; then
                COLOR=$GREEN
                ICON="✓"
            elif [ "$COV_INT" -ge 60 ]; then
                COLOR=$YELLOW
                ICON="◐"
            else
                COLOR=$RED
                ICON="✗"
            fi

            printf "  ${COLOR}%-5s${NC} %-40s ${BOLD}%6.1f%%${NC}\n" "$ICON" "$FUNC" "$COV" | tee -a "$TEMP_OUTPUT"
        done | head -20

        output ""
    fi
fi

# Recommendations
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output "${BOLD}  RECOMMENDATIONS${NC}"
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output ""

if [ "$TEST_RESULT" -ne 0 ]; then
    output "  ${RED}✗${NC} Fix failing tests before proceeding"
fi

if [ "$OVERALL_COVERAGE_INT" -lt 70 ]; then
    output "  ${YELLOW}◐${NC} Increase test coverage to at least 70%"
fi

if [ "$OVERALL_COVERAGE_INT" -ge 80 ]; then
    output "  ${GREEN}✓${NC} Excellent test coverage! Keep it up!"
fi

if [ "$FAILED_TESTS" -eq 0 ]; then
    output "  ${GREEN}✓${NC} All tests passing!"
fi

# Check for packages with low coverage
LOW_COVERAGE=$(go test -cover ./... 2>&1 | grep -E "coverage: [0-9]" | awk '{if ($4 != "" && $4+0 < 60 && $4+0 > 0) print $2}' | head -3)
if [ -n "$LOW_COVERAGE" ]; then
    output ""
    output "  ${YELLOW}Packages needing attention:${NC}"
    echo "$LOW_COVERAGE" | while read -r pkg; do
        SHORT=$(echo "$pkg" | sed 's/github.com\/alenon\/gokanon\///')
        output "    • $SHORT"
    done
fi

output ""

# Footer
output "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
output ""

if [ "$TEST_RESULT" -eq 0 ] && [ "$OVERALL_COVERAGE_INT" -ge 60 ]; then
    output "${BOLD}${GREEN}✓ Tests passed with good coverage!${NC}"
    output ""
    FINAL_EXIT=0
else
    if [ "$TEST_RESULT" -ne 0 ]; then
        output "${BOLD}${RED}✗ Some tests failed!${NC}"
    fi
    if [ "$OVERALL_COVERAGE_INT" -lt 60 ]; then
        output "${BOLD}${YELLOW}⚠ Test coverage below recommended threshold (60%)${NC}"
    fi
    output ""
    FINAL_EXIT=1
fi

# Save to file if requested
if [ -n "$OUTPUT_FILE" ]; then
    cat "$TEMP_OUTPUT" > "$OUTPUT_FILE"
    echo -e "${CYAN}Report saved to: $OUTPUT_FILE${NC}"
fi

# Cleanup
rm -f "$TEST_OUTPUT" "$COVERAGE_FILE" "$TEMP_OUTPUT"

exit $FINAL_EXIT
