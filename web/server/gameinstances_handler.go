package server

import (
	"encoding/json"
	"net/http"
)

// GameInstancesHandler handles gameinstance showcase pages
type GameInstancesHandler struct {
	VC *ViewContext
}

// NewGameInstancesHandler creates a new gameinstances handler
func NewGameInstancesHandler(vc *ViewContext) *GameInstancesHandler {
	return &GameInstancesHandler{VC: vc}
}

// Handler returns an HTTP handler for gameinstances routes
func (h *GameInstancesHandler) Handler() http.Handler {
	mux := http.NewServeMux()

	// GameInstance listing page
	mux.HandleFunc("/gameinstances", h.handleGameInstanceListing)
	mux.HandleFunc("/gameinstances/", h.handleGameInstanceListing)

	// GameInstance details page
	mux.HandleFunc("/gameinstance/", h.handleGameInstanceDetails)

	return mux
}

// handleGameInstanceListing renders the gameinstance listing page
func (h *GameInstancesHandler) handleGameInstanceListing(w http.ResponseWriter, r *http.Request) {
	/*
		// Get all gameinstances from catalog
		gameinstances := h.catalog.ListGameInstances()

		// Prepare template data
		data := map[string]any{
			"Title":    "GameInstance Examples",
			"PageType": "gameinstance-listing",
			"GameInstances": gameinstances,
			"PageDataJSON": toJSON(map[string]any{
				"pageType": "gameinstance-listing",
			}),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("gameinstances/listing.html", "")

		// Render the template
		if err := h.templateGroup.RenderHtmlTemplate(w, templates[0], "", data, nil); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
			return
		}
	*/
}

// handleGameInstanceDetails renders the gameinstance details page
func (h *GameInstancesHandler) handleGameInstanceDetails(w http.ResponseWriter, r *http.Request) {
	/*
		// Extract gameinstance ID from path
		// Path format: /gameinstance/bitly
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		gameinstanceID := parts[2]

		// Get gameinstance from catalog
		gameinstance := h.catalog.GetGameInstance(gameinstanceID)
		if gameinstance == nil {
			http.NotFound(w, r)
			return
		}

		// Get mode from query params (default to server mode)
		mode := "server"
		if r.URL.Query().Get("mode") == "wasm" {
			mode = "wasm"
		}

		// Get version (default to gameinstance's default version)
		version := r.URL.Query().Get("version")
		if version == "" {
			version = gameinstance.DefaultVersion
		}

		// Get SDL and recipe content for the version
		versionData := gameinstance.Versions[version]

		// Prepare minimal page data for the client (content will be loaded via API)
		pageData := map[string]any{
			"gameinstanceId": gameinstance.ID,
			"mode":      mode,
		}

		// Prepare template data
		data := map[string]any{
			"Title":        gameinstance.Name + " - SDL GameInstance",
			"PageType":     "gameinstance-details",
			"GameInstance":      gameinstance,
			"Mode":         mode,
			"PageDataJSON": toJSON(pageData),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("gameinstances/details.html", "")

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
