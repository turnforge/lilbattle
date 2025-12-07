package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	goal "github.com/panyam/goapplib"
	protos "github.com/turnforge/weewar/gen/go/weewar/v1/models"
)

// WorldsGroup implements goal.PageGroup for /worlds routes.
type WorldsGroup struct {
	weewarApp *WeewarApp
	goalApp   *goal.App[*WeewarApp]
}

// RegisterRoutes registers all world-related routes using goal.Register.
func (g *WorldsGroup) RegisterRoutes(app *goal.App[*WeewarApp]) *http.ServeMux {
	mux := http.NewServeMux()

	// Register pages using goal's generic registration
	goal.Register[*WorldListingPage](app, mux, "/")
	goal.Register[*SelectWorldPage](app, mux, "/select")
	goal.Register[*WorldViewerPage](app, mux, "/{worldId}/view")
	goal.Register[*WorldEditorPage](app, mux, "/{worldId}/edit")

	// Custom handlers that don't fit the View pattern
	mux.HandleFunc("/new", createNewWorldHandler(app))
	mux.HandleFunc("/{worldId}/start", worldStartHandler)
	mux.HandleFunc("/{worldId}/copy", worldCopyHandler)
	// Screenshot handler delegates to WeewarApp's ViewsRoot method
	mux.HandleFunc("/{worldId}/screenshot/live", g.weewarApp.ViewsRoot.handleWorldScreenshotLive)
	mux.HandleFunc("/{worldId}", worldActionsHandler(app))

	return mux
}

// Custom handlers - these use closures to access the App

func createNewWorldHandler(app *goal.App[*WeewarApp]) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := app.Context
		loggedInUserId := ctx.AuthMiddleware.GetLoggedInUserId(req)
		client := ctx.ClientMgr.GetWorldsSvcClient()

		createReq := &protos.CreateWorldRequest{
			World: &protos.World{
				Name:        "Untitled World",
				Description: "",
				CreatorId:   loggedInUserId,
				Tags:        []string{},
				Difficulty:  "",
			},
		}

		resp, err := client.CreateWorld(context.Background(), createReq)
		if err != nil {
			log.Printf("Failed to create world: %v", err)
			http.Error(w, "Failed to create world", http.StatusInternalServerError)
			return
		}

		editURL := fmt.Sprintf("/worlds/%s/edit", resp.World.Id)
		http.Redirect(w, req, editURL, http.StatusFound)
	}
}

func worldStartHandler(w http.ResponseWriter, r *http.Request) {
	worldId := r.PathValue("worldId")
	redirectURL := fmt.Sprintf("/games/new?worldId=%s", worldId)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func worldCopyHandler(w http.ResponseWriter, r *http.Request) {
	worldId := r.PathValue("worldId")
	http.Redirect(w, r, fmt.Sprintf("/worlds/new?copyFrom=%s", worldId), http.StatusFound)
}

func worldActionsHandler(app *goal.App[*WeewarApp]) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodDelete:
			deleteWorldHandler(app, w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func deleteWorldHandler(app *goal.App[*WeewarApp], w http.ResponseWriter, req *http.Request) {
	worldId := req.PathValue("worldId")
	if worldId == "" {
		http.Error(w, "World ID is required", http.StatusBadRequest)
		return
	}

	ctx := app.Context
	loggedInUserId := ctx.AuthMiddleware.GetLoggedInUserId(req)
	log.Printf("Delete world request: worldId=%s, userId=%s", worldId, loggedInUserId)

	client := ctx.ClientMgr.GetWorldsSvcClient()
	deleteReq := &protos.DeleteWorldRequest{Id: worldId}

	_, err := client.DeleteWorld(context.Background(), deleteReq)
	if err != nil {
		log.Printf("Failed to delete world %s: %v", worldId, err)
		http.Error(w, "Failed to delete world", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted world: %s", worldId)
	http.Redirect(w, req, "/worlds/", http.StatusFound)
}
