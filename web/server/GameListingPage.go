package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

type GameListingPage struct {
	BasePage
	Header Header

	GameListView GameListView
}

func (m *GameListingPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	m.Title = "Games"
	m.ActiveTab = "games"
	m.DisableSplashScreen = true
	m.Header.Load(r, w, app)
	return m.GameListView.Load(r, w, app)
}
