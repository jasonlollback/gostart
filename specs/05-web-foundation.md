# Specification GS-040--043 --- Templ web foundation

Status: proposed\
Depends on: GS-020, GS-011\
Related decisions: ADR-007

## Outcome

When frontend support is selected, the generated application contains a
working Templ/Tailwind/shadcn-templ web UI.

## Requirements

-   Templ for server-rendered HTML.
-   Tailwind CSS 4.
-   shadcn-templ as the component foundation.
-   HTMX for suitable server-driven interactions.
-   Minimal JavaScript.
-   Responsive layout and accessible navigation.
-   Templ generation and CSS build tasks in Taskfile.
-   Tool versions pinned or reproducibly managed in the separate tools
    module where applicable.
-   UI templates are omitted when frontend is disabled.

## Acceptance criteria

-   [ ] Clean checkout can generate Templ Go code.
-   [ ] Tailwind CSS builds from a clean checkout.
-   [ ] shadcn-templ setup is documented.
-   [ ] A minimal page renders.
-   [ ] Generated project compiles.
-   [ ] Disabled frontend produces no frontend-specific dependencies or
    assets.
-   [ ] No alternative UI framework is added.
-   [ ] State-changing HTMX requests have a documented CSRF strategy.
