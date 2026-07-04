#!/usr/bin/env bash
# e2e-full.sh: boot a dev server, run the e2e replay harness against it,
# tear the server down cleanly on exit (including on failure or Ctrl-C).
#
# Config knobs (env overrides):
#   LILBATTLE_E2E_HTTP_PORT   (default 8090)  — HTTP port for the dev server
#   LILBATTLE_E2E_GRPC_PORT   (default 9091)  — gRPC port for the dev server
#   LILBATTLE_E2E_STARTUP_TIMEOUT (default 30) — seconds to wait for /ready
#
# Anything after `--` is passed to `go test`, so you can target a single
# replay: `scripts/e2e-full.sh -- -run TestReplayScripts/29146`.
#
# CI: cheapest possible orchestration for the full suite. Local dev: use
# `make e2e` against a server you already have running instead.
set -euo pipefail

HTTP_PORT="${LILBATTLE_E2E_HTTP_PORT:-8090}"
GRPC_PORT="${LILBATTLE_E2E_GRPC_PORT:-9091}"
TIMEOUT="${LILBATTLE_E2E_STARTUP_TIMEOUT:-30}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "[e2e-full] booting server on :$HTTP_PORT (gRPC :$GRPC_PORT)"
LILBATTLE_WEB_PORT=":$HTTP_PORT" \
LILBATTLE_GRPC_PORT=":$GRPC_PORT" \
DISABLE_API_AUTH=true \
go run main.go -games_service_be=local -worlds_service_be=local >/tmp/e2e-server.log 2>&1 &
SERVER_PID=$!

# Kill the server on any exit — success, failure, or interrupt.
cleanup() {
    if kill -0 "$SERVER_PID" 2>/dev/null; then
        echo "[e2e-full] stopping server ($SERVER_PID)"
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT INT TERM

# Wait for the server to answer HTTP. `curl -sf` treats non-2xx as
# failure; the game viewer route returns 200 or 404 depending on the
# game ID, both of which prove the listener is bound.
echo "[e2e-full] waiting for server to become ready..."
deadline=$(( SECONDS + TIMEOUT ))
until curl -s -o /dev/null "http://localhost:$HTTP_PORT/games/probe/view"; do
    if (( SECONDS >= deadline )); then
        echo "[e2e-full] server did not become ready within ${TIMEOUT}s" >&2
        echo "[e2e-full] last 20 lines of server log:" >&2
        tail -20 /tmp/e2e-server.log >&2 || true
        exit 1
    fi
    sleep 1
done
echo "[e2e-full] server ready"

# Run the e2e tests. Any extra args after `--` land here so a single
# replay can be targeted from the command line.
extra_args=()
if (( $# > 0 )); then
    while [[ "${1:-}" != "--" && $# -gt 0 ]]; do shift; done
    if [[ "${1:-}" == "--" ]]; then shift; fi
    extra_args=("$@")
fi

echo "[e2e-full] running tests (extra args: ${extra_args[*]:-<none>})"
LILBATTLE_E2E_SERVER="http://localhost:$HTTP_PORT/api" \
go test -tags=e2e ./tests/e2e/ -v -count=1 "${extra_args[@]}"
