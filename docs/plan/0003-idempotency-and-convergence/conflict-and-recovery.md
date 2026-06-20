# Conflict And Recovery

## Why This Beats Blind LWW

Blind last-writer-wins makes retries easy to code, but it can hide real
problems:

- duplicate mutation may reapply
- stale device may overwrite newer state
- user loses data without a clear reason

The stronger model should be:

- duplicate retry becomes the same receipt with no double write
- stale mutation becomes explicit conflict with no silent overwrite
- successful mutation returns a stable canonical change id

## Recommended Entity Semantics

### DDL

- whole-document replace is acceptable in phase 1
- guarded by `base_change_id`

### Habit

- whole-document replace remains acceptable in phase 1
- guarded by carrier `ddl_uid` and `base_change_id`

### Habit Records

- do not implement free-form incremental counters first
- keep replace-whole-document semantics for safety

## Local Recovery Flow

When client push gets `conflict`:

1. keep local unsynced content
2. pull latest server state
3. mark the local object as needing reconciliation
4. let the client decide whether to auto-merge or ask the user later
