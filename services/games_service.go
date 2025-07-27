package services

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

var GAMES_STORAGE_DIR = weewar.DevDataPath("storage/games")

// GamesServiceImpl implements the GamesService gRPC interface
type GamesServiceImpl struct {
	BaseGamesServiceImpl
	storage *FileStorage // Storage area for all files
}

// NewGamesService creates a new GamesService implementation for server mode
func NewGamesService() *GamesServiceImpl {
	service := &GamesServiceImpl{
		BaseGamesServiceImpl: BaseGamesServiceImpl{
			WorldsService: NewWorldsServiceImpl(),
		},
		storage: NewFileStorage(GAMES_STORAGE_DIR),
	}

	return service
}

// ListGames returns all available games (metadata only for performance)
func (s *GamesServiceImpl) ListGames(ctx context.Context, req *v1.ListGamesRequest) (resp *v1.ListGamesResponse, err error) {
	resp = &v1.ListGamesResponse{
		Items: []*v1.Game{},
		Pagination: &v1.PaginationResponse{
			HasMore:      false,
			TotalResults: 0,
		},
	}
	resp.Items, err = ListFSEntities[*v1.Game](s.storage, nil)
	resp.Pagination.TotalResults = int32(len(resp.Items))
	return resp, nil
}

// DeleteGame deletes a game
func (s *GamesServiceImpl) DeleteGame(ctx context.Context, req *v1.DeleteGameRequest) (resp *v1.DeleteGameResponse, err error) {
	resp = &v1.DeleteGameResponse{}
	err = s.storage.DeleteEntity(req.Id)
	return
}

// CreateWorld creates a new world
func (s *GamesServiceImpl) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (resp *v1.CreateGameResponse, err error) {
	if req.Game == nil {
		return nil, fmt.Errorf("game data is required")
	}

	req.Game.Id, err = s.storage.CreateEntity(req.Game.Id)
	if err != nil {
		return resp, err
	}

	now := time.Now()
	req.Game.CreatedAt = tspb.New(now)
	req.Game.UpdatedAt = tspb.New(now)

	// Save game metadta
	if err := s.storage.SaveArtifact(req.Game.Id, "metadata", req.Game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	world, err := s.WorldsService.GetWorld(ctx, &v1.GetWorldRequest{Id: req.Game.WorldId})
	if err != nil {
		return nil, fmt.Errorf("Error loading world: %w", err)
	}

	// Save a new empty game state and a new move list
	gs := &v1.GameState{
		GameId:    req.Game.Id,
		WorldData: world.WorldData,
	}
	if err := s.storage.SaveArtifact(req.Game.Id, "state", gs); err != nil {
		log.Printf("Failed to create state for game %s: %v", req.Game.Id, err)
	}

	// Save a new empty game history and a new move list
	if err := s.storage.SaveArtifact(req.Game.Id, "history", &v1.GameMoveHistory{GameId: req.Game.Id}); err != nil {
		log.Printf("Failed to create state for game %s: %v", req.Game.Id, err)
	}

	resp = &v1.CreateGameResponse{
		Game:      req.Game,
		GameState: gs,
	}

	return resp, nil
}

// GetGame returns a specific game with complete data including tiles and units
func (s *GamesServiceImpl) GetGame(ctx context.Context, req *v1.GetGameRequest) (resp *v1.GetGameResponse, err error) {
	if req.Id == "" {
		return nil, fmt.Errorf("game ID is required")
	}

	game, err := LoadFSArtifact[*v1.Game](s.storage, req.Id, "metadata")
	if err != nil {
		return nil, fmt.Errorf("game metadata not found: %w", err)
	}

	gameState, err := LoadFSArtifact[*v1.GameState](s.storage, req.Id, "state")
	if err != nil {
		return nil, fmt.Errorf("game state not found: %w", err)
	}

	gameHistory, err := LoadFSArtifact[*v1.GameMoveHistory](s.storage, req.Id, "history")
	if err != nil {
		return nil, fmt.Errorf("game state not found: %w", err)
	}

	resp = &v1.GetGameResponse{
		Game:        game,
		State:       gameState,
		MoveHistory: gameHistory,
	}

	return resp, nil
}

// ProcessMoves processes moves for an existing game
func (s *GamesServiceImpl) ProcessMoves(ctx context.Context, req *v1.ProcessMovesRequest) (resp *v1.ProcessMovesResponse, err error) {
	if len(req.Moves) == 0 {
		return nil, fmt.Errorf("at least one move is required")
	}

	gameresp, err := s.GetGame(ctx, &v1.GetGameRequest{Id: req.GameId})
	if err != nil || gameresp.Game == nil {
		return nil, err
	}

	// Given the game we want to get the game runtime (for its active state)
	rtGame, gameState, err := s.Storage.LoadGameState(gameresp.Game.Id)
	if err != nil {
		return nil, err
	}

	// Get the moves validted by the move processor, it is upto the move processor
	// to decide how "transactional" it wants to be - ie fail after  N moves,
	// success only if all moves succeeds etc.  Note that at this point the game
	// state has not changed and neither has the Runtime Game object.  Both the
	// GameState and the Runtime Game are checkpointed at before the moves started
	var dmp weewar.DefaultMoveProcessor
	results, err := dmp.ProcessMoves(rtGame, req.Moves)
	if err != nil {
		return nil, err
	}
	resp.MoveResults = results

	// Now that we have the results, we want to update our gamestate by applying the
	// results - this would also set the next "checkoint" to after the reuslts.
	// It is upto the storage to see how the runtime game is also updated.  For example
	// a storage that persists the gameState may just not do anythign and let it be
	// reconstructed on the next load
	err = s.Storage.SaveGameState(gameState, results)

	return resp, err
}

// GetMovementOptions returns available movement options for a unit
func (s *GamesServiceImpl) GetMovementOptions(ctx context.Context, req *v1.GetMovementOptionsRequest) (resp *v1.GetMovementOptionsResponse, err error) {
	// TODO: Implement actual movement calculation logic
	// For now, return mock data
	resp = &v1.GetMovementOptionsResponse{
		Options: []*v1.MovementOption{
			// Mock movement options - replace with actual game logic
		},
	}
	return resp, nil
}

// GetAttackOptions returns available attack options for a unit
func (s *GamesServiceImpl) GetAttackOptions(ctx context.Context, req *v1.GetAttackOptionsRequest) (resp *v1.GetAttackOptionsResponse, err error) {
	// TODO: Implement actual attack calculation logic
	// For now, return mock data
	resp = &v1.GetAttackOptionsResponse{
		Options: []*v1.AttackOption{
			// Mock attack options - replace with actual game logic
		},
	}
	return resp, nil
}

// CanSelectUnit determines if a unit can be selected by the current player
func (s *GamesServiceImpl) CanSelectUnit(ctx context.Context, req *v1.CanSelectUnitRequest) (resp *v1.CanSelectUnitResponse, err error) {
	// TODO: Implement actual unit selection logic
	// For now, return mock data
	resp = &v1.CanSelectUnitResponse{
		CanSelect: true, // Mock - replace with actual ownership/turn logic
		Reason:    "",
	}
	return resp, nil
}
