# LilBattle Developer Guide

A guide for developing, testing, and running the LilBattle turn-based strategy game.

## Quick Start (zero configuration)

A fresh clone runs against local filesystem storage with username/password auth. No Postgres, no OAuth setup, no `.env` file required.

```bash
# Clone and setup
git clone <repository-url>
cd lilbattle

# Install dependencies
go mod download
cd web && pnpm install && cd ..

# Generate proto code
cd protos && make && cd ..

# Terminal 1: Backend (local FS, no env file needed)
make servelocal

# Terminal 2: Frontend build (watches for changes)
cd web && pnpm run watch
```

Open browser at `http://localhost:8080`. Sign up with any email and password — email verification is off by default, so you're logged in immediately.

**What you get:**

- Games / worlds / users stored under `~/dev-app-data/lilbattle/storage/`.
- Local username/password auth via `/login` and `/signup`.
- Email flows (verification, password reset) print to the server console instead of sending real messages.
- No OAuth buttons on the login page until you configure providers.
- No database or cloud service needed.

**Multi-user testing.** Signup is instant, so a second identity is one browser incognito window and a different signup away. You can also enable the `?dev_user=<handle>` shortcut for impersonation without going through signup at all — see "Dev-mode fake login" below.

## Configuring beyond defaults

Copy `configs/.env.example` to `configs/.env.dev` (dev) or `configs/.env` (production) and uncomment the variables you want to change. The example file has one comment block per concern (auth, backends, email, OAuth, filestore, feature flags). Every variable is optional in dev mode — omitting the whole file falls back to the same defaults you got from Quick Start.

### Custom JWT secret

For anything beyond throwaway experiments, override the CLI-token signing secret. The default is a shared dev value baked into the code, fine for local play but no good for real deployments.

```bash
# configs/.env.dev
JWT_CLI_SECRET=$(openssl rand -base64 32)
```

### OAuth providers (optional)

Google, GitHub, and X/Twitter logins each need a registered OAuth app with `http://localhost:8080/auth/<provider>/callback/` as an allowed redirect. Add the credentials to `configs/.env.dev`:

```bash
OAUTH2_GOOGLE_CLIENT_ID=...
OAUTH2_GOOGLE_CLIENT_SECRET=...
OAUTH2_GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback/
```

The login page conditionally shows a provider button only when its credentials are configured. Twitter requires PKCE; the built-in handler at `web/server/twitter_oauth2.go` covers it.

### Postgres backend

For gameplay against a real database instead of the filesystem, spin up a local Postgres and point the backends at it:

```bash
# 1. Start Postgres (docker-compose has a preconfigured service)
make up

# 2. Set backend + endpoint in configs/.env.dev
GAMES_SERVICE_BE=pg
WORLDS_SERVICE_BE=pg
LILBATTLE_DB_ENDPOINT=postgres://postgres:password@localhost:5432/lilbattledb

# 3. Restart the server
make servepg
```

`make servepg` passes `-games_service_be=pg -worlds_service_be=pg` on the command line for you.

### GAE / Datastore backend

For running against Google Cloud Datastore (App Engine / Firestore-in-Datastore mode):

```bash
# configs/.env.dev
GAMES_SERVICE_BE=gae
WORLDS_SERVICE_BE=gae
GAE_PROJECT=your-gcp-project
GAE_NAMESPACE=your-namespace
GOOGLE_APPLICATION_CREDENTIALS=./configs/service-account.json
```

Then `make servegae`. Requires a service account JSON with Datastore access.

### R2 / S3 filestore

For serving world screenshots / uploaded assets from Cloudflare R2 (or any S3-compatible bucket) instead of local disk:

```bash
# configs/.env.dev
FILESTORE_BE=r2
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
```

The `local` filestore backend (default) writes under `~/dev-app-data/lilbattle/storage/files/`.

### Real email (Resend)

Verification and password-reset emails print to the console by default. To send real email in dev, set:

```bash
RESEND_API_KEY=...
RESEND_FROM_EMAIL=LilBattle <noreply@yourdomain.com>
```

Empty `RESEND_API_KEY` keeps the console sender.

### Production deployment

Production mode is opt-in via `LILBATTLE_ENV=production`. It:

- Loads `configs/.env` instead of `configs/.env.dev`.
- Defaults backends to `pg` (must be overridden if you want anything else).
- Requires `configs/.env` to exist — missing it is fatal (dev mode logs a warning and continues).

Deploy via `make deploy` (GAE) after populating `configs/.env` with production secrets.

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

### Recorded replay harness

`tests/e2e/` runs the committed `.sh` replay scripts under
`tests/e2etests/` against an **already-running server**. The harness
assumes the worlds each script references already exist on the target.
Seeding is a separate step, handled by `scripts/seed-worlds.sh` (which
drives `ww worlds ensure`) — see the "Seeding worlds" subsection below.

The harness targets a URL, never a filesystem — same code path against
local dev, staging, or any other server the CI can reach.

Gated behind the `e2e` build tag while the recorded scripts drift from
current game rules (tracked in issue 183). Default `go test ./...`
compiles the `doc.go` stub and reports "no tests to run"; the harness
only runs when requested.

Quick paths:

```bash
# One command — boots a local dev server, runs tests, tears down.
make cli && make e2e-full

# Target one replay:
make cli && make e2e-full ARGS='-run TestReplayScripts/29146'

# Against a server you already have running (faster iteration):
export LILBATTLE_E2E_SERVER=http://localhost:8090/api
make e2e
make e2e-run REPLAY=29146
make e2e-watch          # auto-opens the game URL in a browser
```

Overrides:

- `LILBATTLE_E2E_SERVER` (required for `make e2e*`) — target server URL,
  including the `/api` suffix. Skipped by `make e2e-full`, which sets it
  for you after starting a background server.
- `LILBATTLE_WW_BIN=/path/to/ww` — override which ww binary the scripts
  invoke. Defaults to the first `ww` on `PATH`.
- `LILBATTLE_E2E_WATCH=true` — after `ww new`, invoke `open` (macOS) or
  `xdg-open` (Linux) on the game's viewer URL so you can watch the replay
  play out in a browser tab.
- `LILBATTLE_E2E_HTTP_PORT` / `LILBATTLE_E2E_GRPC_PORT` (for `e2e-full`) —
  override the server's ports if `:8090` / `:9091` conflict.

Fixture worlds live under `tests/e2e/fixtures/worlds/<worldID>/` as
`data.json` + `metadata.json` (both protojson). The scripts themselves
are generated in the sibling `weemaps` repo from upstream game dumps —
see `weemaps/scripts/history.py`. To add a new replay, generate the
`.sh` there, drop it into `tests/e2etests/`, and commit the world's
fixture files under `tests/e2e/fixtures/worlds/`.

#### Seeding worlds

Worlds get onto the target server via `ww worlds ensure`, which:

- Probes the target with `GetWorld`.
- Missing world + fixture supplied → creates the world.
- Existing world + fixture supplied → hash-compares content
  (tiles + units). Mismatch is a **hard error**, never an overwrite —
  operator picks the fix (update the target, update the fixture, or
  use a fresh world ID).
- Existing world + no fixture → success (basic presence probe).

For the whole replay-fixture set:

```bash
export LILBATTLE_SERVER=http://localhost:8090/api
bash scripts/seed-worlds.sh
```

`make e2e-full` calls this automatically after the server comes up.
For staging / prod, seeding is the operator's responsibility — the
tests will fail loudly if a required world is missing.

The Go entry point (`lib.EnsureWorldExists`, `lib.HashWorldData`) is
callable from any Go program that needs the same guarantees.

#### Manually driving a replay for drift diagnosis

When a recorded script fails, `make e2e-full` reports the line that
failed but doesn't tell you WHY the state at that point diverged from
the recording. This recipe drives one script by hand so you can inspect
state between moves:

```bash
# 1. Server on 8090 with DISABLE_API_AUTH so no login is needed.
LILBATTLE_WEB_PORT=:8090 LILBATTLE_GRPC_PORT=:9091 DISABLE_API_AUTH=true \
    go run main.go -games_service_be=local -worlds_service_be=local

# 2. In another shell — seed the world the script needs.
export LILBATTLE_SERVER=http://localhost:8090/api
ww worlds ensure 7e5016a4 --data-dir tests/e2e/fixtures/worlds/7e5016a4/
# (or run scripts/seed-worlds.sh for all fixtures)

# 3. Create a game, note the ID from the `export` line ww prints.
ww new 7e5016a4

# 4. Drive it manually, inspecting between moves.
export LILBATTLE_GAME_ID=<id>
export LILBATTLE_CONFIRM=false
ww status                  # turn, current player, coins per player
ww tiles                   # owned tiles per player
ww units                   # units per player
ww options t:0,3           # what can I do at this tile? (empty = not buildable / not owned)
ww options A1              # what can unit A1 do? (moves, attacks, captures)
ww map                     # visual snapshot (inline iTerm2 image or PNG)
ww build t:2,1 1           # try a specific move
ww move A1 R               # move unit A1 right
```

Inspect the raw state on disk any time:

```bash
jq '{current_player, turn_counter, finished}' \
    ~/dev-app-data/lilbattle/storage/games/<id>/state.json

# Tiles owned by player 1
jq '.world_data.tiles_map | to_entries[] | select(.value.player == 1)' \
    ~/dev-app-data/lilbattle/storage/games/<id>/state.json

# All units, with their shortcuts and remaining movement
jq '.world_data.units_map[] | {q, r, player, shortcut, health: .available_health, moves: .distance_left}' \
    ~/dev-app-data/lilbattle/storage/games/<id>/state.json
```

Common drift signatures:

- `ww build` rejected with "tile does not belong to player X" → the
  recording assumed the tile was captured earlier in the run; walk
  backwards through the `.sh` to find the capture that didn't fire.
- `ww assert unit ... [health eq N]` off by a few HP → combat math has
  changed since recording. Cross-check against `lib/combat.go` history.
- `ww capture` fails silently or produces wrong ownership → capture
  eligibility rules changed (which unit types can capture, or which
  terrain is capturable).

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
