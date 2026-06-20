# Settings And Migration

## Product Boundary Versus WebDAV

The distinction should be explicit in both code and UX.

### WebDAV Mode

- remote file snapshots
- client-led merge
- distributed synchronization flavor

### Deadliner Account Sync Mode

- authenticated account session
- mutation push and cursor pull
- server-led canonical ordering

## Recommended Settings IA

### Sync

- `Sync Mode`
- `Local Only`
- `WebDAV`
- `Deadliner Account`

### Deadliner Account

- sign in and sign out
- device list
- last sync status
- force resync

### WebDAV

- server URL
- username
- password or token
- test connection

## Migration Strategy

Do not force users from WebDAV to the new backend.

Recommended path:

1. user selects `Deadliner Account`
2. app offers import from current local data
3. first upload creates the canonical remote state
4. other devices sign in and pull from server

For WebDAV users specifically:

- migration should be explicit and one-time
- avoid dual-write to both systems in phase 1
- after switching, one sync mode should be active at a time

## Why Not Dual Sync At Once

Supporting WebDAV and centralized backend at the same time for one dataset
would create:

- two remote authorities
- ambiguous conflict semantics
- much harder support burden
