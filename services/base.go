package services

import v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"

type BaseGamesServiceImpl struct {
	v1.UnimplementedGamesServiceServer
	WorldsService v1.WorldsServiceServer
}

type BaseWorldsServiceImpl struct {
	v1.UnimplementedWorldsServiceServer
}
