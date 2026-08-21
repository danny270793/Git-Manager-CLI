# asdf

This project uses [asdf](https://asdf-vm.com/) to manage the Go version, pinned in [.tool-versions](../.tool-versions).

## Install

```sh
brew install asdf
```

Then add asdf to your shell (see the [asdf docs](https://asdf-vm.com/guide/getting-started.html) for your shell of choice), and restart your terminal.

## Add the golang plugin

```sh
asdf plugin add golang
```

## Install a specific Go version

```sh
asdf install golang 1.25.0
```

## Install the version pinned in .tool-versions

Run this from the repo root to install the Go version this project expects:

```sh
asdf install
```
