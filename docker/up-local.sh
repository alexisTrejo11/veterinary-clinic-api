#!/usr/bin/env bash
# Convenience wrapper for the local Docker stack (same as compose.local.yml).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ ! -f .env ]; then
	echo "Missing .env — run: cp .env.example .env" >&2
	exit 1
fi

exec docker compose --env-file .env -f docker/compose.local.yml up "$@"
