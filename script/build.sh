#!/usr/bin/env bash

RUN_NAME="deadliner"
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
OUTPUT_DIR="$ROOT_DIR/output"
GO_CACHE_DIR="${GOCACHE:-/private/tmp/deadlinerserver-gocache}"
GO_MOD_CACHE_DIR="${GOMODCACHE:-/private/tmp/deadlinerserver-gomodcache}"

mkdir -p "$OUTPUT_DIR/bin"
mkdir -p "$GO_CACHE_DIR" "$GO_MOD_CACHE_DIR"
cp "$SCRIPT_DIR/bootstrap.sh" "$OUTPUT_DIR/bootstrap.sh"
chmod +x "$OUTPUT_DIR/bootstrap.sh"

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    (
        cd "$ROOT_DIR" || exit 1
        GOCACHE="$GO_CACHE_DIR" GOMODCACHE="$GO_MOD_CACHE_DIR" go build -o "$OUTPUT_DIR/bin/${RUN_NAME}" ./cmd/deadlinerserver
    )
else
    (
        cd "$ROOT_DIR" || exit 1
        GOCACHE="$GO_CACHE_DIR" GOMODCACHE="$GO_MOD_CACHE_DIR" go test -c -covermode=set -o "$OUTPUT_DIR/bin/${RUN_NAME}" -coverpkg=./... ./cmd/deadlinerserver
    )
fi
