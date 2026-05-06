#!/bin/bash

# This script runs all fuzz tests in the current package sequentially.
# It accepts an optional duration parameter (e.g., ./fuzz.sh 2m).
# Default duration is 1 minute per fuzz target.

set -e

# Default fuzz time per target if not specified
FUZZ_TIME=${1:-1m}

echo "Starting fuzzing session with duration: $FUZZ_TIME per target"

# Find all fuzz tests (functions starting with 'Fuzz')
FUZZ_TESTS=$(go test -list . | grep ^Fuzz)

if [ -z "$FUZZ_TESTS" ]; then
    echo "No fuzz tests found in the current directory."
    exit 0
fi

for test in $FUZZ_TESTS; do
    echo "------------------------------------------------------------"
    echo "Fuzzing $test..."
    # -run=^$ ensures regular tests are skipped
    # -fuzztime sets the duration for this specific target
    go test -fuzz="$test" -fuzztime="$FUZZ_TIME" -run=^$ .
done

echo "------------------------------------------------------------"
echo "All fuzz tests completed successfully!"
