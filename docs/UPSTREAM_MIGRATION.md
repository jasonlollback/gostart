# Goilerplate integration plan

This repository retains the Goilerplate fork's Git history and implementation as a source of reusable code. The product direction and target architecture are defined by the gostart requirements and decisions in this repository.

## Current state

- The repository history and original upstream relationship are retained.
- The root project identity and contributor guidance have been changed to gostart.
- The gostart_v1 product requirements, architecture, decisions, workflow, backlog and specifications have been ported into this repository.
- The existing Goilerplate generator and generated SaaS template are still present. They are not yet the gostart v1 implementation.
- The existing Goilerplate website, content and release configuration have not yet been assessed for reuse or removal.
- No generated project behaviour has been changed in this integration foundation.

## Integration principles

1. Keep useful, tested Goilerplate implementation where it fits the gostart product.
2. Do not preserve a feature merely because it exists upstream. In particular, the Goilerplate matrix of databases, payment providers, mail providers, frontends, workspaces and SaaS capabilities is not the v1 product scope.
3. Prefer the smaller, tested gostart_v1 renderer and safe writer where they are stronger than the current implementation. Reuse concepts and tests; do not copy blindly.
4. Generated output must be independent of the generator and must not receive automatic updates.
5. Preserve Goilerplate license files and notices for inherited code. Keep original attribution intact.
6. Keep upstream as a remote and incorporate upstream changes selectively, never by blindly merging a new SaaS feature matrix.

## Proposed sequence

### 1. Establish canonical identity

- [x] Create an integration branch from the Goilerplate fork.
- [x] Establish gostart README and agent guidance.
- [x] Port product, architecture, decision, backlog and specification documents.
- [ ] Rename the Go module and CLI from Goilerplate to gostart in one coherent change.
- [ ] Review GoReleaser, CI, release workflows, Taskfile and generated documentation for old product references.
- [ ] Decide whether the existing website and content are useful; rebrand, replace or remove them.
- [ ] Review inherited licenses and retain required notices.

### 2. Reconcile generator architecture

- [ ] Compare the existing Goilerplate selection/rendering pipeline with gostart_v1's render-plan-write pipeline.
- [ ] Adopt deterministic in-memory rendering and complete output-path validation before filesystem writes.
- [ ] Ensure writes are confined, exclusive and cleaned up on failure.
- [ ] Retain golden fixture coverage and add tests for path collisions, case-fold collisions and partial-write cleanup.
- [ ] Keep the implementation small; do not add a general plugin runtime.

### 3. Establish the v1 generated application

- [ ] Build or select the minimal reference app from the gostart_v1 baseline.
- [ ] Align generated project structure and conventions with docs/ARCHITECTURE.md.
- [ ] Add the SQLite and startup migration foundation.
- [ ] Add the chosen web stack and Taskfile/tools module.
- [ ] Add optional email OTP authentication and ADMIN_EMAIL initial administrator provisioning.
- [ ] Add configuration and declarative feature composition only as needed by working use cases.
- [ ] Verify each generated fixture independently.

### 4. Remove out-of-scope template capabilities

Only after the minimal v1 path works, remove or isolate unused Goilerplate capabilities such as PostgreSQL, Stripe/Polar billing, SMTP/Resend choices, OAuth, SvelteKit, workspaces, storage, MCP, and broad SaaS administration. Preserve code only when it is part of an explicitly accepted future capability.

### 5. Release readiness

- [ ] Replace Goilerplate naming in CLI help, docs, generated files and release assets.
- [ ] Confirm module path `github.com/jasonlollback/gostart` (no major-version suffix for v0/v1).
- [ ] Establish a clean CI matrix for the generator and generated fixtures.
- [ ] Tag the first gostart release only after the v1 acceptance criteria pass.

## Important

The `done` statuses from the original gostart_v1 backlog describe work completed in that separate repository. They do not mean the work has already been integrated or verified in this canonical repository. The backlog has been reset to make this distinction explicit.
