# AGENTS.md

## Product identity

- The product and CLI are named `gostart`, always lowercase.
- This repository is an independent fork of Goilerplate. Do not describe gostart as the original Goilerplate product.
- Generated applications are self-contained and must not depend on gostart at runtime.
- Do not add automatic template-update functionality.

## Product and architecture

- Read `docs/PRODUCT.md`, `docs/ARCHITECTURE.md` and `docs/DECISIONS.md` before changing generator or template architecture.
- Read the relevant specification in `specs/` before implementing a backlog item.
- Keep `docs/BACKLOG.md` current when work is completed or scope changes.
- The initial generated application profile is Go, `net/http`, `http.ServeMux`, SQLite with `modernc.org/sqlite`, `database/sql`, sqlc, Templ, Tailwind CSS 4, shadcn-templ, HTMX, slog and tint.
- V1 supports SQLite only. Do not add PostgreSQL or alternative frontend frameworks without an explicit product decision.
- Use feature-oriented packages. Do not impose `/internal/platform` and `/internal/modules`.
- Prefer the standard library. Add dependencies only when they materially improve the result.
- Avoid interfaces without a real substitution or testing need.

## Generated application conventions

- Use `modernc.org/sqlite` with `_time_format=sqlite&_timezone=UTC`.
- Use `database/sql`, WAL, foreign keys, STRICT tables, TEXT timestamps and UUIDv7 identifiers.
- Use sqlc for typed query generation.
- Keep development tools in a separate `/tools` Go module.
- Use Taskfile, not Makefiles.
- Keep `.env.example` aligned with the selected generated features.
- Initial administrator provisioning is triggered by `ADMIN_EMAIL` alone, is idempotent, and must not continually re-promote an account on every startup.
- Authentication uses email one-time codes and database-backed sessions; do not introduce password authentication by default.
- Migrations are embedded, sequential, forward-only and checksum-verified. They run before the HTTP server becomes ready.

## Engineering

- Keep packages small and concrete. Separate rendering, validation and filesystem writing.
- Validate the complete generation plan before writing output.
- Generation must be deterministic for identical inputs.
- Never silently overwrite existing user files.
- Do not execute arbitrary scripts or hooks from feature manifests.
- Test generated output, not only generator internals.
- Run formatting, vet, tests and build checks appropriate to the changed code.
- Update or add specifications and acceptance tests when behaviour changes.

## Git and changes

- Use a focused branch and pull request for substantial work.
- Keep commits small and use Conventional Commit prefixes: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`.
- Preserve upstream history and the `upstream` remote relationship.
- Retain original license and copyright notices for inherited code. Do not replace upstream attribution or imply inherited code was authored by gostart.
- Do not add secrets, private customer data or local machine configuration.
