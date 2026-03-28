#!/bin/bash

# ============================================================
# HOLDING HR SYSTEM - FULL TEST RUNNER
# ============================================================

set -e

echo "============================================================"
echo "HOLDING HR SYSTEM - AUTOMATED TEST RUNNER"
echo "============================================================"

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
REPORT_DIR="${REPORT_DIR:-/app/test-results}"
TEST_TIMEOUT="${TEST_TIMEOUT:-300}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "Base URL: $BASE_URL"
echo "Report Directory: $REPORT_DIR"
echo "============================================================"

# Create report directory
mkdir -p "$REPORT_DIR"

# Wait for server to be ready
echo ""
echo "Waiting for server to start..."
max_attempts=30
attempt=0
while ! curl -s "$BASE_URL/login" > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo -e "${RED}Server did not start within $max_attempts attempts${NC}"
        exit 1
    fi
    echo "  Attempt $attempt/$max_attempts - Server not ready yet..."
    sleep 2
done
echo -e "${GREEN}Server is ready!${NC}"

# Run Go integration tests
echo ""
echo "============================================================"
echo "RUNNING GO INTEGRATION TESTS"
echo "============================================================"

cd /app

# Build and run test
export BASE_URL="$BASE_URL"
export REPORT_DIR="$REPORT_DIR"

# Check if test file exists
if [ -f "tests/integration_test.go" ]; then
    echo "Running integration tests..."
    go run tests/integration_test.go 2>&1 | tee "$REPORT_DIR/test_output.log"
    TEST_EXIT_CODE=$?
else
    echo -e "${RED}Test file not found: tests/integration_test.go${NC}"
    exit 1
fi

echo ""
echo "============================================================"
echo "TEST RESULTS"
echo "============================================================"

# Display results
if [ -f "$REPORT_DIR/test_report.txt" ]; then
    cat "$REPORT_DIR/test_report.txt"
fi

if [ -f "$REPORT_DIR/test_report.json" ]; then
    echo ""
    echo "JSON report available at: $REPORT_DIR/test_report.json"
fi

# Check for failures
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo ""
    echo -e "${GREEN}============================================================${NC}"
    echo -e "${GREEN}ALL TESTS PASSED! ✓${NC}"
    echo -e "${GREEN}============================================================${NC}"
else
    echo ""
    echo -e "${RED}============================================================${NC}"
    echo -e "${RED}SOME TESTS FAILED! ✗${NC}"
    echo -e "${RED}============================================================${NC}"
fi

exit $TEST_EXIT_CODE
