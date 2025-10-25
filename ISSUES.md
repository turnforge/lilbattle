# WeeWar - Issues and Bug Fixes

This document tracks all bugs, issues, and their resolutions in chronological order.

## Current Session (2025-10-25)

### Issue #1: Unit Shortcuts Being Lost on Move/Turn Changes
**Status**: ✅ FIXED
**Severity**: High
**Date Fixed**: 2025-10-25

**Problem**: Unit shortcuts (A1, B1, C1, C2) were being corrupted or reset when units moved or turns ended. Running `ww units` would show shortcuts like "C?", "C6" instead of proper sequential labels.

**Root Cause**:
- ProcessMoveUnit, ProcessEndTurn, and ProcessAttackUnit created unit copies for history recording
- These manual copies didn't include the Shortcut field
- When ApplyChangeResults updated game state from WorldChange, it used incomplete unit objects
- Units got saved without shortcuts, which then got regenerated with incorrect sequential numbers

**Evidence**:
```json
// state.json showed corrupted shortcuts
{"q": 0, "r": 1, "shortcut": "", ...},      // Empty!
{"q": 1, "r": 1, "shortcut": "C6", ...}     // Wrong number - should be C2
```

**Solution**:
- Created `copyUnit()` helper function in services/moves.go that copies all unit fields
- Refactored all manual unit copying to use copyUnit()
- For damage/kill cases: copyUnit() then override specific fields (e.g., AvailableHealth)

**Files Changed**:
- services/moves.go: Added copyUnit(), refactored ProcessMoveUnit, ProcessEndTurn, ProcessAttackUnit

**Code Reduction**: 90 lines of repetitive code → simple copyUnit() calls

**Commit**: eaa6376d "Refactor unit copying with helper function to preserve all fields"

---

### Issue #2: FSGamesService Cache Causing Stale Browser Data
**Status**: ✅ FIXED
**Severity**: High
**Date Fixed**: 2025-10-25 (previous session continuation)

**Problem**: Browser showed stale unit positions and game state even after CLI made changes and state.json was updated on disk. Refreshing the browser didn't help.

**Root Cause**:
- FSGamesService has in-memory caches (gameCache, stateCache, historyCache)
- CLI and gRPC server each have separate FSGamesService instances with separate caches
- When CLI updates disk via its FSGamesService, gRPC server's cache never gets invalidated
- GetGame() checks cache first before reading disk

**Solution**: Disabled cache check in FSGamesService.GetGame() method (lines 150-160 in fsgames_service.go). GetGame() now always reads fresh from disk.

**Files Changed**:
- services/fsgames_service.go: Commented out cache check in GetGame()

**Future Consideration**: If cache is re-enabled, implement file watching or cache invalidation mechanism.

**Commit**: "Disable FSGamesService cache to fix stale data issue"

---

## Previous Sessions

### Issue #3: Lazy Top-Up Bug - Units Not Moving Despite Having Options
**Status**: ✅ FIXED
**Severity**: High
**Date Fixed**: 2025-10-24

**Problem**: Units showing DistanceLeft=0 couldn't move even when GetOptionsAt showed valid movement options. CLI `ww options B1` would show moves, but `ww move B1 R` would fail.

**Root Cause**: Units weren't being "topped up" (movement points refreshed) before move validation. The lazy top-up pattern requires explicit topUpUnitIfNeeded() calls.

**Solution**:
- Added topUpUnitIfNeeded() calls in ProcessMoveUnit before validation
- Added topUpUnitIfNeeded() calls in ProcessAttackUnit before validation
- Units now properly refresh movement points based on current turn

**Files Changed**:
- services/moves.go: Added topUpUnitIfNeeded() calls in ProcessMoveUnit and ProcessAttackUnit

---

### Issue #4: Movement Points Lost on Game Load
**Status**: ✅ FIXED
**Severity**: High
**Date Fixed**: 2025-10-23

**Problem**: Units lost movement points when loading saved game state. Fresh units would show full movement, but saved units would show 0 movement even if they hadn't moved.

**Root Cause**: `NewGame()` always called `initializeStartingUnits()` which reset all unit stats to maximum values, ignoring saved state.

**Solution**:
- Created `NewGameFromState()` function that preserves unit stats from saved state
- FSGamesService now uses NewGameFromState() when loading existing games
- NewWorld() and conversion utilities preserve DistanceLeft from protobuf

**Files Changed**:
- services/utils.go: Preserve DistanceLeft field in conversions
- services/game.go: Added NewGameFromState() function

---

### Issue #5: Unit Duplication Bug - Units Appearing at Two Positions
**Status**: ✅ FIXED
**Severity**: Critical
**Date Fixed**: 2025-10-22

**Problem**: After moving a unit, it would appear at both the old position and new position. This broke the game state completely.

**Root Cause**:
- Transaction layer shared unit object references with parent layer
- When MoveUnit modified unit coordinates, it corrupted both transaction and parent layer
- ProcessMoveUnit captured original unit before move instead of moved copy

**Solution**:
- Implemented copy-on-write semantics in World.MoveUnit()
- Transaction layers now create unit copies before modification
- ProcessMoveUnit captures moved unit from World.UnitAt(destination)
- Parent layer objects remain immutable during transaction processing

**Files Changed**:
- services/world.go: Copy-on-write in MoveUnit, AddUnit, RemoveUnit
- services/moves.go: ProcessMoveUnit captures moved unit correctly

**Testing**:
- Created comprehensive World operation tests (basic moves, replacements, transactions)
- End-to-end ProcessMoves integration tests using WasmGamesService
- Transaction flow simulation tests validating copy-on-write semantics

---

### Issue #6: WorldEditorPage Theme Query Parameter Ignored
**Status**: ✅ FIXED
**Severity**: Medium
**Date Fixed**: 2025-10-21

**Problem**: WorldEditorPage button panel always showed fantasy-themed assets regardless of `?theme=` query parameter.

**Root Cause**: Theme was hardcoded to "fantasy" in WorldEditorPage.go:83, query parameter was never read.

**Solution**:
- Added `Theme` field to WorldEditorPage struct
- Read `?theme=` query parameter in Load() method with "fantasy" default fallback
- Use `v.Theme` in SetupDefaults() instead of hardcoded value

**Files Changed**:
- web/server/WorldEditorPage.go

**Result**: Button panel icons now correctly match theme query parameter.

---

### Issue #7: StartGamePage Nil Pointer Error
**Status**: ✅ FIXED
**Severity**: Medium
**Date Fixed**: 2025-10-21

**Problem**: Template render error when clicking "Create new" from ListGames page - panic on nil pointer dereference.

**Root Cause**: Template accessed `.World.Name` without checking if `.World` is nil when no worldId provided.

**Solution**: Changed all `.World.Name` checks to `and .World .World.Name` pattern in template (5 locations).

**Files Changed**:
- web/templates/StartGamePage.html

**Result**: Page correctly displays "Select a world" UI when no worldId provided.

---

### Issue #8: GetPlayerUnits Index Out of Range in ProcessEndTurn
**Status**: ✅ FIXED
**Severity**: High
**Date Fixed**: 2025-10-20

**Problem**: ProcessEndTurn caused panic: index out of range when trying to reset units for new turn.

**Root Cause**:
- World.Push() creates empty transaction layers
- GetPlayerUnits() only checked current layer's unitsByPlayer array
- Transaction layer had empty array, causing index out of range

**Solution**:
- Fixed GetPlayerUnits() to check parent layer when transaction layer is empty
- Transaction layers now properly fall back to parent data

**Files Changed**:
- services/world.go: GetPlayerUnits() fallback logic

---

## Known Issues (Not Yet Fixed)

### Issue #9: CLI Options Command Async Panel Updates
**Status**: 🚧 KNOWN ISSUE
**Severity**: Medium

**Problem**: Options command calls SceneClicked which runs in goroutines (async), returns immediately before panels are populated. TurnOptionsPanel not yet updated when CLI tries to read it.

**Proposed Solution**: Create Cmd panel versions (CmdTurnOptionsPanel, CmdGameState) with channel-based callbacks for synchronous communication.

**Files**: cmd/cli/cmd/presenter.go, services/panels.go

---

### Issue #10: Visual Updates Use Full Scene Reload
**Status**: 🚧 MINOR POLISH
**Severity**: Low

**Problem**: After moves, the entire scene is reloaded instead of using targeted updates for specific units/tiles.

**Impact**: Slightly slower rendering, but functionally correct.

**Proposed Solution**: Implement incremental SetUnitAt/SetTileAt updates in renderer.

---

### Issue #11: Missing Move Animations
**Status**: 🚧 MINOR POLISH
**Severity**: Low

**Problem**: Units teleport instantly to new positions instead of smoothly animating movement.

**Impact**: User experience, not functionality.

**Proposed Solution**: Add Phaser tween animations for unit movement along path.

---

## Fixed Issues Summary

**Critical Issues Fixed**: 2 (Unit Duplication, FSGamesService Cache)
**High Severity Fixed**: 4 (Shortcuts Lost, Lazy Top-Up, Movement Points, GetPlayerUnits)
**Medium Severity Fixed**: 2 (Theme Query, StartGamePage Nil)
**Low Severity Fixed**: 0

**Total Issues Resolved**: 8
**Known Open Issues**: 3

---

## Debugging Tips for Common Issues

**Unit shortcuts corrupted?**
```bash
jq '.world_data.units[] | {shortcut, q, r}' ~/dev-app-data/weewar/storage/games/{gameId}/state.json
```

**Units not moving?**
```bash
ww --verbose options B1 | grep "DistanceLeft"
jq '.world_data.units[] | select(.shortcut == "B1") | {distance_left, last_topped_up_turn}' state.json
```

**Browser showing stale data?**
- Check if FSGamesService cache is enabled in services/fsgames_service.go
- Verify state.json on disk has correct data
- Check browser console for errors

**Move history not recording?**
```bash
jq '.groups[-1]' ~/dev-app-data/weewar/storage/games/{gameId}/history.json
```
