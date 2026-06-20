package picker

import (
	"math/rand"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// RandomPicker picks uniformly at random from the supplied options. It is the
// baseline Picker used by autoplay before any heuristic / learned policy
// lands. Deterministic under a fixed RNG seed — same seed, same picks.
type RandomPicker struct{}

// NewRandomPicker returns a stateless uniform Picker.
func NewRandomPicker() *RandomPicker {
	return &RandomPicker{}
}

// Pick returns options[rng.Intn(len(options))]. Returns nil if options is
// empty or nil.
func (p *RandomPicker) Pick(options []*v1.GameOption, rng *rand.Rand) *v1.GameOption {
	if len(options) == 0 {
		return nil
	}
	return options[rng.Intn(len(options))]
}
