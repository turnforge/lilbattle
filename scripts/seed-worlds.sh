#!/usr/bin/env bash
# seed-worlds.sh: idempotently seed every fixture world under
# $LILBATTLE_E2E_DATA_DIR/fixtures/worlds/ onto the target server via
# `ww worlds ensure`.
#
# Requires:
#   LILBATTLE_SERVER       — target API endpoint (e.g. http://localhost:8090/api)
#   LILBATTLE_E2E_DATA_DIR — test-data root (contains fixtures/worlds/)
#   ww                     — on PATH (`make cli` installs to $GOBIN)
#
# `ww worlds ensure` treats existing-and-matching as success and content
# mismatch as a hard error — so this script never overwrites production
# worlds. If a mismatch surfaces, the operator picks the fix path.
set -euo pipefail

if [[ -z "${LILBATTLE_SERVER:-}" ]]; then
    echo "seed-worlds: LILBATTLE_SERVER not set" >&2
    exit 1
fi
if [[ -z "${LILBATTLE_E2E_DATA_DIR:-}" ]]; then
    echo "seed-worlds: LILBATTLE_E2E_DATA_DIR not set (e.g. ~/projects/weemaps/e2e)" >&2
    exit 1
fi

FIXTURE_ROOT="$LILBATTLE_E2E_DATA_DIR/fixtures/worlds"
if [[ ! -d "$FIXTURE_ROOT" ]]; then
    echo "seed-worlds: fixture root $FIXTURE_ROOT missing" >&2
    exit 1
fi

for fixture_dir in "$FIXTURE_ROOT"/*/; do
    [[ -d "$fixture_dir" ]] || continue
    world_id="$(basename "$fixture_dir")"
    echo "[seed-worlds] ensuring $world_id from $fixture_dir"
    ww worlds ensure "$world_id" --data-dir "$fixture_dir"
done

echo "[seed-worlds] done"
