# Specification GS-001 --- Repository foundation

Status: done\
Depends on: none\
Related decisions: ADR-008, ADR-011, ADR-012, ADR-013

## 1. Problem

The project needs a minimal, reproducible repository foundation before
generator behaviour is implemented.

## 2. Outcome

A new Go repository builds, tests and exposes a minimal urfave/cli v3
command. Repository instructions and development commands are committed.

## 3. Scope

Included: - Go module. - `cmd/gostart/main.go`. - Minimal CLI root
command and version command. - Taskfile. - `.gitignore`. - README. -
`AGENTS.md` and initial docs. - Baseline unit test.

Excluded: - Template rendering. - Interactive Huh prompts. - Feature
registry. - Project generation. - Generated application dependencies.

## 4. Requirements

-   Use the current stable urfave/cli v3 API, not v2 examples.
-   CLI execution must pass the process arguments correctly.
-   Keep command construction separate from business logic where useful.
-   Version output must come from the build information the Go
    toolchain embeds in the binary (`runtime/debug.ReadBuildInfo`): the
    main module version (git tag or pseudo-version), `vcs.revision` and
    `vcs.modified`. Do not inject versions with `-ldflags`. Values that
    are not available are reported as `unavailable`.
-   Do not add Cobra or Viper.
-   Avoid bringing Templ, sqlc or shadcn-templ into the generator's own
    module unless a later, justified requirement appears.

## 5. Acceptance criteria

-   [x] `go test ./...` succeeds.
-   [x] `go build ./cmd/gostart` succeeds.
-   [x] `gostart --help` displays the root command and version command.
-   [x] `gostart version` displays version information.
-   [x] Invalid commands return a non-zero exit status and useful error.
-   [x] Taskfile provides format-check, test and build tasks.
-   [x] Root agent instructions point to the product, architecture,
    decision and SDD documents.
-   [x] Dependency versions are recorded reproducibly.

## 6. Testing

-   Unit test version extraction from build info, including tagged,
    dirty, `(devel)` and unavailable build information.
-   Test command construction or command execution for expected
    help/version behaviour.
-   Test error propagation and exit handling where practical.

## 7. Verification

``` bash
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/gostart
```

## 8. Implementation plan

Agent must propose exact files and dependency versions before editing.
Confirm the urfave/cli v3 API against its installed module documentation
rather than copying v2 examples.
