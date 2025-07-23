package weewar

import "encoding/json"

type UnitStats struct {
	Cost       int  `json:"cost"`
	Health     int  `json:"health"`
	Movement   int  `json:"movement"`
	Attack     int  `json:"attack"`
	Defense    int  `json:"defense"`
	SightRange int  `json:"sightRange"`
	CanCapture bool `json:"canCapture"`
}

type UnitData struct {
	BaseStats      UnitStats `json:"baseStats"`
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	MovementPoints int       `json:"movementPoints"`
	AttackRange    int       `json:"attackRange"`
	Health         int       `json:"health"`
	Properties     []string  `json:"properties,omitempty"`
	// Note: Movement costs and attack data managed separately by RulesEngine
}

// DamageDistribution represents combat damage in canonical format
type DamageDistribution struct {
	MinDamage      int            `json:"minDamage"`
	MaxDamage      int            `json:"maxDamage"`
	DamageBuckets  []DamageBucket `json:"damageBuckets"`
	ExpectedDamage float64        `json:"expectedDamage"`
}

// DamageBucket represents a damage value with its probability weight
type DamageBucket struct {
	Damage int     `json:"damage"`
	Weight float64 `json:"weight"`
}

// Unit represents a runtime unit instance in the game
type Unit struct {
	UnitType int // Reference to UnitData by ID

	// Runtime state
	DistanceLeft    int // Movement points remaining this turn
	AvailableHealth int // Current health points
	TurnCounter     int // Which turn this unit was created/last acted

	// Position on the map
	Coord AxialCoord `json:"coord"` // Cube coordinate position

	// Player ownership
	Player int
}

// MarshalJSON implements custom JSON marshaling for Unit
func (t *Unit) MarshalJSON() ([]byte, error) {
	// Convert cube map to tile list for JSON
	out := map[string]any{
		"q":         t.Coord.Q,
		"r":         t.Coord.R,
		"unit_type": t.UnitType,
		"player":    t.Player,
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements custom JSON unmarshaling for Tiile
func (t *Unit) UnmarshalJSON(data []byte) error {
	// First try to unmarshal with new bounds format
	type mapJSON struct {
		Q        int `json:"q"`
		R        int `json:"r"`
		UnitType int `json:"unit_type"`
		Player   int `json:"player"`
	}

	var dict mapJSON

	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}

	t.Coord = AxialCoord{dict.Q, dict.R}
	t.UnitType = dict.UnitType
	t.Player = dict.Player
	return nil
}

// GetPosition returns the unit's cube coordinate position
func (u *Unit) GetPosition() AxialCoord {
	return u.Coord
}

// SetPosition sets the unit's position using cube coordinates
func (u *Unit) SetPosition(coord AxialCoord) {
	u.Coord = coord
}

// NewUnit creates a new unit instance
func NewUnit(unitType, playerID int) *Unit {
	return &Unit{
		UnitType:        unitType,
		Player:          playerID,
		DistanceLeft:    0, // Will be set based on UnitData
		AvailableHealth: 0, // Will be set based on UnitData
		TurnCounter:     0,
	}
}

func (u *Unit) Clone() *Unit {
	return &Unit{
		UnitType:        u.UnitType,
		DistanceLeft:    u.DistanceLeft,
		AvailableHealth: u.AvailableHealth,
		TurnCounter:     u.TurnCounter,
		Coord:           u.Coord,
		Player:          u.Player,
	}
}

// Basic unit data map - matches weewar-data.json
var unitDataMap = map[int]UnitData{
	1:  {ID: 1, Name: "Soldier (Basic)"},
	2:  {ID: 2, Name: "Soldier (Advanced)"},
	3:  {ID: 3, Name: "Tank (Basic)"},
	4:  {ID: 4, Name: "Tank (Advanced)"},
	5:  {ID: 5, Name: "Striker"},
	6:  {ID: 6, Name: "Anti-aircraft (Basic)"},
	7:  {ID: 7, Name: "Hovercraft"},
	8:  {ID: 8, Name: "Artillery (Basic)"},
	9:  {ID: 9, Name: "Artillery (Advanced)"},
	10: {ID: 10, Name: "Speedboat"},
	11: {ID: 11, Name: "Capturing"},
	12: {ID: 12, Name: "Battleship"},
	13: {ID: 13, Name: "Destroyer"},
	14: {ID: 14, Name: "Jetfighter"},
	15: {ID: 15, Name: "Crop Duster"},
	16: {ID: 16, Name: "Engineer"},
	17: {ID: 17, Name: "Medic"},
	18: {ID: 18, Name: "Miner"},
	19: {ID: 19, Name: "Paratrooper"},
	20: {ID: 20, Name: "Helicopter"},
	21: {ID: 21, Name: "Zeppelin"},
	22: {ID: 22, Name: "Submarine"},
	24: {ID: 24, Name: "Anti-aircraft (Advanced)"},
	25: {ID: 25, Name: "Artillery (Mega)"},
	26: {ID: 26, Name: "Artillery (Quick)"},
	27: {ID: 27, Name: "Mech"},
	28: {ID: 28, Name: "Goliath RC"},
	29: {ID: 29, Name: "Stratotanker"},
}

// GetUnitData returns unit data for the given type
func GetUnitData(unitType int) *UnitData {
	if data, exists := unitDataMap[unitType]; exists {
		return &data
	}
	return nil
}
