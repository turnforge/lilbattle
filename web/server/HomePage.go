package server

import (
	"net/http"
)

type BasePage struct {
	Title              string
	BodyClass          string
	CustomHeader       bool
	BodyDataAttributes string
}

type HomePage struct {
	BasePage
	Header Header

	// Add any other components here to reflect what you want to show in your home page
	// Note that you would also update your HomePage templates to reflect these
	AppItemListView AppItemListView
}

func (p *HomePage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.Title = "Home"
	p.Header.Load(r, w, vc)
	err, finished = p.AppItemListView.Load(r, w, vc)
	if err != nil || finished {
		return
	}
	return
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
