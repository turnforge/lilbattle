package picker

import (
	"math/rand"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// TestRandomPicker_EmptyOptions_ReturnsNil pins the empty-input → nil
// contract. Callers use the nil return as their "end turn" signal.
func TestRandomPicker_EmptyOptions_ReturnsNil(t *testing.T) {
	p := NewRandomPicker()
	rng := rand.New(rand.NewSource(1))

	if got := p.Pick(nil, rng); got != nil {
		t.Errorf("Pick(nil) = %v; want nil", got)
	}
	if got := p.Pick([]*v1.GameOption{}, rng); got != nil {
		t.Errorf("Pick([]) = %v; want nil", got)
	}
}

// TestRandomPicker_UniformDistribution runs 10k picks across 4 distinct
// options and asserts each option's empirical frequency is within tolerance
// of the uniform 25% mark. A picker that drifts from uniform (e.g. due to
// off-by-one in indexing) fails here even when individual picks still
// return a valid option.
func TestRandomPicker_UniformDistribution(t *testing.T) {
	const (
		picks     = 10000
		tolerance = 0.03 // ±3 percentage points; well outside binomial noise for n=10k
	)

	options := []*v1.GameOption{
		{OptionType: &v1.GameOption_EndTurn{EndTurn: &v1.EndTurnAction{}}},
		{OptionType: &v1.GameOption_EndTurn{EndTurn: &v1.EndTurnAction{}}},
		{OptionType: &v1.GameOption_EndTurn{EndTurn: &v1.EndTurnAction{}}},
		{OptionType: &v1.GameOption_EndTurn{EndTurn: &v1.EndTurnAction{}}},
	}

	p := NewRandomPicker()
	rng := rand.New(rand.NewSource(42))
	counts := make(map[*v1.GameOption]int, len(options))

	for range picks {
		chosen := p.Pick(options, rng)
		if chosen == nil {
			t.Fatal("Pick returned nil on non-empty input")
		}
		counts[chosen]++
	}

	expectedFreq := 1.0 / float64(len(options))
	for i, opt := range options {
		freq := float64(counts[opt]) / float64(picks)
		if diff := freq - expectedFreq; diff < -tolerance || diff > tolerance {
			t.Errorf("option[%d] frequency %.4f deviates from uniform %.4f by more than %.2f",
				i, freq, expectedFreq, tolerance)
		}
	}
}
