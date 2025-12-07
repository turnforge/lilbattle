package server

import (
	"context"

	goal "github.com/panyam/goapplib"
	"github.com/turnforge/weewar/services"
)

type WebAppServer struct {
	goal.WebAppServer
}

func (s *WebAppServer) Start(ctx context.Context, srvErr chan error, stopChan chan bool) error {
	cm := services.NewClientMgr(s.GrpcAddress)
	weewarApp, _, _ := NewWeewarApp(cm)
	return s.StartWithHandler(ctx, weewarApp.Handler(), srvErr, stopChan)
}

type IndexerAppServer struct {
	goal.WebAppServer
}
