# Production Readiness Audit

**Date**: 2026-02-22
**Scope**: Full flow coverage, test coverage, and game/move logic testability

---

## Executive Summary

LilBattle has strong foundations: 379 test cases (368 passing, 4 failing), a well-separated `lib/` vs `services/` architecture, a fluent `GameBuilder` for test setup, and CI that gates builds on test results. However, several production-critical flows have **zero automated test coverage**, and the game logic tests focus on happy-path individual actions rather than multi-step scenario sequences. This audit catalogs every flow, grades its current coverage, and proposes what to add.

### Current Test Inventory

| Package | Tests | Status |
|---------|-------|--------|
| `tests/` (game logic) | ~320 | All passing |
| `lib/` (action sequences, rules) | ~40 | 4 failing (fix adjacency) |
| `cmd/cli/` | ~10 | All passing |
| `services/r2/` | 2 files | Not in default CI |
| `services/authz/` | 1 file | Not in CI |
| `web/server/` | 2 files | Not in CI |
| `web/tests/` (TS) | 2 files | Not in CI |
| Playwright e2e | 0 files | Configured but empty |

### Failing Tests (Pre-existing)

4 tests fail in `lib/action_sequence_test.go` — all related to "move then fix" for support units (Engineer, Medic, Carrier, Support). The fix action requires adjacency but the test places the target 2 hexes away. These should be fixed as part of this effort.

---

## Part 1: Flow-by-Flow Readiness

### 1.1 Authentication Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Email sign-up | `web/server/auth.go`, oneauth | `connect_auth_integration_test.go` covers Bearer token pipeline | Partial |
| Email sign-in | oneauth `UsernameStore` | Integration test covers token → gRPC user ID | Partial |
| OAuth (Google) | `web/server/auth.go` | None | **Gap** |
| OAuth (GitHub) | `web/server/auth.go` | None | **Gap** |
| OAuth (Twitter/PKCE) | `web/server/twitter_oauth2.go` | None | **Gap** |
| Session management | scs SessionManager | None | **Gap** |
| Password reset | `web/server/password_pages.go` | None | **Gap** |
| Auth middleware (gRPC) | `injectAuthMetadata` | Unit test exists | OK |
| Rate limiting | goapplib middleware | None | **Gap** |

**Recommendation**: OAuth flows are hard to unit-test without mocking providers. Add integration tests that verify the middleware chain (mock token → context propagation → service receives correct user ID). Password reset needs at minimum a handler-level test.

### 1.2 World Management Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Create world | `services/worlds_service.go` | `worlds_service_test.go` | OK |
| List worlds | WorldsService | `worlds_service_test.go` | OK |
| Edit world (editor) | WorldEditorPage + backend | None (frontend-only) | **Gap** |
| Delete world | Backend handler | None | **Gap** |
| World viewer (readonly) | WorldViewerPage | None | **Gap** |

**Recommendation**: Backend CRUD is covered. Frontend editor interactions need Playwright e2e tests (paint terrain, save, reload, verify).

### 1.3 Game Creation Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| StartGamePage config | Frontend + CreateGame RPC | None | **Gap** |
| CreateGame RPC | GamesService | `singleton/games_service_test.go` | Partial |
| Join game | JoinGame RPC | `join_game_test.go` | OK |
| Game initialization | `lib.NewGame()` | Used in every test via GameBuilder | OK |
| Player config (teams, colors) | GameConfiguration proto | Partial (builder sets players) | Partial |

**Recommendation**: Add a service-level test that exercises CreateGame → GetGame → verify initial state. Frontend config flow needs e2e coverage.

### 1.4 Gameplay Flow (Core Loop)

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Select unit / get options | `GetOptionsAt()` | `get_options_at_test.go`, `action_options_test.go` | OK |
| Movement | `ProcessMoveUnit()` | `moves_test.go` (9 cases) | OK |
| Attack | `ProcessAttackUnit()` | `lib/attack_test.go`, `tests/combat_formula_test.go` | OK |
| Build | `ProcessBuildUnit()` | `build_coins_test.go` (11 cases) | OK |
| Capture | `ProcessCaptureBuilding()` | `capture_test.go`, `capture_highlights_test.go` | OK |
| Heal | `ProcessHealUnit()` | `heal_test.go` | OK |
| Fix/Repair | `ProcessFixUnit()` | **4 failing tests** — adjacency issue | **Broken** |
| End turn | `ProcessEndTurn()` | Tested in `controller_test.go` | OK |
| Victory check | `checkVictoryConditions()` | Tested indirectly | Partial |
| Splash damage | Splash in `ProcessAttackUnit` | `splash_damage_test.go` | OK |
| Lazy top-up | `TopUpUnitIfNeeded()` | Tested via movement/attack tests | Implicit |
| Action progression | `GetAllowedActionsForUnit()` | `action_progression_test.go`, `lib/action_sequence_test.go` | OK |

**Recommendation**: Fix the 4 failing fix tests. Add explicit top-up tests. Victory conditions should have dedicated tests (see Part 2).

### 1.5 Presenter / UI Sync Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| SceneClicked dispatch | `gameview_presenter.go` | None | **Gap** |
| applyIncrementalChanges | Presenter | `apply_changes_test.go` | OK |
| refreshExhaustedHighlights | Presenter | Indirect in apply_changes | Partial |
| refreshCapturingHighlights | Presenter | `capture_highlights_test.go` | OK |
| ApplyRemoteChanges (sync) | Presenter | None | **Gap** |
| Highlight clearing | clearHighlightsAndSelection | None | **Gap** |
| Build modal trigger | Presenter | None | **Gap** |

**Recommendation**: The presenter is the glue between game logic and UI. `SceneClicked` (the main interaction dispatcher) has zero tests. Add presenter-level integration tests that mock the GameScene/panels and verify the correct sequence of calls for click → highlight → move → animate.

### 1.6 Multiplayer Sync Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| GameSyncManager connect | Frontend TS | None | **Gap** |
| Move publishing | Sync service | None | **Gap** |
| Remote change application | `ApplyRemoteChanges()` | None | **Gap** |
| Sequence tracking | GameSyncManager | None | **Gap** |

**Recommendation**: Multiplayer sync is the highest-risk untested area. At minimum needs a Go-level integration test that: publishes moves on one service instance → receives them on another → verifies state convergence.

### 1.7 CLI Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Position parsing | `position_parser.go` | Used implicitly in CLI tests | Partial |
| ww status | `cmd/status.go` | None directly | **Gap** |
| ww move/attack/build | Various cmd files | `worlds_test.go`, `assert_test.go` | Partial |
| ww map rendering | `cmd/map.go` | None | **Gap** |
| Dry run mode | --dryrun flag | None | **Gap** |
| JSON output mode | --json flag | None | **Gap** |

**Recommendation**: Add position parser unit tests covering all formats (unit shortcuts, Q/R, row/col, directions, t: prefix). CLI commands should have table-driven tests for each subcommand.

### 1.8 Storage / Persistence Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| FS backend (load/save) | `fsbe/games_service.go` | Used by all integration tests | Implicit |
| GORM backend | `gormbe/games_service.go` | None | **Gap** |
| Datastore backend (GAE) | `gaebe/` | None | **Gap** |
| R2/S3 file store | `services/r2/` | `filestore_test.go`, `r2_integration_test.go` | OK but not in CI |
| Cache behavior | `backend_games_service.go` | None | **Gap** |
| History append | `SaveMoves()` | Implicit in integration tests | Partial |

**Recommendation**: GORM backend needs schema migration tests and CRUD tests. R2 tests should be added to CI (or at least the non-integration one). The FS backend works because it's exercised everywhere, but an explicit round-trip test would help.

### 1.9 Deployment Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| CI build check | `ci.yml` | Runs on every push | OK |
| Go test suite in CI | `ci.yml` | `tests/`, `cmd/cli/`, `services/r2/`, `web/assets/themes/` | OK |
| Frontend build in CI | `ci.yml` | `pnpm run buildprod` | OK |
| WASM compile check | `ci.yml` | `go build ./cmd/wasm` | OK |
| Local replace guard | `ci.yml` | Checks go.mod and package.json | OK |
| GAE deploy | `deploy.yml` | Deploy + verify | OK |
| TS unit tests in CI | **Missing** | Not run | **Gap** |
| Playwright e2e in CI | **Missing** | Not configured | **Gap** |

**Recommendation**: Add `cd web && pnpm test` to CI. Configure Playwright with a headed-less runner in CI.

### 1.10 Authorization Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Authz rules | `services/authz/authz.go` | `authz_test.go` | OK but not in CI |
| Game ownership checks | In-service checks | None | **Gap** |
| World ownership checks | In-service checks | None | **Gap** |

**Recommendation**: Add authz tests to CI. Add tests for "player X cannot modify player Y's game" scenarios.

---

## Part 2: Game/Move Logic — Variable Testing

This section focuses on areas where the game logic needs **combinatorial / variable testing** — testing across different unit types, terrain types, player counts, and game states.

### 2.1 Movement Across Terrain Types

**Current state**: `moves_test.go` tests basic movement on grass. `movement_test.go` in lib tests a few terrain costs.

**What's missing** — movement cost varies by unit type x terrain type (44 x 26 matrix):

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Land unit on grass/plains/desert/mountain | High | Partial |
| Land unit blocked by water | High | **Gap** |
| Naval unit on shallow/regular/deep water | High | **Gap** |
| Naval unit blocked by land | High | **Gap** |
| Air unit ignoring terrain costs | Medium | **Gap** |
| Amphibious unit (hovercraft) crossing land↔water | High | **Gap** |
| Road/bridge crossing modifiers | Medium | **Gap** |
| Movement through occupied tiles (pass-through) | High | 1 test exists |
| Movement blocked by enemy units | High | **Gap** |

**Recommendation**: Create a table-driven `TestMovementCosts` that loads the rules JSON and verifies movement costs for representative unit/terrain pairs. Use subtests: `t.Run(fmt.Sprintf("%s_on_%s", unitName, terrainName), ...)`.

### 2.2 Combat Across Unit Matchups

**Current state**: `combat_formula_test.go` tests the damage formula math. `attack_test.go` tests basic attack flow.

**What's missing** — damage varies by attacker x defender x terrain:

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Infantry vs Infantry (baseline) | High | Covered |
| Tank vs Infantry (heavy vs light) | High | **Gap** |
| Anti-Air vs Air unit (specialization) | High | **Gap** |
| Air unit vs Naval unit | Medium | **Gap** |
| Ranged unit attack at max range | High | **Gap** |
| Ranged unit attack at min range (blocked) | High | **Gap** |
| Counter-attack eligibility matrix | High | **Gap** |
| Wound bonus accumulation (3+ attacks) | Medium | Partial |
| Damage with terrain defense bonuses | High | **Gap** |
| Unit killed → removed from world | High | Covered |
| Both units killed simultaneously | Medium | **Gap** |

**Recommendation**: Create `TestCombatMatchups` with a table of `{attacker, defender, terrain, expectedDamageRange}`. Use deterministic seeds for reproducible results. Test both the damage calculation and the state changes (health updated, unit removed if killed, splash applied).

### 2.3 Splash Damage Scenarios

**Current state**: `splash_damage_test.go` exists with basic coverage.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Splash hits adjacent enemy units | High | Covered |
| Splash does NOT hit air units | High | Covered |
| Splash friendly fire | Medium | **Gap** |
| Splash kills adjacent unit | Medium | **Gap** |
| Multiple units in splash radius | Medium | **Gap** |
| Splash damage values by unit type | Medium | **Gap** |

### 2.4 Build System Variations

**Current state**: `build_coins_test.go` covers basic build, coin deduction, wrong turn, non-buildable terrain.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Build each unit type (cost varies) | High | **Gap** |
| Build from landbase vs airport vs naval base | High | **Gap** |
| Allowed units filter restricts build options | High | `allowed_units_test.go` |
| Build at airport only allows air units | Medium | **Gap** |
| Build at naval base only allows naval units | Medium | **Gap** |
| One build per tile per turn | High | Covered |
| Build cost matches rules JSON | Medium | **Gap** |
| New unit has zero movement | High | Covered |

**Recommendation**: Table-driven test with `{baseTileType, unitType, expectedCost, shouldSucceed}` rows covering the buildability matrix.

### 2.5 Capture Mechanics

**Current state**: `capture_test.go` and `capture_highlights_test.go` cover start/complete flow.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Start capture on enemy building | High | Covered |
| Capture completes next turn | High | Covered |
| Capture interrupted by unit death | High | **Gap** |
| Capture interrupted by unit moving away | Medium | **Gap** |
| Only certain units can capture | High | **Gap** |
| Capture neutral building | Medium | **Gap** |
| Capture changes tile ownership | High | Covered |
| Multiple captures in same turn | Medium | **Gap** |

### 2.6 Action Progression Sequences

**Current state**: `action_progression_test.go` and `action_sequence_test.go` cover the state machine.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Move → Attack sequence | High | Covered |
| Move → Capture sequence | High | Covered |
| Move → Fix sequence | High | **Broken** (4 failing tests) |
| Attack exhausts turn | High | Covered |
| Pipe-separated alternatives | Medium | Covered |
| ChosenAlternative persistence | Medium | Covered |
| Unit type-specific progression | Medium | Partial |

### 2.7 End Turn / Turn Cycling

**Current state**: `controller_test.go` exercises multi-turn sequences.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| 2-player turn cycling | High | Covered |
| 3+ player turn cycling | Medium | **Gap** |
| Income calculation per player | High | `build_coins_test.go` (income tests) |
| Capture completion at turn start | High | Covered |
| Lazy top-up reset at turn boundary | High | Implicit |
| Victory condition: last player standing | High | **Gap** |
| Victory condition: all bases captured | Medium | **Gap** |
| Player elimination mid-game | Medium | **Gap** |

**Recommendation**: Add explicit victory condition tests. A 3-player game where P1 eliminates P2, then P3 eliminates P1.

### 2.8 Heal / Fix Mechanics

**Current state**: `heal_test.go` covers terrain-based healing. Fix tests are broken.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Heal on city terrain | High | Covered |
| Heal on non-healing terrain (blocked) | High | Covered |
| Fix adjacent friendly unit | High | **Broken** |
| Fix formula verification | Medium | **Gap** |
| Fix from different support unit types | Medium | **Gap** |
| Medic fix range vs Stratotanker range | Medium | **Gap** |
| Fix cannot target enemy units | Medium | **Gap** |

### 2.9 Pathfinding Edge Cases

**Current state**: Rules engine Dijkstra tested via movement, but edge cases not explicitly covered.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Path around impassable terrain | High | **Gap** |
| Path through narrow corridor | Medium | **Gap** |
| No valid path exists | High | **Gap** |
| Multiple equal-cost paths | Low | **Gap** |
| Pass-through occupied friendly tile | High | 1 test |
| Blocked by enemy unit on path | High | **Gap** |
| Movement exactly exhausting budget | Medium | **Gap** |

### 2.10 World Transaction Layer

**Current state**: `world_test.go` and `world_move_test.go` cover basic operations.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Push/Pop preserves parent state | High | Covered |
| Copy-on-write isolation | High | Covered |
| Move unit across layers | Medium | Covered |
| Concurrent reads on different layers | Medium | **Gap** |
| Deep nesting (3+ layers) | Low | **Gap** |

### 2.11 Seed / RNG Determinism

**Current state**: `seed_test.go` exists.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Same seed → same combat outcomes | High | Covered |
| Different seeds → different outcomes | Medium | Covered |
| Seed preserved across save/load | High | **Gap** |
| Deterministic replay from history | Medium | **Gap** |

---

## Part 3: Priority Recommendations

### P0 — Fix Before Launch

1. **Fix the 4 failing tests** in `lib/action_sequence_test.go` (fix adjacency setup)
2. **Add TS tests to CI** (`cd web && pnpm test`)
3. **Add `services/authz/` tests to CI**
4. **Add `web/server/` tests to CI** (connect auth integration)
5. **Victory condition tests** — a game without verified win conditions is incomplete

### P1 — High Priority (Production Quality)

6. **Presenter SceneClicked tests** — the main interaction path has zero coverage
7. **Movement terrain matrix** — table-driven tests for unit/terrain cost combinations
8. **Combat matchup matrix** — table-driven tests for attacker/defender/terrain damage
9. **Build eligibility matrix** — which units can be built at which bases
10. **Capture interruption** — capture broken by death or movement
11. **3+ player turn cycling and elimination**
12. **Pathfinding edge cases** — blocked paths, enemy blocking, narrow corridors

### P2 — Important (Robustness)

13. **Multiplayer sync integration test** — publish/receive/converge
14. **GORM backend tests** — production database layer untested
15. **Position parser unit tests** — all input formats
16. **CLI subcommand tests** — status, map, options output
17. **Splash damage variations** — friendly fire, multi-target, kill
18. **Seed determinism across save/load**
19. **Rate limiting verification**

### P3 — Nice to Have (Polish)

20. **Playwright e2e tests** — game creation → move → end turn → verify
21. **OAuth integration tests** with mocked providers
22. **Cache behavior tests** for backend service
23. **World editor e2e** — paint, save, reload
24. **Fix formula verification** — probabilistic correctness testing
25. **Session/password reset handler tests**

---

## Part 4: Testing Infrastructure Improvements

### 4.1 CI Coverage Gaps

The current CI command is:
```
go test ./tests/... ./cmd/cli/... ./services/r2/... ./web/assets/themes/...
```

This **misses**:
- `./lib/...` (action sequence tests, rules loader)
- `./services/authz/...`
- `./web/server/...` (connect auth integration)
- TypeScript unit tests (`web/tests/`)

Proposed CI test command:
```
go test ./tests/... ./cmd/cli/... ./lib/... ./services/authz/... ./services/r2/... ./web/server/... ./web/assets/themes/...
```

Note: Adding `./lib/...` will make CI fail until the 4 fix tests are repaired.

### 4.2 Test Coverage Reporting

Add coverage collection to CI:
```
go test -coverprofile=coverage.out ./tests/... ./lib/... ./cmd/cli/...
go tool cover -func=coverage.out | tail -1  # Total coverage %
```

### 4.3 GameBuilder Enhancements

The `GameBuilder` is solid. Suggested additions:
- `WithCrossing(fromQ, fromR, toQ, toR, terrainType)` for road/bridge tests
- `WithAllowedUnits(unitTypes...)` for build restriction tests
- `UnitDamaged(q, r, player, unitType, health)` shorthand for combat tests
- `MultiPlayer(n)` with auto-generated bases for N-player tests

### 4.4 Table-Driven Test Template

For game logic variable testing, establish this pattern:
```go
func TestMovementCosts(t *testing.T) {
    cases := []struct {
        name     string
        unitType int32
        terrain  int32
        expected float64
    }{
        {"soldier on grass", UnitTypeSoldierBasic, lib.TileTypeGrass, 1.0},
        {"soldier on mountain", UnitTypeSoldierBasic, lib.TileTypeMountain, 2.0},
        {"tank on road", UnitTypeTank, TileTypeRoad, 0.5},
        // ... load from rules JSON for completeness
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // ... verify movement cost
        })
    }
}
```

### 4.5 Playwright E2E Setup

Create `web/playwright.config.ts` and initial smoke tests:
1. **Auth smoke**: Load login page → verify form renders
2. **Game creation smoke**: Login → create game → verify redirect to viewer
3. **Gameplay smoke**: Open existing game → click unit → see highlights → click to move
4. **World editor smoke**: Open editor → paint tile → save → reload → verify

---

## Summary Scorecard

| Area | Coverage | Grade |
|------|----------|-------|
| Game engine (lib/) | 368/372 passing, comprehensive | **A** |
| Move processing | Well-tested, but missing cross-type coverage | **B+** |
| Combat formula | Formula math covered, matchups not | **B** |
| Build/Capture/Heal | Happy paths covered | **B** |
| Fix/Repair | Broken tests | **D** |
| Action progression | Good coverage | **A-** |
| Turn management | Basic coverage | **B** |
| Victory conditions | No dedicated tests | **D** |
| Presenter layer | Minimal | **D** |
| Auth/authz | Partial | **C** |
| Multiplayer sync | None | **F** |
| Frontend (TS unit) | 2 test files, not in CI | **D** |
| E2E tests | None | **F** |
| CLI commands | Minimal | **C-** |
| Storage backends | FS implicit, GORM/GAE untested | **C** |
| CI pipeline | Good for Go, missing TS and several packages | **B-** |
