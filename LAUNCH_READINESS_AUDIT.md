# LilBattle Launch Readiness Audit

**Audit Date**: January 9, 2026
**Last Updated**: June 15, 2026
**Overall Status**: READY for public launch (core), test gaps remain
**Estimated Completion**: 95% (features), ~60% (test coverage)

---

## June 2026 Refresh — What Changed

Targeted refresh after the stack-bump sweep landed today. Game-logic
analysis (sections 3, 4) untouched — the gameplay code hasn't moved
since February.

- **`lib/` tests now pass — and were already passing in Feb.** PR 109
  squash-merged three logical commits: the audit doc, a consolidation
  pass, and the actual fix-adjacency repair (moved the damaged
  friendly from `(-1, 0)` to `(1, -1)` so it's adjacent both before
  and after the fixer moves). The audit text was drafted between the
  first and third commits, so the "4 failing" line shipped stale. P0
  "fix 4 failing tests" is dropped and adding `./lib/...` to CI is
  risk-free.
- **Stack drift mostly resolved.** Issue 111 ("6 of 8 stack components
  behind") is closeable after PRs 116 and 118: `oneauth`, `goapplib`,
  `templar`, `servicekit`, and (transitively) `gocurrent` are now
  current. `protoc-gen-dal` v0.0.10 → v0.0.13 is the residual gap.
- **New oneauth security capabilities now available.** oneauth v0.1.x
  ships OAuth 2.1 alignment, `client_secret_jwt` token-endpoint client
  auth, encrypted private PEMs at rest, and a `core` / `httpauth` /
  `apiauth` / `localauth` / `federatedauth` split. Lilbattle uses only
  the basic surface today; the threat model didn't move, but the room
  to harden it did — see section 1 for what to consider.
- **Issue 113 is more attractive now.** "Check if `alexedwards/scs`
  is redundant with oneauth sessions" was filed before the v0.1.x
  session/middleware surface was as rich as it now is. Worth re-reading.

---

## Executive Summary

LilBattle has a solid technical foundation with production-ready core gameplay, multi-backend persistence, and clean architecture. All critical security and legal blockers have been addressed. The February 2026 audit adds a detailed analysis of test coverage across every flow and game logic path, identifying where combinatorial/variable testing is needed.

### Test Inventory (June 2026)

| Package | Tests | Status |
|---------|-------|--------|
| `tests/` (game logic) | ~320 | All passing |
| `lib/` (action sequences, rules) | ~40 | All passing (Feb audit's "4 failing" line was stale — fix shipped in same PR 109) |
| `cmd/cli/` | ~10 | All passing |
| `services/r2/` | 2 files | Not in default CI |
| `services/authz/` | 1 file | Not in CI |
| `web/server/` | 2 files | Not in CI (Connect-auth tests rewrote in PR 118 for oneauth Subject rename) |
| `web/tests/` (TS) | 2 files | Not in CI |
| Playwright e2e | 0 files | Configured but empty |

**Total**: 417 test runs, 417 passing, 0 failing across
`./tests/... ./cmd/cli/... ./lib/... ./services/authz/...
./services/r2/... ./web/server/...`. CI today covers only the first
two — see section 6.1 for the gap.

### Critical Blockers Status (Jan 2026 — all resolved)

| Area | Issue | Status |
|------|-------|--------|
| Security | API layer authentication | ✅ COMPLETED (#70) |
| Security | Rate limiting | ✅ COMPLETED (#70) |
| Security | Test credentials conditional on env var | ✅ COMPLETED (#70) |
| Legal | LICENSE file | ✅ COMPLETED |
| Docs | About page | ✅ COMPLETED (#66) |
| Docs | Contact/support page | ✅ COMPLETED (#66) |
| Persistence | UsersService multi-backend | ✅ COMPLETED (#71) |
| Security | Authorization on game/world ops | ✅ COMPLETED (#72) |
| Security | Security headers middleware | ✅ COMPLETED (#72) |
| Security | Authorization unit tests | ✅ COMPLETED (#72) |
| Persistence | No backup/disaster recovery strategy | 🟡 DEFERRED (cloud storage) |

---

## 1. Security & Authentication

### Completed ✅

**API Authentication** (PR #70, refreshed PR 118)
- gRPC/Connect endpoints have authentication via metadata
- Subject (formerly UserID) passed from HTTP session to gRPC context
- Auth interceptors enabled in grpcserver.go
- Uses oneauth v0.1.29 (split into core/httpauth/apiauth/localauth/
  federatedauth sub-modules); gRPC metadata key migrated `x-user-id`
  → `x-subject`, session cookie key migrated `loggedInUserId` →
  `loggedInSubject`

**Rate Limiting** (PR #70)
- Auth endpoints: 10 requests per 15 minutes
- API endpoints: 100 requests per minute
- IP-based limiting with proper headers

**Test Credentials Secured** (PR #70)
- Test auth conditional on `ENABLE_TEST_AUTH=true`
- User switching requires `ENABLE_SWITCH_AUTH=true`
- Auth disabled only with explicit `DISABLE_API_AUTH=true`

**Authorization Checks** (PR #72)
- Owner validation on UpdateGame/DeleteGame, UpdateWorld/DeleteWorld
- Player validation on ProcessMoves (must be game player AND current turn)
- Services: `services/authz/authz.go` with 17 unit test cases

**Security Headers** (PR #72)
- CSP, X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy
- Middleware: `web/server/securityheaders.go`

**OAuth Providers**
- Google, GitHub, Twitter/X (with PKCE)
- Session management via SCS library — see [issue 113][i113] for
  the question of whether oneauth's session surface (now richer
  in v0.1.x) makes SCS redundant

[i113]: https://github.com/turnforge/lilbattle/issues/113

### Remaining Security Items

| Priority | Task | Status |
|----------|------|--------|
| P1 | Fix insecure gRPC connections (TLS) | 🟡 TODO |
| P1 | Input validation framework | 🟡 TODO |
| P2 | CSRF tokens on all forms | TODO |
| P2 | Audit logging | TODO |
| P2 | Decide on SCS vs oneauth session surface (#113) | 🟡 TODO (more attractive after oneauth v0.1.x) |

### New oneauth Capabilities Available (June 2026)

Unlocked by the v0.1.x bump in PR 118 but not yet adopted in lilbattle.
None of these are launch blockers, but worth evaluating as the threat
model evolves:

- **OAuth 2.1 alignment** (audit/admin/apiauth) — stricter token
  endpoint behaviour, deprecates implicit/ROPC patterns we're not
  using. Reduces surface area for token-related vulns.
- **`client_secret_jwt` client auth** at the token endpoint —
  alternative to `client_secret_basic`/`_post` for CLI/confidential
  clients (the lilbattle CLI uses `client_secret_basic` today).
- **Encrypted private PEMs at rest** + Ed25519 SSH helper — relevant
  if we ever issue per-user signing keys.
- **`apiauth.OneAuth` composition root** — newer entry point that
  funnels grants/middleware through gRPC-shape `Issuer()`/`Validator()`
  accessors (oneauth issue 218). PR 118 already migrated lilbattle's
  one consumer (the JWT-cookie verifier) to `Validator().ValidateToken`
  since the legacy `APIAuth.VerifyTokenFunc` was removed.

---

## 2. Gameplay Features

### Complete ✅ (Production Ready)
- **Movement System**: Dijkstra pathfinding, 44 unit types, terrain costs
- **Combat System**: Probabilistic damage, 1.2MB authentic matrices, wound bonus, splash
- **Building Units**: Cost validation, terrain checks, shortcut generation
- **Capture System**: Multi-turn mechanic, ownership transfer
- **Healing System**: Terrain-based, unit-type restrictions
- **Fix/Repair System**: Adjacent friendly unit repair (4 test failures — see below)
- **Turn Management**: Player rotation, income generation, lazy top-up
- **Persistence**: File, PostgreSQL, Datastore backends
- **Replay/History**: Move groups, world changes, save/load
- **Unit Balance**: 44 types with authentic LilBattle data
- **Maps/Worlds**: Hex coords, 26 terrains, dynamic sizing
- **Multiplayer Infrastructure**: Sync broadcasting, transactions

### Known Issues

**Fix/Repair Tests Broken** — 4 tests fail in `lib/action_sequence_test.go`:
- `engineer_move_then_fix`, `support_move_then_fix`, `medic_move_then_fix`, `carrier_move_then_fix`
- Root cause: After moving from (0,0) to an adjacent tile, the friendly unit at (-1,0) is distance 2 from the new position, but fix requires adjacency.

**Victory Conditions** — Simplistic (last player with units wins). No alternate win conditions, elimination detection, or scoring.

**Damage Estimate** — Hardcoded at 50 in `game.go:755`. UI previews show incorrect values.

**AI Opponents** — Library exists in `.attic/lib/ai/` but not integrated with web UI.

---

## 3. Flow-by-Flow Test Coverage

### 3.1 Authentication Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Email sign-up | `web/server/auth.go`, oneauth | `connect_auth_integration_test.go` | Partial |
| Email sign-in | oneauth `UsernameStore` | Integration test covers token pipeline | Partial |
| OAuth (Google/GitHub/Twitter) | `web/server/auth.go` | None | **Gap** |
| Session management | scs SessionManager | None | **Gap** |
| Password reset | `web/server/password_pages.go` | None | **Gap** |
| Auth middleware (gRPC) | `injectAuthMetadata` | Unit test exists | OK |
| Rate limiting | goapplib middleware | None | **Gap** |

### 3.2 World Management Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Create world | `services/worlds_service.go` | `worlds_service_test.go` | OK |
| List worlds | WorldsService | `worlds_service_test.go` | OK |
| Edit world (editor) | WorldEditorPage + backend | None (frontend-only) | **Gap** |
| Delete world | Backend handler | None | **Gap** |
| World viewer (readonly) | WorldViewerPage | None | **Gap** |

### 3.3 Game Creation Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| StartGamePage config | Frontend + CreateGame RPC | None | **Gap** |
| CreateGame RPC | GamesService | `singleton/games_service_test.go` | Partial |
| Join game | JoinGame RPC | `join_game_test.go` | OK |
| Game initialization | `lib.NewGame()` | Used in every test via GameBuilder | OK |
| Player config (teams, colors) | GameConfiguration proto | Partial (builder sets players) | Partial |

### 3.4 Gameplay Flow (Core Loop)

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Select unit / get options | `GetOptionsAt()` | `get_options_at_test.go`, `action_options_test.go` | OK |
| Movement | `ProcessMoveUnit()` | `moves_test.go` (9 cases) | OK |
| Attack | `ProcessAttackUnit()` | `lib/attack_test.go`, `tests/combat_formula_test.go` | OK |
| Build | `ProcessBuildUnit()` | `build_coins_test.go` (11 cases) | OK |
| Capture | `ProcessCaptureBuilding()` | `capture_test.go`, `capture_highlights_test.go` | OK |
| Heal | `ProcessHealUnit()` | `heal_test.go` | OK |
| Fix/Repair | `ProcessFixUnit()` | **4 failing tests** — adjacency issue | **Broken** |
| End turn | `ProcessEndTurn()` | `controller_test.go` | OK |
| Victory check | `checkVictoryConditions()` | Tested indirectly | Partial |
| Splash damage | Splash in `ProcessAttackUnit` | `splash_damage_test.go` | OK |
| Lazy top-up | `TopUpUnitIfNeeded()` | Tested via movement/attack tests | Implicit |
| Action progression | `GetAllowedActionsForUnit()` | `action_progression_test.go`, `action_sequence_test.go` | OK |

### 3.5 Presenter / UI Sync Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| SceneClicked dispatch | `gameview_presenter.go` | None | **Gap** |
| applyIncrementalChanges | Presenter | `apply_changes_test.go` | OK |
| refreshExhaustedHighlights | Presenter | Indirect in apply_changes | Partial |
| refreshCapturingHighlights | Presenter | `capture_highlights_test.go` | OK |
| ApplyRemoteChanges (sync) | Presenter | None | **Gap** |
| Highlight clearing | clearHighlightsAndSelection | None | **Gap** |
| Build modal trigger | Presenter | None | **Gap** |

### 3.6 Multiplayer Sync Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| GameSyncManager connect | Frontend TS | None | **Gap** |
| Move publishing | Sync service | None | **Gap** |
| Remote change application | `ApplyRemoteChanges()` | None | **Gap** |
| Sequence tracking | GameSyncManager | None | **Gap** |

### 3.7 CLI Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Position parsing | `position_parser.go` | Used implicitly in CLI tests | Partial |
| ww status | `cmd/status.go` | None directly | **Gap** |
| ww move/attack/build | Various cmd files | `worlds_test.go`, `assert_test.go` | Partial |
| ww map rendering | `cmd/map.go` | None | **Gap** |
| Dry run mode | --dryrun flag | None | **Gap** |
| JSON output mode | --json flag | None | **Gap** |

### 3.8 Storage / Persistence Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| FS backend (load/save) | `fsbe/games_service.go` | Used by all integration tests | Implicit |
| GORM backend | `gormbe/games_service.go` | None | **Gap** |
| Datastore backend (GAE) | `gaebe/` | None | **Gap** |
| R2/S3 file store | `services/r2/` | `filestore_test.go`, `r2_integration_test.go` | OK but not in CI |
| Cache behavior | `backend_games_service.go` | None | **Gap** |
| History append | `SaveMoves()` | Implicit in integration tests | Partial |

### 3.9 Deployment Flow

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

### 3.10 Authorization Flow

| Step | Implementation | Tests | Verdict |
|------|---------------|-------|---------|
| Authz rules | `services/authz/authz.go` | `authz_test.go` | OK but not in CI |
| Game ownership checks | In-service checks | None | **Gap** |
| World ownership checks | In-service checks | None | **Gap** |

---

## 4. Game/Move Logic — Variable Testing

Areas where the game logic needs combinatorial testing across unit types, terrain types, player counts, and game states.

### 4.1 Movement Across Terrain Types

**Current**: `moves_test.go` tests basic movement on grass. `movement_test.go` tests a few terrain costs.

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Land unit on grass/plains/desert/mountain | High | Partial |
| Land unit blocked by water | High | **Gap** |
| Naval unit on shallow/regular/deep water | High | **Gap** |
| Naval unit blocked by land | High | **Gap** |
| Air unit ignoring terrain costs | Medium | **Gap** |
| Amphibious unit (hovercraft) crossing land/water | High | **Gap** |
| Road/bridge crossing modifiers | Medium | **Gap** |
| Movement through occupied tiles (pass-through) | High | 1 test exists |
| Movement blocked by enemy units | High | **Gap** |

### 4.2 Combat Across Unit Matchups

**Current**: `combat_formula_test.go` tests damage formula math. `attack_test.go` tests basic flow.

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
| Unit killed and removed from world | High | Covered |
| Both units killed simultaneously | Medium | **Gap** |

### 4.3 Splash Damage Scenarios

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Splash hits adjacent enemy units | High | Covered |
| Splash does NOT hit air units | High | Covered |
| Splash friendly fire | Medium | **Gap** |
| Splash kills adjacent unit | Medium | **Gap** |
| Multiple units in splash radius | Medium | **Gap** |
| Splash damage values by unit type | Medium | **Gap** |

### 4.4 Build System Variations

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

### 4.5 Capture Mechanics

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

### 4.6 Action Progression Sequences

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Move then Attack sequence | High | Covered |
| Move then Capture sequence | High | Covered |
| Move then Fix sequence | High | **Broken** (4 failing tests) |
| Attack exhausts turn | High | Covered |
| Pipe-separated alternatives | Medium | Covered |
| ChosenAlternative persistence | Medium | Covered |
| Unit type-specific progression | Medium | Partial |

### 4.7 End Turn / Turn Cycling

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| 2-player turn cycling | High | Covered |
| 3+ player turn cycling | Medium | **Gap** |
| Income calculation per player | High | `build_coins_test.go` |
| Capture completion at turn start | High | Covered |
| Lazy top-up reset at turn boundary | High | Implicit |
| Victory condition: last player standing | High | **Gap** |
| Victory condition: all bases captured | Medium | **Gap** |
| Player elimination mid-game | Medium | **Gap** |

### 4.8 Heal / Fix Mechanics

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Heal on city terrain | High | Covered |
| Heal on non-healing terrain (blocked) | High | Covered |
| Fix adjacent friendly unit | High | **Broken** |
| Fix formula verification | Medium | **Gap** |
| Fix from different support unit types | Medium | **Gap** |
| Medic fix range vs Stratotanker range | Medium | **Gap** |
| Fix cannot target enemy units | Medium | **Gap** |

### 4.9 Pathfinding Edge Cases

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Path around impassable terrain | High | **Gap** |
| Path through narrow corridor | Medium | **Gap** |
| No valid path exists | High | **Gap** |
| Multiple equal-cost paths | Low | **Gap** |
| Pass-through occupied friendly tile | High | 1 test |
| Blocked by enemy unit on path | High | **Gap** |
| Movement exactly exhausting budget | Medium | **Gap** |

### 4.10 World Transaction Layer

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Push/Pop preserves parent state | High | Covered |
| Copy-on-write isolation | High | Covered |
| Move unit across layers | Medium | Covered |
| Concurrent reads on different layers | Medium | **Gap** |
| Deep nesting (3+ layers) | Low | **Gap** |

### 4.11 Seed / RNG Determinism

| Test Scenario | Priority | Status |
|---------------|----------|--------|
| Same seed produces same combat outcomes | High | Covered |
| Different seeds produce different outcomes | Medium | Covered |
| Seed preserved across save/load | High | **Gap** |
| Deterministic replay from history | Medium | **Gap** |

---

## 5. Priority Recommendations

### P0 — Fix Now

1. ~~**Fix the 4 failing tests** in `lib/action_sequence_test.go`~~ —
   ✅ they were fixed in the same PR that added the audit (PR 109
   squashed three commits; audit text drafted before the fix
   commit). Permanent — re-verified June 2026.
2. **Add `lib/` tests to CI** — unblocked; safe to land.
3. **Add TS tests to CI** (`cd web && pnpm test`)
4. **Add `services/authz/` and `web/server/` tests to CI**
5. **Victory condition tests** — a game without verified win conditions is incomplete
6. **Fix damage estimate** — still hardcoded at 50 in `lib/game.go:755` (audit said `game.go` — actual location is `lib/game.go`)

### P1 — High Priority (Production Quality)

7. **Presenter SceneClicked tests** — main interaction path, zero coverage
8. **Movement terrain matrix** — table-driven tests for unit/terrain cost combinations
9. **Combat matchup matrix** — table-driven tests for attacker/defender/terrain damage
10. **Build eligibility matrix** — which units can be built at which bases
11. **Capture interruption** — capture broken by death or movement
12. **3+ player turn cycling and elimination**
13. **Pathfinding edge cases** — blocked paths, enemy blocking, narrow corridors

### P2 — Important (Robustness)

14. **Multiplayer sync integration test** — publish/receive/converge
15. **GORM backend tests** — production database layer untested
16. **Position parser unit tests** — all input formats
17. **CLI subcommand tests** — status, map, options output
18. **Splash damage variations** — friendly fire, multi-target, kill
19. **Seed determinism across save/load**
20. **Rate limiting verification**
21. **Input validation framework**

### P3 — Nice to Have

22. **Playwright e2e tests** — game creation, move, end turn, verify
23. **OAuth integration tests** with mocked providers
24. **Cache behavior tests** for backend service
25. **World editor e2e** — paint, save, reload
26. **Fix formula verification** — probabilistic correctness testing
27. **Session/password reset handler tests**
28. **API documentation** for developers
29. **FAQ/Help page** for users
30. **Browser-based game tutorial**

---

## 6. Testing Infrastructure

### 6.1 CI Coverage Gaps

The current CI command:
```
go test ./tests/... ./cmd/cli/... ./services/r2/... ./web/assets/themes/...
```

**Misses**:
- `./lib/...` (action sequence tests, rules loader)
- `./services/authz/...`
- `./web/server/...` (connect auth integration)
- TypeScript unit tests (`web/tests/`)

Proposed CI test command:
```
go test ./tests/... ./cmd/cli/... ./lib/... ./services/authz/... ./services/r2/... ./web/server/... ./web/assets/themes/...
cd web && pnpm test
```

Note: As of June 2026, the proposed command above passes locally
(417 runs, all green). The Feb-era caveat that `./lib/...` would break
CI no longer applies.

### 6.2 Test Coverage Reporting

```
go test -coverprofile=coverage.out ./tests/... ./lib/... ./cmd/cli/...
go tool cover -func=coverage.out | tail -1
```

### 6.3 GameBuilder Enhancements

The `GameBuilder` in `tests/game_builder.go` is solid. Suggested additions:
- `WithCrossing(fromQ, fromR, toQ, toR, terrainType)` for road/bridge tests
- `WithAllowedUnits(unitTypes...)` for build restriction tests
- `UnitDamaged(q, r, player, unitType, health)` shorthand for combat tests
- `MultiPlayer(n)` with auto-generated bases for N-player tests

### 6.4 Table-Driven Test Template

For game logic variable testing:
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
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // ... verify movement cost
        })
    }
}
```

### 6.5 Playwright E2E Setup

Create `web/playwright.config.ts` and initial smoke tests:
1. **Auth smoke**: Load login page, verify form renders
2. **Game creation smoke**: Login, create game, verify redirect to viewer
3. **Gameplay smoke**: Open game, click unit, see highlights, click to move
4. **World editor smoke**: Open editor, paint tile, save, reload, verify

---

## 7. Documentation & Legal

### Completed ✅
- LICENSE file (MIT License)
- About page (PR #66)
- Contact/Support page (PR #66)
- README.md with architecture diagram and CLI examples
- CLI User Guide (404 lines)
- Architecture docs (2000 lines)
- Terms of Service and Privacy Policy (generic)
- Profile page (account management)

### Remaining

| Priority | Task | Status |
|----------|------|--------|
| P1 | API documentation (OpenAPI or README) | 🟡 TODO |
| P1 | Customize Terms/Privacy for LilBattle practices | TODO |
| P1 | Help/FAQ page | 🟡 TODO |
| P2 | Browser game tutorial | TODO |
| P2 | CONTRIBUTING.md | TODO |
| P3 | CHANGELOG.md | TODO |

---

## 8. Monetization

See [MONETIZATION.md](./MONETIZATION.md) for full strategy.

| Phase | Scope | Status |
|-------|-------|--------|
| Phase 1 | Footer banners, homepage ads, feature flags | ✅ Implemented |
| Phase 2 | Game end screen, turn transition ads | Planned |
| Phase 3 | Rewarded video ads | Planned |
| Phase 4 | Premium ad-free tier ($3-5/month) | Planned |

Phase 1 ad infrastructure (AdSlot component, AdScript loader, feature flags, CSP updates) is implemented. Google AdSense account setup still needed.

---

## 9. Risk Assessment

### Mitigated ✅
1. ~~Data Breach: API has no auth~~ → API authentication (#70)
2. ~~Authorization Bypass~~ → Authorization checks (#72)
3. ~~Data Loss: No backups~~ → Cloud storage with built-in redundancy
4. ~~Legal Liability: No LICENSE~~ → MIT License added
5. ~~Denial of Service~~ → Rate limiting (#70)
6. ~~Security Headers Missing~~ → Security headers middleware (#72)

### Medium Risk
1. **Test gaps in game logic** — combat matchups, terrain movement, victory conditions untested
2. ~~**Fix/repair broken** — 4 failing tests~~ → fixed in PR 109
   alongside the audit itself; permanent
3. **Multiplayer untested** — sync flow has zero coverage
4. **Poor retention** — no tutorial, users may churn
5. **Damage estimates wrong** — hardcoded at 50 in `lib/game.go`

### Low Risk
1. AI not integrated but game works in hotseat mode
2. Simple win condition but functional
3. Missing API docs (internal-facing for now)

---

## 10. Summary Scorecard

| Area | Coverage | Grade |
|------|----------|-------|
| Security (auth, authz, headers) | All critical items completed; oneauth v0.1.x unlocks more (#113) | **A** |
| Game engine (lib/) | All passing (June 2026), comprehensive | **A** |
| Move processing | Well-tested, missing cross-type coverage | **B+** |
| Combat formula | Formula math covered, matchups not | **B** |
| Build/Capture/Heal | Happy paths covered | **B** |
| Fix/Repair | Previously D (4 broken tests); now passing | **B** |
| Action progression | Good coverage | **A-** |
| Turn management | Basic coverage | **B** |
| Victory conditions | No dedicated tests | **D** |
| Presenter layer | Minimal | **D** |
| Multiplayer sync | None | **F** |
| Frontend (TS unit) | 2 test files, not in CI | **D** |
| E2E tests | None | **F** |
| CLI commands | Minimal | **C-** |
| Storage backends | FS implicit, GORM/GAE untested | **C** |
| CI pipeline | Good for Go, missing TS and several packages | **B-** |
| Documentation & Legal | Core items done, API docs missing | **B** |
| Monetization | Phase 1 infrastructure done | **B+** |
