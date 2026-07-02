//go:build !wasm
// +build !wasm

package singleton

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// recordingPersister captures the args of the last Save call and can be
// configured to return a specific error. Used to prove SaveMoveGroup wires
// the args through verbatim and propagates errors.
type recordingPersister struct {
	saved     int
	gotGameID string
	gotState  *v1.GameState
	gotGroup  *v1.GameMoveGroup
	err       error
}

func (r *recordingPersister) Save(_ context.Context, gameID string, state *v1.GameState, group *v1.GameMoveGroup) error {
	r.saved++
	r.gotGameID = gameID
	r.gotState = state
	r.gotGroup = group
	return r.err
}

// TestSaveMoveGroup_DelegatesToPersister pins the contract the browser JS
// bridge and any future persister impl rely on: SaveMoveGroup forwards
// (gameId, state, group) verbatim to the injected Persister and returns its
// error. Regression here would silently drop browser writes or swallow
// server auth errors.
func TestSaveMoveGroup_DelegatesToPersister(t *testing.T) {
	svc := NewSingletonGamesService()
	rec := &recordingPersister{}
	svc.Persister = rec

	state := &v1.GameState{TurnCounter: 7}
	group := &v1.GameMoveGroup{GroupNumber: 3}
	err := svc.SaveMoveGroup(context.Background(), "game-abc", state, group)

	if err != nil {
		t.Fatalf("SaveMoveGroup returned unexpected error: %v", err)
	}
	if rec.saved != 1 {
		t.Errorf("persister.saved = %d, want 1", rec.saved)
	}
	if rec.gotGameID != "game-abc" {
		t.Errorf("gameId = %q, want %q", rec.gotGameID, "game-abc")
	}
	if rec.gotState != state {
		t.Errorf("state pointer mismatch (persister must receive the caller's pointer verbatim)")
	}
	if rec.gotGroup != group {
		t.Errorf("group pointer mismatch (persister must receive the caller's pointer verbatim)")
	}
}

// TestSaveMoveGroup_PropagatesPersisterError pins the error path: server
// auth rejection (401 becomes non-nil error at the JS bridge) must reach the
// caller of SaveMoveGroup so BaseGamesService.ProcessMoves can bubble it up
// and the FE can log / react.
func TestSaveMoveGroup_PropagatesPersisterError(t *testing.T) {
	svc := NewSingletonGamesService()
	sentinel := errors.New("simulated 401 unauthorized")
	svc.Persister = &recordingPersister{err: sentinel}

	err := svc.SaveMoveGroup(context.Background(), "game-x", &v1.GameState{}, &v1.GameMoveGroup{})
	if !errors.Is(err, sentinel) {
		t.Errorf("SaveMoveGroup = %v, want %v", err, sentinel)
	}
}

// TestSaveMoveGroup_DefaultNoopPreservesLegacyBehavior pins the
// backward-compat contract for autoplay smoke tests and any non-browser
// caller: a freshly-constructed SingletonGamesService (no persister set)
// treats SaveMoveGroup as a no-op returning nil. Without the NoopPersister
// default, calling SaveMoveGroup would nil-dereference svc.Persister and
// break autoplay in a hard-to-diagnose way.
func TestSaveMoveGroup_DefaultNoopPreservesLegacyBehavior(t *testing.T) {
	svc := NewSingletonGamesService()

	err := svc.SaveMoveGroup(context.Background(), "any", &v1.GameState{}, &v1.GameMoveGroup{})
	if err != nil {
		t.Errorf("default NoopPersister.Save returned %v, want nil", err)
	}
}
