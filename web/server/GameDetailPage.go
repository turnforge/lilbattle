package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type GameDetailPage struct {
	BasePage
	Header  Header
	Game *protos.Game
	GameId string
}

func (p *GameDetailPage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.GameId = r.PathValue("appItemId")
	if p.GameId == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return nil, true
	}

	p.Title = "Game Details"
	p.Header.Load(r, w, vc)

	// Fetch the Game using the client manager
	client, err := vc.ClientMgr.GetGamesSvcClient()
	if err != nil {
		log.Printf("Error getting Games client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, true
	}

	req := &protos.GetGameRequest{
		Id: p.GameId,
	}

	resp, err := client.GetGame(context.Background(), req)
	if err != nil {
		log.Printf("Error fetching Game %s: %v", p.GameId, err)
		http.Error(w, "Game not found", http.StatusNotFound)
		return nil, true
	}

	if resp.Appitem != nil {
		// Convert from GameProject to Game (assuming we need the basic info)
		p.Game = &protos.Game{
			Id:          resp.Appitem.Id,
			Name:        resp.Appitem.Name,
			Description: resp.Appitem.Description,
			CreatedAt:   resp.Appitem.CreatedAt,
			UpdatedAt:   resp.Appitem.UpdatedAt,
		}
		p.Title = p.Game.Name
	}

	return nil, false
}

func (p *GameDetailPage) Copy() View { 
	return &GameDetailPage{} 
}