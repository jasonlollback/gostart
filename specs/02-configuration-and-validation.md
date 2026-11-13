# Specification GS-020/021 --- gostart.yml configuration

Status: proposed\
Depends on: GS-011\
Related decisions: ADR-009

## Outcome

A versioned, validated configuration model drives project generation and
is persisted into generated projects.

## Requirements

-   Define a typed Go representation for project, application, database,
    frontend, authentication, features and deployment.
-   `version` is mandatory.
-   Validate project name and Go module path.
-   Define defaults centrally and document them.
-   Reject unsupported schema versions.
-   Report field-specific errors.
-   Support YAML read/write round trips.
-   Persist the resolved configuration in generated `gostart.yml`.
-   Do not load `gostart.yml` at generated application runtime.
-   Do not silently drop unknown fields if doing so could cause a
    materially different generation result; define and test the chosen
    strictness policy.

## Initial fields

``` yaml
version: 1
project:
  name: example
  module: github.com/example/example
application:
  type: web
database:
  engine: sqlite
  driver: modernc
frontend:
  enabled: true
  framework: templ
  components: shadcn-templ
authentication:
  enabled: false
features:
  admin: false
  registration: false
  workspaces: false
  billing: false
deployment:
  docker: true
  compose: true
  github_actions: true
```

This is illustrative. The implementation should avoid fields that have
no consumer yet.

## Acceptance criteria

-   [ ] Valid configuration loads and round-trips.
-   [ ] Unsupported versions fail.
-   [ ] Invalid project/module values fail with actionable messages.
-   [ ] Defaults are deterministic.
-   [ ] Validation has no filesystem side effects.
-   [ ] Persisted config reflects resolved feature selections.
-   [ ] Tests cover malformed YAML, missing fields, unknown values and
    defaults.
