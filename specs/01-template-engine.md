# Specification GS-011/012/013 --- Template engine and safe generation

Status: done\
Depends on: GS-010\
Related decisions: ADR-001, ADR-009, ADR-010

## Outcome

The generator renders a minimal application from embedded templates into
a new destination directory, deterministically and without overwriting
user data.

## Requirements

-   Embed template source with `embed.FS`.
-   Use `text/template`.
-   Keep template files readable in the repository.
-   Use explicit typed template data.
-   Separate rendering from filesystem writing.
-   Reject path traversal and invalid output paths.
-   Detect duplicate output paths.
-   Refuse to overwrite existing files by default.
-   Return actionable errors containing the source template and
    destination path.
-   Include generated `gostart.yml` when configuration support lands.
-   No executable hooks in template manifests.

## Acceptance criteria

-   [ ] Identical inputs render identical file contents.
-   [ ] Template parse and execution errors identify the failing
    template.
-   [ ] Path traversal attempts are rejected.
-   [ ] Duplicate output paths are rejected.
-   [ ] Existing files are not silently overwritten.
-   [ ] Partial writes are cleaned up or reported in a way that permits
    safe recovery.
-   [ ] A generated minimal fixture formats, tests and builds.
-   [ ] Tests cover empty files, nested directories, permissions and
    filesystem errors.

## Suggested API shape

Do not treat this as a mandatory API until the implementation plan is
reviewed:

``` go
type Renderer interface {
    Render(ctx context.Context, config Config) ([]GeneratedFile, error)
}

type GeneratedFile struct {
    Path    string
    Content []byte
}
```

Returning an in-memory file set before writing makes collision detection
and validation easier. Avoid a public interface if a concrete type is
sufficient.

## Non-goals

-   Template updates.
-   Arbitrary template functions with side effects.
-   Runtime template loading from untrusted locations.
-   Automatic merging into an existing project.

## GS-011 --- Embedded template renderer

Status: done

Scope: in-memory rendering only. Path safety, duplicate detection,
overwrite protection, writing and golden fixtures belong to GS-012 and
GS-013.

### Decisions

-   **D1 Location.** One package, `internal/templates`. Template sets
    live beneath it (`internal/templates/minimal/`) and are embedded
    with `//go:embed all:minimal` so dot-files are included. No
    top-level `templates/` package.
-   **D2 Literal `{{`.** Default delimiters are kept. A literal `{{` in
    template source is written `{{"{{"}}`.
-   **D3 Every file is a template.** Every file under a set's root must
    end in `.tmpl`. The suffix is stripped to form the output path, and
    it keeps `.go` and `go.mod` sources out of the generator's build.
    Any other file (for example a stray `.DS_Store`) is an error, not
    skipped.
-   **D4 Module path guard.** `Data.Validate` runs before any template
    executes. It accepts only non-empty paths of ASCII letters, digits
    and `- . _ ~ /`, with no empty elements and no element that begins
    or ends with a dot. This only stops the value from breaking out of
    generated Go source or `go.mod`; full validation is GS-021. No
    `golang.org/x/mod` dependency.
-   **D5 Reference retained.** (Superseded by GS-013: the reference
    moved to `testdata/golden/minimal` and these checks now run through
    the writer.) `testdata/reference/minimal` stays as the
    expected output for `ModulePath: "example.com/app"`. A test asserts
    that rendering reproduces it byte for byte, so templates and
    reference cannot drift.
-   **D6 Template data.** `Data{ModulePath string}` only. The binary
    name and Go version stay fixed in the templates.

### API

``` go
type Data struct{ ModulePath string }
func (d Data) Validate() error

type File struct {
    Path    string // slash-separated, relative to the output directory
    Content []byte
}

// Render executes every template under root in fsys, sorted by path.
func Render(fsys fs.FS, root string, data Data) ([]File, error)

// Minimal renders the embedded minimal set.
func Minimal(data Data) ([]File, error)
```

Templates are parsed with `missingkey=error` and no custom functions.
Errors are `render template <path in fsys>: ...`. The template is
named after its source path (for example
`minimal/cmd/server/main.go.tmpl`), so `text/template` line numbers
identify the source too. On error no files are returned.

### Completion evidence

-   `task check`, `task golden:check`, `task lint`,
    `task staticcheck`, `task vuln`: pass.
-   `go test -race -count=1 ./internal/templates/...`: pass, including
    the reference equivalence test and a standalone gofmt/vet/test/build
    of the set rendered as `github.com/acme/widget`.
-   Deliberately removing the `{{"{{"}}` escape and adding a stray
    `.DS_Store` each made the equivalence test fail with a message
    naming the file.
-   The generator's `go.mod` and `go.sum` are unchanged.

## GS-012 --- Safe output writer

Status: done

Scope: writing a rendered file set to disk safely. Golden fixtures and
building a written minimal fixture are GS-013. Printing the created
files is the CLI's job (GS-060/061); the writer returns them.

### Decisions

-   **D1 Destination.** The destination either does not exist, in
    which case it is created and its parent must already exist, or is
    an existing empty directory. A non-empty directory (even one holding
    only `.git`) or a file is refused before anything is written. There
    is no overwrite or merge mode. Whether overwrite is ever supported
    remains an open decision in `docs/DECISIONS.md`.
-   **D2 No parent creation.** A missing parent is an error, so a typo
    cannot create a deep tree.
-   **D3 Write strategy.** Writes go through `os.OpenRoot(dir)`, which
    confines them to the destination. Directories are created one level
    at a time with `Mkdir` and never reused. Files are created with
    `O_CREATE|O_EXCL`, so a file or symlink that appears concurrently is
    never overwritten or followed. On failure, everything the call
    created is removed newest first, including the destination if it
    created it. If cleanup fails too, the error lists every path left
    behind. No staging directory and rename: that would not work for an
    existing empty destination. The cost is that a process killed
    mid-write leaves partial files.
-   **D4 Path rules.** Before the filesystem is touched, every path must
    pass `fs.ValidPath`, must not be `.` and must not contain a
    backslash or NUL. Exact duplicates, paths or directories differing
    only in letter case, and a path that is both a file and a directory
    are rejected on every platform. All problems are reported together.
-   **D5 Modes.** Files 0644 and directories 0755, before umask.
    `templates.File` has no mode field until a template needs one.
-   **D6 No context.** Cancellation is handled by GS-060.
-   **D7 API.** `generator.Write(dir string, files []templates.File)
    ([]string, error)` returns the created file paths sorted, or nil
    and an error. Input order does not matter and the input is not
    modified.

### Completion evidence

-   `task check`, `task golden:check`, `task lint`,
    `task staticcheck`, `task vuln`: pass.
-   `go test -race -count=1 ./internal/generator/`: pass, with no tests
    skipped when not run as root.
-   With cleanup disabled, the cleanup tests failed. With `O_EXCL`
    replaced by `O_TRUNC`, the symlink tests failed.
-   The generator's `go.mod` and `go.sum` are unchanged.

## GS-013 --- Golden fixture tests

Status: done

Scope: the minimal fixture only. Web, admin and SaaS fixtures are
GS-074; `gostart.yml` in the output is GS-024.

### Decisions

-   **D1 One golden fixture.** The GS-010 reference moved to
    `testdata/golden/minimal`. It is now the reviewed expected output
    of the minimal set rendered with `example.com/app`, not hand-written
    source. `internal/reference/reference_test.go` and
    `internal/templates/minimal_test.go` were removed; their checks
    moved into `internal/generator/golden_test.go`, which uses the real
    path: `templates.Minimal`, then `generator.Write`.
-   **D2 Updating.** `task golden:update` runs the golden test with
    `-update`. It deletes and rewrites `testdata/golden/minimal` (a fixed
    path in the test) with `generator.Write`. The change is reviewed as
    a diff. Golden fixtures are never edited by hand.
-   **D3 Module path.** The golden fixture uses `example.com/app`. A
    second render as `github.com/acme/widget` checks that `go.mod`
    declares it and that `example.com/app` appears nowhere. That check is
    textual and does not build.
-   **D4 Unselected features.** The minimal fixture must not contain
    `gostart.yml`, `sqlc.yml`, `Dockerfile`, `compose.yml`, `db/`,
    `migrations/`, `assets/`, `.github/`,
    `internal/{database,migrate,auth,users,admin}`, `*.templ` or `*.sql`
    files, nor import gostart, `database/sql`, `modernc.org/sqlite`,
    `github.com/a-h/templ` or `html/template`. These checks are
    independent of the golden fixture, so a careless update cannot hide
    them. The lists live in the test; later fixtures get their own.
-   **D5 No race detector in the test.** The golden test runs plain
    `go test`; `task golden:check` runs `-race` locally.

### Tests

`TestGoldenMinimal` generates the fixture twice in temporary
directories with `generator.Write`, then:

-   **repeatable:** both copies have the same paths, bytes and modes;
-   **matches golden:** the file list and bytes equal
    `testdata/golden/minimal`;
-   **no unselected features:** checks the paths and imports in D4;
-   **formats and builds:** runs `gofmt -l`, `go vet`, `go test` and
    `go build ./cmd/server` with `GOWORK=off` (skipped with `-short` or
    without `go`).

`TestGoldenMinimalModulePath` covers D3.

### Completion evidence

-   `task check`, `task golden:check`, `task lint`, `task staticcheck`,
    `task vuln`: pass.
-   `go test -race -count=1 ./internal/generator/...`: pass.
-   `task golden:update` with unchanged templates leaves the golden
    fixture unchanged.
-   Each of these deliberate breakages failed the test:
    -   a hand edit to a golden file;
    -   a `db/*.sql` template, including after a careless
        `golden:update`;
    -   a hard-coded `example.com/app` in a template.
-   The generator's `go.mod` and `go.sum` are unchanged.
