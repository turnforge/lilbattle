package tests

import (
	"github.com/turnforge/weewar/services"
	"github.com/turnforge/weewar/services/fsbe"
	"github.com/turnforge/weewar/services/singleton"
)

type AxialCoord = services.AxialCoord
type CombatContext = services.CombatContext
type Game = services.Game
type World = services.World
type MoveProcessor = services.MoveProcessor
type SingletonGamesService = singleton.SingletonGamesService

var DevDataPath = fsbe.DevDataPath
var UnitSetCoord = services.UnitSetCoord
var NewWorld = services.NewWorld
var NewGame = services.NewGame
var NewTile = services.NewTile
var NewUnit = services.NewUnit
var CubeDistance = services.CubeDistance
var NewSingletonGamesService = singleton.NewSingletonGamesService
var DefaultRulesEngine = services.DefaultRulesEngine
var ParseActionAlternatives = services.ParseActionAlternatives
var LoadRulesEngineFromFile = services.LoadRulesEngineFromFile
