package server

import (
	"encoding/json"
	"net/http"
)

// GamesHandler handles game showcase pages
type GamesHandler struct {
	VC *ViewContext
}

// NewGamesHandler creates a new games handler
func NewGamesHandler(vc *ViewContext) *GamesHandler {
	return &GamesHandler{VC: vc}
}

// Handler returns an HTTP handler for games routes
func (h *GamesHandler) Handler() http.Handler {
	mux := http.NewServeMux()

	// Game listing page
	mux.HandleFunc("/games", h.handleGameListing)
	mux.HandleFunc("/games/", h.handleGameListing)

	// Game details page
	mux.HandleFunc("/game/", h.handleGameDetails)

	return mux
}

// handleGameListing renders the game listing page
func (h *GamesHandler) handleGameListing(w http.ResponseWriter, r *http.Request) {
	/*
		// Get all games from catalog
		games := h.catalog.ListGames()

		// Prepare template data
		data := map[string]any{
			"Title":    "Game Examples",
			"PageType": "game-listing",
			"Games": games,
			"PageDataJSON": toJSON(map[string]any{
				"pageType": "game-listing",
			}),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("games/listing.html", "")

		// Render the template
		if err := h.templateGroup.RenderHtmlTemplate(w, templates[0], "", data, nil); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
			return
		}
	*/
}

// handleGameDetails renders the game details page
func (h *GamesHandler) handleGameDetails(w http.ResponseWriter, r *http.Request) {
	/*
		// Extract game ID from path
		// Path format: /game/bitly
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		gameID := parts[2]

		// Get game from catalog
		game := h.catalog.GetGame(gameID)
		if game == nil {
			http.NotFound(w, r)
			return
		}

		// Get mode from query params (default to server mode)
		mode := "server"
		if r.URL.Query().Get("mode") == "wasm" {
			mode = "wasm"
		}

		// Get version (default to game's default version)
		version := r.URL.Query().Get("version")
		if version == "" {
			version = game.DefaultVersion
		}

		// Get SDL and recipe content for the version
		versionData := game.Versions[version]

		// Prepare minimal page data for the client (content will be loaded via API)
		pageData := map[string]any{
			"gameId": game.ID,
			"mode":      mode,
		}

		// Prepare template data
		data := map[string]any{
			"Title":        game.Name + " - SDL Game",
			"PageType":     "game-details",
			"Game":      game,
			"Mode":         mode,
			"PageDataJSON": toJSON(pageData),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("games/details.html", "")

		// Render the template
		if err := h.templateGroup.RenderHtmlTemplate(w, templates[0], "", data, nil); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
			return
		}
	*/
}

// toJSON converts data to JSON string for template use
func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
