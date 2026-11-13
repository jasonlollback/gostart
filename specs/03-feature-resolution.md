# Specification GS-022/023/024 --- Feature manifests and composition

Status: proposed\
Depends on: GS-020\
Related decisions: ADR-001, ADR-010, ADR-011

## Outcome

A requested configuration is converted into a validated, deterministic
set of features and templates before any output is written.

## Requirements

Each feature declares: - Stable ID. - Name and description. - Required
features. - Conflicts. - Template paths. - Optional config validation.

The resolver must: - Reject unknown features. - Detect dependency
cycles. - Resolve required dependencies. - Explain automatically enabled
features. - Reject conflicts. - Avoid guessing where user intent is
ambiguous. - Produce deterministic output ordering. - Complete before
filesystem mutation.

Manifests are data only. No executable hooks.

## Acceptance criteria

-   [ ] A feature with a missing dependency is handled according to
    documented auto-enable policy.
-   [ ] Cyclic dependencies fail clearly.
-   [ ] Conflicting features fail clearly.
-   [ ] Unknown feature IDs fail.
-   [ ] Repeated resolution yields the same result.
-   [ ] The user receives an explanation of auto-enabled features.
-   [ ] Invalid resolution writes no files.
-   [ ] Tests cover simple, nested, cyclic and conflicting graphs.
