//go:build !wasm
// +build !wasm

package services

import (
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// TestBuildGameStatusRequest_MidGame pins the in-progress payload: end-state
// fields stay at their zero values while turn/player track gameState. This is
// the every-turn case that runs on every PlayerChanged / CoinsChanged event,
// so it must not accidentally announce victory.
func TestBuildGameStatusRequest_MidGame(t *testing.T) {
	state := &v1.GameState{
		CurrentPlayer: 2,
		TurnCounter:   7,
		Finished:      false,
		WinningPlayer: 0,
		Status:        v1.GameStatus_GAME_STATUS_PLAYING,
	}
	got := buildGameStatusRequest(state)

	if got.CurrentPlayer != 2 || got.TurnCounter != 7 {
		t.Errorf("turn/player mismatch: got cp=%d tc=%d", got.CurrentPlayer, got.TurnCounter)
	}
	if got.Finished {
		t.Errorf("expected Finished=false mid-game; got true")
	}
	if got.WinningPlayer != 0 {
		t.Errorf("expected WinningPlayer=0 mid-game; got %d", got.WinningPlayer)
	}
	if got.Status != v1.GameStatus_GAME_STATUS_PLAYING {
		t.Errorf("expected Status=IN_PROGRESS mid-game; got %v", got.Status)
	}
}

// TestBuildGameStatusRequest_GameEnded pins the victory payload: the three
// end-state fields propagate verbatim from gameState. This is the contract
// the FE relies on to trigger the GameEndedModal — if either Finished or
// the winner gets dropped here, the modal never fires.
func TestBuildGameStatusRequest_GameEnded(t *testing.T) {
	state := &v1.GameState{
		CurrentPlayer: 1,
		TurnCounter:   42,
		Finished:      true,
		WinningPlayer: 1,
		Status:        v1.GameStatus_GAME_STATUS_ENDED,
	}
	got := buildGameStatusRequest(state)

	if !got.Finished {
		t.Errorf("expected Finished=true after victory; got false")
	}
	if got.WinningPlayer != 1 {
		t.Errorf("expected WinningPlayer=1; got %d", got.WinningPlayer)
	}
	if got.Status != v1.GameStatus_GAME_STATUS_ENDED {
		t.Errorf("expected Status=ENDED; got %v", got.Status)
	}
	if got.CurrentPlayer != 1 || got.TurnCounter != 42 {
		t.Errorf("turn/player mismatch on ended payload: got cp=%d tc=%d", got.CurrentPlayer, got.TurnCounter)
	}
}
