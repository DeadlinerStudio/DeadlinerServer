#!/usr/bin/env bash

RUN_NAME="deadliner"
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
OUTPUT_DIR="$ROOT_DIR/output"
TMP_ROOT="${TMPDIR:-/tmp}"
TMP_ROOT="${TMP_ROOT%/}"

DEFAULT_GO_CACHE_DIR=$(go env GOCACHE 2>/dev/null || true)
DEFAULT_GO_MOD_CACHE_DIR=$(go env GOMODCACHE 2>/dev/null || true)

resolve_writable_dir() {
    local candidate="$1"
    local fallback="$2"
    local probe_file=""

    if [ -n "$candidate" ] && mkdir -p "$candidate" 2>/dev/null; then
        probe_file="$candidate/.deadlinerserver-write-test"
        if ( : > "$probe_file" ) 2>/dev/null; then
            rm -f "$probe_file"
            printf '%s\n' "$candidate"
            return 0
        fi
    fi

    mkdir -p "$fallback"
    printf '%s\n' "$fallback"
}

GO_CACHE_DIR=$(resolve_writable_dir "${GOCACHE:-$DEFAULT_GO_CACHE_DIR}" "${TMP_ROOT}/deadlinerserver-gocache")
GO_MOD_CACHE_DIR=$(resolve_writable_dir "${GOMODCACHE:-$DEFAULT_GO_MOD_CACHE_DIR}" "${TMP_ROOT}/deadlinerserver-gomodcache")

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
