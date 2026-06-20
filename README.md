# DeadlinerServer

DeadlinerServer is the open-source, self-hostable Go backend for Deadliner.

This repository is being built as a Kitex-based, multi-user backend with strict
account isolation, `app / domain / infra / service` layered architecture, and a
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

- [docs/plan/0001-foundation.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation.md)
- [docs/plan/0002-sync-product-shape.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0002-sync-product-shape.md)
- [docs/plan/0003-idempotency-and-convergence.md](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence.md)

## Layout

```text
main.go / handler.go       Kitex-generated service bootstrap
build.sh / script/         Kitex runtime packaging scripts
conf/                      JSON runtime configuration
db/migrations/             driver-specific SQL schema
docs/plan/                 SDD-style planning artifacts
idl/                       Kitex thrift contracts
internal/app/              app-layer commands and DTOs
internal/domain/           models, services, repository contracts
internal/infra/            GORM and storage adapters
internal/config/           local config defaults
internal/utils/            shared utility helpers
kitex_gen/                 generated thrift and service stubs
```

## Local Commands

```bash
make test
make test-rpc
make run
make generate
```

Notes:

- `make test` verifies the foundation packages under `internal/...`
- `make test-rpc` verifies the generated Kitex service packages too
- first-time Kitex dependency resolution may still require outbound module access
- runtime config is loaded from [conf/config.json](/Users/aritxonly/Codes/Golang/DeadlinerServer/conf/config.json)
- `infra/gorm` currently needs GORM dependency checksums in `go.sum` before full package tests can pass
- `make generate` requires the Kitex CLI to be installed locally
