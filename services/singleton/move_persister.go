package singleton

import (
	"context"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// MovePersister is the hook SingletonGamesService.SaveMoveGroup delegates to
// when the singleton needs to write moves durably. The singleton itself only
// holds an in-memory copy of the game — persistence is out-of-band and
// pluggable so browser (WASM) code can push writes back to the server, tests
// can capture the moves, and future 1P / vs-AI modes can persist to
// IndexedDB or a local-first mesh without changing the game logic.
//
// Implementations must be safe to call from any goroutine.
type MovePersister interface {
	// Save persists a completed move group and the resulting game state.
	// Called from BaseGamesService.ProcessMoves after the moves have been
	// applied to the in-memory runtime. A non-nil error surfaces to the
	// caller of ProcessMoves; the FE decides whether to log, roll back, or
	// reload from canonical state.
	Save(ctx context.Context, gameId string, state *v1.GameState, group *v1.GameMoveGroup) error
}

// NoopPersister discards every Save call. This is SingletonGamesService's
// default so autoplay smoke tests and any non-browser caller that doesn't
// wire a persister keep today's in-memory-only behavior.
type NoopPersister struct{}

// Save always returns nil.
func (NoopPersister) Save(context.Context, string, *v1.GameState, *v1.GameMoveGroup) error {
	return nil
}
