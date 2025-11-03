//go:build !wasm
// +build !wasm

// This file is excluded from WASM builds.
// It contains gRPC client code that requires net/http packages
// which are not supported by TinyGo's WASM target.

package server

import (
	"context"
	"errors"
	"log"

	v1s "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1/services"
	"github.com/panyam/turnengine/games/weewar/services/fsbe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const APP_ID = "weewar"

var ErrNoSuchEntity = errors.New("entity not found")

type ClientMgr struct {
	svcAddr         string
	gamesSvcClient  v1s.GamesServiceClient
	worldsSvcClient v1s.WorldsServiceClient
	authSvc         *fsbe.AuthService
	// We may need an auth svc at some point
}

func NewClientMgr(svc_addr string) *ClientMgr {
	if svc_addr == "" {
		panic("Service Address is nil")
	}
	return &ClientMgr{svcAddr: svc_addr}
}

func (c *ClientMgr) Address() string {
	return c.svcAddr
}

func (c *ClientMgr) ClientContext(ctx context.Context, loggedInUserId string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return metadata.AppendToOutgoingContext(context.Background(), "LoggedInUserId", loggedInUserId)
}

func (c *ClientMgr) GetAuthService() *fsbe.AuthService {
	if c.authSvc == nil {
		c.authSvc = &fsbe.AuthService{
			// clients: c
		}
	}
	return c.authSvc
}

// We will have one client per service here
func (c *ClientMgr) GetWorldsSvcClient() (out v1s.WorldsServiceClient, err error) {
	if c.worldsSvcClient == nil {
		worldsSvcConn, err := grpc.NewClient(c.svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("cannot connect with server %v", err)
			return nil, err
		}

		c.worldsSvcClient = v1s.NewWorldsServiceClient(worldsSvcConn)
	}
	return c.worldsSvcClient, nil
}

// We will have one client per service here
func (c *ClientMgr) GetGamesSvcClient() (out v1s.GamesServiceClient, err error) {
	if c.gamesSvcClient == nil {
		gamesSvcConn, err := grpc.NewClient(c.svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("cannot connect with server %v", err)
			return nil, err
		}

		c.gamesSvcClient = v1s.NewGamesServiceClient(gamesSvcConn)
	}
	return c.gamesSvcClient, nil
}
