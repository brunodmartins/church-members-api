# Copilot / AI Assistant Instructions

This file gives immediate, actionable context for AI coding agents working in this repository.

Overview
- Language: Go (modules). Entry: `main.go` and the `cmd/` package.
- Architecture: small monolith split into HTTP handlers (`internal/handler`), domain models (`internal/domain`), DTOs and converters (`internal/dto`, `internal/constants/dto`), module services (`modules/*`), and infra/platform code (`infra/`, `platform/`).

Key patterns and conventions
- Domain objects live in `internal/domain` (IDs are string-based; check `internal/domain/id.go`).
- DTOs and HTTP shapes live under `internal/constants/dto/response.go` and `internal/dto/converter.go` — use the converter utilities when translating domain -> response types.
- Handlers follow the `internal/handler/*` layout (see `internal/handler/api/routes.go` and `internal/handler/api/*_handler.go`). Prefer reusing existing validation and middleware located in `internal/handler/validator.go` and `internal/handler/middleware/`.
- Modules (business logic) live under `modules/<area>/service.go` and `modules/<area>/repository.go`. Tests use mocks under `modules/<area>/mock`.
- Jobs/cron-like tasks live in `modules/jobs` and `cmd/job.go`.
- Platform and AWS-specific wrappers are under `platform/aws` and `cdi/` for container/DI helpers.

Build / test / run (discoverable commands)
- Run unit tests: `go test ./...` (or a package: `go test ./internal/... -run TestName`).
- Run the app locally: `go run main.go` or use Docker: `docker-compose up` (Dockerfile and docker-compose.yml are in repository root).
- Formatting: use `gofmt` / `go fmt ./...` before committing.

Project-specific advice for code changes
- When adding or changing public HTTP responses, update `internal/constants/dto/response.go` and the converter in `internal/dto/converter.go`.
- Follow existing error handling in `infra/errors` and reuse common error shapes rather than introducing new ad-hoc HTTP formats.
- Authentication & JWT: check `infra/jwt.go` and `infra/auth_service.go` for token construction/validation.
- Persistence: parts of the project use DynamoDB (see `terraform/storage/dynamodb` and `test/dynamodbhelper`). When adding repository code, follow existing repository interfaces in `modules/*/repository.go`.

Testing notes
- Unit tests are standard `*_test.go` files. Use existing test helpers in `test/dynamodbhelper` for DB-related tests.
- Many tests use table-driven style and rely on small mock packages under `modules/*/mock` — follow that pattern.

Integration points / external systems
- DynamoDB: terraform and test helpers present. Keep AWS-specific code scoped to `platform/aws` and DI wrappers.
- Email/notification/storage services have adapters under `services/` — prefer adding adapters there and injecting them into modules.

When editing code, prefer minimal, focused changes
- Keep API-compatible changes to DTOs and document them in `docs/specs/swagger.yaml` when altering public endpoints.

Files to inspect for examples
- `internal/constants/dto/response.go` — canonical response structs.
- `internal/dto/converter.go` — domain -> dto mapping helpers.
- `internal/handler/api/routes.go` & `internal/handler/api/*_handler.go` — routing and handler patterns.
- `modules/*/service.go` and `modules/*/repository.go` — typical service/repository separation.

If anything here is unclear or you want more detail (examples of converters, test harnesses, or common Makefile targets), ask and I'll iteratively refine this file.
