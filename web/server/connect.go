package server

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	"github.com/panyam/turnengine/games/weewar/services"
)

// ConnectGamesServiceAdapter adapts the gRPC GamesService to Connect's interface
type ConnectGamesServiceAdapter struct {
	svc *services.GamesServiceImpl
}

func NewConnectGamesServiceAdapter(svc *services.GamesServiceImpl) *ConnectGamesServiceAdapter {
	return &ConnectGamesServiceAdapter{svc: svc}
}

func (a *ConnectGamesServiceAdapter) CreateGame(ctx context.Context, req *connect.Request[v1.CreateGameRequest]) (*connect.Response[v1.CreateGameResponse], error) {
	resp, err := a.svc.CreateGame(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectGamesServiceAdapter) ListGames(ctx context.Context, req *connect.Request[v1.ListGamesRequest]) (*connect.Response[v1.ListGamesResponse], error) {
	resp, err := a.svc.ListGames(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectGamesServiceAdapter) GetGame(ctx context.Context, req *connect.Request[v1.GetGameRequest]) (*connect.Response[v1.GetGameResponse], error) {
	resp, err := a.svc.GetGame(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectGamesServiceAdapter) GetGames(ctx context.Context, req *connect.Request[v1.GetGamesRequest]) (*connect.Response[v1.GetGamesResponse], error) {
	resp, err := a.svc.GetGames(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectGamesServiceAdapter) DeleteGame(ctx context.Context, req *connect.Request[v1.DeleteGameRequest]) (*connect.Response[v1.DeleteGameResponse], error) {
	resp, err := a.svc.DeleteGame(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectGamesServiceAdapter) UpdateGame(ctx context.Context, req *connect.Request[v1.UpdateGameRequest]) (*connect.Response[v1.UpdateGameResponse], error) {
	resp, err := a.svc.UpdateGame(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

/** If you had a streamer than you can use this to act as a bridge between websocket and grpc streams
func (a *ConnectGameServiceAdapter) StreamSomeThing(ctx context.Context, req *connect.Request[v1.StreamSomeThingRequest], stream *connect.ServerStream[v1.StreamSomeThingResponse]) error {
	// Create a custom stream implementation that bridges to Connect
	bridgeStream := &ConnectStreamBridge[v1.StreamSomeThingResponse]{
		connectStream: stream,
		ctx:           ctx,
	}

	// Call your existing gRPC streaming method
	return a.svc.StreamSomeThing(req.Msg, bridgeStream)
}
*/

// ConnectMapsServiceAdapter adapts the gRPC MapsService to Connect's interface
type ConnectMapsServiceAdapter struct {
	svc *services.MapsServiceImpl
}

func NewConnectMapsServiceAdapter(svc *services.MapsServiceImpl) *ConnectMapsServiceAdapter {
	return &ConnectMapsServiceAdapter{svc: svc}
}

func (a *ConnectMapsServiceAdapter) CreateMap(ctx context.Context, req *connect.Request[v1.CreateMapRequest]) (*connect.Response[v1.CreateMapResponse], error) {
	resp, err := a.svc.CreateMap(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectMapsServiceAdapter) ListMaps(ctx context.Context, req *connect.Request[v1.ListMapsRequest]) (*connect.Response[v1.ListMapsResponse], error) {
	resp, err := a.svc.ListMaps(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectMapsServiceAdapter) GetMap(ctx context.Context, req *connect.Request[v1.GetMapRequest]) (*connect.Response[v1.GetMapResponse], error) {
	resp, err := a.svc.GetMap(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectMapsServiceAdapter) GetMaps(ctx context.Context, req *connect.Request[v1.GetMapsRequest]) (*connect.Response[v1.GetMapsResponse], error) {
	resp, err := a.svc.GetMaps(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectMapsServiceAdapter) DeleteMap(ctx context.Context, req *connect.Request[v1.DeleteMapRequest]) (*connect.Response[v1.DeleteMapResponse], error) {
	resp, err := a.svc.DeleteMap(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *ConnectMapsServiceAdapter) UpdateMap(ctx context.Context, req *connect.Request[v1.UpdateMapRequest]) (*connect.Response[v1.UpdateMapResponse], error) {
	resp, err := a.svc.UpdateMap(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

/** If you had a streamer than you can use this to act as a bridge between websocket and grpc streams
func (a *ConnectMapServiceAdapter) StreamSomeThing(ctx context.Context, req *connect.Request[v1.StreamSomeThingRequest], stream *connect.ServerStream[v1.StreamSomeThingResponse]) error {
	// Create a custom stream implementation that bridges to Connect
	bridgeStream := &ConnectStreamBridge[v1.StreamSomeThingResponse]{
		connectStream: stream,
		ctx:           ctx,
	}

	// Call your existing gRPC streaming method
	return a.svc.StreamSomeThing(req.Msg, bridgeStream)
}
*/
