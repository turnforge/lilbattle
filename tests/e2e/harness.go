//go:build e2e
// +build e2e

// Package e2e provides an integration-test harness for the recorded ww
// replay scripts under tests/e2etests/. The harness runs the scripts
// against an ALREADY-RUNNING server pointed at via LILBATTLE_E2E_SERVER
// and assumes the referenced worlds already exist on that target. World
// seeding is a separate step — for local dev flows, `make e2e-full`
// invokes `scripts/seed-worlds.sh` which drives `ww worlds ensure`.
//
// Gated behind the `e2e` build tag while the recorded scripts drift
// from current game rules (tracked in issue 183). Run with:
//
//	LILBATTLE_E2E_SERVER=http://localhost:8090/api \
//	  go test -tags=e2e ./tests/e2e/
//
// Or via the Makefile:
//
//	make e2e            # requires a server up + worlds seeded
//	make e2e-run REPLAY=29146
//	make e2e-full       # boots a local server, seeds, runs, tears down
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// serverURL resolves the target server URL from LILBATTLE_E2E_SERVER.
// Every test call goes here first — a missing var is treated as an
// actionable configuration error, not a silent skip, since the whole
// suite depends on it.
func serverURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("LILBATTLE_E2E_SERVER")
	if url == "" {
		t.Fatalf("LILBATTLE_E2E_SERVER not set — point at a running server (e.g. http://localhost:8090/api) or run via `make e2e-full`")
	}
	return url
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

// repoRoot walks up from the current test's CWD to find the go.mod.
// Script paths are anchored here rather than at CWD so the tests work
// regardless of what dir `go test` is invoked from.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find go.mod (no repo root)")
		}
		dir = parent
	}
}
