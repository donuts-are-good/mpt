#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_PATH="$ROOT_DIR/tmp/mpt"
TMP_DIR="$ROOT_DIR/tmp"
RUN_DIR="$TMP_DIR/test-suite"
CONTENT_DIR="$ROOT_DIR/content"

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

echo "[single] json -> msgpack"
"$BIN_PATH" "$RUN_DIR/view-source.json" "$RUN_DIR/view-source.msgpack"

echo "[view] msgpack -> stdout (json)"
"$BIN_PATH" --view "$RUN_DIR/view-source.msgpack" > "$RUN_DIR/view-output.json"

echo "[single] msgpack -> json"
"$BIN_PATH" "$RUN_DIR/view-source.msgpack" "$RUN_DIR/view-source.roundtrip.json"

echo "[single] yaml -> msgpack"
"$BIN_PATH" "$RUN_DIR/demo1.yaml" "$RUN_DIR/demo1.msgpack"

echo "[single] msgpack -> yaml"
"$BIN_PATH" "$RUN_DIR/demo1.msgpack" "$RUN_DIR/demo1.roundtrip.yaml"

echo "[override] json (weird extension) -> msgpack"
"$BIN_PATH" --from json --to msgpack "$RUN_DIR/weird.bin" "$RUN_DIR/weird-output.bin"

echo "[override] msgpack (weird extension) -> json"
"$BIN_PATH" --from msgpack --to json "$RUN_DIR/weird-output.bin" "$RUN_DIR/weird-output.json"

echo "[stdout] json -> stdout (json)"
"$BIN_PATH" "$RUN_DIR/stdout-json-source.json" --json > "$RUN_DIR/stdout-json-output.json"

echo "[stdout] json -> stdout (yaml)"
"$BIN_PATH" "$RUN_DIR/stdout-yaml-source.json" --yaml > "$RUN_DIR/stdout-yaml-output.yaml"

echo "[batch] preparing batch directories"
BATCH_DIR="$RUN_DIR/batch"
mkdir -p "$BATCH_DIR/json" "$BATCH_DIR/msgpack"
cp "$CONTENT_DIR/json/demo0.json" "$BATCH_DIR/json/demo0.json"
cp "$CONTENT_DIR/json/demo1.json" "$BATCH_DIR/json/demo1.json"

echo "[batch] *.json -> *.msgpack"
"$BIN_PATH" "$BATCH_DIR/json"/*.json --to-msgpack

echo "[batch] copying msgpack inputs"
cp "$BATCH_DIR/json"/*.msgpack "$BATCH_DIR/msgpack/"

echo "[batch] *.msgpack -> *.json"
"$BIN_PATH" "$BATCH_DIR/msgpack"/*.msgpack --to-json

echo "[batch] *.msgpack -> *.yaml"
"$BIN_PATH" "$BATCH_DIR/msgpack"/*.msgpack --to-yaml

echo "all commands completed"

