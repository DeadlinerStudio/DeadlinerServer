# API Docs

This directory contains the HTTP API contract for DeadlinerServer.

Files:

- [openapi.yaml](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/api/openapi.yaml): machine-readable OpenAPI 3.0 description

## Current Scope

The documented surface currently covers the first public HTTP endpoints:

- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `POST /v1/auth/refresh`
- `GET /v1/sync/pull`
- `POST /v1/sync/push`
- `GET /healthz`
- `GET /admin/config`
- `GET /admin/api/config`
- `PUT /admin/api/config`

## Security Defaults

The HTTP server now applies several protections before requests reach business
handlers:

- access token enforcement on `/v1/sync/*`
- `application/json` enforcement on HTTP write endpoints
- request body size limits
- per-client in-memory rate limits
- request IDs on every response
- generic 5xx responses that avoid leaking internal details

## Config Admin

The first config backend is intentionally small:

- an HTML page at `/admin/config`
- a protected JSON API at `/admin/api/config`
- only non-sensitive runtime config is editable
- secret values stay in `conf/secret.json`
- server restart is required after saving updates

## Notes

- Preferred client authentication is `Authorization: Bearer <access-token>`.
- Error responses include `request_id` to help correlate client failures with
  server logs.
- Habit `color` is intentionally not part of synchronized habit documents.
