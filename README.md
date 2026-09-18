# gitmanager

A CLI that clones and syncs GitLab repositories into a local folder tree,
mirroring the group/subgroup structure from GitLab.

[![CI](https://github.com/danny270793/Git-Manager-CLI/actions/workflows/ci.yml/badge.svg)](https://github.com/danny270793/Git-Manager-CLI/actions/workflows/ci.yml)

- `sync --group=<url>` walks a GitLab group recursively (including
  subgroups), discovers every project, and clones or updates each one.
- `sync --repo=<url>` clones or updates a single repository.
- `pull --destination=<dir>` walks a local folder recursively and updates
  every repository already cloned there, without talking to GitLab (no
  discovery of new repositories).
- Repositories that already exist locally are updated by stashing any
  local changes, switching to `main`, and pulling.
- A progress bar tracks how many repositories have been synced.

## Requirements

- `git` on `PATH`
- For `--method=ssh`: an SSH key registered with GitLab and available to
  your SSH agent
- For `--group` on private groups: a GitLab personal access token with
  `read_api` scope

## Install

### Prebuilt binary

Downloads the latest release for your OS/architecture and installs it to
`~/.local/bin/gitmanager`:

```sh
curl -fsSL https://raw.githubusercontent.com/danny270793/Git-Manager-CLI/main/scripts/install.sh | bash
```

Make sure `~/.local/bin` is on your `PATH` — the script will tell you if it isn't.

### From source

Requires Go, managed via [asdf](https://asdf-vm.com/) — see [docs/asdf.md](docs/asdf.md) for setup.

```sh
git clone git@github.com:danny270793/Git-Manager-CLI.git
cd Git-Manager-CLI
asdf install
```

## Usage

Run without building, straight from source:

```sh
./scripts/start.sh sync --group=https://gitlab.com/danny270793 --method=ssh --destination=.
```

Or build a binary first:

```sh
./scripts/build.sh
./build/gitmanager sync --group=https://gitlab.com/danny270793 --method=ssh --destination=.
```

A project at `https://gitlab.com/danny270793/group/project` ends up at
`./group/project` (the group's own path segment, `danny270793`, is
stripped since `.` already represents the group root).

Sync a single repository instead of a whole group:

```bash
gitmanager sync --repo=https://gitlab.com/danny270793/group/project --method=ssh --destination=.
```

Update every repository already cloned under a folder, without checking
GitLab for new repositories:

```bash
gitmanager pull --destination=.
```

Other commands:

```bash
gitmanager --help
gitmanager --version
```

### `sync` flags

| Flag            | Description                                                                          | Default |
| --------------- | ------------------------------------------------------------------------------------- | ------- |
| `--group`       | GitLab group URL to sync recursively                                                  | -       |
| `--repo`        | Single GitLab repository URL to sync                                                  | -       |
| `--method`      | Clone method: `ssh` or `https`                                                        | `ssh`   |
| `--destination` | Local destination folder                                                              | `.`     |
| `--token`       | GitLab personal access token (falls back to `GITLAB_TOKEN` env var); only used by `--group` on private groups | - |

`--group` and `--repo` are mutually exclusive; exactly one is required.

### `pull` flags

| Flag            | Description                                    | Default |
| --------------- | ----------------------------------------------- | ------- |
| `--destination` | Local folder to walk for existing repositories | `.`     |

## Project layout

```
main.go                     entrypoint
cmd/                        CLI commands (root, sync, pull)
internal/giturl/            GitLab URL parsing and clone URL construction
internal/gitlabapi/         recursive group project listing via the GitLab API
internal/gitops/            git clone / stash+checkout+pull via the git CLI, local repo discovery
scripts/                    build, run, test, format and install scripts
```

## Development

```sh
./scripts/build.sh    # build the binary into ./build
./scripts/start.sh    # run from source, arguments are passed through
./scripts/test.sh     # run the test suite
./scripts/format.sh   # format the source with gofmt
```

Override the build output path or embedded version with the `OUTPUT` and
`VERSION` environment variables:

```sh
OUTPUT=/usr/local/bin/gitmanager VERSION=1.0.0 ./scripts/build.sh
```

See [AGENTS.md](AGENTS.md) for contribution conventions (Conventional Commits, branch naming, etc.) — the single source of truth for both human and AI contributors.
