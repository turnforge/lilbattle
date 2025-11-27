package services

import (
	"context"
	"fmt"

	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
	"github.com/turnforge/weewar/lib"
	"github.com/turnforge/weewar/web/assets/themes"
)

// Data-Only panel implementations

type PanelBase struct {
	Theme       themes.Theme
	RulesEngine *v1.RulesEngine
}

func (p *PanelBase) SetTheme(t themes.Theme) {
	p.Theme = t
}

func (p *PanelBase) SetRulesEngine(r *v1.RulesEngine) {
	p.RulesEngine = r
}

// BaseGameState is a non-UI implementation of GameState interface
// Used for CLI and testing - stores game state without rendering
type BaseGameState struct {
	Game  *v1.Game
	State *v1.GameState
}

func (b *BaseGameState) SetGameState(_ context.Context, req *v1.SetGameStateRequest) (*v1.SetGameStateResponse, error) {
	b.Game = req.Game
	b.State = req.State
	return nil, nil
}

func (b *BaseGameState) SetUnitAt(_ context.Context, req *v1.SetUnitAtRequest) (*v1.SetUnitAtResponse, error) {
	if b.State == nil || b.State.WorldData == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	// Initialize map if needed
	if b.State.WorldData.UnitsMap == nil {
		b.State.WorldData.UnitsMap = make(map[string]*v1.Unit)
	}

	// Set unit at coordinate using map-based storage
	key := lib.CoordKey(req.Q, req.R)
	b.State.WorldData.UnitsMap[key] = req.Unit

	return nil, nil
}

func (b *BaseGameState) RemoveUnitAt(_ context.Context, req *v1.RemoveUnitAtRequest) (*v1.RemoveUnitAtResponse, error) {
	if b.State == nil || b.State.WorldData == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	// Remove unit at coordinate using map-based storage
	if b.State.WorldData.UnitsMap != nil {
		key := lib.CoordKey(req.Q, req.R)
		delete(b.State.WorldData.UnitsMap, key)
	}

	return nil, nil
}

func (b *BaseGameState) UpdateGameStatus(_ context.Context, req *v1.UpdateGameStatusRequest) (*v1.UpdateGameStatusResponse, error) {
	if b.State == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	b.State.CurrentPlayer = req.CurrentPlayer
	b.State.TurnCounter = req.TurnCounter

	return nil, nil
}

type BaseUnitPanel struct {
	PanelBase
	Unit *v1.Unit
}

type BaseTilePanel struct {
	PanelBase
	Tile *v1.Tile
}

type BaseGameScene struct {
	PanelBase
	CurrentPathsRequest      *v1.ShowPathRequest
	CurrentHighlightsRequest *v1.ShowHighlightsRequest
}

func (b *BaseGameScene) ClearPaths(context.Context) {
	b.CurrentPathsRequest = nil
}

func (b *BaseGameScene) ClearHighlights(_ context.Context, req *v1.ClearHighlightsRequest) {
	// Only clear CurrentHighlightsRequest if clearing all or clearing specific interactive types
	if req == nil || len(req.Types) == 0 {
		b.CurrentHighlightsRequest = nil
	}
}

func (b *BaseGameScene) ShowPath(_ context.Context, p *v1.ShowPathRequest) {
	b.CurrentPathsRequest = p
}

func (b *BaseGameScene) ShowHighlights(_ context.Context, h *v1.ShowHighlightsRequest) {
	b.CurrentHighlightsRequest = h
}

// Animation methods - no-ops for CLI
func (b *BaseGameScene) MoveUnit(_ context.Context, _ *v1.MoveUnitRequest) (*v1.MoveUnitResponse, error) {
	return &v1.MoveUnitResponse{}, nil
}

func (b *BaseGameScene) ShowAttackEffect(_ context.Context, _ *v1.ShowAttackEffectRequest) (*v1.ShowAttackEffectResponse, error) {
	return &v1.ShowAttackEffectResponse{}, nil
}

func (b *BaseGameScene) ShowHealEffect(_ context.Context, _ *v1.ShowHealEffectRequest) (*v1.ShowHealEffectResponse, error) {
	return &v1.ShowHealEffectResponse{}, nil
}

func (b *BaseGameScene) ShowCaptureEffect(_ context.Context, _ *v1.ShowCaptureEffectRequest) (*v1.ShowCaptureEffectResponse, error) {
	return &v1.ShowCaptureEffectResponse{}, nil
}

func (b *BaseGameScene) SetUnitAt(_ context.Context, _ *v1.SetUnitAtRequest) (*v1.SetUnitAtResponse, error) {
	return &v1.SetUnitAtResponse{}, nil
}

func (b *BaseGameScene) RemoveUnitAt(_ context.Context, _ *v1.RemoveUnitAtRequest) (*v1.RemoveUnitAtResponse, error) {
	return &v1.RemoveUnitAtResponse{}, nil
}

type BaseTurnOptionsPanel struct {
	BaseUnitPanel
	Options *v1.GetOptionsAtResponse
}

func (b *BaseTurnOptionsPanel) CurrentOptions() *v1.GetOptionsAtResponse {
	return b.Options
}

func (b *BaseTurnOptionsPanel) SetCurrentUnit(_ context.Context, unit *v1.Unit, options *v1.GetOptionsAtResponse) {
	b.Unit = unit
	if options == nil {
		options = &v1.GetOptionsAtResponse{}
	}
	b.Options = options
}

func (b *BaseUnitPanel) CurrentUnit() *v1.Unit {
	return b.Unit
}

func (b *BaseUnitPanel) SetCurrentUnit(_ context.Context, u *v1.Unit) {
	b.Unit = u
}

func (b *BaseTilePanel) CurrentTile() *v1.Tile {
	return b.Tile
}

func (b *BaseTilePanel) SetCurrentTile(_ context.Context, u *v1.Tile) {
	b.Tile = u
}

type BaseBuildOptionsModal struct {
	PanelBase
	BuildOptions []*v1.BuildUnitAction
	Tile         *v1.Tile
	PlayerCoins  int32
}

func (b *BaseBuildOptionsModal) Show(_ context.Context, tile *v1.Tile, buildOptions []*v1.BuildUnitAction, playerCoins int32) {
	b.Tile = tile
	b.BuildOptions = buildOptions
	b.PlayerCoins = playerCoins
}

func (b *BaseBuildOptionsModal) Hide(_ context.Context) {
	b.Tile = nil
	b.BuildOptions = nil
	b.PlayerCoins = 0
}

type BaseCompactSummaryCardPanel struct {
	PanelBase
	Tile *v1.Tile
	Unit *v1.Unit
}

func (b *BaseCompactSummaryCardPanel) SetCurrentData(_ context.Context, tile *v1.Tile, unit *v1.Unit) {
	b.Tile = tile
	b.Unit = unit
}
