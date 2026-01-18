package server

import (
	"log"
	"net/http"

	goal "github.com/panyam/goapplib"
	protos "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

type GameDetailPage struct {
	BasePage
	Header Header
	Game   *protos.Game
	GameId string
}

func (p *GameDetailPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*LilBattleApp]) (err error, finished bool) {
	p.GameId = r.PathValue("appItemId")
	if p.GameId == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return nil, true
	}

	p.Title = "Game Details"
	p.Header.Load(r, w, app)

	ctx := app.Context
	loggedInUserId := ctx.AuthMiddleware.GetLoggedInUserId(r)
	client := ctx.ClientMgr.GetGamesSvcClient()
	req := &protos.GetGameRequest{Id: p.GameId}

	resp, err := client.GetGame(GrpcAuthContext(loggedInUserId), req)
	if err != nil {
		log.Printf("Error fetching Game %s: %v", p.GameId, err)
		return HandleGRPCError(err, w, r, app)
	}

	if resp.Game != nil {
		p.Game = &protos.Game{
			Id:          resp.Game.Id,
			Name:        resp.Game.Name,
			Description: resp.Game.Description,
			CreatedAt:   resp.Game.CreatedAt,
			UpdatedAt:   resp.Game.UpdatedAt,
		}
		p.Title = p.Game.Name
		p.SetCanonicalFromRequest(app, r)
		p.MetaTitle = p.Game.Name + " - LilBattle Game"
		if p.Game.Description != "" {
			p.MetaDescription = p.Game.Description
		} else {
			p.MetaDescription = "View " + p.Game.Name + " - a LilBattle game."
		}
	}

	return nil, false
}
