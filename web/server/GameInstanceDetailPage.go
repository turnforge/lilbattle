package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type GameInstanceDetailPage struct {
	BasePage
	Header  Header
	GameInstance *protos.GameInstance
	GameInstanceId string
}

func (p *GameInstanceDetailPage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.GameInstanceId = r.PathValue("appItemId")
	if p.GameInstanceId == "" {
		http.Error(w, "GameInstance ID is required", http.StatusBadRequest)
		return nil, true
	}

	p.Title = "GameInstance Details"
	p.Header.Load(r, w, vc)

	// Fetch the GameInstance using the client manager
	client, err := vc.ClientMgr.GetGameInstancesSvcClient()
	if err != nil {
		log.Printf("Error getting GameInstances client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, true
	}

	req := &protos.GetGameInstanceRequest{
		Id: p.GameInstanceId,
	}

	resp, err := client.GetGameInstance(context.Background(), req)
	if err != nil {
		log.Printf("Error fetching GameInstance %s: %v", p.GameInstanceId, err)
		http.Error(w, "GameInstance not found", http.StatusNotFound)
		return nil, true
	}

	if resp.Appitem != nil {
		// Convert from GameInstanceProject to GameInstance (assuming we need the basic info)
		p.GameInstance = &protos.GameInstance{
			Id:          resp.Appitem.Id,
			Name:        resp.Appitem.Name,
			Description: resp.Appitem.Description,
			CreatedAt:   resp.Appitem.CreatedAt,
			UpdatedAt:   resp.Appitem.UpdatedAt,
		}
		p.Title = p.GameInstance.Name
	}

	return nil, false
}

func (p *GameInstanceDetailPage) Copy() View { 
	return &GameInstanceDetailPage{} 
}