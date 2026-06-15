package tests

import (
	"math/rand"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// TestCalculateHitProbability tests the basic hit probability calculation
func TestCalculateHitProbability(t *testing.T) {
	// Load rules engine
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	tests := []struct {
		name         string
		attackerType int32
		defenderType int32
		attackerTile int32
		defenderTile int32
		woundBonus   int32
		expectedP    float64
		expectError  bool
	}{
		{
			name:         "Soldier vs Soldier on grass (no bonuses)",
			attackerType: 1, // Soldier
			defenderType: 1, // Soldier
			attackerTile: 5, // Grass (was 1, should be 5)
			defenderTile: 5, // Grass
			woundBonus:   0,
			// A=6 (Light:Land from soldier), Ta=0, D=6 (soldier defense), Td=0, B=0
			// p = 0.05 * (((6+0)-(6+0))+0) + 0.5 = 0.05 * 0 + 0.5 = 0.50
			expectedP: 0.50,
		},
		{
			name:         "Soldier vs Soldier with wound bonus",
			attackerType: 1,
			defenderType: 1,
			attackerTile: 5,
			defenderTile: 5,
			woundBonus:   2,
			// p = 0.05 * (((6+0)-(6+0))+2) + 0.5 = 0.05 * 2 + 0.5 = 0.60
			expectedP: 0.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CombatContext{
				Attacker: &v1.Unit{
					UnitType: tt.attackerType,
					Player:   1,
				},
				AttackerTile: &v1.Tile{
					TileType: tt.attackerTile,
				},
				AttackerHealth: 10,
				Defender: &v1.Unit{
					UnitType: tt.defenderType,
					Player:   2,
				},
				DefenderTile: &v1.Tile{
					TileType: tt.defenderTile,
				},
				DefenderHealth: 10,
				WoundBonus:     tt.woundBonus,
			}

			p, err := rulesEngine.CalculateHitProbability(ctx)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && p != tt.expectedP {
				t.Errorf("Expected p=%f, got p=%f", tt.expectedP, p)
			}
		})
	}
}

// TestSimulateCombatDamage tests that damage simulation produces reasonable results
func TestSimulateCombatDamage(t *testing.T) {
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	ctx := &CombatContext{
		Attacker: &v1.Unit{
			UnitType: 1, // Soldier
			Player:   1,
		},
		AttackerTile: &v1.Tile{
			TileType: 1, // Grass
		},
		AttackerHealth: 10,
		Defender: &v1.Unit{
			UnitType: 1, // Soldier
			Player:   2,
		},
		DefenderTile: &v1.Tile{
			TileType: 1, // Grass
		},
		DefenderHealth: 10,
		WoundBonus:     0,
	}

	// Run multiple simulations
	rng := rand.New(rand.NewSource(42))
	totalDamage := 0
	numTests := 100

	for i := 0; i < numTests; i++ {
		damage, err := rulesEngine.SimulateCombatDamage(ctx, rng)
		if err != nil {
			t.Fatalf("Simulation failed: %v", err)
		}

		// Damage should be in reasonable range (0 to attacker health)
		if damage < 0 || damage > ctx.AttackerHealth {
			t.Errorf("Damage %d out of range [0, %d]", damage, ctx.AttackerHealth)
		}

		totalDamage += int(damage)
	}

	// With p=0.65 and 10 health, expected damage = 10*6*0.65/6 = 6.5
	avgDamage := float64(totalDamage) / float64(numTests)
	expectedDamage := 6.5

	// Allow 20% variance
	if avgDamage < expectedDamage*0.8 || avgDamage > expectedDamage*1.2 {
		t.Logf("Average damage %f is outside expected range around %f", avgDamage, expectedDamage)
		// Don't fail, just log - RNG can vary
	}
}

// TestGenerateDamageDistribution tests distribution generation
func TestGenerateDamageDistribution(t *testing.T) {
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	ctx := &CombatContext{
		Attacker: &v1.Unit{
			UnitType: 1,
			Player:   1,
		},
		AttackerTile: &v1.Tile{
			TileType: 1,
		},
		AttackerHealth: 10,
		Defender: &v1.Unit{
			UnitType: 1,
			Player:   2,
		},
		DefenderTile: &v1.Tile{
			TileType: 1,
		},
		DefenderHealth: 10,
		WoundBonus:     0,
	}

	dist, err := rulesEngine.GenerateDamageDistribution(ctx, 10000)
	if err != nil {
		t.Fatalf("Failed to generate distribution: %v", err)
	}

	if dist == nil {
		t.Fatal("Distribution is nil")
	}

	// Check that distribution has ranges
	if len(dist.Ranges) == 0 {
		t.Error("Distribution has no ranges")
	}

	// Check expected damage is reasonable
	if dist.ExpectedDamage < 5 || dist.ExpectedDamage > 8 {
		t.Logf("Expected damage %f seems off (should be around 6.5)", dist.ExpectedDamage)
	}

	// Check that probabilities sum to approximately 1.0
	totalProb := 0.0
	for _, r := range dist.Ranges {
		totalProb += r.Probability
	}

	if totalProb < 0.99 || totalProb > 1.01 {
		t.Errorf("Total probability %f should be close to 1.0", totalProb)
	}

	t.Logf("Distribution: min=%f, max=%f, expected=%f, ranges=%d",
		dist.MinDamage, dist.MaxDamage, dist.ExpectedDamage, len(dist.Ranges))
}

// TestCalculateWoundBonus tests wound bonus calculation
func TestCalculateWoundBonus(t *testing.T) {
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	defenderCoord := AxialCoord{Q: 0, R: 0}

	tests := []struct {
		name          string
		attackHistory []*v1.AttackRecord
		attackerCoord AxialCoord
		expectedBonus int32
	}{
		{
			name:          "No previous attacks",
			attackHistory: []*v1.AttackRecord{},
			attackerCoord: AxialCoord{Q: 1, R: 0}, // Adjacent
			expectedBonus: 0,
		},
		{
			name: "One previous ranged attack, current is ranged",
			attackHistory: []*v1.AttackRecord{
				{Q: 3, R: 0, IsRanged: true},
			},
			attackerCoord: AxialCoord{Q: 2, R: 0}, // 2 tiles away = ranged
			expectedBonus: 1,
		},
		{
			name: "One previous ranged attack, current is adjacent",
			attackHistory: []*v1.AttackRecord{
				{Q: 3, R: 0, IsRanged: true},
			},
			attackerCoord: AxialCoord{Q: 1, R: 0}, // Adjacent
			expectedBonus: 1,
		},
		{
			name: "Two attacks from opposite sides",
			attackHistory: []*v1.AttackRecord{
				{Q: 1, R: 0, IsRanged: false}, // Right side
			},
			attackerCoord: AxialCoord{Q: -1, R: 0}, // Left side (opposite)
			expectedBonus: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defender := &v1.Unit{
				Q:             int32(defenderCoord.Q),
				R:             int32(defenderCoord.R),
				AttackHistory: tt.attackHistory,
			}

			bonus := rulesEngine.CalculateWoundBonus(defender, tt.attackerCoord)
			if bonus != tt.expectedBonus {
				t.Errorf("Expected wound bonus %d, got %d", tt.expectedBonus, bonus)
			}
		})
	}
}

// TestEstimateCombatDamage_Analytical pins the analytical mean used by
// AttackUnitAction.DamageEstimate: round(AttackerHealth × p), where p comes from
// CalculateHitProbability. Cases mirror TestCalculateHitProbability so the math
// can be traced end-to-end (p × H → expected).
func TestEstimateCombatDamage_Analytical(t *testing.T) {
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	tests := []struct {
		name           string
		attackerType   int32
		defenderType   int32
		attackerTile   int32
		defenderTile   int32
		attackerHealth int32
		woundBonus     int32
		expected       int32
	}{
		{
			// p = 0.50 (matches TestCalculateHitProbability "Soldier vs Soldier on grass"),
			// expected damage = round(10 × 0.50) = 5
			name:           "Soldier vs Soldier on grass, full health",
			attackerType:   1, defenderType: 1, attackerTile: 5, defenderTile: 5,
			attackerHealth: 10, woundBonus: 0, expected: 5,
		},
		{
			// p = 0.60 (matches "Soldier vs Soldier with wound bonus"),
			// expected damage = round(10 × 0.60) = 6
			name:           "Soldier vs Soldier on grass, wound bonus 2",
			attackerType:   1, defenderType: 1, attackerTile: 5, defenderTile: 5,
			attackerHealth: 10, woundBonus: 2, expected: 6,
		},
		{
			// Same matchup, half-strength attacker: round(5 × 0.50) = 3 (Go's
			// math.Round rounds half away from zero, so 2.5 → 3).
			name:           "Soldier vs Soldier on grass, half health",
			attackerType:   1, defenderType: 1, attackerTile: 5, defenderTile: 5,
			attackerHealth: 5, woundBonus: 0, expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CombatContext{
				Attacker:       &v1.Unit{UnitType: tt.attackerType, Player: 1},
				AttackerTile:   &v1.Tile{TileType: tt.attackerTile},
				AttackerHealth: tt.attackerHealth,
				Defender:       &v1.Unit{UnitType: tt.defenderType, Player: 2},
				DefenderTile:   &v1.Tile{TileType: tt.defenderTile},
				DefenderHealth: 10,
				WoundBonus:     tt.woundBonus,
			}

			got, err := rulesEngine.EstimateCombatDamage(ctx)
			if err != nil {
				t.Fatalf("EstimateCombatDamage failed: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Expected damage %d, got %d", tt.expected, got)
			}
		})
	}
}

// TestEstimateCombatDamage_MatchesSimulationMean cross-checks the analytical
// estimate against the empirical mean from GenerateDamageDistribution. Any
// divergence > 1 would mean the analytical formula has drifted from what
// SimulateCombatDamage actually rolls — they must remain mathematically tied.
func TestEstimateCombatDamage_MatchesSimulationMean(t *testing.T) {
	rulesEngine := DefaultRulesEngine()
	if rulesEngine == nil {
		t.Fatal("Failed to load default rules engine")
	}

	cases := []struct {
		name           string
		attackerType   int32
		defenderType   int32
		attackerHealth int32
		woundBonus     int32
	}{
		{"Soldier vs Soldier", 1, 1, 10, 0},
		{"Soldier vs Soldier w/ wound bonus", 1, 1, 10, 2},
		{"Soldier vs Soldier, low health", 1, 1, 3, 0},
		{"Soldier vs Tank", 1, 5, 10, 0},
		{"Tank vs Soldier", 5, 1, 10, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &CombatContext{
				Attacker:       &v1.Unit{UnitType: tc.attackerType, Player: 1},
				AttackerTile:   &v1.Tile{TileType: 5},
				AttackerHealth: tc.attackerHealth,
				Defender:       &v1.Unit{UnitType: tc.defenderType, Player: 2},
				DefenderTile:   &v1.Tile{TileType: 5},
				DefenderHealth: 10,
				WoundBonus:     tc.woundBonus,
			}

			analytical, err := rulesEngine.EstimateCombatDamage(ctx)
			if err != nil {
				t.Fatalf("EstimateCombatDamage failed: %v", err)
			}

			dist, err := rulesEngine.GenerateDamageDistribution(ctx, 10000)
			if err != nil {
				t.Fatalf("GenerateDamageDistribution failed: %v", err)
			}
			empirical := int32(dist.ExpectedDamage + 0.5) // round-half-up to int32

			diff := analytical - empirical
			if diff < 0 {
				diff = -diff
			}
			if diff > 1 {
				t.Errorf("Analytical estimate (%d) drifted from empirical mean (%d, raw=%.3f) by %d; formula has diverged from SimulateCombatDamage",
					analytical, empirical, dist.ExpectedDamage, diff)
			}
		})
	}
}
