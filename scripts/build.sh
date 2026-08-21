#!/usr/bin/env bash
# Builds the gitmanager binary into ./build.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

OUTPUT="${OUTPUT:-build/gitmanager}"
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"

mkdir -p "$(dirname "$OUTPUT")"

go build -ldflags "-X danny270793/gitmanager/cmd.Version=${VERSION}" -o "$OUTPUT" .

echo "Built ${OUTPUT} (version ${VERSION})"
