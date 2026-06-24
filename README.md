# DeadlinerServer

DeadlinerServer is the open-source, self-hostable Go backend for Deadliner.

This repository is being built as a Kitex + Hertz, multi-user backend with
strict account isolation, `app / domain / infra` layered architecture, and a
backend-first sync model for iOS, Android, and HarmonyOS.

## Current Scope

The current foundation includes:

- an SDD-style implementation plan
- a first-pass Kitex thrift contract
- a first-pass Hertz HTTP client API
- a DDD-oriented service skeleton
- a multi-database-ready persistence skeleton with MySQL as the default

## Why This Server Exists

The iOS app already has a stable local sync shape with:

- `uid`
- logical versions `Ver(ts, ctr, dev)`
- tombstones independent from business state
- habit documents anchored by a carrier DDL uid

The server should preserve those domain rules where they are valuable, but it
should not remain a file-based client-coordinated sync system forever. This
repository moves Deadliner toward a backend-first model where:

- the server becomes the canonical source of truth
- clients submit mutations instead of overwriting remote snapshots
- the server assigns canonical change versions and pull cursors
- account data is isolated by tenant boundaries

## Plan

The initial architecture and delivery plan lives at:

- [docs/plan/README.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/README.md)
- [docs/plan/0001-foundation/README.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation/README.md)
- [docs/plan/0002-sync-product-shape/README.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0002-sync-product-shape/README.md)
- [docs/plan/0003-idempotency-and-convergence/README.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence/README.md)
- [docs/api/README.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/api/README.md)

## Layout

```text
cmd/deadlinerserver/       process entrypoint
conf/                      JSON runtime configuration
db/migrations/             driver-specific SQL schema
docs/plan/                 SDD-style planning artifacts
idl/                       Kitex thrift contracts
internal/app/              bootstrap, auth context, use cases, Kitex/Hertz transport
internal/domain/           models, services, repository and provider contracts
internal/infra/            GORM and provider adapters
internal/config/           local config defaults
internal/utils/            shared utility helpers
script/                    packaging and runtime bootstrap scripts
kitex_gen/                 generated thrift and service stubs
```

The repository root intentionally keeps only the conventional Go / Kitex /
open-source entry files, such as `go.mod`, `go.sum`, `Makefile`, `README.md`,
`LICENSE`, `.gitignore`, and `kitex_info.yaml`.

Within `docs/plan/`, each numbered topic now uses a small directory package
with a `README` and focused subdocuments so planning can evolve without
creating another 300-line monolith.

## Local Commands

```bash
make test
make test-rpc
make run
make package
make generate
```

## Config Split

Non-sensitive runtime settings stay in [conf/config.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/config.json).

Sensitive settings are now centralized in:

- `conf/secret.json`

The repository only tracks the example file:

- [conf/secret.example.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/secret.example.json)

The current sensitive config set is:

- `auth.accessTokenSecret`
- `database.dsn`

You can also override the secret config path with `DEADLINER_SECRET_CONFIG`.

For Docker Compose, the repository also includes:

- [conf/secret.docker.example.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/secret.docker.example.json)

Create your own untracked `conf/secret.docker.json` from that example before
starting the stack.

## Interface Logs

The server now emits lightweight interface access logs for both transports:

- `HTTP ...` for Hertz HTTP requests
- `KITEX ...` for Kitex RPC requests

The logs intentionally record only metadata such as method, route, status,
latency, caller IP / address, and payload size. They do not print request
bodies, passwords, or tokens.

## API Docs

HTTP API documentation now lives in:

- [docs/api/openapi.yaml](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/api/openapi.yaml)

The OpenAPI document currently covers health, auth, and sync endpoints.

## HTTP Security

The Hertz server now enables a small default security middleware chain:

- request IDs on every response
- generic 5xx error bodies with `request_id`
- bearer token enforcement on `/v1/sync/*`
- `application/json` enforcement on write endpoints
- request body size limits
- in-memory per-client rate limits
- conservative security headers for API responses

## Docker

The repository now includes:

- [Dockerfile](/Users/aritxonly/Codes/Golang/DeadlinerServer/Dockerfile)
- [compose.yaml](/Users/aritxonly/Codes/Golang/DeadlinerServer/compose.yaml)
- [conf/secret.docker.example.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/secret.docker.example.json)

Typical local flow:

```bash
cp conf/secret.docker.example.json conf/secret.docker.json
docker compose up --build
```

This starts:

- `deadliner` on `8080` for HTTP and `8888` for Kitex
- `mysql` on `3306`

## CI

GitHub Actions now runs [`.github/workflows/ci.yml`](/Users/aritxonly/Codes/Golang/DeadlinerServer/.github/workflows/ci.yml) on pull requests and pushes to `main`.

The workflow currently does three things:

- verifies `gofmt` cleanliness under `cmd/` and `internal/`
- runs the full Go test suite
- packages the Linux server binary and uploads it as a workflow artifact

Notes:

- `make test` verifies the foundation packages under `internal/...`
- `make test-rpc` verifies the generated Kitex service packages too
- `make run` now starts [cmd/deadlinerserver/main.go](/Users/aritxonly/Codes/Golang/DeadlinerServer/cmd/deadlinerserver/main.go)
- `make package` now runs [script/build.sh](/Users/aritxonly/Codes/Golang/DeadlinerServer/script/build.sh)
- runtime config is loaded from [conf/config.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/config.json) and sensitive settings are merged from `conf/secret.json`
- Kitex RPC listens on `service.address` and Hertz HTTP listens on `http.address`
- the first HTTP endpoints are `/v1/auth/register`, `/v1/auth/login`, `/v1/auth/refresh`, `/v1/sync/pull`, and `/v1/sync/push`
- auth config now includes access token secret, token TTLs, and password hash cost
- `make generate` requires the Kitex CLI to be installed locally
