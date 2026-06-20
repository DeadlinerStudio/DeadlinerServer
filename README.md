# DeadlinerServer

DeadlinerServer is the open-source, self-hostable Go backend for Deadliner.

This repository is being built as a Kitex-based, multi-user backend with strict
account isolation, `app / domain / infra` layered architecture, and a
backend-first sync model for iOS, Android, and HarmonyOS.

## Current Scope

The current foundation includes:

- an SDD-style implementation plan
- a first-pass Kitex thrift contract
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

## Layout

```text
cmd/deadlinerserver/       process entrypoint
conf/                      JSON runtime configuration
db/migrations/             driver-specific SQL schema
docs/plan/                 SDD-style planning artifacts
idl/                       Kitex thrift contracts
internal/app/              bootstrap, auth context, use cases, transport mappers
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
- runtime config is loaded from [conf/config.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/config.json)
- auth config now includes access token secret, token TTLs, and password hash cost
- `make generate` requires the Kitex CLI to be installed locally
