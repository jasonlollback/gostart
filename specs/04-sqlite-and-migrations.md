# Specification GS-030--034 --- Generated SQLite and migration system

Status: proposed\
Depends on: GS-010, GS-011\
Related decisions: ADR-002, ADR-003, ADR-004

## Outcome

Generated database-enabled applications use SQLite and run safe,
forward-only SQL migrations before the HTTP server becomes ready.

## SQLite requirements

-   `database/sql` with `modernc.org/sqlite`.
-   DSN includes `_time_format=sqlite&_timezone=UTC`.
-   Enable and verify foreign-key enforcement.
-   WAL mode.
-   Bounded connection settings and documented busy handling.
-   STRICT tables.
-   TEXT timestamps.
-   UUIDv7 identifiers.
-   Database file path is configurable and persistent in Docker.
-   Tests cover timestamp precision and round trips.

## Migration file format

``` text
migrations/sqlite/
  000001_create_users.sql
  000002_create_sessions.sql
```

Plain SQL only. No down files or rollback support.

Metadata table:

``` sql
CREATE TABLE schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    checksum   TEXT NOT NULL,
    applied_at TEXT NOT NULL
) STRICT;
```

## Runner requirements

The runner receives explicit dependencies: database connection,
migration filesystem and dialect identifier. It must not use global
application state.

-   Discover and sort migrations.
-   Validate names, versions and duplicates.
-   Acquire concurrency protection.
-   Check applied checksums.
-   Execute pending migrations in order.
-   Record a migration only after successful execution.
-   Release protection on success and failure.
-   Return actionable errors.
-   Startup fails if migration execution fails.

The implementation must verify actual `modernc.org/sqlite` transaction
and locking semantics. Do not assume a lock strategy based on another
SQLite driver.

## Startup lifecycle

1.  Load and validate config.
2.  Open database.
3.  Verify connection.
4.  Run migrations.
5.  Provision initial administrator if auth/admin are selected.
6.  Initialise handlers/services.
7.  Mark readiness and start serving.

## Concurrency and failure requirements

-   Multiple instances starting against the same database must not apply
    the same migration concurrently or corrupt metadata.
-   Lock acquisition must have bounded waiting/retry behaviour.
-   A failed migration must not be marked successful.
-   A checksum mismatch must stop startup.
-   Applied migrations are immutable; fix mistakes with a new migration.

## Acceptance criteria

-   [ ] Fresh database migrates to latest version.
-   [ ] Existing database applies only pending migrations.
-   [ ] Repeated startup is safe.
-   [ ] Checksum mismatch fails.
-   [ ] Failed migration is not recorded as successful.
-   [ ] Concurrent startup test passes reliably.
-   [ ] Migration failure prevents readiness.
-   [ ] No down migration mechanism exists.
-   [ ] Migration runner can be tested with an in-memory or temporary
    SQLite database.
-   [ ] SQLite driver-specific locking assumptions are documented and
    tested.
