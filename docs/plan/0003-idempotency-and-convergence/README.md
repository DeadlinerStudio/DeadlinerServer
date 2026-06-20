# 0003 Idempotency And Convergence

This package defines how Deadliner Account Sync should become more stable than
WebDAV under retries, weak network, and concurrent device activity.

## Subdocuments

- [Client Rules](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence/client-rules.md)
- [Server Rules](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence/server-rules.md)
- [Conflict And Recovery](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence/conflict-and-recovery.md)
- [Phase 1 Guarantees](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0003-idempotency-and-convergence/phase1-guarantees.md)

## Intent

The unit of sync should be a durable mutation, not a remote snapshot rewrite.
