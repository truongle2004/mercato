# Mercato agent instructions

These instructions are the canonical repository guidance for coding agents. Keep tool-specific entry files thin and put shared guidance under `.agents/`.

## Project

Mercato is a Go e-commerce service with REST and gRPC transports, PostgreSQL, Redis, Stripe payments, notifications, and a React/Vite storefront under `web/`.

Read `.agents/rules/project-context.md` when a task needs the detailed domain map or shipped-feature context.

## Commands

```bash
go build -o main ./cmd/api
make unittest
go test ./internal/product/service/... -v -run TestProductServiceTestSuite/TestCreateSuccess
golangci-lint run
make mock
make doc
cd proto && buf generate
make migrate-up
make migrate-new name=add_index
```

Before committing Go changes, run the narrowest relevant tests, then `go vet ./...`, `golangci-lint run`, `go test ./...`, and `go mod tidy` when dependencies changed. Integration tests require Docker and run with `make integration`.

## Architecture invariants

- Preserve the handler → service → repository dependency direction.
- Keep business logic out of transports and repositories.
- Depend on interfaces at boundaries and update generated mocks when interfaces change.
- Pass `context.Context` as the first parameter for I/O-capable operations.
- Map domain errors consistently at HTTP and gRPC boundaries.
- Manage schema changes only through paired, versioned SQL files in `migrations/`; do not introduce `AutoMigrate`.
- Do not change public HTTP, gRPC, or protobuf contracts without explicit approval.
- Never log credentials, tokens, payment secrets, or sensitive personal data.

## Relevant rules

Read only the files relevant to the task:

- Go code: `.agents/rules/go-conventions.md`
- Architecture or domain boundaries: `.agents/rules/architecture.md`
- HTTP or gRPC APIs: `.agents/rules/api-design.md`
- Repositories, models, or migrations: `.agents/rules/database.md`
- Tests or mocks: `.agents/rules/testing.md`
- Authentication, authorization, or sensitive data: `.agents/rules/security.md`
- Performance-sensitive work: `.agents/rules/performance.md`
- Logging, tracing, or metrics: `.agents/rules/logging-observability.md`
- Configuration: `.agents/rules/configuration.md`
- Protobuf definitions or generation: `.agents/rules/grpc-proto.md`
- Shared `pkg/` code: `.agents/rules/shared-packages.md`
- Reviews: `.agents/rules/code-review.md`
- Git or commits: `.agents/rules/git-workflow.md`
- Complete command reference: `.agents/rules/commands.md`

## Portable skills

Reusable workflows live in `.agents/skills/*/SKILL.md`. Load the matching skill when the request is specifically about API scaffolding, architecture, debugging, API documentation, migrations, performance, refactoring, review, security auditing, or test generation.

## Commits

Use `type(scope): concise description` with one concern per commit. Commit using the configured Git identity and do not append agent co-author trailers.
