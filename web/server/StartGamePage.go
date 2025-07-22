package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

var DefaultRulesEngine *weewar.RulesEngine

func init() {
	var err error
	// TODO - only for dev
	DefaultRulesEngine, err = weewar.LoadRulesEngineFromFile(weewar.DevDataPath("data/rules-data.json"))
	if err != nil {
		panic(fmt.Sprintf("Failed to load rules engine: %v", err))
	}
}

type StartGamePage struct {
	BasePage
	Header    Header
	Map       *protos.Map
	MapId     string
	UnitTypes []UnitType
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

	// Load unit types for unit restrictions UI
	p.loadUnitTypes()

	return nil, false
}

// loadUnitTypes populates the UnitTypes field for the unit restrictions UI
func (p *StartGamePage) loadUnitTypes() {
	// Load unit types with icons from rules engine
	p.UnitTypes = []UnitType{}

	// Get all available unit types from the rules engine
	rulesEngine := DefaultRulesEngine

	// If rules engine is not populated, fall back to GetUnitData function which uses the unitDataMap
	unitIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 24, 25, 26, 27, 28, 29}

	// If rules engine has units loaded, use those; otherwise use the static list
	if rulesEngine.GetLoadedUnitsCount() > 0 {
		// Use units from rules engine
		for unitID := range rulesEngine.Units {
			unitData := rulesEngine.Units[unitID]
			if unitData != nil {
				// Use web-accessible static URL path for the unit asset
				iconDataURL := fmt.Sprintf("/static/assets/v1/Units/%d/0.png", unitID)

				p.UnitTypes = append(p.UnitTypes, UnitType{
					ID:          unitData.ID,
					Name:        unitData.Name,
					IconDataURL: iconDataURL,
				})
			}
		}
	} else {
		// Fall back to static unit data map
		for _, unitID := range unitIDs {
			unitData := weewar.GetUnitData(unitID)
			if unitData != nil {
				// Use web-accessible static URL path for the unit asset
				iconDataURL := fmt.Sprintf("/static/assets/v1/Units/%d/0.png", unitID)

				p.UnitTypes = append(p.UnitTypes, UnitType{
					ID:          unitData.ID,
					Name:        unitData.Name,
					IconDataURL: iconDataURL,
				})
			}
		}
	}
}

func (p *StartGamePage) Copy() View {
	return &StartGamePage{}
}
