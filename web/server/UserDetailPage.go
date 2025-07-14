package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type UserDetailPage struct {
	BasePage
	Header  Header
	User *protos.User
	UserId string
}

func (p *UserDetailPage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.UserId = r.PathValue("appItemId")
	if p.UserId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return nil, true
	}

	p.Title = "User Details"
	p.Header.Load(r, w, vc)

	// Fetch the User using the client manager
	client, err := vc.ClientMgr.GetUsersSvcClient()
	if err != nil {
		log.Printf("Error getting Users client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, true
	}

	req := &protos.GetUserRequest{
		Id: p.UserId,
	}

	resp, err := client.GetUser(context.Background(), req)
	if err != nil {
		log.Printf("Error fetching User %s: %v", p.UserId, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return nil, true
	}

	if resp.Appitem != nil {
		// Convert from UserProject to User (assuming we need the basic info)
		p.User = &protos.User{
			Id:          resp.Appitem.Id,
			Name:        resp.Appitem.Name,
			Description: resp.Appitem.Description,
			CreatedAt:   resp.Appitem.CreatedAt,
			UpdatedAt:   resp.Appitem.UpdatedAt,
		}
		p.Title = p.User.Name
	}

	return nil, false
}

func (p *UserDetailPage) Copy() View { 
	return &UserDetailPage{} 
}