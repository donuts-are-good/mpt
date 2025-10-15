#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_PATH="$ROOT_DIR/tmp/mpt"
TMP_DIR="$ROOT_DIR/tmp"
RUN_DIR="$TMP_DIR/test-suite"
CONTENT_DIR="$ROOT_DIR/content"

PASS=0
FAIL=0

pass() { echo "✓ $1"; ((PASS++)); }
fail() { echo "✗ $1"; ((FAIL++)); }

check_file_exists() {
    if [[ -f "$1" ]]; then
        pass "$2"
    else
        fail "$2: file missing"
    fi
}

check_valid_json() {
    if jq empty "$1" 2>/dev/null; then
        pass "$2"
    else
        fail "$2: invalid json"
    fi
}

check_json_equal() {
    if diff <(jq -S . "$1") <(jq -S . "$2") >/dev/null 2>&1; then
        pass "$3"
    else
        fail "$3: content mismatch"
    fi
}

echo "[setup] building mpt binary"
(cd "$ROOT_DIR" && go build -o "$BIN_PATH" .)

echo "[setup] preparing workspace at $RUN_DIR"
rm -rf "$RUN_DIR"
mkdir -p "$RUN_DIR"

echo "[setup] staging sample data"
cp "$CONTENT_DIR/json/demo1.json" "$RUN_DIR/view-source.json"
cp "$CONTENT_DIR/json/demo2.json" "$RUN_DIR/weird.bin"
cp "$CONTENT_DIR/json/demo3.json" "$RUN_DIR/stdout-json-source.json"
cp "$CONTENT_DIR/json/demo4.json" "$RUN_DIR/stdout-yaml-source.json"
cp "$CONTENT_DIR/yaml/demo1.yaml" "$RUN_DIR/demo1.yaml"

set +e

echo ""
echo "[test] json -> msgpack -> json roundtrip"
"$BIN_PATH" "$RUN_DIR/view-source.json" "$RUN_DIR/view-source.msgpack"
check_file_exists "$RUN_DIR/view-source.msgpack" "msgpack created"
"$BIN_PATH" "$RUN_DIR/view-source.msgpack" "$RUN_DIR/view-source.roundtrip.json"
check_file_exists "$RUN_DIR/view-source.roundtrip.json" "roundtrip json created"
check_json_equal "$RUN_DIR/view-source.json" "$RUN_DIR/view-source.roundtrip.json" "roundtrip preserves data"

echo ""
echo "[test] --view outputs valid json"
"$BIN_PATH" --view "$RUN_DIR/view-source.msgpack" > "$RUN_DIR/view-output.json"
check_valid_json "$RUN_DIR/view-output.json" "view output is valid json"
check_json_equal "$RUN_DIR/view-source.json" "$RUN_DIR/view-output.json" "view matches original"

echo ""
echo "[test] --view with json file"
"$BIN_PATH" --view "$RUN_DIR/view-source.json" > "$RUN_DIR/view-json.json"
check_valid_json "$RUN_DIR/view-json.json" "view json file produces valid json"
check_json_equal "$RUN_DIR/view-source.json" "$RUN_DIR/view-json.json" "view json preserves data"

echo ""
echo "[test] yaml roundtrip"
"$BIN_PATH" "$RUN_DIR/demo1.yaml" "$RUN_DIR/demo1.msgpack"
check_file_exists "$RUN_DIR/demo1.msgpack" "yaml to msgpack created"
"$BIN_PATH" "$RUN_DIR/demo1.msgpack" "$RUN_DIR/demo1.roundtrip.yaml"
check_file_exists "$RUN_DIR/demo1.roundtrip.yaml" "msgpack to yaml created"

echo ""
echo "[test] format override"
"$BIN_PATH" --from json --to msgpack "$RUN_DIR/weird.bin" "$RUN_DIR/weird-output.bin"
check_file_exists "$RUN_DIR/weird-output.bin" "override created msgpack"
"$BIN_PATH" --from msgpack --to json "$RUN_DIR/weird-output.bin" "$RUN_DIR/weird-output.json"
check_valid_json "$RUN_DIR/weird-output.json" "override roundtrip is valid json"
check_json_equal "$RUN_DIR/weird.bin" "$RUN_DIR/weird-output.json" "override preserves data"

echo ""
echo "[test] stdout conversions"
"$BIN_PATH" "$RUN_DIR/stdout-json-source.json" --json > "$RUN_DIR/stdout-json-output.json"
check_valid_json "$RUN_DIR/stdout-json-output.json" "stdout json is valid"
check_json_equal "$RUN_DIR/stdout-json-source.json" "$RUN_DIR/stdout-json-output.json" "stdout json preserves data"

"$BIN_PATH" "$RUN_DIR/stdout-yaml-source.json" --yaml > "$RUN_DIR/stdout-yaml-output.yaml"
check_file_exists "$RUN_DIR/stdout-yaml-output.yaml" "stdout yaml created"

echo ""
echo "[test] batch conversion"
BATCH_DIR="$RUN_DIR/batch"
mkdir -p "$BATCH_DIR/json" "$BATCH_DIR/msgpack"
cp "$CONTENT_DIR/json/demo0.json" "$BATCH_DIR/json/demo0.json"
cp "$CONTENT_DIR/json/demo1.json" "$BATCH_DIR/json/demo1.json"

"$BIN_PATH" "$BATCH_DIR/json"/*.json --to-msgpack
check_file_exists "$BATCH_DIR/json/demo0.msgpack" "batch created demo0.msgpack"
check_file_exists "$BATCH_DIR/json/demo1.msgpack" "batch created demo1.msgpack"

cp "$BATCH_DIR/json"/*.msgpack "$BATCH_DIR/msgpack/"
"$BIN_PATH" "$BATCH_DIR/msgpack"/*.msgpack --to-json
check_file_exists "$BATCH_DIR/msgpack/demo0.json" "batch reconvert demo0.json"
check_file_exists "$BATCH_DIR/msgpack/demo1.json" "batch reconvert demo1.json"
check_json_equal "$BATCH_DIR/json/demo0.json" "$BATCH_DIR/msgpack/demo0.json" "batch roundtrip demo0"
check_json_equal "$BATCH_DIR/json/demo1.json" "$BATCH_DIR/msgpack/demo1.json" "batch roundtrip demo1"

"$BIN_PATH" "$BATCH_DIR/msgpack"/*.msgpack --to-yaml
check_file_exists "$BATCH_DIR/msgpack/demo0.yaml" "batch created demo0.yaml"
check_file_exists "$BATCH_DIR/msgpack/demo1.yaml" "batch created demo1.yaml"

echo ""
echo "=========================================="
echo "passed: $PASS"
echo "failed: $FAIL"
echo "=========================================="

if [[ $FAIL -gt 0 ]]; then
    exit 1
fi

