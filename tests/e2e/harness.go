//go:build e2e
// +build e2e

// Package e2e provides an integration-test harness for the recorded ww
// replay scripts under tests/e2etests/. The harness runs the scripts
// against an ALREADY-RUNNING server pointed at via LILBATTLE_E2E_SERVER
// (with the /api suffix included) and uses the WorldsService.CreateWorld
// RPC to seed fixtures — same code path against local dev, staging, or
// any other target.
//
// Gated behind the `e2e` build tag while the recorded scripts drift from
// current game rules (tracked in issue 183). Run with:
//
//	LILBATTLE_E2E_SERVER=http://localhost:8090/api \
//	  go test -tags=e2e ./tests/e2e/
//
// Or via the Makefile which sets sensible defaults:
//
//	make e2e            # requires a server up + LILBATTLE_E2E_SERVER set
//	make e2e-run REPLAY=29146
//	make e2e-full       # boots a local server, runs, tears down
package e2e

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	"github.com/turnforge/lilbattle/services/connectclient"
	pj "google.golang.org/protobuf/encoding/protojson"
)

// serverURL resolves the target server URL from LILBATTLE_E2E_SERVER.
// Every test call goes here first — a missing var is treated as an
// actionable configuration error, not a silent skip, since the whole
// suite depends on it and skipping would hide misconfiguration in CI.
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
// Fixture paths are anchored here rather than at CWD so the tests work
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

// ensureFixtureWorld seeds a fixture world onto the target server via
// the WorldsService RPC. Idempotent — checks GetWorld first, and only
// calls CreateWorld if the target reports the world missing.
//
// Using CreateWorld (rather than dropping files into a specific storage
// dir) keeps the harness backend-agnostic: works against FS, gorm, gae,
// or a remote server we don't have shell access to.
func ensureFixtureWorld(t *testing.T, serverBase, worldID string) {
	t.Helper()
	ctx := context.Background()
	client := connectclient.NewConnectWorldsClient(serverBase)

	// Probe: if the world already exists, we're done.
	resp, err := client.GetWorld(ctx, &v1.GetWorldRequest{Id: worldID})
	if err == nil && resp != nil && resp.World != nil {
		return
	}

	// Load fixture protos from the repo.
	fixtureDir := filepath.Join(repoRoot(t), "tests", "e2e", "fixtures", "worlds", worldID)
	world := loadWorld(t, filepath.Join(fixtureDir, "metadata.json"))
	worldData := loadWorldData(t, filepath.Join(fixtureDir, "data.json"))

	_, err = client.CreateWorld(ctx, &v1.CreateWorldRequest{
		World:     world,
		WorldData: worldData,
	})
	if err != nil {
		t.Fatalf("CreateWorld %s: %v", worldID, err)
	}
}

// loadWorld parses a fixture metadata.json into a *v1.World via
// protojson. The files were captured from a live server's storage layer
// (protojson-formatted), so no schema translation is needed.
func loadWorld(t *testing.T, path string) *v1.World {
	t.Helper()
	data := mustReadFile(t, path)
	w := &v1.World{}
	if err := pj.Unmarshal(data, w); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	return w
}

func loadWorldData(t *testing.T, path string) *v1.WorldData {
	t.Helper()
	data := mustReadFile(t, path)
	wd := &v1.WorldData{}
	if err := pj.Unmarshal(data, wd); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	return wd
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

