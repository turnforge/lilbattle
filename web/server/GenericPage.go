package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

type GenericPage struct {
	BasePage
	Header Header
}

func (v *GenericPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	return
}
