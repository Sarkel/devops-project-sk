#!/bin/bash
set -e

COVERAGE_FILE="${BASE_DIR}/coverage-summary.json"

echo "=== Web Coverage Validation ==="
echo "Base Directory: $BASE_DIR"
echo "Report file: $COVERAGE_FILE"
echo "Required threshold: $THRESHOLD%"

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "ERROR: coverage-summary.json file not found!"
    echo "Ensure tests generated 'json-summary' report and artifact was downloaded correctly."
    echo "Debug: Listing contents of '$BASE_DIR':"

    ls -R "$BASE_DIR" || echo "Directory $BASE_DIR does not exist."
    exit 1
fi

ACTUAL_COVERAGE=$(jq -r '.total.statements.pct // empty' "$COVERAGE_FILE")

if [ -z "$ACTUAL_COVERAGE" ]; then
     echo "ERROR: Failed to parse coverage percentage from JSON."
     exit 1
fi

echo "Current coverage (Statements): $ACTUAL_COVERAGE%"

if [ $(echo "$ACTUAL_COVERAGE < $THRESHOLD" | bc) -eq 1 ]; then
    echo "FAILURE: Code coverage is below threshold ($ACTUAL_COVERAGE% < $THRESHOLD%)"
    exit 1
else
    echo "SUCCESS: Code coverage is sufficient."
    exit 0
fi