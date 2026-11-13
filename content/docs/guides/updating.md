---
title: "Update a Project"
description: "Bring a generated project forward without uploading your code"
order: 4
---

## TLDR

Commit your work. Install the newest CLI. Then run:

```bash
go install github.com/axadrn/gostart/v4/cmd/gostart@latest
gostart update
```

gostart creates a separate update branch. Your current branch and working files stay untouched.

## ELI5

Git compares three trees:

1. The version gostart originally generated.
2. Your current project with your changes.
3. The newest version gostart would generate with the same choices.

Git already knows how to merge those trees. gostart does not implement its own merge algorithm.

gostart renders trees 1 and 3 on your machine. Tree 1 comes from the release named in `gostart.lock`, downloaded from the public Go module proxy. Tree 3 is the template inside your CLI. Tree 2 never leaves your computer.

`gostart new` already created the repository and its first commit. Before an update, you only need to commit your own current work.

## Requirements

- A generated project with `gostart.lock`.
- A Git repository. `gostart new` creates it automatically.
- A clean worktree, including no untracked files.
- A project made with gostart v4 or newer. Projects made with v3 or older keep working, but do not update to v4.
- Network access to the Go module proxy, or a `GOPROXY` that serves gostart releases.
- Git 2.38 or newer. Older versions lack `git merge-tree --write-tree`, and `gostart update` stops with a message naming your version.

## Review the result

The new branch is named after the target version:

```text
gostart-update-v4.1.0
```

Switch to it and review the commit:

```bash
git switch gostart-update-v4.1.0
git show --stat
```

When Git finds conflicts, the update branch contains normal conflict markers. Resolve and commit them like any other Git conflict.

## Configuration changes

New source files such as `internal/brand/brand.go` or `internal/content/content.go` usually merge cleanly. If your project already owns the same path, Git shows an explicit conflict instead of guessing.

The same rule applies when an update changes generated config or `.env.example`. Your private `.env` is ignored and never replaced. Review the new example, copy over newly required runtime values, and keep architecture choices in source.

## Abort

Stay on the original branch and delete the update branch:

```bash
git branch -D gostart-update-v4.1.0
```

No restore command is needed because gostart never checked out or changed your original branch.

## Current boundary

The first update version keeps the original generator config. Changing providers or adding Workspaces later belongs to the planned `gostart add` command family.
