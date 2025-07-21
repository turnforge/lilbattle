package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type StartGamePage struct {
	BasePage
	Header Header
	Map    *protos.Map
	MapId  string
}

func (p *StartGamePage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	// Get mapId from query parameter (optional)
	p.MapId = r.URL.Query().Get("mapId")
	
	p.Title = "New Game"
	p.Header.Load(r, w, vc)

	// If a mapId is provided, fetch the map data
	if p.MapId != "" {
		// Fetch the Map using the client manager
		client, err := vc.ClientMgr.GetMapsSvcClient()
		if err != nil {
			log.Printf("Error getting Maps client: %v", err)
			// Don't fail the page, just log the error
			p.MapId = ""
		} else {
			req := &protos.GetMapRequest{
				Id: p.MapId,
			}

			resp, err := client.GetMap(context.Background(), req)
			if err != nil {
				log.Printf("Error fetching Map %s: %v", p.MapId, err)
				// Don't fail the page, just clear the mapId
				p.MapId = ""
			} else if resp.Map != nil {
				p.Map = resp.Map
				p.Title = "New Game - " + p.Map.Name
			}
		}
	}

	return nil, false
}

func (p *StartGamePage) Copy() View {
	return &StartGamePage{}
}