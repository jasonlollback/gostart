# Goilerplate customisation plan

## Direction

The Goilerplate fork is the canonical implementation of gostart. We will customise the existing codebase rather than porting the separate `gostart_v1` implementation or rebuilding the generator around its architecture.

The developer is already familiar with Goilerplate's codebase. Retaining its generator, template-selection mechanism, package structure and proven workflows is an explicit project preference.

## Role of gostart_v1

The former `gostart_v1` repository is a reference for product preferences, design ideas and selected generated-application defaults. Its code and architecture are not being merged. Its specifications are not binding requirements; use them as input when making individual feature decisions.

## Principles

1. Start from the current Goilerplate implementation and understand its existing conventions.
2. Rebrand the product as gostart without unnecessary structural changes.
3. Keep existing features that are useful and maintainable.
4. Remove features that do not fit the intended product or create disproportionate maintenance cost.
5. Add and customise capabilities incrementally, following existing patterns where possible.
6. Avoid parallel architectures, wholesale rewrites and premature generalisation.
7. Test both the generator and representative generated projects after template changes.
8. Preserve Goilerplate license files, notices, attribution and Git history.
9. Keep the upstream remote and selectively incorporate useful upstream fixes.

## Work sequence

### 1. Establish gostart identity

- [x] Create an integration branch from the Goilerplate fork.
- [x] Establish a gostart README and product identity on the integration branch.
- [ ] Rename the Go module, CLI, commands, workflows and release assets coherently.
- [ ] Review website, docs, content and generated-project references for old branding.
- [ ] Retain license and copyright notices.

### 2. Audit existing features

- [ ] Inventory generator options and their dependencies.
- [ ] Map each option to the template files, tests, documentation and generated dependencies it affects.
- [ ] Identify features to retain, simplify, change or remove.
- [ ] Identify gaps where a new capability is needed.
- [ ] Record decisions before making broad template changes.

### 3. Customise incrementally

For each feature change:

- [ ] Define the user-visible behaviour and defaults.
- [ ] Identify the smallest changes to the existing generator and template tree.
- [ ] Update validation and feature interactions.
- [ ] Update generated configuration, docs and examples.
- [ ] Add generator tests and generated-project fixture coverage.
- [ ] Verify unaffected feature combinations.

### 4. Release and maintain

- [ ] Establish gostart release naming and versioning.
- [ ] Ensure CI exercises representative generated combinations.
- [ ] Document supported feature combinations and prerequisites.
- [ ] Review upstream changes periodically and selectively.

## Existing gostart_v1 documents

The `docs/PRODUCT.md`, `docs/ARCHITECTURE.md`, `docs/DECISIONS.md`, `docs/BACKLOG.md` and `specs/` files were ported from the earlier gostart_v1 work. Treat them as historical/reference material until each item is explicitly reconciled with the Goilerplate implementation and current product direction. Do not use their completion statuses as evidence that corresponding work exists in this repository.
