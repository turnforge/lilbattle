package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	protos "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// XML sitemap structures per sitemaps.org protocol

type SitemapIndex struct {
	XMLName  xml.Name         `xml:"sitemapindex"`
	XMLNS    string           `xml:"xmlns,attr"`
	Sitemaps []SitemapIndexEntry `xml:"sitemap"`
}

type SitemapIndexEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type URLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// RegisterSEORoutes registers robots.txt and sitemap handlers
func (n *RootViewsHandler) RegisterSEORoutes() {
	n.mux.HandleFunc("/robots.txt", n.handleRobotsTxt)
	n.mux.HandleFunc("/sitemap.xml", n.handleSitemapIndex)
	n.mux.HandleFunc("/sitemap-static.xml", n.handleSitemapStatic)
	n.mux.HandleFunc("/sitemap-worlds.xml", n.handleSitemapWorlds)
	n.mux.HandleFunc("/sitemap-games.xml", n.handleSitemapGames)
}

func (n *RootViewsHandler) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	baseUrl := n.LilBattleApp.BaseUrl
	if baseUrl == "" {
		baseUrl = "https://lilbattle.com"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, `# LilBattle robots.txt
User-agent: *

# Allow public pages
Allow: /worlds/
Allow: /games/
Allow: /privacy/
Allow: /terms/
Allow: /about
Allow: /contact

# Disallow private/dynamic pages
Disallow: /login
Disallow: /profile
Disallow: /auth/
Disallow: /api/
Disallow: /views/
Disallow: /worlds/new
Disallow: /worlds/create
Disallow: /worlds/select
Disallow: /worlds/*/edit
Disallow: /games/new
Disallow: /games/*/edit

# Sitemap
Sitemap: %s/sitemap.xml
`, baseUrl)
}

func (n *RootViewsHandler) handleSitemapIndex(w http.ResponseWriter, r *http.Request) {
	baseUrl := n.LilBattleApp.BaseUrl
	if baseUrl == "" {
		baseUrl = "https://lilbattle.com"
	}

	today := time.Now().Format("2006-01-02")

	index := SitemapIndex{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: []SitemapIndexEntry{
			{Loc: baseUrl + "/sitemap-static.xml", LastMod: today},
			{Loc: baseUrl + "/sitemap-worlds.xml", LastMod: today},
			{Loc: baseUrl + "/sitemap-games.xml", LastMod: today},
		},
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(index)
}

func (n *RootViewsHandler) handleSitemapStatic(w http.ResponseWriter, r *http.Request) {
	baseUrl := n.LilBattleApp.BaseUrl
	if baseUrl == "" {
		baseUrl = "https://lilbattle.com"
	}

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []SitemapURL{
			{Loc: baseUrl + "/worlds/", ChangeFreq: "daily", Priority: "1.0"},
			{Loc: baseUrl + "/games/", ChangeFreq: "daily", Priority: "0.9"},
			{Loc: baseUrl + "/privacy/", ChangeFreq: "monthly", Priority: "0.3"},
			{Loc: baseUrl + "/terms/", ChangeFreq: "monthly", Priority: "0.3"},
			{Loc: baseUrl + "/about", ChangeFreq: "monthly", Priority: "0.5"},
			{Loc: baseUrl + "/contact", ChangeFreq: "monthly", Priority: "0.5"},
		},
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(urlset)
}

func (n *RootViewsHandler) handleSitemapWorlds(w http.ResponseWriter, r *http.Request) {
	baseUrl := n.LilBattleApp.BaseUrl
	if baseUrl == "" {
		baseUrl = "https://lilbattle.com"
	}

	ctx := n.LilBattleApp
	client := ctx.ClientMgr.GetWorldsSvcClient()

	// Fetch all public worlds
	resp, err := client.ListWorlds(GrpcAuthContext(""), &protos.ListWorldsRequest{
		Pagination: &protos.Pagination{PageSize: 1000},
	})

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []SitemapURL{},
	}

	if err == nil && resp != nil {
		for _, world := range resp.Items {
			url := SitemapURL{
				Loc:        fmt.Sprintf("%s/worlds/%s/view", baseUrl, world.Id),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			}
			if world.UpdatedAt != nil {
				url.LastMod = world.UpdatedAt.AsTime().Format("2006-01-02")
			}
			urlset.URLs = append(urlset.URLs, url)
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(urlset)
}

func (n *RootViewsHandler) handleSitemapGames(w http.ResponseWriter, r *http.Request) {
	baseUrl := n.LilBattleApp.BaseUrl
	if baseUrl == "" {
		baseUrl = "https://lilbattle.com"
	}

	ctx := n.LilBattleApp
	client := ctx.ClientMgr.GetGamesSvcClient()

	// Fetch all public games
	resp, err := client.ListGames(GrpcAuthContext(""), &protos.ListGamesRequest{
		Pagination: &protos.Pagination{PageSize: 1000},
	})

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []SitemapURL{},
	}

	if err == nil && resp != nil {
		for _, game := range resp.Items {
			url := SitemapURL{
				Loc:        fmt.Sprintf("%s/games/%s", baseUrl, game.Id),
				ChangeFreq: "weekly",
				Priority:   "0.6",
			}
			if game.UpdatedAt != nil {
				url.LastMod = game.UpdatedAt.AsTime().Format("2006-01-02")
			}
			urlset.URLs = append(urlset.URLs, url)
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(urlset)
}
