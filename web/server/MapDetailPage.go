package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type MapDetailPage struct {
	BasePage
	Header Header
	Map    *protos.Map
	MapId  string
}

func (p *MapDetailPage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.MapId = r.PathValue("appItemId")
	if p.MapId == "" {
		http.Error(w, "Map ID is required", http.StatusBadRequest)
		return nil, true
	}

	p.Title = "Map Details"
	p.Header.Load(r, w, vc)

	// Fetch the Map using the client manager
	client, err := vc.ClientMgr.GetMapsSvcClient()
	if err != nil {
		log.Printf("Error getting Maps client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, true
	}

	req := &protos.GetMapRequest{
		Id: p.MapId,
	}

	resp, err := client.GetMap(context.Background(), req)
	if err != nil {
		log.Printf("Error fetching Map %s: %v", p.MapId, err)
		http.Error(w, "Map not found", http.StatusNotFound)
		return nil, true
	}

	if resp.Map != nil {
		// Convert from MapProject to Map (assuming we need the basic info)
		p.Map = &protos.Map{
			Id:          resp.Map.Id,
			Name:        resp.Map.Name,
			Description: resp.Map.Description,
			CreatedAt:   resp.Map.CreatedAt,
			UpdatedAt:   resp.Map.UpdatedAt,
		}
		p.Title = p.Map.Name
	}

	return nil, false
}

func (p *MapDetailPage) Copy() View {
	return &MapDetailPage{}
}
