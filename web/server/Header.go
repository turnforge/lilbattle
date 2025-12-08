package server

import (
	"html/template"
	"log"
	"net/http"

	goal "github.com/panyam/goapplib"
)

type HeaderMenuItem struct {
	Id            string
	Title         string
	Link          string
	Class         string
	LowWidthClass string
	Raw           string
	OnClick       string
	Attributes    map[string]string
}

type Header struct {
	AppName             string
	PageData            any
	HeaderLogo          func() template.HTML
	LogoTitle           string
	Styles              map[string]any
	Width               string
	HeaderStyleLink     string
	MenuStyleLink       string
	FixedHeader         bool
	ShowHomeButton      bool
	ShowComposeButton   bool
	ShowLoginButton     bool
	ShowLogoutButton    bool
	HomeButtonImage     string
	CenterMenuItems     []HeaderMenuItem
	SpecialMenuItems    []HeaderMenuItem
	RightMenuItems      []HeaderMenuItem
	HideCenterMenuItems bool
	IsLoggedIn          bool
	LoggedInUserId      string
	Username            string
}

func (h *Header) SetupDefaults() {
	h.AppName = "WeeWar"
	h.LogoTitle = h.AppName
	h.HeaderStyleLink = "/static/css/Header.css"
	h.MenuStyleLink = "/static/css/Menu.css"
	h.Width = "max-w-7xl"
	h.ShowHomeButton = false
	h.ShowComposeButton = true
	h.ShowLoginButton = true
	h.ShowLogoutButton = true
	h.HomeButtonImage = "/static/icons/homebutton.jpg"
	h.HeaderStyleLink = "/static/css/Header.css"
	h.MenuStyleLink = "/static/css/Menu.css"
	h.HideCenterMenuItems = false
	h.SpecialMenuItems = []HeaderMenuItem{
		{Title: "Compose", Link: "/compose"},
	}
	h.RightMenuItems = []HeaderMenuItem{
		{Title: "My Music", Link: "/mymusic"},
		// {Title: "Tutorial", Link: "/tutorial"},
	}
	h.Styles = map[string]any{
		"FixedHeightHeader":          false,
		"HeaderHeightIfFixed":        "80px",
		"MinWidthForFullWidthMenu":   "24em",
		"MinWidthForHamburgerMenu":   "48em",
		"MinWidthForCompressingLogo": "24em",
	}
}

func (v *Header) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	v.SetupDefaults()
	ctx := app.Context
	v.LoggedInUserId = ctx.AuthMiddleware.GetLoggedInUserId(r)
	v.IsLoggedIn = v.LoggedInUserId != ""

	// Load username if logged in
	if v.IsLoggedIn && ctx.AuthService != nil {
		user, userErr := ctx.AuthService.GetUserById(v.LoggedInUserId)
		if userErr != nil {
			log.Printf("Header: Error loading user %s: %v", v.LoggedInUserId, userErr)
		} else if user != nil {
			profile := user.Profile()
			log.Printf("Header: User profile for %s: %+v", v.LoggedInUserId, profile)
			if username, ok := profile["username"].(string); ok {
				v.Username = username
				log.Printf("Header: Set username to: %s", v.Username)
			} else {
				log.Printf("Header: No username in profile or wrong type")
			}
		}
	}

	return
}
