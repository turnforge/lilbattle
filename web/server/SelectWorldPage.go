package server

import "net/http"

type SelectWorldPage struct {
	BasePage
	Header Header

	WorldListView WorldListView
}

func (m *SelectWorldPage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	m.Title = "Select a World"
	m.Header.Load(r, w, vc)

	// Set action mode to "select" to show Play buttons instead of action menus
	m.WorldListView.ActionMode = "select"

	err, finished = m.WorldListView.Load(r, w, vc)
	if err != nil || finished {
		return
	}
	return
}

func (m *SelectWorldPage) Copy() View {
	return &SelectWorldPage{}
}
