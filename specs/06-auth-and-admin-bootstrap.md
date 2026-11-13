# Specification GS-050--055 --- Authentication and initial administrator

Status: proposed\
Depends on: GS-030, GS-040\
Related decisions: ADR-005, ADR-006

## Outcome

Optional authentication uses email one-time codes and database-backed
sessions. When an administrator feature is selected, the first
administrator can be provisioned from `ADMIN_EMAIL` without an
interactive container command.

## Email one-time codes

Requirements: - Cryptographically secure code generation. - Store only a
hash of the code. - Short, configurable validity period. - Single use. -
Attempt limits and resend throttling. - Avoid account enumeration in
public responses. - Never log codes or code hashes. - Email provider
behind a small interface. - Tests use a fake sender.

Exact code length, expiry, throttling values and initial provider are
open decisions and must be resolved in a dedicated design task before
production implementation.

## Sessions

Requirements: - Secure random session token. - Store token hash where
practical. - Database-backed expiry and revocation. - HTTP-only
cookie. - Secure and SameSite attributes configured appropriately. -
Logout revokes session. - Session expiry and revocation tested.

## ADMIN_EMAIL bootstrap

`.env.example` must contain an explanation similar to:

``` dotenv
# Email address for the initial system administrator.
# If no administrator exists on first startup, this account
# will be created with administrator privileges.
# Remove or leave blank after initial provisioning.
ADMIN_EMAIL=admin@example.org
```

Behaviour: 1. If an administrator already exists, do not change user
roles based on `ADMIN_EMAIL`. 2. If no administrator exists and
`ADMIN_EMAIL` is configured, normalise the email consistently with
authentication. 3. Locate the matching user or create it. 4. Grant the
initial administrator role. 5. Commit user and role changes atomically.
6. If no admin and no `ADMIN_EMAIL`, log a warning and continue startup.
7. Do not generate a password or send a one-time code as part of
bootstrap. 8. The administrator signs in through the normal email-code
flow.

The operation must be safe when multiple instances start simultaneously.
Use a database transaction and a uniqueness/locking strategy that has
been tested against SQLite.

## Acceptance criteria

-   [ ] Auth-disabled projects contain no auth routes, tables or
    dependencies.
-   [ ] Expired, reused and invalid codes are rejected.
-   [ ] Rate limits and attempt limits are tested.
-   [ ] Codes and session tokens are not logged.
-   [ ] Session expiry, logout and revocation work.
-   [ ] `ADMIN_EMAIL` provisions an administrator when none exists.
-   [ ] Existing admin prevents subsequent environment-driven promotion.
-   [ ] Existing non-admin matching user is promoted only as part of
    initial bootstrap when no admin exists.
-   [ ] Missing variable with no admin warns but does not prevent
    startup.
-   [ ] Concurrent bootstrap attempts create exactly one initial
    administrator.
-   [ ] `.env.example` explains removal after first successful
    provisioning.
