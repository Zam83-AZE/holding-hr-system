#!/bin/bash

# Test nəticələrini GitHub-a push etmək üçün skript

REPORT_DIR="test-results"
BRANCH_NAME="test-results-$(date +%Y%m%d-%H%M%S)"

echo "=========================================="
echo "Pushing Test Results to GitHub"
echo "=========================================="

# Git config
git config --global user.email "test@holding-hr.local"
git config --global user.name "Test Runner"

# Check if there are test results
if [ ! -d "$REPORT_DIR" ] || [ -z "$(ls -A $REPORT_DIR 2>/dev/null)" ]; then
    echo "No test results found in $REPORT_DIR"
    exit 1
fi

# Add test results
echo "Adding test results..."
git add $REPORT_DIR/

# Check if there are changes to commit
if git diff --staged --quiet; then
    echo "No changes to commit"
    exit 0
fi

# Commit
echo "Committing test results..."
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
git commit -m "Test results - $TIMESTAMP

Auto-generated test report from integration tests."

# Push
echo "Pushing to origin main..."
git push origin main

echo "=========================================="
echo "Test results pushed successfully!"
echo "=========================================="
