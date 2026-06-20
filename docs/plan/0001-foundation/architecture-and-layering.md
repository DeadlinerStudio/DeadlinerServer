# Architecture And Layering

## Backend-First Adaptation

The current mobile model assumes each device can author canonical writes. The
server changes that assumption:

- client logical versions are accepted as mutation metadata
- server commit order becomes the canonical conflict boundary
- the server emits monotonic change ids as pull cursors
- the server stores mutation provenance for debugging and idempotency

This gives us a practical bridge:

1. keep `uid` stable across platforms
2. keep tombstones independent
3. keep habit carrier semantics
4. move canonical ordering authority to the server

## Transport Choice

Use Kitex with thrift as the primary internal and mobile-facing RPC layer.

Reasons:

- explicit contracts for all three clients
- good fit for strongly typed Go services
- easy future split into auth and sync services
- consistent scaffolding for open-source maintenance

## Service Boundary

Phase 1 starts with one Kitex service:

- `DeadlinerService`

Initial RPC groups:

- account registration
- login and session refresh
- mutation push
- incremental pull

## Repository Layering

Recommended repository shape:

- `internal/app`: application composition, transport mapping, auth boundary
- `internal/domain`: business models, services, and contracts
- `internal/infra`: GORM adapters, SQL drivers, token and time providers
- `internal/config`: runtime configuration loading and defaults
- `internal/utils`: small shared helpers that are not domain logic

Inside the code, domain-facing abstractions should be named by capability, not
by storage technology. SQL and GORM details stay on the infra side.
