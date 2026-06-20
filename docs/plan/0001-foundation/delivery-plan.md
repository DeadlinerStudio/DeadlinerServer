# Delivery Plan

## Mobile Client Migration Strategy

### iOS

- map current `Ver(ts, ctr, dev)` to client mutation metadata
- keep current DDL and habit payload shapes as close as possible
- replace WebDAV snapshot reads and writes with push and pull calls

### Android

- align local model with canonical `state`
- embed subtasks in the DDL payload
- use the same carrier DDL habit model as iOS

### HarmonyOS

- follow the same thrift contract as Android and iOS
- keep compatibility logic at the client boundary only

## Milestones

### M1 Foundation

- plan documents
- thrift contract
- repository skeleton
- multi-database skeleton with MySQL schema

### M2 Auth And Session

- account registration
- password login
- refresh token rotation
- device registration

### M3 Sync Write Path

- push deadline mutations
- push habit mutations
- idempotent mutation handling
- change feed append

### M4 Sync Read Path

- pull by cursor
- pagination
- deadline and habit change serialization

### M5 Client Adoption

- iOS migration
- Android migration
- HarmonyOS migration

## Acceptance Criteria

The foundation is complete when:

- the repository clearly explains the target architecture
- the thrift contract reflects backend-first sync
- the default MySQL schema supports multi-user isolation
- Go packages exist for continued implementation without reorganizing later

## Open Questions

1. Should account identity be email-only or support username login too?
2. Should session auth stay token-based only, or also expose a proxy-friendly
   cookie mode later?
3. Do we want a public REST gateway after the Kitex RPC layer stabilizes?
4. Should the canonical server version remain only `change_id`, or should we
   expose both `change_id` and RFC3339 `committed_at` everywhere?
