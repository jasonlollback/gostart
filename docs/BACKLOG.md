# Implementation backlog

This backlog is deliberately ordered around proving the smallest useful
generator before adding optional capabilities.

## Delivery strategy

Build a reference generated application first, then extract its stable
structure into templates. Avoid building a sophisticated feature
registry before a real generated project exists.

Use one branch/PR per vertical slice where practical. Keep the assigned
specification attached to each task.

## Milestones

### M0 --- Repository and engineering foundation

**Outcome:** A clean generator repository with documented decisions,
reproducible Go tooling and a passing baseline CI workflow.

-   `GS-001` Repository foundation and conventions --- `ported in gostart_v1; integration pending`
-   `GS-002` Generator CLI skeleton --- `ported in gostart_v1; integration pending`
-   `GS-003` CI for generator repository --- `ported in gostart_v1; integration pending`

Exit criteria: generator builds, tests run, CLI help works, CI is green.

### M1 --- Reference application and template engine

**Outcome:** Generate a minimal, compilable Go application from embedded
templates.

-   `GS-010` Minimal generated application reference --- `ported in gostart_v1; integration pending`
-   `GS-011` Embedded template renderer --- `ported in gostart_v1; integration pending`
-   `GS-012` Safe output writer and collision detection --- `ported in gostart_v1; integration pending`
-   `GS-013` Golden fixture tests --- `ported in gostart_v1; integration pending`

Exit criteria: minimal output is deterministic, formatted, tested and
buildable.

### M2 --- Configuration and feature composition

**Outcome:** Saved configuration describes a validated feature
composition.

-   `GS-020` Versioned gostart.yml schema --- depends on M1
-   `GS-021` Config load, defaults and validation --- depends on GS-020
-   `GS-022` Feature manifests --- depends on GS-020
-   `GS-023` Dependency and conflict resolver --- depends on GS-022
-   `GS-024` Config-driven generation --- depends on GS-021 and GS-023

Exit criteria: invalid configurations fail before output is written;
valid configurations generate reproducibly.

### M3 --- SQLite and migration foundation

**Outcome:** Generated application opens SQLite and safely migrates on
startup.

-   `GS-030` SQLite database package --- depends on M1
-   `GS-031` Migration metadata and discovery --- depends on GS-030
-   `GS-032` Migration execution and checksums --- depends on GS-031
-   `GS-033` Concurrent startup and failure tests --- depends on GS-032
-   `GS-034` sqlc integration --- depends on GS-030

Exit criteria: fresh and existing databases initialise correctly;
checksum mismatch and migration failure stop startup; concurrent starts
are safe.

### M4 --- Web foundation

**Outcome:** A generated server-rendered web application using the
chosen UI stack.

-   `GS-040` Templ integration --- depends on M2
-   `GS-041` Tailwind CSS 4 and shadcn-templ integration --- depends on
    GS-040
-   `GS-042` Layout, navigation and health pages --- depends on GS-041
-   `GS-043` HTMX conventions and example interaction --- depends on
    GS-042

Exit criteria: generated web fixture builds and serves a usable page
with documented asset generation.

### M5 --- Authentication and administration

**Outcome:** Optional email-code auth and secure initial administrator
provisioning.

-   `GS-050` Email delivery boundary --- depends on M3
-   `GS-051` One-time-code lifecycle --- depends on GS-050
-   `GS-052` Database-backed sessions --- depends on GS-051
-   `GS-053` ADMIN_EMAIL bootstrap --- depends on GS-052
-   `GS-054` User management --- depends on GS-053
-   `GS-055` Admin UI --- depends on M4 and GS-054

Exit criteria: auth abuse cases are tested; initial administrator is
provisioned once; no public registration is generated unless selected.

### M6 --- CLI and project experience

**Outcome:** A polished interactive and non-interactive project
generation workflow.

-   `GS-060` `gostart new` interactive flow --- depends on M2
-   `GS-061` CLI flags and config input --- depends on GS-060
-   `GS-062` `list`, `info`, `validate`, `version` --- depends on GS-022
    and GS-060
-   `GS-063` Presets, if approved --- depends on GS-060
-   `GS-064` Generated README and next steps --- depends on M4 and M6

Exit criteria: a new developer can generate a project without editing
templates or reading generator internals.

### M7 --- Deployment and release

**Outcome:** Generated projects can be tested and deployed reproducibly.

-   `GS-070` Taskfile and tools module templates --- depends on M1
-   `GS-071` Dockerfile and Compose templates --- depends on M3
-   `GS-072` CI workflow template --- depends on M4
-   `GS-073` GHCR publish workflow --- depends on GS-072
-   `GS-074` End-to-end fixture matrix --- depends on M3, M4, M5, M7

Exit criteria: all supported fixtures pass generation, code generation,
tests and build in CI.

### M8 --- Additional capabilities

Do not start until the core generation path is stable.

-   RBAC.
-   Workspaces.
-   Audit events.
-   File storage.
-   Background jobs.
-   Stripe billing.
-   PostgreSQL support.

Each capability requires its own specification, dependency declaration,
generated fixture and migration tests.

## Detailed task definitions

### GS-001 --- Repository foundation and conventions

**Goal:** Establish the gostart repository and engineering conventions.

**Scope:** - Go module and initial `cmd/gostart` entry point. - Root
`AGENTS.md` and `docs/` baseline. - Taskfile with format, test, build
and lint tasks. - `.gitignore`, README and contribution/development
notes. - Dependency policy and tool version strategy.

**Acceptance criteria:** - \[x\] `go test ./...` succeeds. - \[x\]
`go build ./cmd/gostart` succeeds. - \[x\] Taskfile exposes documented
baseline tasks. - \[x\] No application-template dependencies are added
to the generator without need. - \[x\] Architecture decisions are
recorded and discoverable.

### GS-002 --- Generator CLI skeleton

**Goal:** Establish urfave/cli v3 command structure.

**Acceptance criteria:** - \[x\] `gostart --help` works. - \[x\]
`gostart version` reports version information. - \[x\] Unknown commands
and invalid flags return useful errors. - \[x\] CLI logic is separated
from generation logic. - \[x\] CLI tests cover command parsing and error
cases.

### GS-003 --- Generator CI

**Acceptance criteria:** - \[x\] GitHub Actions runs formatting
verification, tests and build. - \[x\] Toolchain version is explicit. -
\[x\] Workflow does not publish images on ordinary commits. - \[x\]
Failure is visible as a required CI check.

Note: the repository is private on a plan without branch protection, so
the `check` job is visible on every pull request but is not enforced as
a required status check. Enforcement is deferred; revisit if the
repository becomes public or the plan changes.

### GS-010 --- Minimal generated application reference

**Goal:** Create the smallest useful generated Go app before
generalising it.

**Acceptance criteria:** - \[x\] Output has its own `go.mod`. - \[x\]
Contains a `cmd/server` entry point. - \[x\] Uses `net/http` and
`http.ServeMux`. - \[x\] Has structured logging and graceful shutdown. -
\[x\] Has a basic health endpoint. - \[x\] Builds and tests without
gostart installed. - \[x\] Contains no database, UI or auth code unless
selected.

### GS-011 --- Embedded template renderer

**Acceptance criteria:** - \[x\] Template source is embedded using
`embed.FS`. - \[x\] Rendering is deterministic for identical input. -
\[x\] Template errors identify the source template. - \[x\] Template data
is explicit and typed. - \[x\] Unit tests cover success and failure
cases.

### GS-012 --- Safe output writer

**Acceptance criteria:** - \[x\] Rejects path traversal and unsafe
output paths. - \[x\] Detects duplicate output paths. - \[x\] Does not
silently overwrite existing user files. - \[x\] Reports created files
clearly. - \[x\] Handles filesystem errors without leaving misleading
success output.

### GS-013 --- Golden fixture tests

**Acceptance criteria:** - \[x\] Minimal fixture is generated in a
temporary directory. - \[x\] Expected files are checked. - \[x\]
Unexpected optional-feature files are absent. - \[x\] Generated Go is
formatted and buildable. - \[x\] Fixture generation is repeatable.

### GS-020 --- Versioned gostart.yml schema

**Acceptance criteria:** - \[ \] Schema has a required version field. -
\[ \] Project/module identifiers are validated. - \[ \] Feature settings
have typed Go representations. - \[ \] Defaults are documented and
deterministic. - \[ \] Unknown schema versions fail with an actionable
error. - \[ \] Round-trip load/write tests exist.

### GS-021 --- Configuration validation

**Acceptance criteria:** - \[ \] Missing required values produce
field-specific errors. - \[ \] Invalid module paths and project names
are rejected. - \[ \] Defaults are applied in one documented place. - \[
\] Validation does not write project files. - \[ \] Tests cover
malformed YAML and invalid combinations.

### GS-022 --- Feature manifests

**Acceptance criteria:** - \[ \] Each feature has a stable ID. - \[ \]
Manifests declare requirements and conflicts. - \[ \] Unknown feature
IDs are rejected. - \[ \] Manifest parsing is deterministic. - \[ \] No
executable hooks are supported.

### GS-023 --- Feature resolver

**Acceptance criteria:** - \[ \] Required dependencies are resolved
before generation. - \[ \] Conflicts are reported before output is
written. - \[ \] Cyclic dependencies are detected. - \[ \] Auto-enabled
features are reported to the user. - \[ \] Ambiguous choices are not
silently guessed. - \[ \] Unit tests cover dependency graphs and cycles.

### GS-024 --- Config-driven generation

**Acceptance criteria:** - \[ \] Valid config produces the selected
files. - \[ \] Disabled features produce no feature-specific files or
dependencies. - \[ \] `gostart.yml` is included in output. - \[ \]
Repeating generation into a clean directory yields identical output. -
\[ \] Invalid configuration fails before filesystem mutation.

### GS-030 --- SQLite package

**Acceptance criteria:** - \[ \] Uses `database/sql` and
`modernc.org/sqlite`. - \[ \] Uses
`_time_format=sqlite&_timezone=UTC`. - \[ \] Enables and verifies
foreign-key enforcement. - \[ \] Configures WAL and a bounded connection
pool. - \[ \] Database path and permissions are documented. - \[ \]
Tests verify timestamp and UUID round trips. - \[ \] Tests exercise
expected concurrent access behaviour.

### GS-031 --- Migration discovery and metadata

**Acceptance criteria:** - \[ \] Migration files are ordered by
sequential version. - \[ \] Metadata table records version, name,
checksum and applied time. - \[ \] Duplicate versions and invalid names
fail clearly. - \[ \] Migration source is supplied through `fs.FS`. - \[
\] Discovery is deterministic.

### GS-032 --- Migration execution

**Acceptance criteria:** - \[ \] Pending migrations execute in order. -
\[ \] Successful migrations are recorded. - \[ \] Applied migration
checksum mismatch fails. - \[ \] Failed migration does not get recorded
as successful. - \[ \] No down migration mechanism is present. - \[ \]
Transaction semantics are documented and tested. - \[ \] Startup does
not proceed after a migration error.

### GS-033 --- Concurrent startup tests

**Acceptance criteria:** - \[ \] Two or more runners against the same
SQLite database do not corrupt migration state. - \[ \] Only one
successful application of each migration is recorded. - \[ \] Lock
acquisition has bounded waiting/retry behaviour. - \[ \] A failed or
interrupted migration leaves a diagnosable state. - \[ \] Test is
repeatable and does not rely on timing sleeps alone.

### GS-034 --- sqlc integration

**Acceptance criteria:** - \[ \] Generated `sqlc.yml` contains agreed
defaults. - \[ \] SQLite type overrides are tested. - \[ \] sqlc
generation runs from Taskfile. - \[ \] Generated code compiles. - \[ \]
SQL query files are organised by feature or another documented
convention.

### GS-040 --- Templ integration

**Acceptance criteria:** - \[ \] Templ source is included only when UI
is selected. - \[ \] Templ generation is available through Taskfile. -
\[ \] Generated Go compiles. - \[ \] A minimal page and layout are
included.

### GS-041 --- Tailwind and shadcn-templ

**Acceptance criteria:** - \[ \] Tailwind CSS 4 build works from a clean
checkout. - \[ \] shadcn-templ setup and component generation are
documented. - \[ \] Tool versions are pinned or otherwise
reproducible. - \[ \] No unrelated UI framework is added. - \[ \] CSS
output is included in the generated build/deployment path.

### GS-042 --- Web layout and health pages

**Acceptance criteria:** - \[ \] Responsive application layout exists. -
\[ \] Accessible navigation and page structure exist. - \[ \]
Health/readiness behaviour is documented. - \[ \] Server errors do not
leak internal details.

### GS-043 --- HTMX conventions

**Acceptance criteria:** - \[ \] Example interaction is server-driven. -
\[ \] Progressive enhancement behaviour is documented where
applicable. - \[ \] No frontend framework is introduced. - \[ \] CSRF
considerations are addressed for state-changing requests.

### GS-050 --- Email delivery boundary

**Acceptance criteria:** - \[ \] Auth logic depends on a small
email-sender interface. - \[ \] Provider configuration is
environment-based. - \[ \] Delivery errors are handled without exposing
codes in logs. - \[ \] A test/fake sender is available. - \[ \] Provider
choice is an explicit open decision before production implementation.

### GS-051 --- One-time-code lifecycle

**Acceptance criteria:** - \[ \] Code is generated using a
cryptographically secure source. - \[ \] Only a hash of the code is
stored. - \[ \] Expiry, attempt limits and resend throttles are
enforced. - \[ \] Codes are single-use. - \[ \] Codes and hashes are
never logged. - \[ \] Tests cover expiry, reuse, brute-force and
rate-limit cases.

### GS-052 --- Database-backed sessions

**Acceptance criteria:** - \[ \] Session tokens are generated
securely. - \[ \] Only a token hash is stored where practical. - \[ \]
Cookies are HTTP-only and configured Secure/SameSite appropriately. - \[
\] Session expiry and revocation work. - \[ \] Logout invalidates the
session. - \[ \] Tests cover expiry, revocation and cookie attributes.

### GS-053 --- ADMIN_EMAIL bootstrap

**Acceptance criteria:** - \[ \] `.env.example` explains initial
administrator provisioning and removal after first use. - \[ \] If no
administrator exists and `ADMIN_EMAIL` is set, the matching user is
created or located and granted initial admin privileges. - \[ \] If an
administrator already exists, startup does not promote or modify
accounts based on the variable. - \[ \] Email normalisation is
consistent with authentication. - \[ \] Provisioning is atomic and safe
under concurrent startup. - \[ \] No password or one-time code is
generated as part of bootstrap. - \[ \] Missing `ADMIN_EMAIL` with no
admin produces a warning but does not prevent startup. - \[ \] Tests
cover all cases.

### GS-054 --- User management

**Acceptance criteria:** - \[ \] Admin can create, view, disable and
re-enable users as specified. - \[ \] Public registration is absent
unless selected. - \[ \] User status and session revocation behaviour
are defined. - \[ \] Sensitive actions are authorised and tested.

### GS-055 --- Admin UI

**Acceptance criteria:** - \[ \] Admin routes require appropriate
authorisation. - \[ \] Non-admin users cannot access admin pages or
actions. - \[ \] UI uses the selected shadcn-templ component
foundation. - \[ \] Empty, loading and error states are handled.

### GS-060 --- Interactive `gostart new`

**Acceptance criteria:** - \[ \] Huh prompts gather required project and
feature choices. - \[ \] Defaults are visible and sensible. - \[ \] User
can review selections before writing. - \[ \] Cancellation leaves no
partial project. - \[ \] Resulting `gostart.yml` matches selected
options.

### GS-061 --- Non-interactive generation

**Acceptance criteria:** - \[ \] CLI flags can specify project and
feature choices. - \[ \] A config file can be supplied. - \[ \]
Non-interactive mode never waits for input. - \[ \] Missing required
values produce errors. - \[ \] Exit codes are suitable for CI use.

### GS-062 --- Discovery and validation commands

**Acceptance criteria:** - \[ \] `list` displays available features. -
\[ \] `info` describes a feature and dependencies. - \[ \] `validate`
validates a configuration without writing output. - \[ \] `version`
reports build version information.

### GS-063 --- Presets (optional)

**Acceptance criteria:** - \[ \] Presets are named shortcuts over
ordinary feature configuration. - \[ \] Users can override preset
choices. - \[ \] Presets do not create a separate template fork. - \[ \]
Include only if real usage shows they improve the experience.

### GS-064 --- Generated README

**Acceptance criteria:** - \[ \] Describes the generated application and
selected capabilities. - \[ \] Includes local setup and Taskfile
commands. - \[ \] Documents required environment variables. - \[ \]
Documents first administrator provisioning. - \[ \] Documents Docker
deployment and persistent data. - \[ \] Does not describe disabled
features.

### GS-070 --- Taskfile and tools module

**Acceptance criteria:** - \[ \] Taskfile loads `.env` for development
tasks. - \[ \] Application still loads its own environment config
through caarlos0/env. - \[ \] Tools are isolated in `/tools/go.mod`. -
\[ \] Tool versions are pinned. - \[ \] Tasks cover format, test, build,
lint, generation and frontend assets as applicable.

### GS-071 --- Docker and Compose

**Acceptance criteria:** - \[ \] Multi-stage image build. - \[ \]
Non-root runtime. - \[ \] Persistent SQLite volume. - \[ \] Correct data
directory ownership. - \[ \] Environment file documented. - \[ \] Health
check and graceful shutdown work. - \[ \] No reverse proxy is forced
into the default Compose stack.

### GS-072 --- CI template

**Acceptance criteria:** - \[ \] Generated project CI runs format
verification, generation, tests, vet and build. - \[ \] Optional checks
are included when corresponding tools/features are enabled. - \[ \]
Workflow works on a clean GitHub checkout.

### GS-073 --- GHCR publishing template

**Acceptance criteria:** - \[ \] Tag-triggered publishing. - \[ \] Uses
GitHub Actions permissions minimally. - \[ \] Publishes version tags and
`latest` according to documented rules. - \[ \] Does not publish secrets
into build output.

### GS-074 --- End-to-end fixture matrix

**Acceptance criteria:** - \[ \] Minimal, web, admin and SaaS fixtures
are generated. - \[ \] Each fixture is independently built and tested. -
\[ \] Disabled features are checked for absence. - \[ \] Database
migrations are tested against real temporary SQLite files. - \[ \] CI
runtime remains practical; expensive checks can be separated into jobs.

## Backlog operating rule

Do not mark a task `done` until its acceptance criteria are met,
verification evidence is recorded and the diff has been reviewed. If an
implementation reveals an architecture issue, stop and update the
decision/specification before continuing.
