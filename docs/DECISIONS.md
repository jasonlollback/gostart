# Architecture decision records

Status: decisions below are agreed baseline decisions unless marked
otherwise. Changes require an explicit proposal and review.

## ADR-001 --- Generated projects are self-contained

**Decision:** Generated projects are ordinary Go applications with their
own source, `go.mod`, build and deployment configuration. They do not
depend on gostart at runtime.

**Reason:** Developers must be free to modify generated code without
depending on a generator framework.

## ADR-002 --- SQLite is the only database in v1

**Decision:** V1 generates SQLite applications using
`modernc.org/sqlite` and `database/sql`.

**Reason:** Keep the first release focused. PostgreSQL can be added
after the SQLite path is proven.

## ADR-003 --- Custom forward-only migrations

**Decision:** Implement a small internal SQL migration package. Do not
use Goose or golang-migrate. No down migrations or rollback support.

**Reason:** Keep migration behaviour explicit and aligned with the
generated application's needs.

## ADR-004 --- Migrations run at application startup

**Decision:** The generated application applies pending migrations
during initialisation, before serving requests. A separate migration
executable is not required.

**Requirements:** Checksum verification, concurrency protection,
failure-safe startup and tests for multiple simultaneous starts.

## ADR-005 --- Email one-time-code authentication

**Decision:** When authentication is selected, use email one-time codes,
database-backed sessions and secure HTTP-only cookies. No password
authentication in the default implementation.

## ADR-006 --- Initial administrator via ADMIN_EMAIL

**Decision:** `ADMIN_EMAIL` alone triggers initial administrator
provisioning when no administrator exists. No additional bootstrap flag,
CLI command or `/bootstrap` endpoint is required in v1.

**Requirements:** Idempotent, atomic provisioning; normalise email
consistently; do not continually promote or modify users once an
administrator exists. Document in `.env.example` that the variable can
be removed after initial provisioning.

## ADR-007 --- Templ and shadcn-templ

**Decision:** Use Templ, Tailwind CSS 4 and shadcn-templ as the UI
foundation. HTMX is the preferred minimal-interactivity option.

## ADR-008 --- urfave/cli v3 and Huh

**Decision:** Use urfave/cli v3 for CLI commands and flags; use Huh for
interactive terminal prompts.

## ADR-009 --- Persist gostart.yml

**Decision:** Each generated project includes a versioned `gostart.yml`
recording generation choices. It is not loaded by the generated
application at runtime.

## ADR-010 --- No automatic template updates

**Decision:** Generated projects do not automatically update when
gostart templates change. Future update tooling, if any, must be a
separately designed opt-in capability.

## ADR-011 --- Feature-oriented structure

**Decision:** Use feature-oriented packages and explicit dependencies.
Do not require `/internal/platform` and `/internal/modules`.

## ADR-012 --- Separate tools module

**Decision:** Generated projects keep development tools in a separate
`/tools` Go module to avoid polluting the application's module
dependencies.

## ADR-013 --- Taskfile

**Decision:** Use Taskfile rather than Makefiles. Taskfile loads `.env`
for development tasks; application configuration is loaded separately by
`caarlos0/env`.

## Open decisions

Do not silently decide these during implementation:

-   Exact email provider(s) included in the first auth release.
-   Exact one-time-code length, expiry, attempt limits and resend
    throttles.
-   Exact session lifetime and rotation policy.
-   Whether the first generated CLI release offers presets or only
    explicit feature selection.
-   Whether a manual migration status command is included in v1.
-   Exact dependency-resolution behaviour when a feature requires
    another feature.
-   Whether generated file overwrite is ever supported.
-   PostgreSQL design and release timing.

When one of these blocks a task, propose options with trade-offs and
wait for a decision.
