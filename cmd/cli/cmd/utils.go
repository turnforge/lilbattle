package cmd

import (
	"context"
	"fmt"

	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
	"github.com/turnforge/weewar/lib"
)

func GetGame() (pc *PresenterContext, game *v1.Game, gameState *v1.GameState, gameHistory *v1.GameMoveHistory, rtGame *lib.Game, err error) {
	// Get game ID
	gameID, err = getGameID()
	if err != nil {
		return
	}

	// Create presenter
	pc, err = createPresenter(gameID)
	if err != nil {
		return
	}

	ctx := context.Background()
	getGameResp, err := pc.Presenter.GetGame(ctx, gameID)
	game, gameState, gameHistory = getGameResp.Game, getGameResp.State, getGameResp.History

	// Get runtime game for parsing positions
	rtGame, err = pc.Presenter.GamesService.GetRuntimeGame(getGameResp.Game, getGameResp.State)
	if err != nil {
		err = fmt.Errorf("failed to get runtime game: %w", err)
		return
	}
	return
}
