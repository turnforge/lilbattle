package server

import (
	"context"
	"log"
	"net/http"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

type UserListView struct {
	Users  []*protos.User
	Paginator Paginator
}

func (g *UserListView) Copy() View { return &UserListView{} }

func (p *UserListView) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	userId := vc.AuthMiddleware.GetLoggedInUserId(r)

	// if we are an independent view then read its params from the query params
	// otherwise those will be passed in
	_, _ = p.Paginator.Load(r, w, vc)

	client, _ := vc.ClientMgr.GetUsersSvcClient()

	req := protos.ListUsersRequest{
		Pagination: &protos.Pagination{
			PageOffset: int32(p.Paginator.CurrentPage * p.Paginator.PageSize),
			PageSize:   int32(p.Paginator.PageSize),
		},
		OwnerId: userId,
		// CollectionId: p.CollectionId,
	}
	resp, err := client.ListUsers(context.Background(), &req)
	if err != nil {
		log.Println("error getting notations: ", err)
		return err, false
	}
	log.Println("Found Users: ", resp.Items)
	p.Users = resp.Items
	p.Paginator.HasPrevPage = p.Paginator.CurrentPage > 0
	if resp.Pagination != nil {
		p.Paginator.HasNextPage = resp.Pagination.HasMore
		p.Paginator.EvalPages(p.Paginator.CurrentPage*p.Paginator.PageSize + int(resp.Pagination.TotalResults))
	}
	return nil, false
}
