# 0001 Foundation

This package defines the baseline architecture for the open-source,
self-hostable Deadliner backend.

## Subdocuments

- [Goals And Scope](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation/goals-and-scope.md)
- [Architecture And Layering](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation/architecture-and-layering.md)
- [Sync And Persistence](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation/sync-and-persistence.md)
- [Delivery Plan](/Users/aritxonly/Codes/Golang/DeadlinerServer/docs/plan/0001-foundation/delivery-plan.md)

## Intent

The foundation plan exists to keep the codebase aligned around a few stable
decisions:

- centralized backend authority with local-first clients
- thrift + Kitex as the RPC contract and Go service scaffold
- `app`, `domain`, `infra`, `config`, and `utils` as the main repository shape
- multi-database persistence adapters with MySQL as the first production path
