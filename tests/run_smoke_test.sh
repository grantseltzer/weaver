#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

TEST_FUNCTIONS_FILE="./tests/test_data/functions_file.txt"
REFERENCE_FILE="./tests/test_data/tester_output_sorted_slurped.json"
OUTPUT_FILE=$(mktemp /tmp/weaver-smoke-test.XXXXXX)

# Check needed files exists
test -f ./bin/weaver
test -f ./bin/tester
test -f $REFERENCE_FILE
test -f $TEST_FUNCTIONS_FILE

echo "[*] Running Weaver with $TEST_FUNCTIONS_FILE"

./bin/weaver --json -f $TEST_FUNCTIONS_FILE ./bin/tester > $OUTPUT_FILE&

OSTER_PID=$!

# Set cleanup function
function cleanup {
    echo "[*] Killing weaver process..."
    kill -9 $OSTER_PID > /dev/null 2>&1
}
trap cleanup EXIT

echo "[*] Weaver running with pid $OSTER_PID, writing to $OUTPUT_FILE"

echo "[*] Waiting for uprobes to be installed"

sleep 10

echo "[*] Running test program"

./bin/tester
sleep 1

cat $OUTPUT_FILE | jq -s -c 'sort_by(.FunctionName)' | jq | tee $OUTPUT_FILE > /dev/null

echo "[*] Checking output"

diff $REFERENCE_FILE $OUTPUT_FILE

echo "[*] Looks good :)"
