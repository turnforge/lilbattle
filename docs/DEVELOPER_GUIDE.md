# LilBattle Developer Guide

A guide for developing, testing, and running the LilBattle turn-based strategy game.

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd lilbattle

# Install dependencies
go mod download
cd web && npm install && cd ..

# Generate proto code
buf generate

# Start development (uses devloop for live reload)
devloop

# Or manually:
# Terminal 1: Backend
go run main.go serve

# Terminal 2: Frontend build (watches for changes)
cd web && npm run watch
```

Open browser at `http://localhost:8080`

## Architecture Overview

LilBattle uses a modern web architecture:

```
Browser
├── Phaser.js (WebGL rendering)
├── TypeScript UI Layer
└── WASM Game Logic (Go compiled)
    ↕ gRPC
Go Backend
├── Web Server (Templar templates)
├── Services (gRPC)
└── Rules Engine (data-driven)
```

### Key Components

- **Backend (`services/`)**: Core game logic, move processing, rules engine
- **Frontend (`web/src/`)**: TypeScript pages with Phaser.js rendering
- **Templates (`web/templates/`)**: Templar engine with goapplib integration
- **Protos (`protos/`)**: Protocol Buffers for all data structures
- **CLI (`cmd/cli/`)**: Command-line interface for headless gameplay

### Template System (Templar)

Templates use the templar engine with namespace/include/extend directives:

```html
{{# namespace "lilbattle" #}}
{{# include "goapplib/BasePage.html" #}}
{{# extend "goapplib/BasePage.html" #}}

{{ define "Header" }}
  {{# include "Header.html" #}}
{{ end }}

{{ define "Body" }}
  <!-- Page content -->
{{ end }}
```

Component templates (`.templar.html`) are rendered by presenters for dynamic panels.

## CLI Interface

Build and use the CLI for command-line gameplay:

```bash
# Build CLI
make cli

# Basic commands
export LILBATTLE_GAME_ID=<gameId>

ww status                    # Show game state
ww units                     # List all units
ww options A1                # Show moves for unit A1
ww options t:A1              # Show build options for tile A1
ww move A1 R                 # Move unit right (L/R/TL/TR/BL/BR)
ww move A1 0,-3             # Move to coordinates
ww attack A1 B2             # Attack unit
ww build t:A1 trooper       # Build unit at tile
ww endturn                  # End current turn

# Flags
ww --verbose units          # Debug output
ww --dryrun move A1 R      # Preview without saving
ww --json status            # JSON output
```

### Position Format Support

- **Unit shortcuts**: `A1`, `B2` (references a unit)
- **Q,R coordinates**: `0,-3`, `5,2` (axial hex coordinates)
- **Row,Col coordinates**: `r4,5` (offset coordinates)
- **Direction shortcuts**: `L`, `R`, `TL`, `TR`, `BL`, `BR` (relative)
- **Tile prefix**: `t:A1` (forces tile lookup instead of unit)

## Development with devloop

The `devloop` tool handles continuous builds:

```bash
devloop config              # Get configuration
devloop paths               # List watched file patterns
devloop trigger <rulename>  # Trigger rule execution
devloop logs <rulename>     # Stream logs
devloop status <rulename>   # Get rule status
```

Builds for frontend, WASM, and backend run continuously. Do NOT manually run:
- `npm run build` (web module auto-builds)
- `buf generate` (protos auto-regenerate)

### Pre-push hook

A pre-push hook lives at `.githooks/pre-push` and mirrors the CI test
job. Install it once per clone:

```bash
make setup-hooks
```

That runs `git config core.hooksPath .githooks`. From then on, every
`git push` first runs:

1. `go build` over the production package set (excludes `cmd/wasm`,
   `cmd/repl`, `cmd/indexer`, and `tests/`).
2. `go test` over the CI-covered Go packages.
3. WASM build (`GOOS=js GOARCH=wasm`) — required because
   `web/wasmLoading.test.ts` loads the real binary.
4. `pnpm test` inside `web/`.

Skipped: `pnpm install --frozen-lockfile` and `pnpm run buildprod`.
Stale node_modules surfaces loudly via the jest run; the production
webpack bundle isn't needed for jest's own transpilation.

A failed step aborts the push. Bypass for a deliberate WIP push:
`git push --no-verify`. Use sparingly — CI is the only other gate.

### Dev-mode fake login (`?dev_user=`)

For multi-client testing without registering N real accounts, the server
supports a fake-login query parameter — opt-in via env var. Disabled by
default; the middleware isn't even installed unless the gate is on.

```bash
# Start the dev server with the gate on
ENABLE_DEV_FAKE_LOGIN=true devloop
```

Then open multiple browser windows, each pointing at the same game with
a different identity:

```
http://localhost:8080/?dev_user=alice
http://localhost:8080/?dev_user=bob
```

Each window now acts as a different logged-in user. The middleware writes
the standard session/cookie/JWT triple via `oneauth.SetLoggedInSubject`,
so every downstream auth check (page handlers, gRPC auth, sync) sees the
identity exactly as it would after a real login. Subject sticks via cookie
until you swap with `?dev_user=<other>` or clear cookies.

The handle goes through verbatim as the session subject — meaning a game
created with players `UserId="alice"` and `UserId="bob"` pairs naturally
with two windows opened as those handles.

> **Never enable `ENABLE_DEV_FAKE_LOGIN=true` in production.** The server
> logs a startup banner and a per-event `slog.Warn` line on every
> fake-login request so an accidental prod enablement is impossible to
> miss in logs. The middleware itself is only added to the handler chain
> when the env var is set at server boot — production handlers never
> hold a reference to it.

## Testing

```bash
# All tests
go test ./...

# Specific package with verbose output
go test ./services/ -v

# With coverage
go test ./services/ -cover

# Specific test
go test ./services/ -run TestActionProgression -v
```

## Game Storage Structure

Games stored in `~/dev-app-data/lilbattle/storage/games/{gameId}/`:
- `metadata.json`: Game configuration
- `state.json`: Current game state
- `history.json`: Move history

Worlds stored in `~/dev-app-data/lilbattle/storage/worlds/{worldId}/`:
- `metadata.json`: World metadata
- `world.json`: Map data

### Debugging with jq

```bash
# Check game status
jq '{current_player, turn_counter, status}' ~/dev-app-data/lilbattle/storage/games/{gameId}/state.json

# List units for player
jq '.world_data.units[] | select(.player == 1) | {shortcut, q, r, moves: .distance_left}' state.json

# View recent moves
jq '.groups[-1]' ~/dev-app-data/lilbattle/storage/games/{gameId}/history.json
```

## Key Files

### Services (`services/`)
- `game.go`: Core game state management
- `world.go`: Hex coordinate system, unit/tile operations
- `moves.go`: Move processing and validation
- `rules_engine.go`: Data-driven game mechanics
- `singleton_gameview_presenter.go`: UI update orchestration

### Frontend (`web/src/pages/`)
- `GameViewerPage/`: Interactive game interface (DockView, Grid, Mobile variants)
- `WorldEditorPage/`: Map editor with tools and panels
- `common/`: Shared code (World, PhaserWorldScene, animations)

### Templates (`web/templates/`)
- `BasePage.html`: Base layout extending goapplib
- `*.templar.html`: Component templates for presenter rendering

## Proto Field Naming

Proto fields use snake_case in JSON but camelCase in Go:
- JSON: `available_health`, `distance_left`, `unit_type`
- Go: `AvailableHealth`, `DistanceLeft`, `UnitType`

## Further Documentation

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Detailed technical architecture
- [PROJECT.md](../PROJECT.md) - Current status and achievements
- [ROADMAP.md](./ROADMAP.md) - Development phases
- [ATTACK.md](./ATTACK.md) - Combat mechanics
- [GAMELOG.md](./GAMELOG.md) - Move history system
