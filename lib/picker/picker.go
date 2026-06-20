// Package picker selects a single GameOption from a slice of valid options.
//
// Picker is the internal abstraction the presenter uses to implement its
// NextMove method. Drivers (CLI, web, RPC) never see picker — they only see
// the presenter's NextMove surface. This lets the policy swap (Random →
// Heuristic → AI) entirely inside the presenter without touching either
// driver.
package picker

import (
	"math/rand"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// Picker chooses one option from a slice of valid options.
//
// Implementations vary in strategy — random, heuristic, learned/RL agent —
// but the contract is identical: given non-empty input, return a non-nil
// element of the input; given empty input, return nil.
type Picker interface {
	// Pick returns one of the supplied options. Returns nil when options
	// is empty or nil (a signal callers translate into "end turn").
	Pick(options []*v1.GameOption, rng *rand.Rand) *v1.GameOption
}
