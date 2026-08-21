# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project

gitmanager — a Go command-line tool that clones and syncs GitLab repositories into a local folder tree.

## Toolchain

- Go version is pinned via [.tool-versions](.tool-versions) and managed with `asdf`. See [docs/asdf.md](docs/asdf.md) for setup.
- Run `asdf install` from the repo root to install the pinned Go version.

## Common commands

- Run locally: `./scripts/start.sh`
- Build binary (outputs to `./build`): `./scripts/build.sh`
- Run tests: `./scripts/test.sh`
- Format source: `./scripts/format.sh`

## Conventions

- Module path: `danny270793/gitmanager`.
- Keep `./build` out of version control (already in `.gitignore`).
- Commit related changes file by file rather than one large commit, unless told otherwise.
- Do not commit, create branches, or open MRs/PRs until the user explicitly asks for it.

## Commits: Conventional Commits (required)

Every commit message **must** follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

### Format

```unknown
<type>[optional scope]: <short description>

[optional body]

[optional footer(s)]
```

- **Description:** imperative mood, lowercase start (no trailing period required but stay consistent).
- **Header max length:** keep the first line ≤ **72** characters when practical.

### Allowed types (common)

Use one of: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.

- **`feat`:** new behavior or capability for users.
- **`fix`:** a bug fix.
- **`docs`:** documentation only.
- **`chore`:** maintenance that is not a user-facing feature or fix (deps, config, tooling).
- **`ci`:** CI/CD pipeline or automation only.

### Scope (optional)

A noun in parentheses after the type, e.g. `fix(gitops): handle detached HEAD before checkout`.

### Breaking changes

Either:

- append **`!`** after the type/scope: `feat(cli)!: rename --group to --namespace`, or
- add a footer: `BREAKING CHANGE: <what changed and what to do>`.

### Examples (valid)

- `feat(sync): add --repo flag for single-repository sync`
- `fix(gitops): skip stash when working tree is clean`
- `docs: document the pull command`
- `chore: bump go-gitlab client`

### Examples (invalid — do not use)

- `Update cli` (missing type)
- `Fixed bug` (not conventional)
- `WIP` / `misc changes`

When proposing or creating commits, **always** use this format. If multiple unrelated changes exist, **split into multiple commits** rather than one vague message.
