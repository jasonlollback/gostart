# Specification GS-060--064 --- CLI and project generation experience

Status: proposed\
Depends on: GS-020, GS-023, GS-024\
Related decisions: ADR-008, ADR-009

## Outcome

Developers can generate projects interactively or non-interactively,
inspect available features and validate saved configurations.

## Commands

Initial command set: - `new` - `list` - `info` - `validate` - `version`

Use urfave/cli v3. Use Huh for interactive terminal prompts.

## Interactive generation

Gather only required choices: - Project name and Go module path. -
Application type. - Database choice. - Frontend choice. - Authentication
choice. - User/admin provisioning choice. - Optional features. -
Deployment choices.

Show a review of resolved features before writing files. Cancellation
must not leave a partial project.

## Non-interactive generation

Support CLI flags and a configuration file. Non-interactive mode must
never wait for terminal input. Missing required values produce
actionable errors and non-zero exit codes.

## Acceptance criteria

-   [ ] Interactive flow generates a valid project.
-   [ ] User can review resolved choices before writing.
-   [ ] Cancellation leaves no partial project.
-   [ ] Config-file generation works without prompts.
-   [ ] CLI flags override config only according to a documented
    precedence rule.
-   [ ] `list` and `info` expose feature descriptions and dependencies.
-   [ ] `validate` does not write files.
-   [ ] `version` includes build metadata where available.
-   [ ] Help and errors are useful and tested.
