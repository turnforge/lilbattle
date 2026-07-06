//go:build e2e
// +build e2e

// Package e2e provides an integration-test harness for recorded ww
// replay scripts. The harness itself lives here (in lilbattle, where
// ww + lib.EnsureWorldExists live), but the scripts + fixture world
// data live outside this repo — the test data root is passed via
// LILBATTLE_E2E_DATA_DIR:
//
//	LILBATTLE_E2E_DATA_DIR/
//	├── replays/*.sh                # what the harness runs
//	└── fixtures/worlds/<id>/       # world data seeded via ww worlds ensure
//	    ├── data.json
//	    └── metadata.json
//
// Run with:
//
//	LILBATTLE_E2E_SERVER=http://localhost:8090/api \
//	LILBATTLE_E2E_DATA_DIR=/path/to/weemaps/e2e \
//	  go test -tags=e2e ./tests/e2e/
//
// Or via the Makefile:
//
//	make e2e DATA_DIR=~/projects/weemaps/e2e
//	make e2e-run REPLAY=29146 DATA_DIR=~/projects/weemaps/e2e
//	make e2e-full DATA_DIR=~/projects/weemaps/e2e
//
// Gated behind the `e2e` build tag while recorded scripts drift from
// current game rules (issue 183). Default `go test ./...` compiles
// doc.go and reports "no tests to run".
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// serverURL resolves the target server URL from LILBATTLE_E2E_SERVER.
// A missing var is treated as an actionable configuration error, not
// a silent skip, since the whole suite depends on it.
func serverURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("LILBATTLE_E2E_SERVER")
	if url == "" {
		t.Fatalf("LILBATTLE_E2E_SERVER not set — point at a running server (e.g. http://localhost:8090/api) or run via `make e2e-full`")
	}
	return url
}

// dataDir resolves the test-data root from LILBATTLE_E2E_DATA_DIR. The
// harness expects `<dir>/replays/*.sh` and (indirectly, via
// scripts/seed-worlds.sh) `<dir>/fixtures/worlds/<id>/`. Kept out of
// this repo so upstream-derived material (game IDs, world data) doesn't
// live in lilbattle's git history.
func dataDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("LILBATTLE_E2E_DATA_DIR")
	if dir == "" {
		t.Fatalf("LILBATTLE_E2E_DATA_DIR not set — point at a test-data root containing replays/ and fixtures/worlds/ (e.g. ~/projects/weemaps/e2e)")
	}
	return dir
}

// replaysDir returns `<data-dir>/replays`.
func replaysDir(t *testing.T) string {
	return filepath.Join(dataDir(t), "replays")
}

// wwBinaryPath resolves the ww binary the replay scripts will invoke.
// Precedence: LILBATTLE_WW_BIN env var (explicit override for CI or
// unusual layouts), then PATH lookup for a bare "ww" (the common dev
// case — `make cli` installs to GOBIN). No auto-build.
func wwBinaryPath(t *testing.T) string {
	t.Helper()
	if override := os.Getenv("LILBATTLE_WW_BIN"); override != "" {
		return override
	}
	path, err := exec.LookPath("ww")
	if err != nil {
		t.Fatalf("ww not on PATH and LILBATTLE_WW_BIN not set — run `make cli` from the repo root")
	}
	return path
}

// wwPathDir wraps wwBinaryPath in a tempdir with a symlink named exactly
// "ww". The .sh replay scripts call `ww ...` unqualified; prepending
// this dir to PATH ensures they resolve to the intended binary even
// when LILBATTLE_WW_BIN points at a differently-named artifact.
func wwPathDir(t *testing.T) string {
	t.Helper()
	src := wwBinaryPath(t)
	dir := t.TempDir()
	link := filepath.Join(dir, "ww")
	if err := os.Symlink(src, link); err != nil {
		t.Fatalf("symlink ww: %v", err)
	}
	return dir
}
