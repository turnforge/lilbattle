package assets

import (
	_ "embed"
)

//go:embed weewar-rules.json
var RulesDataJSON []byte
