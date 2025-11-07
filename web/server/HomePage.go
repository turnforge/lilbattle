package server

import (
	"net/http"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1/models"
)

type BasePage struct {
	Title               string
	BodyClass           string
	CustomHeader        bool
	BodyDataAttributes  string
	DisableSplashScreen bool
	SplashTitle         string
	SplashMessage       string
	ActiveTab           string
}

type HomePage struct {
	BasePage
	Header Header

	// Dashboard data
	RecentGames  []*v1.Game
	RecentWorlds []*v1.World
	TotalGames   int32
	TotalWorlds  int32
}

func (p *HomePage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	// Redirect to the first visible tab
	if !vc.HideWorlds {
		http.Redirect(w, r, "/worlds/", http.StatusFound)
		return nil, true
	} else if !vc.HideGames {
		http.Redirect(w, r, "/games/", http.StatusFound)
		return nil, true
	} else {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return nil, true
	}
}

type PrivacyPolicy struct {
	Header Header
}

func (p *PrivacyPolicy) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	return p.Header.Load(r, w, vc)
}

type TermsOfService struct {
	Header Header
}

func (p *TermsOfService) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	return p.Header.Load(r, w, vc)
}

func (g *TermsOfService) Copy() View { return &TermsOfService{} }
func (g *PrivacyPolicy) Copy() View  { return &PrivacyPolicy{} }
func (g *HomePage) Copy() View       { return &HomePage{} }
