# Architecture

## 1. Generator architecture

The generator is a small Go CLI with these responsibilities:

-   CLI parsing and interactive prompts.
-   Configuration loading and validation.
-   Feature dependency resolution.
-   Template selection and rendering.
-   Safe output creation.
-   Reporting generated files and next steps.

Suggested package boundaries:

``` text
cmd/gostart/
internal/
  cli/
  config/
  features/
  generator/
  manifest/
  templates/
  validation/
templates/
testdata/
```

Do not create packages merely to mirror this diagram if a smaller
implementation is clearer. Keep responsibilities separated without
excessive interfaces.

## 2. Template engine

-   Use `embed.FS` to embed template source.
-   Use `text/template` for controlled substitutions.
-   Store templates as normal repository files.
-   Use explicit template data structures.
-   Do not execute arbitrary hooks or scripts from manifests.
-   Detect output path collisions.
-   Refuse to overwrite existing files unless an explicit, reviewed
    overwrite mode is designed.
-   Escape or validate project/module identifiers before using them in
    paths or generated source.
-   Keep template rendering deterministic.

## 3. Configuration

`gostart.yml` is the persistent generation manifest in each generated
project.

It is not runtime application configuration.

Runtime settings are supplied through environment variables, documented
in `.env.example`, and loaded by the generated application with
`github.com/caarlos0/env/v11`.

Configuration must have a schema version. Unknown or unsupported
versions should produce a useful error. Defaults should be applied
consistently and documented.

## 4. Feature model

Each feature has metadata describing:

-   Stable feature ID.
-   Human-readable name and description.
-   Required features.
-   Conflicting features.
-   Templates and generated assets.
-   Configuration fields, where applicable.

Feature resolution happens before writing output. Missing dependencies
should either be automatically enabled with an explicit explanation or
rejected when user intent is ambiguous. Conflicts must be reported
before generation.

Do not build a general-purpose plugin runtime. Feature manifests are
declarative data.

## 5. Generated application architecture

Use a feature-oriented package layout. Do not require
`/internal/platform` and `/internal/modules`.

Typical packages may include:

``` text
internal/
  app/
  config/
  database/
  migrate/
  health/
  auth/
  users/
  admin/
```

Only generate packages relevant to selected capabilities.

Boundaries:

-   HTTP handlers do not execute SQL directly.
-   Repositories own database access.
-   Services own application operations and business rules where that
    separation adds value.
-   Features must not depend on other features' HTTP handlers.
-   Cross-feature dependencies are explicit.
-   Do not create an interface for every concrete type by default.

## 6. Database

V1 supports SQLite only, using `modernc.org/sqlite` with `database/sql`.

The connection string includes:

``` text
_time_format=sqlite&_timezone=UTC
```

Use WAL, foreign-key enforcement, STRICT tables, TEXT timestamps, UUIDv7
identifiers and bounded connection settings. Test timestamp round trips
and SQLite locking behaviour.

Use sqlc for typed query generation. Keep the agreed sqlc defaults in a
version-controlled template.

## 7. Migrations

Generated applications contain an internal migration package. Migrations
are embedded SQL files and run during application startup before the
HTTP server becomes ready.

Requirements:

-   Forward-only, sequential migrations.
-   No down migrations.
-   No rollback command.
-   Checksum verification.
-   Transactional execution where supported.
-   Protection against concurrent application instances.
-   Fail startup if migration execution fails.
-   Applied migration files are immutable; corrections use a new
    migration.

The migration runner should accept explicit dependencies such as
`*sql.DB`, `fs.FS` and a dialect identifier. Avoid global state.

## 8. Authentication and bootstrap

Authentication, when selected, uses email one-time codes and
database-backed sessions.

Initial administrator provisioning is triggered by `ADMIN_EMAIL`. If no
administrator exists, the application creates or locates the configured
user and grants the initial administrator role. Once an administrator
exists, the variable is not a continuing role-synchronisation mechanism.

The operation must be atomic and safe during concurrent startup.

## 9. Frontend

Use Templ, Tailwind CSS 4 and shadcn-templ. HTMX may be used for
server-driven interactions. Avoid adding a JavaScript framework.

## 10. Development and deployment

Generated projects use:

-   Taskfile, not Makefile.
-   `.env` and `.env.example`.
-   A separate `/tools` Go module for development tools.
-   Docker multi-stage builds.
-   Docker Compose with persistent SQLite storage.
-   Separate CI and Docker publishing workflows.
-   GHCR for published images.

Do not force a reverse proxy into the application's Compose stack.
Provide an optional Caddy example.

## 11. Generator testing

Test both the generator and generated outputs.

Use golden fixtures for at least:

-   Minimal application.
-   Web application.
-   Admin application.
-   SaaS-style composition.

Generated fixtures should be formatted, generated where relevant,
compiled, tested and built in CI.
