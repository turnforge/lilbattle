//go:build e2e
// +build e2e

// Package e2e provides an integration-test harness for the recorded ww
// replay scripts under tests/e2etests/. Each script drives real ww against
// a real server; the harness spins up an ephemeral server + tempdir per
// test so runs are hermetic.
//
// Gated behind the `e2e` build tag because the recorded scripts have
// drifted from current game rules — the harness itself is green, but the
// assertions inside the .sh files fail against today's rules engine. Run
// with `go test -tags=e2e ./tests/e2e/`. Drift fixes are tracked as
// separate follow-up issues.
package e2e

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// pickFreePort asks the kernel for an unused TCP port and returns it as a
// string. Two random ports (one for HTTP, one for gRPC) per server; the
// listener is closed immediately so main.go can rebind. Cheaper than
// running the whole server bootstrap twice to detect port conflicts.
func pickFreePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("pick free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	_ = lis.Close()
	return port
}

// wwBinaryPath resolves the ww binary the replay scripts will invoke.
// Precedence: LILBATTLE_WW_BIN env var (explicit override for CI or
// unusual layouts), then PATH lookup for a bare "ww" (the common dev
// case — the CLI is installed to GOBIN by `make cli` and picked up
// globally, per CLAUDE.md). No auto-build: if ww isn't findable the test
// fails with an actionable message pointing at `make cli`.
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

// repoRoot walks up from the current test binary's directory to find the
// go.mod. Needed because `go test` sets CWD to the test package's dir,
// but we need to `go build` from the module root.
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

// TestServer describes a running ephemeral lilbattle server.
type TestServer struct {
	URL           string // http://127.0.0.1:<port>
	StorageDir    string // parent dir holding games/ + worlds/
	GamesStorage  string // <StorageDir>/games
	WorldsStorage string // <StorageDir>/worlds
	cmd           *exec.Cmd
}

// startTestServer boots main.go with local backends, random ports, and a
// per-test storage dir via LILBATTLE_{GAMES,WORLDS}_STORAGE_DIR. Waits for
// the HTTP endpoint to answer before returning. t.Cleanup tears it down.
//
// Boot cost is dominated by `go run` compiling the server binary on first
// call (~5s cold). If this proves too slow for CI, cache the binary with
// t.TempDir shared across subtests — but Go's testing package makes a
// per-test tempdir hard to share cleanly, and the current 3-script suite
// wall-clock is dominated by the replays themselves, not boot.
func startTestServer(t *testing.T) *TestServer {
	t.Helper()
	storageDir := t.TempDir()
	gamesDir := filepath.Join(storageDir, "games")
	worldsDir := filepath.Join(storageDir, "worlds")
	if err := os.MkdirAll(gamesDir, 0o755); err != nil {
		t.Fatalf("mkdir games: %v", err)
	}
	if err := os.MkdirAll(worldsDir, 0o755); err != nil {
		t.Fatalf("mkdir worlds: %v", err)
	}

	httpPort := pickFreePort(t)
	grpcPort := pickFreePort(t)

	cmd := exec.Command("go", "run", "main.go",
		"-games_service_be=local",
		"-worlds_service_be=local",
	)
	cmd.Dir = repoRoot(t)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("LILBATTLE_WEB_PORT=:%d", httpPort),
		fmt.Sprintf("LILBATTLE_GRPC_PORT=:%d", grpcPort),
		"LILBATTLE_GAMES_STORAGE_DIR="+gamesDir,
		"LILBATTLE_WORLDS_STORAGE_DIR="+worldsDir,
		// Auth stays enabled; the CLI does its own auth via profile tokens.
		// The replay scripts run ww in local mode (no profile), which
		// hits the endpoint anonymously. `ww new` needs a server but no
		// auth; ProcessMoves DOES need auth. Games created with no user_id
		// (the default in `ww new` without profile) are playable by anyone
		// under current authz rules — see services/authz.CanSubmitMoves.
		"DISABLE_API_AUTH=true",
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatalf("start server: %v", err)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	waitUntilReady(t, url, 30*time.Second)

	server := &TestServer{
		URL:           url,
		StorageDir:    storageDir,
		GamesStorage:  gamesDir,
		WorldsStorage: worldsDir,
		cmd:           cmd,
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	})
	return server
}

// waitUntilReady polls the server's root URL until it responds or the
// timeout elapses. Any HTTP response (even 404) proves the listener is
// bound; we don't need a specific status.
func waitUntilReady(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	client := &http.Client{Timeout: 500 * time.Millisecond}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url + "/")
		if err == nil {
			_ = resp.Body.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("server at %s never became ready within %v", url, timeout)
}

// copyFixtureWorld copies a fixture world into the server's worlds
// storage so `ww new <worldID>` can find it. Fixtures live in
// tests/e2e/fixtures/worlds/<id>/ (data.json + metadata.json).
func copyFixtureWorld(t *testing.T, server *TestServer, worldID string) {
	t.Helper()
	src := filepath.Join(repoRoot(t), "tests", "e2e", "fixtures", "worlds", worldID)
	dst := filepath.Join(server.WorldsStorage, worldID)
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("copy fixture world %s: %v", worldID, err)
	}
}

// copyDir does a shallow (single-level) recursive copy. Fixtures are
// flat data.json + metadata.json — no need for the general-purpose
// io/fs.WalkDir dance.
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // fixtures are file-only
		}
		data, err := os.ReadFile(filepath.Join(src, entry.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, entry.Name()), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
