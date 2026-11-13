# Specification-driven development workflow

## Purpose

Use this workflow with Codex, Claude Code or another coding agent. It is
designed to prevent large, ambiguous implementation prompts and
uncontrolled architecture drift.

## Source of truth

Use this precedence:

1.  `docs/DECISIONS.md` for agreed architecture decisions.
2.  `docs/PRODUCT.md` for product scope and goals.
3.  `docs/ARCHITECTURE.md` for structural conventions.
4.  The assigned feature specification for detailed behaviour and
    acceptance criteria.
5.  The approved implementation plan for the current task.

If sources conflict, stop and report the conflict. Do not choose
whichever instruction is easiest.

## Work-item lifecycle

A work item moves through:

`proposed -> ready -> in-progress -> review -> done`

A work item is `ready` only when:

-   Its purpose and user-visible outcome are clear.
-   Scope and non-goals are explicit.
-   Dependencies are satisfied.
-   Acceptance criteria are testable.
-   Security and data implications are considered.
-   Required tests are identified.

## Before implementation

The agent must:

1.  Read the mandatory repository documents and assigned specification.
2.  Inspect the repository, current branch, working tree and relevant
    code.
3.  Identify existing patterns before proposing new ones.
4.  Summarise understanding and any ambiguity.
5.  Propose files to add/change.
6.  Propose tests and exact commands.
7.  Identify risks and migration implications.
8.  Wait for approval when asked to plan first, or whenever an
    unresolved decision affects architecture or security.

Do not ask the agent to plan the entire product repeatedly. Plan only
the assigned slice.

## Implementation

-   Implement one vertical slice at a time.
-   Prefer small commits or reviewable diffs.
-   Write tests with the implementation.
-   Keep unrelated formatting or refactoring out of the change.
-   Do not update dependencies casually.
-   Do not regenerate large files without checking the diff.
-   Do not add placeholders that make acceptance appear complete.
-   Keep generated templates and their output in sync.

## Verification

The agent must run the commands required by the specification and report
actual results. Typical checks include:

``` bash
gofmt -w .
go test ./...
go vet ./...
staticcheck ./...
golangci-lint run
govulncheck ./...
```

Only run tools that are installed and relevant to the current slice. If
a required tool is unavailable, report it rather than claiming success.

For generated application templates, generate a fixture and run its own
formatting, code generation, tests and build. A unit test of the
template renderer alone is insufficient.

## Review checklist

-   Does the implementation satisfy every acceptance criterion?
-   Is the scope limited to the assigned work item?
-   Are errors handled and surfaced clearly?
-   Are secrets and sensitive values excluded from logs?
-   Are concurrency and failure cases covered?
-   Are generated files correct and reproducible?
-   Are docs and examples updated?
-   Are tests meaningful and deterministic?
-   Has any decision changed? If so, was it approved and recorded?

## Agent completion report

Use this structure:

``` text
Implemented:
- ...

Files changed:
- ...

Acceptance criteria:
- [x] ...
- [ ] ...

Verification:
- `command`: PASS/FAIL/NOT RUN — observed result

Design decisions:
- ...

Known limitations / follow-up:
- ...

Deviation from specification:
- None / explain
```

## Prompt pattern

Use this prompt for an implementation slice:

> Read `AGENTS.md`, the product and architecture documents, the decision
> register, the SDD workflow and the assigned specification. Inspect the
> repository before editing. First return a concise implementation plan,
> files expected to change, tests to add and commands to run. Identify
> any ambiguity or conflict. Do not change agreed decisions. Wait for
> approval. After approval, implement only this slice, run the required
> checks and provide the prescribed completion report.

For a second agent performing review:

> Review this change against the assigned specification and `AGENTS.md`.
> Do not modify files. Identify correctness bugs, security issues,
> missing acceptance criteria, concurrency risks, untested failure paths
> and scope creep. Cite the relevant file and code location. Distinguish
> blocking issues from suggestions.
