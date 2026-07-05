package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	pj "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// WorldsClient is the subset of the WorldsService client surface
// EnsureWorldExists needs. Kept as an interface so tests can substitute
// a fake without pulling in Connect / gRPC transport plumbing.
type WorldsClient interface {
	GetWorld(ctx context.Context, req *v1.GetWorldRequest) (*v1.GetWorldResponse, error)
	CreateWorld(ctx context.Context, req *v1.CreateWorldRequest) (*v1.CreateWorldResponse, error)
}

// WorldFixture identifies the local material EnsureWorldExists uses when
// the world is missing (create) or present (validate match). Exactly one
// of the three source forms should be non-zero; helpers below normalize
// them into World + WorldData protos.
//
// Dir: a directory containing "metadata.json" + "data.json" (the layout
// the FS backend uses in its storage tree — same shape as
// tests/e2e/fixtures/worlds/<id>/).
//
// MetadataJSONPath + DataJSONPath: explicit file paths, useful when the
// caller wants either an unconventional layout or a subset (metadata-
// only checks are rare but harmless).
//
// World + WorldData: pre-loaded protos, useful for programmatic callers
// that already have them in memory.
type WorldFixture struct {
	Dir              string
	MetadataJSONPath string
	DataJSONPath     string
	World            *v1.World
	WorldData        *v1.WorldData
}

// HashWorldData returns a deterministic hex-encoded SHA-256 of the
// world's game-content — tiles + units, everything WorldData carries.
// proto.MarshalOptions{Deterministic: true} sorts map keys and
// stabilizes field ordering, so two logically-equal WorldData values
// always produce the same hash regardless of insertion order or
// runtime state.
//
// Callers compare hashes to decide whether a fixture matches what's on
// the server. A mismatch is a hard error the operator must resolve;
// EnsureWorldExists never overwrites.
//
// If we later add a WorldData.content_hash field on the proto side
// (server-computed at write time), this function stays the client-side
// reference implementation for verifying against it.
func HashWorldData(wd *v1.WorldData) (string, error) {
	if wd == nil {
		return "", nil
	}
	data, err := proto.MarshalOptions{Deterministic: true}.Marshal(wd)
	if err != nil {
		return "", fmt.Errorf("marshal world data: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// EnsureWorldExists probes the target for worldID and, based on whether
// a fixture was supplied, either creates the world, validates it
// matches, or reports that the operator needs to seed it.
//
// Four cases:
//
//   1. Not on target + no fixture → error. Operator must seed manually
//      or pass a fixture. Never silently succeeds.
//   2. Not on target + fixture → CreateWorld with the fixture data.
//      Preserves the fixture's world ID.
//   3. On target + no fixture → returns nil. Basic presence check
//      succeeded; no content validation possible.
//   4. On target + fixture → compare HashWorldData(server) vs
//      HashWorldData(fixture). Match → nil. Mismatch → error. Never
//      overwrites; the operator decides whether to fix the target, fix
//      the fixture, or use a fresh world ID.
func EnsureWorldExists(ctx context.Context, client WorldsClient, worldID string, fixture *WorldFixture) error {
	if worldID == "" {
		return fmt.Errorf("worldID is required")
	}

	// Probe. Any error from GetWorld other than "not found" (which the
	// server signals via a nil World and NotFound status) surfaces
	// immediately — network / auth issues shouldn't be masked as
	// "world missing, will create."
	resp, err := client.GetWorld(ctx, &v1.GetWorldRequest{Id: worldID})
	found := err == nil && resp != nil && resp.World != nil
	if err != nil && !isNotFound(err) {
		return fmt.Errorf("GetWorld %s: %w", worldID, err)
	}

	if !found {
		if fixture == nil {
			return fmt.Errorf("world %s not on target and no fixture supplied — seed the target manually or pass a fixture", worldID)
		}
		world, worldData, lerr := loadFixture(fixture)
		if lerr != nil {
			return lerr
		}
		// Preserve the caller's world ID over whatever the fixture files
		// happened to carry — the .sh replay scripts and the operator's
		// mental model both key off the ID they passed in.
		world.Id = worldID
		_, cerr := client.CreateWorld(ctx, &v1.CreateWorldRequest{
			World:     world,
			WorldData: worldData,
		})
		if cerr != nil {
			return fmt.Errorf("CreateWorld %s: %w", worldID, cerr)
		}
		return nil
	}

	// World exists on target. If a fixture was passed, validate content.
	if fixture == nil {
		return nil
	}
	_, fixtureData, lerr := loadFixture(fixture)
	if lerr != nil {
		return lerr
	}
	serverHash, err := HashWorldData(resp.WorldData)
	if err != nil {
		return fmt.Errorf("hash server WorldData for %s: %w", worldID, err)
	}
	fixtureHash, err := HashWorldData(fixtureData)
	if err != nil {
		return fmt.Errorf("hash fixture WorldData for %s: %w", worldID, err)
	}
	if serverHash != fixtureHash {
		return fmt.Errorf("world %s content mismatch (server hash %s != fixture hash %s) — fix the target, update the fixture, or use a fresh world ID; ensure never overwrites",
			worldID, serverHash, fixtureHash)
	}
	return nil
}

// loadFixture normalizes a WorldFixture into (world, worldData) protos.
// Precedence: explicit World/WorldData fields > MetadataJSONPath +
// DataJSONPath > Dir (with the FS-backend layout).
func loadFixture(f *WorldFixture) (*v1.World, *v1.WorldData, error) {
	if f == nil {
		return nil, nil, fmt.Errorf("nil fixture")
	}
	if f.World != nil && f.WorldData != nil {
		return f.World, f.WorldData, nil
	}
	metaPath := f.MetadataJSONPath
	dataPath := f.DataJSONPath
	if f.Dir != "" {
		if metaPath == "" {
			metaPath = filepath.Join(f.Dir, "metadata.json")
		}
		if dataPath == "" {
			dataPath = filepath.Join(f.Dir, "data.json")
		}
	}
	if metaPath == "" || dataPath == "" {
		return nil, nil, fmt.Errorf("fixture must provide Dir OR (MetadataJSONPath + DataJSONPath) OR (World + WorldData)")
	}
	world, err := loadWorldProto[v1.World](metaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load metadata %s: %w", metaPath, err)
	}
	worldData, err := loadWorldProto[v1.WorldData](dataPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load data %s: %w", dataPath, err)
	}
	return world, worldData, nil
}

// loadWorldProto reads a protojson file into a fresh T. Generic wrapper
// so metadata.json and data.json share one implementation without a
// per-type wrapper each.
func loadWorldProto[T any, PT interface {
	*T
	proto.Message
}](path string) (PT, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var msg T
	pt := PT(&msg)
	if err := pj.Unmarshal(data, pt); err != nil {
		return nil, err
	}
	return pt, nil
}

// isNotFound recognizes the several ways WorldsService signals
// "the world you asked for isn't here." The FS backend returns a gRPC
// NotFound status; other backends may wrap it differently. Match the
// substring rather than depending on a specific error type since
// clients receive an already-serialized string over Connect.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "NotFound") ||
		strings.Contains(msg, "404")
}
