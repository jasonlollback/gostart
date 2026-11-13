# Specification GS-070--074 --- Generated project quality and release

Status: proposed\
Depends on: core generation, database and web foundation\
Related decisions: ADR-001, ADR-012, ADR-013

## Outcome

Generated applications have reproducible development commands, container
deployment examples and CI workflows. The generator verifies
representative generated projects.

## Taskfile and tools module

-   Use Taskfile, not Make.
-   Load `.env` for development tasks.
-   Application loads its own config with caarlos0/env.
-   Keep development tools in `/tools/go.mod`.
-   Pin versions for Templ, shadcn-templ CLI, sqlc and golangci-lint v2.
-   Include tasks only when relevant to selected features.

## Docker

-   Multi-stage build.
-   Non-root runtime.
-   Persistent SQLite volume.
-   Correct directory ownership.
-   Graceful shutdown.
-   Health check.
-   No mandatory reverse proxy in Compose.

## GitHub Actions

CI should run applicable formatting, code generation, tests, vet,
staticcheck, golangci-lint, govulncheck and build steps.

Docker publishing should be a separate workflow triggered by tags and
publish to GHCR with documented version/latest tags and least-privilege
permissions.

## Fixture matrix

At minimum: - Minimal Go app. - Web app. - Admin app. - SaaS-style
composition.

Each fixture must be generated and independently formatted,
code-generated, tested and built.

## Acceptance criteria

-   [ ] Taskfile works from a clean checkout.
-   [ ] Tool module is separate from app go.mod.
-   [ ] Docker image runs non-root.
-   [ ] SQLite data survives container recreation.
-   [ ] Generated CI workflow passes on a clean GitHub checkout.
-   [ ] GHCR workflow publishes only on intended tag events.
-   [ ] All supported fixtures pass generator and generated-app tests.
-   [ ] Disabled capabilities are absent from generated output.
