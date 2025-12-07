package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

type WorldListingPage struct {
	BasePage
	Header Header

	WorldListView WorldListView
}

func (m *WorldListingPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	m.DisableSplashScreen = true
	m.Title = "Worlds"
	m.ActiveTab = "worlds"
	m.Header.Load(r, w, app)
	return m.WorldListView.Load(r, w, app)
}
