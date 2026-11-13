# Specification: GS-010 --- Minimal generated application reference

Status: done\
Owner: Jason Lollback\
Depends on: GS-001, GS-002, GS-003 (M0)\
Related decisions: ADR-001, ADR-011, ADR-013

## 1. Problem

gostart has a CLI skeleton but no definition of what it generates. The
backlog says to build a reference generated application first and
extract templates from it later (GS-011), so the template engine is
shaped by a real, working project rather than by speculation.

## 2. Outcome

The repository contains a hand-written, minimal Go application under
`testdata/reference/minimal/`. It is a complete standalone module that
formats, vets, tests and builds with no dependency on gostart. It
serves a health endpoint, logs through slog and shuts down gracefully.
It contains no database, UI or authentication code.

## 3. Scope

### Included

-   A standalone module at `testdata/reference/minimal/` using the
    placeholder module path `example.com/app`.
-   `cmd/server` entry point.
-   `internal/config`, `internal/logging`, `internal/version`,
    `internal/health` and `internal/app` packages.
-   `Taskfile.yml`, `.env.example` and `.gitignore` for the reference.
-   Unit tests for each package.
-   A generator-side check that proves the reference is independent and
    buildable (section 9).
-   A `testdata/` line in `CONTRIBUTING.md` and backlog updates once
    verified.

### Excluded

-   Template rendering, placeholders, `embed.FS`, path parameterisation
    (GS-011).
-   Output writing or collision detection (GS-012).
-   `gostart.yml` (GS-020).
-   Database, migrations, sqlc (M3). Templ, Tailwind, shadcn-templ,
    HTMX (M4). Auth, sessions, admin (M5).
-   The `/tools` module and pinned development tools (GS-070).
-   Dockerfile, Compose and CI templates (GS-071 to GS-073).
-   `/readyz` and readiness behaviour (GS-042).
-   Any change to the generator's CLI or `internal/cli`.

## 4. Behaviour and requirements

### 4.1 Startup

`cmd/server/main.go` is thin. It calls `run() error` and, on error,
writes one line to stderr and exits with status 1. `run`:

1.  Loads and validates configuration (`config.Load(env.ToMap(os.Environ()))`). On failure it
    returns the error; it does not start the server.
2.  Builds the logger (`logging.New`).
3.  Logs `starting` with attributes `build` (`version.Get()`) and
    `addr`.
4.  Creates a context cancelled on `SIGINT` or `SIGTERM`
    (`signal.NotifyContext`).
5.  Calls `app.Run(ctx, cfg, logger)`.

### 4.2 Configuration

Loaded from environment variables by `github.com/caarlos0/env/v11`.
Defaults are production-safe. A `.env` file is never read by the
application; it is loaded only by the Taskfile (ADR-013).

`config.Load(environ map[string]string)` takes the environment
explicitly so tests are hermetic (no `t.Setenv`). `Config` uses nested
structs with `envPrefix` (`Log`, `HTTP`) and `slog.Level` for the log
level, which `env` parses through `encoding.TextUnmarshaler`.

| Variable                    | Type     | Default       | Notes                                                    |
|-----------------------------|----------|---------------|----------------------------------------------------------|
| `HTTP_ADDR`                 | string   | `:8080`       | `host:port`; empty host allowed, port required           |
| `LOG_FORMAT`                | string   | `json`        | `json` or `text`; see 4.3                               |
| `LOG_LEVEL`                 | slog.Level | `info`      | any level `slog.Level` accepts (`debug`, `info`, `warn`, `error`, `info+2`) |
| `HTTP_SHUTDOWN_TIMEOUT`     | duration | `10s`         | must be greater than zero                                |
| `HTTP_READ_HEADER_TIMEOUT`  | duration | `5s`          | must be greater than zero                                |
| `HTTP_READ_TIMEOUT`         | duration | `15s`         | must be greater than zero                                |
| `HTTP_WRITE_TIMEOUT`        | duration | `30s`         | must be greater than zero                                |
| `HTTP_IDLE_TIMEOUT`         | duration | `60s`         | must be greater than zero                                |

`Config.Validate` reports every problem at once using `errors.Join`,
and its messages name the environment variable. Parse failures from
`env` (an unparsable duration or log level) are wrapped as
`parse config: ...` and name the Go field, which is accepted. Values may
appear in messages for non-secret settings only; a future secret setting
must not be echoed in any error.

### 4.3 Logging

`logging.New(cfg, w io.Writer) *slog.Logger`:

-   `LOG_FORMAT=json` (the default): always `slog.NewJSONHandler`,
    whether or not a console is present.
-   `LOG_FORMAT=text`: if `w` is a terminal, the
    `github.com/lmittmann/tint` handler (colour disabled when
    `NO_COLOR` is set); otherwise `slog.NewTextHandler`.
-   There is no `APP_ENVIRONMENT` variable; log selection depends only
    on `LOG_FORMAT` and the terminal check (D2). Local development
    selects tint by setting `LOG_FORMAT=text` in `.env`.
-   Terminal detection uses `github.com/mattn/go-isatty`
    (`IsTerminal` or `IsCygwinTerminal` on the file descriptor). A
    writer that is not an `*os.File` is not a terminal. Detection is
    isolated in one small function so tests can inject the result
    (`newLogger(cfg, w, terminal bool)` behind `New`).
-   The level comes from `LOG_LEVEL`.
-   The logger is passed explicitly. `slog.SetDefault` is not called.

### 4.4 Health endpoint

`GET /healthz` returns `200 OK`, `Content-Type: application/json`, body:

``` json
{"status":"ok"}
```

No build or version information is served. It is reported only in the
startup log (decision D1). `HEAD` is also accepted. Other methods return `405` with an `Allow`
header (Go 1.22+ `ServeMux` method patterns provide this). Unknown
paths return the mux default `404`.

### 4.5 Server and shutdown

`app.Run(ctx, cfg, logger) error`:

-   Listens with `net.Listen("tcp", cfg.HTTPAddr)` so listen errors are
    returned before serving, then serves with an `http.Server` using
    the configured timeouts.
-   Builds routes with `http.NewServeMux` in a function that takes its
    dependencies explicitly (no globals).
-   When `ctx` is cancelled, calls `Server.Shutdown` with a context
    bounded by `HTTP_SHUTDOWN_TIMEOUT`. A clean shutdown returns nil.
    If the shutdown deadline is exceeded, it closes the server and
    returns the error.
-   `http.ErrServerClosed` is not treated as an error.
-   A serve error other than `ErrServerClosed` is returned.
-   Logs `listening` (with the bound address) and `shutdown complete`.
    Request bodies, headers and query strings are not logged.

An `app.Run` variant for tests accepts a pre-made `net.Listener` so
tests can use port 0 and learn the address without sleeping.

### 4.6 Version

`internal/version` is copied from the generator's `internal/version`
(`Info`, `Get`, `LogValue`, with the `unavailable` convention and tests
for tagged, dirty, `(devel)` and unavailable build info). The version
comes only from `runtime/debug.ReadBuildInfo`. `-ldflags` is not used in
the Taskfile or any command (project rule).

## 5. Interfaces and data

``` text
testdata/reference/minimal/
  go.mod                      module example.com/app; go 1.27.0
  go.sum
  Taskfile.yml                fmt, fmt:check, vet, test, build, run, check
  .env.example                documented variables from 4.2
  .gitignore                  /bin/, .env, test output
  cmd/server/main.go
  internal/config/config.go
  internal/config/config_test.go
  internal/logging/logging.go
  internal/logging/logging_test.go
  internal/version/version.go
  internal/version/version_test.go
  internal/health/health.go
  internal/health/health_test.go
  internal/app/app.go
  internal/app/app_test.go
```

Generator-side additions:

-   `internal/reference/reference_test.go`: the independence check
    described in section 9.
-   `Taskfile.yml` (generator): task `reference:check` runs the
    verification commands in section 10; `check` is extended to depend on
    it only if it runs in a reasonable time (see D4).

Dependencies of the reference module (the generator's own `go.mod` is
unchanged):

| Module                         | Why                                  | Stdlib alternative considered            |
|--------------------------------|--------------------------------------|------------------------------------------|
| `github.com/caarlos0/env/v11`  | Agreed env parsing (Architecture 3)  | `os.Getenv` with hand-written parsing    |
| `github.com/lmittmann/tint`    | Agreed console log handler           | `slog.TextHandler` (no colour/formatting)|
| `github.com/mattn/go-isatty`   | Approved terminal detection          | `os.File.Stat` with `ModeCharDevice`, which also matches `/dev/null` and is not a reliable terminal test |

## 6. Architecture constraints

-   ADR-001: the module imports nothing from gostart and has its own
    `go.mod`.
-   ADR-011: feature-oriented packages; no `platform` or `modules`.
-   ADR-013: Taskfile loads `.env` (`dotenv: [".env"]`); the application
    does not.
-   No interfaces unless a test needs a seam. No ORM, no UI framework.
-   No `-ldflags`; version from build info only.
-   The reference must not be wired into the generator's own module:
    `go test ./...` from the repository root skips `testdata`.
-   Package names and the `example.com/app` placeholder are provisional
    until GS-011 decides how they are parameterised; tests must not
    depend on the literal module path.

## 7. Security and reliability

-   Threat assumptions: the app is reachable by untrusted HTTP clients;
    it holds no secrets or user data in this slice.
-   All server timeouts are set, including `ReadHeaderTimeout`, to
    bound slow-client resource use.
-   Error responses do not include internal details.
-   The logger never logs full configuration or the environment.
-   `/healthz` exposes no version or build information (D1).
-   Shutdown does not drop in-flight requests inside the timeout.
-   Signal handling is stopped after the first signal (via
    `NotifyContext`'s stop) so a second signal terminates the process
    if shutdown hangs.

## 8. Acceptance criteria

From the backlog:

-   [x] Output has its own `go.mod` and no import of gostart.
-   [x] Contains a `cmd/server` entry point.
-   [x] Uses `net/http` and `http.ServeMux`.
-   [x] Has structured logging (slog; JSON, or tint on a terminal with `LOG_FORMAT=text`) and graceful
    shutdown.
-   [x] Has a basic health endpoint.
-   [x] Builds and tests without gostart installed.
-   [x] Contains no database, UI or auth code.

Added by this specification:

-   [x] Given no environment variables, when `config.Load` runs, then
    the defaults in 4.2 are returned.
-   [x] Given several invalid values, when `Validate` runs, then one
    error names every offending variable.
-   [x] Given `LOG_FORMAT=json`, then JSON is produced whether or not
    the writer is a terminal.
-   [x] Given `LOG_FORMAT=text` and a terminal writer, then tint is used.
-   [x] Given `LOG_FORMAT=text` and a non-terminal writer, then
    `slog.TextHandler` output is produced and never tint.
-   [x] Given `NO_COLOR` set, then tint output has no colour codes.
-   [x] Given `LOG_LEVEL`, then records below the level are dropped.
-   [x] Given a running server, when `GET /healthz`, then 200 with the
    JSON body `{"status":"ok"}` in 4.4 with no version field; `POST`
    returns 405.
-   [x] Given a running server, when the context is cancelled, then
    `Run` returns nil within the shutdown timeout and the port is
    released.
-   [x] Given the address is already in use, then `Run` returns an
    error without serving.
-   [x] `go.mod` of the reference lists only `env/v11`, `tint`, `go-isatty` and
    their transitive requirements.
-   [x] The generator's own `go.mod` and `go.sum` are unchanged.
-   [x] `.env.example` documents every variable in 4.2.

## 9. Testing plan

Reference unit tests (run inside the reference module):

-   `config`: defaults; overrides; each invalid value; multiple errors
    joined; duration parsing. Passes an explicit environment map.
-   `logging`: handler selection per environment, terminal state and
    format (terminal state injected); `NO_COLOR`; level filtering; output captured in a buffer.
-   `version`: tagged, dirty, `(devel)`, unavailable.
-   `health`: 200 body with no version field, `Content-Type`, `HEAD`, 405 with `Allow`.
-   `app`: routes via `httptest`; `Run` with a pre-made listener on
    `127.0.0.1:0`: request succeeds, cancel, returns nil within a
    deadline; address-in-use error. No `time.Sleep`; synchronise with
    channels and contexts.

Generator-side independence test (`internal/reference/reference_test.go`):

-   Skipped under `testing.Short()` and when `go` is not on `PATH`.
-   Copies `testdata/reference/minimal` to `t.TempDir()`, runs with
    `GOWORK=off` and `GOFLAGS=-mod=mod` unset: `gofmt -l`, `go vet ./...`,
    `go test ./...`, `go build ./cmd/server`. Fails on any output or
    non-zero exit.
-   Walks the copy and fails if any file imports
    `github.com/jasonlollback/gostart`.
-   Asserts the copy contains no `database/sql`, `templ` or auth
    imports (the "no unselected features" criterion).

## 10. Verification commands

``` bash
# Generator baseline (unchanged behaviour)
task check

# Reference application, run in place
cd testdata/reference/minimal
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
go build -o /tmp/gs010-app ./cmd/server

# Smoke test
HTTP_ADDR=127.0.0.1:18080 /tmp/gs010-app &
sleep 1; curl -i http://127.0.0.1:18080/healthz
kill -TERM %1; wait

# Independence test and lint scope
cd ../../..
go test ./internal/reference/...
task lint staticcheck vuln
git diff --stat -- go.mod go.sum        # expect no output
```

Expected: every command exits 0; `/healthz` returns 200 with the body in
4.4; the process exits 0 after `SIGTERM` and logs `shutdown complete`;
`go.mod` and `go.sum` of the generator show no diff.

## 11. Implementation plan

1.  Create the reference module, `config` and `logging` with tests.
2.  Copy `version` from the generator with its tests.
3.  Add `health` and `app` with tests.
4.  Add `cmd/server`, `.env.example`, `.gitignore`, `Taskfile.yml`.
5.  Add `internal/reference` independence test and the generator task.
6.  Run section 10, review the diff, update `CONTRIBUTING.md` and the
    backlog, and report using the completion template.

## 12. Risks and open questions

**D1 --- Version on `/healthz` (resolved).** Removed entirely. The
endpoint returns `{"status":"ok"}`. Build details appear only in the
startup log. Reason: an unauthenticated endpoint reachable by load
balancers and monitors would let an attacker fingerprint the release.

**D2 --- Logging mode selection (resolved).** Simplified: `LOG_FORMAT=json`
always yields JSON; `LOG_FORMAT=text` yields tint on a terminal and
plain text otherwise. `APP_ENVIRONMENT` is dropped from this slice
because nothing else uses it; it can be added by the feature that needs
it. This supersedes the earlier "tint when `APP_ENVIRONMENT` is
development" instruction. The default stays `json` so an unconfigured
deployment logs machine-readable output.

**D3 --- `LOG_LEVEL` and timeout variables (resolved).** Kept.

**D4 --- Independence test (resolved).** Part of `go test ./...`,
skipped with `-short` or when `go` is not on `PATH`. If it proves too
slow for `task check` or CI, move it to a separate task in a later
change.

**D7 --- Terminal detection (resolved).** `github.com/mattn/go-isatty`
is approved. `golang.org/x/term` is not used.

**D5 --- File name.** The spec is numbered `09` to avoid renumbering
existing specs; it could be renamed.

**D6 --- Placeholder module path.** `example.com/app` is provisional
and is replaced when GS-011 defines template data.

Other risks: first-run network access to fetch `env`, `tint` and `go-isatty`;
shutdown tests becoming flaky (mitigated by listener injection and
contexts); terminal detection can differ on Windows and in CI, so tests
inject the terminal result and only one test touches real detection;
`go-isatty` and its `golang.org/x/sys` requirement add to the
reference's dependency graph.

## 13. Completion evidence

Implemented on branch `gs-010-minimal-reference-app` and merged to `main` via pull request #3.

Files: `testdata/reference/minimal/**` (new module),
`internal/reference/reference_test.go`, `Taskfile.yml`
(`reference:check`), `CONTRIBUTING.md`, this spec.

Commands run (all observed to pass):

-   `task reference:check` (gofmt, vet, `go test -race`, build): PASS.
-   `task check`, `task lint`, `task staticcheck`, `task vuln`: PASS.
-   `go test ./internal/reference/`: PASS (imports and standalone
    build).
-   govulncheck on the reference module: no vulnerabilities after
    raising indirect `golang.org/x/sys` from v0.28.0 to v0.44.0
    (GO-2026-5024, Windows only, not called).
-   Smoke test of the built binary: `GET /healthz` 200
    `{"status":"ok"}`; `POST` 405 with `Allow: GET, HEAD`; SIGTERM
    exits 0 logging `shutdown complete`; `HTTP_ADDR=bad` exits 1 with a
    useful message; `git diff go.mod go.sum` of the generator empty.

Deviations and notes:

-   `config` follows the project's standard pattern (explicit
    environment map, nested `envPrefix` structs, `slog.Level`), adopted
    after review. An earlier draft rewrote env's parse errors to name
    the variable; that code was removed.
-   `golang.org/x/sys` is pinned above go-isatty's minimum (see above).
-   The real tint-on-a-terminal path is covered by tests with an injected
    terminal flag; it was not exercised in an interactive terminal.
-   Backlog entry GS-010 is marked done after review.
