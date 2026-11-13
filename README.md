# gostart

**A Go project generator for small, composable, production-minded applications.**

gostart generates self-contained Go applications from composable templates and optional capabilities. Generated projects are ordinary Go source with their own `go.mod`, build, tests and deployment configuration. They do not depend on gostart at runtime.

**Status:** early development. The repository is being reshaped around a focused generator and a reference Go application. See [product requirements](docs/PRODUCT.md), [architecture](docs/ARCHITECTURE.md) and [implementation backlog](docs/BACKLOG.md).

## Product principles

- Go-first and idiomatic. Prefer the standard library when it is sufficient.
- Generate small, comprehensible applications with explicit feature composition.
- Generated projects are self-contained and independently maintainable.
- No automatic template updates.
- Secure defaults and reproducible generation.
- Test generated projects, not just the generator.

## Initial application profile

The initial reference output is a web application using Go, `net/http`, `http.ServeMux`, SQLite, `database/sql`, sqlc, Templ, Tailwind CSS 4, shadcn-templ, HTMX, slog and tint. Development uses Taskfile and a separate tools module. Deployment examples use Docker, Compose, GitHub Actions and GHCR.

The first release is deliberately focused. SQLite is the only planned database for v1; authentication is email one-time-code based when selected. See the product requirements for scope and non-goals.

## Development

The generator is a Go CLI. Consult [AGENTS.md](AGENTS.md) and [the specification-driven development workflow](docs/SDD_WORKFLOW.md) before making changes.

## Documentation

- [Product requirements](docs/PRODUCT.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Architecture decisions](docs/DECISIONS.md)
- [Implementation backlog](docs/BACKLOG.md)
- [Specification-driven development workflow](docs/SDD_WORKFLOW.md)
- [Specification template](docs/SPEC_TEMPLATE.md)
- [Implementation specifications](specs/)

## Upstream

This repository is a fork of [axadrn/goilerplate](https://github.com/axadrn/goilerplate). The upstream relationship is retained so useful fixes can be reviewed and selectively incorporated. gostart is an independent project; changes do not need to be submitted upstream.
