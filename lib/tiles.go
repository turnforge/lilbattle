package weewar

import (
	"encoding/json"
)

// TerrainType represents whether terrain is nature or player-controllable
type TerrainType int

const (
	TerrainNature TerrainType = iota // Natural terrain (grass, mountains, water, etc.)
	TerrainPlayer                    // Player-controllable structures (bases, cities, etc.)
)

// TerrainData represents terrain type information
type TerrainData struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	BaseMoveCost float64     `json:"baseMoveCost"` // Base movement cost for this terrain
	DefenseBonus float64     `json:"defenseBonus"`
	Type         TerrainType `json:"type"` // Nature or Player terrain
	Properties   []string    `json:"properties,omitempty"`
	// Note: Unit-specific movement costs in RulesEngine can override base cost
}

// Tile represents a single hex tile on the map
type Tile struct {
	Coord AxialCoord `json:"coord"`

	TileType int `json:"tileType"` // Reference to TerrainData by ID

	// Optional: Player this tile belongs to if it is a city tile
	Player int `json:"player"`
}

// NewTile creates a new tile at the specified position
func NewTile(coord AxialCoord, tileType int) *Tile {
	return &Tile{
		Coord:    coord,
		TileType: tileType,
	}
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		Coord:    t.Coord,
		TileType: t.TileType,
		Player:   t.Player, // Units are cloned separately
	}
}

// MarshalJSON implements custom JSON marshaling for Tile
func (t *Tile) MarshalJSON() ([]byte, error) {
	// Convert cube map to tile list for JSON
	out := map[string]any{
		"q":         t.Coord.Q,
		"r":         t.Coord.R,
		"tile_type": t.TileType,
		"player":    t.Player,
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements custom JSON unmarshaling for Tiile
func (t *Tile) UnmarshalJSON(data []byte) error {
	// First try to unmarshal with new bounds format
	type mapJSON struct {
		Q        int `json:"q"`
		R        int `json:"r"`
		TileType int `json:"tile_type"`
		Player   int `json:"player"`
	}

	var dict mapJSON

	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}

	t.Coord = AxialCoord{dict.Q, dict.R}
	t.TileType = dict.TileType
	t.Player = dict.Player
	return nil
}
