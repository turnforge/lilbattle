package fsbe

import (
	"math"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/panyam/goutils/utils"
)

const WEEWAR_DATA_ROOT = "~/dev-app-data/weewar"

// For dev
func DevDataPath(path string) string {
	return filepath.Join(utils.ExpandUserPath(WEEWAR_DATA_ROOT), path)
}

// shortRandSuffix generates a 4-character random suffix for ID suggestions
func shortRandSuffix() string {
	max_id := int64(math.Pow(36, 4))
	randval := rand.Int63() % max_id
	return strconv.FormatInt(randval, 36)
}
