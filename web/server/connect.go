package server

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	"github.com/panyam/turnengine/games/weewar/services"
)

// ConnectAppItemsServiceAdapter adapts the gRPC AppItemsService to Connect's interface
type ConnectAppItemsServiceAdapter struct {
	svc *services.AppItemsServiceImpl
}

func NewConnectAppItemsServiceAdapter(svc *services.AppItemsServiceImpl) *ConnectAppItemsServiceAdapter {
	return &ConnectAppItemsServiceAdapter{svc: svc}
}

func (a *ConnectAppItemsServiceAdapter) CreateAppItem(ctx context.Context, req *connect.Request[v1.CreateAppItemRequest]) (*connect.Response[v1.CreateAppItemResponse], error) {
	resp, err := a.svc.CreateAppItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectAppItemsServiceAdapter) ListAppItems(ctx context.Context, req *connect.Request[v1.ListAppItemsRequest]) (*connect.Response[v1.ListAppItemsResponse], error) {
	resp, err := a.svc.ListAppItems(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectAppItemsServiceAdapter) GetAppItem(ctx context.Context, req *connect.Request[v1.GetAppItemRequest]) (*connect.Response[v1.GetAppItemResponse], error) {
	resp, err := a.svc.GetAppItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectAppItemsServiceAdapter) GetAppItems(ctx context.Context, req *connect.Request[v1.GetAppItemsRequest]) (*connect.Response[v1.GetAppItemsResponse], error) {
	resp, err := a.svc.GetAppItems(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectAppItemsServiceAdapter) DeleteAppItem(ctx context.Context, req *connect.Request[v1.DeleteAppItemRequest]) (*connect.Response[v1.DeleteAppItemResponse], error) {
	resp, err := a.svc.DeleteAppItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectAppItemsServiceAdapter) UpdateAppItem(ctx context.Context, req *connect.Request[v1.UpdateAppItemRequest]) (*connect.Response[v1.UpdateAppItemResponse], error) {
	resp, err := a.svc.UpdateAppItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

/** If you had a streamer than you can use this to act as a bridge between websocket and grpc streams
func (a *ConnectAppItemServiceAdapter) StreamSomeThing(ctx context.Context, req *connect.Request[v1.StreamSomeThingRequest], stream *connect.ServerStream[v1.StreamSomeThingResponse]) error {
	// Create a custom stream implementation that bridges to Connect
	bridgeStream := &ConnectStreamBridge[v1.StreamSomeThingResponse]{
		connectStream: stream,
		ctx:           ctx,
	}

	// Call your existing gRPC streaming method
	return a.svc.StreamSomeThing(req.Msg, bridgeStream)
}
*/
