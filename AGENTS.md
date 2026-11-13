# AGENTS.md

## Product identity

- The product and CLI are named `gostart`, always lowercase.
- This repository is a customised fork of Goilerplate. Do not describe gostart as the original Goilerplate product.
- Goilerplate's existing generator, template-selection mechanism and project structure are the implementation foundation.
- Generated applications remain independently owned Go projects. Do not add a runtime dependency on gostart unless explicitly designed and approved.

## Working with the existing codebase

- Understand and follow the existing Goilerplate implementation before proposing architectural changes.
- Prefer incremental customisation of existing packages, commands and template conventions over wholesale rewrites.
- Preserve working capabilities unless there is a clear product, maintenance, security or complexity reason to change or remove them.
- Add features using the existing generator's established patterns where practical.
- Remove unwanted capabilities deliberately, including their template files, dependencies, configuration, tests, documentation and UI references.
- Avoid introducing a second template engine, feature registry or parallel architecture without a demonstrated need.
- Keep changes focused and preserve existing tests. Add regression tests for changed generator behaviour and generated output.
- Do not port the gostart_v1 implementation or impose its architecture as a requirement.

## Product preferences and reference material

- The former `gostart_v1` repository contains useful product preferences and design notes, not a binding implementation specification.
- Its documents may inform decisions about generated app defaults, including Go/net/http, SQLite, sqlc, Templ, Tailwind CSS 4, shadcn-templ, HTMX, Taskfile, a separate tools module, Docker and GitHub Actions.
- These preferences should be evaluated against the existing Goilerplate implementation and actual product requirements. Do not force them all into the codebase.
- Do not assume the v1 non-goals or feature scope apply to the canonical repository; decide feature-by-feature.

## Generated application conventions

- Prefer idiomatic Go and the standard library when appropriate.
- Keep generated code understandable and independently maintainable.
- Avoid dependencies and abstractions that do not provide clear value.
- Apply secure defaults to authentication, sessions, database access, deployment and secrets.
- Where SQLite with modernc.org/sqlite is used, the preferred connection parameters are `_time_format=sqlite&_timezone=UTC`.
- If initial administrator provisioning is implemented using `ADMIN_EMAIL`, it should be idempotent and must not continually re-promote an account on every startup.

## Engineering

- Keep changes small, testable and consistent with repository conventions.
- Run formatting, tests, vet and relevant generation/build checks.
- Test generated projects as well as generator internals when templates or generation behaviour change.
- Update user-facing documentation and examples when behaviour changes.
- Do not add secrets, private customer data or local machine configuration.

## Git and upstream

- Preserve the original Goilerplate Git history and upstream relationship.
- Retain original license files, copyright notices and attribution for inherited code.
- Review upstream changes selectively; do not blindly merge features that conflict with gostart's direction.
- Use focused branches and pull requests for substantial changes.
