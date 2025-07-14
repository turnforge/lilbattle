package services

import (
	"context"
	"errors"
	"log"

	protos "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const APP_ID = "weewar"

var ErrNoSuchEntity = errors.New("entity not found")

type ClientMgr struct {
	svcAddr           string
	appitemsSvcClient protos.AppItemsServiceClient
	authSvc           *AuthService
	// We may need an auth svc at some point
}

func NewClientMgr(svc_addr string) *ClientMgr {
	log.Println("Client Mgr Svc Addr: ", svc_addr)
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

// We will have one client per service here
func (c *ClientMgr) GetAppItemsSvcClient() (out protos.AppItemsServiceClient, err error) {
	if c.appitemsSvcClient == nil {
		log.Println("Addr: ", c.svcAddr)
		appitemsSvcConn, err := grpc.NewClient(c.svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("cannot connect with server %v", err)
			return nil, err
		}

		c.appitemsSvcClient = protos.NewAppItemsServiceClient(appitemsSvcConn)
	}
	return c.appitemsSvcClient, nil
}

func (c *ClientMgr) GetAuthService() *AuthService {
	if c.authSvc == nil {
		c.authSvc = &AuthService{clients: c}
	}
	return c.authSvc
}
