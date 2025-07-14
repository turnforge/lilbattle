package server

import (
	"net/http"
)

// MapsHandler handles map showcase pages
type MapsHandler struct {
	VC *ViewContext
}

// NewMapsHandler creates a new maps handler
func NewMapsHandler(vc *ViewContext) *MapsHandler {
	return &MapsHandler{VC: vc}
}

// Handler returns an HTTP handler for maps routes
func (h *MapsHandler) Handler() http.Handler {
	mux := http.NewServeMux()

	// Map listing page
	mux.HandleFunc("/maps", h.handleMapListing)
	mux.HandleFunc("/maps/", h.handleMapListing)

	// Map details page
	mux.HandleFunc("/map/", h.handleMapDetails)

	return mux
}

// handleMapListing renders the map listing page
func (h *MapsHandler) handleMapListing(w http.ResponseWriter, r *http.Request) {
	/*
		// Get all maps from catalog
		maps := h.catalog.ListMaps()

		// Prepare template data
		data := map[string]any{
			"Title":    "Map Examples",
			"PageType": "map-listing",
			"Maps": maps,
			"PageDataJSON": toJSON(map[string]any{
				"pageType": "map-listing",
			}),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("maps/listing.html", "")

		// Render the template
		if err := h.templateGroup.RenderHtmlTemplate(w, templates[0], "", data, nil); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
			return
		}
	*/
}

// handleMapDetails renders the map details page
func (h *MapsHandler) handleMapDetails(w http.ResponseWriter, r *http.Request) {
	/*
		// Extract map ID from path
		// Path format: /map/bitly
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		mapID := parts[2]

		// Get map from catalog
		map := h.catalog.GetMap(mapID)
		if map == nil {
			http.NotFound(w, r)
			return
		}

		// Get mode from query params (default to server mode)
		mode := "server"
		if r.URL.Query().Get("mode") == "wasm" {
			mode = "wasm"
		}

		// Get version (default to map's default version)
		version := r.URL.Query().Get("version")
		if version == "" {
			version = map.DefaultVersion
		}

		// Get SDL and recipe content for the version
		versionData := map.Versions[version]

		// Prepare minimal page data for the client (content will be loaded via API)
		pageData := map[string]any{
			"mapId": map.ID,
			"mode":      mode,
		}

		// Prepare template data
		data := map[string]any{
			"Title":        map.Name + " - SDL Map",
			"PageType":     "map-details",
			"Map":      map,
			"Mode":         mode,
			"PageDataJSON": toJSON(pageData),
		}

		// Load and render template
		templates := h.templateGroup.MustLoad("maps/details.html", "")

		// Render the template
		if err := h.templateGroup.RenderHtmlTemplate(w, templates[0], "", data, nil); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
			return
		}
	*/
}
