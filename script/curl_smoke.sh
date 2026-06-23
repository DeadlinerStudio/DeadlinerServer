#!/usr/bin/env bash

set -euo pipefail

if ! command -v curl >/dev/null 2>&1; then
    echo "error: curl is required" >&2
    exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
    echo "error: jq is required for JSON parsing" >&2
    exit 1
fi

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
EMAIL="${EMAIL:-deadliner-smoke-$(date +%s)@example.com}"
PASSWORD="${PASSWORD:-deadliner-smoke-pass}"
DISPLAY_NAME="${DISPLAY_NAME:-Deadliner Smoke User}"
DEVICE_UID="${DEVICE_UID:-device-smoke-$(date +%s)}"
DEVICE_NAME="${DEVICE_NAME:-Local Curl Smoke}"
PLATFORM="${PLATFORM:-ios}"
PULL_LIMIT="${PULL_LIMIT:-50}"

WORK_DIR="$(mktemp -d "${TMPDIR:-/tmp}/deadliner-curl-smoke.XXXXXX")"
trap 'rm -rf "$WORK_DIR"' EXIT

LAST_RESPONSE_FILE=""

log_step() {
    printf '\n==> %s\n' "$1"
}

log_info() {
    printf '%s\n' "$1"
}

run_curl_json() {
    local label="$1"
    local method="$2"
    local path="$3"
    local body_file="${4:-}"
    local auth_token="${5:-}"

    local response_file="$WORK_DIR/${label}.json"
    local status_file="$WORK_DIR/${label}.status"
    local curl_args=(
        --silent
        --show-error
        --output "$response_file"
        --write-out "%{http_code}"
        -X "$method"
        -H "Content-Type: application/json"
    )

    if [ -n "$auth_token" ]; then
        curl_args+=(-H "Authorization: Bearer $auth_token")
    fi
    if [ -n "$body_file" ]; then
        curl_args+=(--data @"$body_file")
    fi

    local status_code
    status_code="$(curl "${curl_args[@]}" "${BASE_URL}${path}")"
    printf '%s' "$status_code" > "$status_file"

    LAST_RESPONSE_FILE="$response_file"

    printf 'HTTP %s\n' "$status_code"
    jq . "$response_file"

    if [ "$status_code" -lt 200 ] || [ "$status_code" -ge 300 ]; then
        echo "error: ${label} failed with status ${status_code}" >&2
        exit 1
    fi
}

read_json() {
    local jq_expr="$1"
    jq -r "$jq_expr" "$LAST_RESPONSE_FILE"
}

cat > "$WORK_DIR/register.json" <<EOF
{
  "email": "${EMAIL}",
  "password": "${PASSWORD}",
  "display_name": "${DISPLAY_NAME}",
  "device_uid": "${DEVICE_UID}",
  "device_name": "${DEVICE_NAME}",
  "platform": "${PLATFORM}"
}
EOF

cat > "$WORK_DIR/login.json" <<EOF
{
  "email": "${EMAIL}",
  "password": "${PASSWORD}",
  "device_uid": "${DEVICE_UID}",
  "device_name": "${DEVICE_NAME}",
  "platform": "${PLATFORM}"
}
EOF

DDL_UID="ddl-smoke-$(date +%s)"
HABIT_DDL_UID="ddl-habit-smoke-$(date +%s)"
DEADLINE_MUTATION_ID="mutation-deadline-smoke-$(date +%s)"
HABIT_MUTATION_ID="mutation-habit-smoke-$(date +%s)"
BUSINESS_TS="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

cat > "$WORK_DIR/push.json" <<EOF
{
  "device_uid": "${DEVICE_UID}",
  "base_cursor": "",
  "mutations": [
    {
      "mutation_id": "${DEADLINE_MUTATION_ID}",
      "device_uid": "${DEVICE_UID}",
      "entity_uid": "${DDL_UID}",
      "base_change_id": 0,
      "deadline": {
        "deleted": false,
        "doc": {
          "uid": "${DDL_UID}",
          "legacy_id": 0,
          "name": "Curl Smoke Task",
          "start_time": "",
          "end_time": "",
          "state": "active",
          "complete_time": "",
          "note": "created by script/curl_smoke.sh",
          "is_stared": false,
          "type": "task",
          "habit_count": 0,
          "habit_total_count": 0,
          "calendar_event": -1,
          "timestamp": "${BUSINESS_TS}",
          "sub_tasks": []
        }
      }
    },
    {
      "mutation_id": "${HABIT_MUTATION_ID}",
      "device_uid": "${DEVICE_UID}",
      "entity_uid": "${HABIT_DDL_UID}",
      "base_change_id": 0,
      "habit": {
        "deleted": false,
        "doc": {
          "ddl_uid": "${HABIT_DDL_UID}",
          "habit": {
            "name": "Curl Smoke Habit",
            "description": "created by script/curl_smoke.sh",
            "icon_key": "figure.walk",
            "period": "DAILY",
            "times_per_period": 1,
            "goal_type": "PER_PERIOD",
            "total_target": 0,
            "created_at": "${BUSINESS_TS}",
            "updated_at": "${BUSINESS_TS}",
            "status": "ACTIVE",
            "sort_order": 0,
            "alarm_time": ""
          },
          "records": [
            {
              "date": "$(date -u +"%Y-%m-%d")",
              "count": 1,
              "status": "COMPLETED",
              "created_at": "${BUSINESS_TS}"
            }
          ]
        }
      }
    }
  ]
}
EOF

log_info "Base URL: ${BASE_URL}"
log_info "Email: ${EMAIL}"
log_info "Device UID: ${DEVICE_UID}"

log_step "Register"
run_curl_json "register" "POST" "/v1/auth/register" "$WORK_DIR/register.json"
REGISTER_ACCOUNT_UID="$(read_json '.session.account_uid')"
REGISTER_ACCESS_TOKEN="$(read_json '.session.access_token')"
REGISTER_REFRESH_TOKEN="$(read_json '.session.refresh_token')"

log_step "Login"
run_curl_json "login" "POST" "/v1/auth/login" "$WORK_DIR/login.json"
LOGIN_ACCESS_TOKEN="$(read_json '.session.access_token')"
LOGIN_REFRESH_TOKEN="$(read_json '.session.refresh_token')"

log_step "Push"
run_curl_json "push" "POST" "/v1/sync/push" "$WORK_DIR/push.json" "$LOGIN_ACCESS_TOKEN"
NEXT_CURSOR="$(read_json '.next_cursor')"
PUSH_DEADLINE_COUNT="$(read_json '.deadline_changes | length')"
PUSH_HABIT_COUNT="$(read_json '.habit_changes | length')"

log_step "Pull"
run_curl_json "pull" "GET" "/v1/sync/pull?device_uid=${DEVICE_UID}&cursor=&limit=${PULL_LIMIT}&include_deleted=true" "" "$LOGIN_ACCESS_TOKEN"
PULL_NEXT_CURSOR="$(read_json '.next_cursor')"
PULL_DEADLINE_COUNT="$(read_json '.deadline_changes | length')"
PULL_HABIT_COUNT="$(read_json '.habit_changes | length')"

cat > "$WORK_DIR/refresh.json" <<EOF
{
  "refresh_token": "${LOGIN_REFRESH_TOKEN}",
  "device_uid": "${DEVICE_UID}"
}
EOF

log_step "Refresh"
run_curl_json "refresh" "POST" "/v1/auth/refresh" "$WORK_DIR/refresh.json"
REFRESH_ACCESS_TOKEN="$(read_json '.session.access_token')"

printf '\n==> Summary\n'
printf 'account_uid: %s\n' "$REGISTER_ACCOUNT_UID"
printf 'register_access_token: %s\n' "$REGISTER_ACCESS_TOKEN"
printf 'register_refresh_token: %s\n' "$REGISTER_REFRESH_TOKEN"
printf 'login_access_token: %s\n' "$LOGIN_ACCESS_TOKEN"
printf 'next_cursor_after_push: %s\n' "$NEXT_CURSOR"
printf 'push_deadline_changes: %s\n' "$PUSH_DEADLINE_COUNT"
printf 'push_habit_changes: %s\n' "$PUSH_HABIT_COUNT"
printf 'next_cursor_after_pull: %s\n' "$PULL_NEXT_CURSOR"
printf 'pull_deadline_changes: %s\n' "$PULL_DEADLINE_COUNT"
printf 'pull_habit_changes: %s\n' "$PULL_HABIT_COUNT"
printf 'refresh_access_token: %s\n' "$REFRESH_ACCESS_TOKEN"
printf 'ddl_uid: %s\n' "$DDL_UID"
printf 'habit_ddl_uid: %s\n' "$HABIT_DDL_UID"
