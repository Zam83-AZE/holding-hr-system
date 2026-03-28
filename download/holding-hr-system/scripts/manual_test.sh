#!/bin/bash

# ============================================================
# HOLDING HR SYSTEM - MANUAL TEST RUNNER
# ============================================================
# Usage: ./scripts/manual_test.sh [BASE_URL]
# Example: ./scripts/manual_test.sh http://localhost:8080
# ============================================================

BASE_URL="${1:-http://localhost:8080}"
REPORT_DIR="./test-results"

echo "============================================================"
echo "HOLDING HR SYSTEM - MANUAL TEST RUNNER"
echo "============================================================"
echo "Base URL: $BASE_URL"
echo "Report Directory: $REPORT_DIR"
echo "============================================================"

# Create report directory
mkdir -p "$REPORT_DIR"

# Export environment variables
export BASE_URL
export REPORT_DIR

# Run tests
echo ""
echo "Running integration tests..."
go run tests/integration_test.go

echo ""
echo "Test completed. Check $REPORT_DIR for results."
