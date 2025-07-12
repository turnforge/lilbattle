package weewar

import (
	"github.com/panyam/turnengine/internal/turnengine"
)

type PositionComponent struct {
	X, Y, Z int `json:"x,y,z"`
}

func (p PositionComponent) Type() string { return "position" }

type HealthComponent struct {
	Current, Max int `json:"current,max"`
}

func (h HealthComponent) Type() string { return "health" }

type MovementComponent struct {
	Range     int `json:"range"`
	MovesLeft int `json:"movesLeft"`
}

func (m MovementComponent) Type() string { return "movement" }

type CombatComponent struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
}

func (c CombatComponent) Type() string { return "combat" }

type UnitTypeComponent struct {
	UnitType string `json:"unitType"`
	Cost     int    `json:"cost"`
}

func (u UnitTypeComponent) Type() string { return "unitType" }

type TeamComponent struct {
	TeamID int `json:"teamId"`
}

func (t TeamComponent) Type() string { return "team" }

type TerrainComponent struct {
	TerrainType    string `json:"terrainType"`
	DefenseBonus   int    `json:"defenseBonus"`
	MovementCost   int    `json:"movementCost"`
	CapturePoints  int    `json:"capturePoints"`
	MaxCapture     int    `json:"maxCapture"`
	Owner          int    `json:"owner"`
}

func (t TerrainComponent) Type() string { return "terrain" }

type ProductionComponent struct {
	Income      int      `json:"income"`
	CanProduce  []string `json:"canProduce"`
	Producing   string   `json:"producing"`
	TurnsLeft   int      `json:"turnsLeft"`
}

func (p ProductionComponent) Type() string { return "production" }

type StatusComponent struct {
	Status    string `json:"status"`
	TurnsLeft int    `json:"turnsLeft"`
}

func (s StatusComponent) Type() string { return "status" }

func RegisterWeeWarComponents(registry *turnengine.ComponentRegistry) {
	registry.Register(PositionComponent{})
	registry.Register(HealthComponent{})
	registry.Register(MovementComponent{})
	registry.Register(CombatComponent{})
	registry.Register(UnitTypeComponent{})
	registry.Register(TeamComponent{})
	registry.Register(TerrainComponent{})
	registry.Register(ProductionComponent{})
	registry.Register(StatusComponent{})
}