package lib

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	pj "google.golang.org/protobuf/encoding/protojson"
)

// fakeWorldsClient records the calls made against it and returns
// pre-configured responses. Lets the four EnsureWorldExists cases run
// without spinning up an actual server.
type fakeWorldsClient struct {
	getResp *v1.GetWorldResponse
	getErr  error

	createResp *v1.CreateWorldResponse
	createErr  error

	getCalls    []string
	createCalls []*v1.CreateWorldRequest
}

func (f *fakeWorldsClient) GetWorld(_ context.Context, req *v1.GetWorldRequest) (*v1.GetWorldResponse, error) {
	f.getCalls = append(f.getCalls, req.Id)
	return f.getResp, f.getErr
}

func (f *fakeWorldsClient) CreateWorld(_ context.Context, req *v1.CreateWorldRequest) (*v1.CreateWorldResponse, error) {
	f.createCalls = append(f.createCalls, req)
	return f.createResp, f.createErr
}

// writeFixtureDir emits a WorldFixture-compatible dir (metadata.json +
// data.json) at t.TempDir()/<worldID>/. Returns the dir path so callers
// can pass it as WorldFixture.Dir.
func writeFixtureDir(t *testing.T, worldID string, tilesMap map[string]*v1.Tile) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), worldID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	world := &v1.World{Id: worldID, Name: worldID + " test"}
	worldData := &v1.WorldData{TilesMap: tilesMap}
	worldBytes, err := pj.Marshal(world)
	if err != nil {
		t.Fatalf("marshal world: %v", err)
	}
	dataBytes, err := pj.Marshal(worldData)
	if err != nil {
		t.Fatalf("marshal worldData: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "metadata.json"), worldBytes, 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "data.json"), dataBytes, 0o644); err != nil {
		t.Fatalf("write data: %v", err)
	}
	return dir
}

func tilesMap(q, r int32) map[string]*v1.Tile {
	return map[string]*v1.Tile{
		"0,0": {Q: q, R: r, TileType: 1},
	}
}

// notFoundErr mirrors what the Connect client sees when the FS backend
// returns a gRPC NotFound status — a stringy error containing the
// backend's message. isNotFound matches on substring, so any of these
// three variants is valid.
type notFoundErr struct{ msg string }

func (e *notFoundErr) Error() string { return e.msg }

// TestHashWorldData_DeterministicAcrossMapOrder pins the hash's stability
// contract: two logically-equal WorldData values with maps built in
// different insertion orders must produce the same hash. Without
// proto.MarshalOptions{Deterministic: true} this is the first thing that
// breaks; the constant is the whole reason HashWorldData exists.
func TestHashWorldData_DeterministicAcrossMapOrder(t *testing.T) {
	a := &v1.WorldData{TilesMap: map[string]*v1.Tile{
		"1,0": {Q: 1, R: 0, TileType: 1},
		"0,1": {Q: 0, R: 1, TileType: 2},
	}}
	b := &v1.WorldData{TilesMap: map[string]*v1.Tile{
		"0,1": {Q: 0, R: 1, TileType: 2},
		"1,0": {Q: 1, R: 0, TileType: 1},
	}}
	ha, err := HashWorldData(a)
	if err != nil {
		t.Fatalf("hash a: %v", err)
	}
	hb, err := HashWorldData(b)
	if err != nil {
		t.Fatalf("hash b: %v", err)
	}
	if ha != hb {
		t.Errorf("expected same hash for logically-equal WorldData; got %q vs %q", ha, hb)
	}
}

// TestHashWorldData_DifferentForDifferentContent guards the other side:
// a real content change (different tile type at same coord) MUST change
// the hash, otherwise the match check in EnsureWorldExists is a
// rubber stamp.
func TestHashWorldData_DifferentForDifferentContent(t *testing.T) {
	a := &v1.WorldData{TilesMap: tilesMap(1, 0)}
	b := &v1.WorldData{TilesMap: tilesMap(2, 0)}
	ha, _ := HashWorldData(a)
	hb, _ := HashWorldData(b)
	if ha == hb {
		t.Errorf("expected different hashes for different tile coords; both were %q", ha)
	}
}

// TestEnsureWorldExists_MissingNoFixture is the "must not silently
// succeed" case. Missing world + no fixture → hard error pointing the
// operator at seeding. No CreateWorld attempt.
func TestEnsureWorldExists_MissingNoFixture(t *testing.T) {
	client := &fakeWorldsClient{getErr: &notFoundErr{msg: "world not found: xyz"}}

	err := EnsureWorldExists(context.Background(), client, "xyz", nil)
	if err == nil {
		t.Fatal("expected error for missing world + no fixture")
	}
	if !strings.Contains(err.Error(), "no fixture") {
		t.Errorf("error should mention missing fixture; got: %v", err)
	}
	if len(client.createCalls) != 0 {
		t.Errorf("must not call CreateWorld when no fixture; got %d calls", len(client.createCalls))
	}
}

// TestEnsureWorldExists_MissingWithFixtureCreates pins the seed path.
// Missing world + fixture dir → CreateWorld with the fixture data, and
// the caller's world ID overrides whatever the fixture files carry.
func TestEnsureWorldExists_MissingWithFixtureCreates(t *testing.T) {
	client := &fakeWorldsClient{getErr: &notFoundErr{msg: "NotFound"}}
	fixtureDir := writeFixtureDir(t, "fixture-id", tilesMap(3, 4))

	err := EnsureWorldExists(context.Background(), client, "requested-id", &WorldFixture{Dir: fixtureDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(client.createCalls) != 1 {
		t.Fatalf("expected 1 CreateWorld call; got %d", len(client.createCalls))
	}
	got := client.createCalls[0].World
	if got.Id != "requested-id" {
		t.Errorf("CreateWorld id = %q, want %q (caller's ID must override fixture)", got.Id, "requested-id")
	}
}

// TestEnsureWorldExists_ExistsMatchesFixture is the "everything is
// aligned" case. Content matches → no CreateWorld, no error.
func TestEnsureWorldExists_ExistsMatchesFixture(t *testing.T) {
	wd := &v1.WorldData{TilesMap: tilesMap(5, 6)}
	client := &fakeWorldsClient{
		getResp: &v1.GetWorldResponse{
			World:     &v1.World{Id: "w1"},
			WorldData: wd,
		},
	}
	fixtureDir := writeFixtureDir(t, "w1", tilesMap(5, 6))

	err := EnsureWorldExists(context.Background(), client, "w1", &WorldFixture{Dir: fixtureDir})
	if err != nil {
		t.Fatalf("expected no error on match; got %v", err)
	}
	if len(client.createCalls) != 0 {
		t.Errorf("must not call CreateWorld on match; got %d calls", len(client.createCalls))
	}
}

// TestEnsureWorldExists_ExistsMismatchIsError is the safety property
// this whole exercise was about: server content diverges from the
// fixture → hard error, no overwrite.
func TestEnsureWorldExists_ExistsMismatchIsError(t *testing.T) {
	serverData := &v1.WorldData{TilesMap: tilesMap(5, 6)}
	client := &fakeWorldsClient{
		getResp: &v1.GetWorldResponse{
			World:     &v1.World{Id: "w1"},
			WorldData: serverData,
		},
	}
	fixtureDir := writeFixtureDir(t, "w1", tilesMap(99, 99)) // deliberately different

	err := EnsureWorldExists(context.Background(), client, "w1", &WorldFixture{Dir: fixtureDir})
	if err == nil {
		t.Fatal("expected error on content mismatch")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Errorf("error should mention mismatch; got: %v", err)
	}
	if len(client.createCalls) != 0 {
		t.Errorf("must not call CreateWorld on mismatch (never overwrite); got %d calls", len(client.createCalls))
	}
}

// TestEnsureWorldExists_ExistsNoFixture is the "presence probe only"
// path. World exists on target, caller didn't pass a fixture → nil,
// nothing to validate.
func TestEnsureWorldExists_ExistsNoFixture(t *testing.T) {
	client := &fakeWorldsClient{
		getResp: &v1.GetWorldResponse{
			World: &v1.World{Id: "w1"},
		},
	}
	err := EnsureWorldExists(context.Background(), client, "w1", nil)
	if err != nil {
		t.Errorf("expected no error for exists+no-fixture; got %v", err)
	}
}
