# Server Rules

## Receipt Table For Deduplication

The server persists a mutation receipt keyed by:

- `account_id`
- `device_uid`
- `mutation_id`

On duplicate push:

- do not reapply the mutation
- return the previously recorded result

## Append-Only Canonical Ordering

Every accepted mutation gets a canonical `change_id`.

Rules:

- globally ordered within one account
- monotonic
- used as the pull cursor

## Entity-Level Version Check

Each entity stores its current `server_change_id`.

When a mutation arrives:

- if `base_change_id` matches current entity version, apply normally
- if it does not match, treat the mutation as stale and resolve
  deterministically

Phase-1 recommendation:

- reject stale mutation with a conflict result
- return latest authoritative entity state

## Push Contract

Recommended push behavior:

1. client sends pending mutations in local sequence order
2. server processes each mutation transactionally
3. server returns `applied`, `replayed`, `conflict`, or `rejected`
4. client marks queue entries accordingly
5. client performs a follow-up pull to the returned cursor

## Pull Contract

Pull by cursor is naturally idempotent:

- `PullChanges(cursor=X)` can be retried safely
- server returns all changes after `X` in canonical order
- applying the same pulled change twice must be safe locally
