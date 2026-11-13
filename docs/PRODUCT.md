# Product requirements

## Product

`gostart` is a Go CLI project generator for producing self-contained Go
applications from a small set of composable templates and optional
capabilities.

The generated application is the product of the generator. It must be
ordinary, understandable Go source with its own `go.mod`, build process,
tests and deployment configuration. It must not need gostart at runtime.

## Target user

A Go developer who wants a repeatable, production-minded starting point
without inheriting a large, opinionated boilerplate or unwanted
application capabilities.

## Primary outcomes

1.  Generate a minimal Go application quickly.
2.  Add optional capabilities such as a web UI, authentication, admin,
    workspaces or billing only when selected.
3.  Preserve the generation choices in `gostart.yml`.
4.  Produce projects that can be developed and deployed independently.
5.  Make generated code understandable to both developers and coding
    agents.
6.  Validate generated combinations through automated fixture builds.

## Product principles

-   Go-first and idiomatic.
-   Prefer the standard library when it is sufficient.
-   Explicit feature composition and dependencies.
-   Small, comprehensible generated applications.
-   No runtime generator dependency.
-   No automatic template updates.
-   Security-sensitive defaults.
-   Reproducible generation.
-   Test the generated project, not only the generator.

## Initial application profile

The initial reference output is a web application using:

-   Go and `net/http`.
-   `http.ServeMux`.
-   SQLite with `modernc.org/sqlite`.
-   `database/sql`.
-   sqlc.
-   Custom forward-only SQL migrations executed at application startup.
-   Templ, Tailwind CSS 4 and shadcn-templ.
-   HTMX where server-driven interactivity is useful.
-   slog and tint.
-   Taskfile and a separate tools Go module.
-   Docker and Docker Compose.
-   GitHub Actions and GHCR publishing examples.

The generator itself uses urfave/cli v3, Huh, `embed.FS` and
`text/template`.

## V1 scope

-   Go application generation.
-   SQLite database option.
-   Base HTTP application.
-   Optional Templ/shadcn-templ frontend.
-   Optional email one-time-code authentication.
-   Optional user management and admin functionality.
-   Initial administrator bootstrap using `ADMIN_EMAIL`.
-   Feature manifests and dependency validation.
-   `gostart.yml` read/write/validation.
-   Docker, Compose and GitHub Actions templates.
-   Generator tests and generated fixture builds.

The first implementation should prioritise a small, working vertical
slice. Do not attempt to implement every optional capability before the
generator's core is proven.

## Explicit non-goals for v1

-   PostgreSQL support.
-   Automatic updates to generated projects.
-   A runtime framework or hosted service.
-   Arbitrary user-defined executable template hooks.
-   Multiple frontend frameworks.
-   ORM support.
-   Down migrations or automated rollback.
-   Mandatory public self-registration.
-   Mandatory Stripe billing.
-   A separate migration executable.
-   A requirement to run an interactive CLI command inside a deployment
    container.

## Success criteria

A developer can generate a minimal application and a selected web
application, inspect the resulting source, run its tests, build it and
deploy it without gostart being installed. Invalid feature combinations
are rejected before files are written. A saved `gostart.yml` accurately
records the generated composition.
