package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

type SelectWorldPage struct {
	BasePage
	Header        Header
	WorldListView WorldListView
}

func (m *SelectWorldPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	m.Title = "Select a World"
	m.DisableSplashScreen = true
	m.Header.Load(r, w, app)
	m.WorldListView.ActionMode = "select"
	return m.WorldListView.Load(r, w, app)
}
