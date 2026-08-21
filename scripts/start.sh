#!/usr/bin/env bash
# Runs gitmanager from source, forwarding all arguments.
# Example: ./scripts/start.sh sync --group=https://gitlab.com/sofiinc --method=ssh --destination=.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

go run . "$@"
