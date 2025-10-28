# Proposed Proto Changes for Attack Formula Implementation

## Overview

Currently, combat uses a pre-calculated damage distribution lookup table. To implement the actual attack formula from `docs/ATTACK.md`, we need to track additional state and change how combat properties are stored.

## Attack Formula (from ATTACK.md)

```
p = 0.05 * ( ( ( A + Ta ) - ( D + Td ) ) + B ) + 0.5
```

Where:
- **A** = Attack value of attacking unit (base attack)
- **Ta** = Terrain attack bonus for attacker
- **D** = Defense value of defending unit
- **Td** = Terrain defense bonus for defender
- **B** = Wound bonus (accumulated attacks this turn)
- **Ha** = Health of attacker (number of dice rolls = Ha * 6)

## Current State

### What We Have
- ✅ `UnitDefinition.defense` - base defense (D)
- ✅ `UnitDefinition.splash_damage` - splash damage amount (S)
- ✅ `TerrainUnitProperties.attack_bonus` - terrain attack bonus (Ta)
- ✅ `TerrainUnitProperties.defense_bonus` - terrain defense bonus (Td)
- ❌ `UnitUnitProperties.damage` - damage distribution (doesn't use formula)

### What We're Missing
- ❌ **Base attack value (A)** - not in UnitDefinition
- ❌ **Wound bonus tracking (B)** - no state for attacks received this turn
- ❌ **Attack history** - no way to calculate wound bonus based on attack positions

## Proposed Proto Changes

### 1. Add Base Attack to UnitDefinition

```protobuf
message UnitDefinition {
  int32 id = 1;
  string name = 2;
  string description = 3;
  int32 health = 4;
  int32 coins = 5;
  double movement_points = 6;

  int32 attack = 7;              // NEW: Base attack value (A)
  int32 defense = 8;             // MOVED: Was field 7
  int32 attack_range = 9;        // MOVED: Was field 8
  int32 min_attack_range = 10;   // MOVED: Was field 9
  int32 splash_damage = 11;      // MOVED: Was field 10

  map<int32, TerrainUnitProperties> terrain_properties = 12;  // MOVED: Was field 11
  repeated string properties = 13;  // MOVED: Was field 12
}
```

**Rationale:** The attack value (A) is fundamental to the attack formula. Every unit needs a base attack stat.

### 2. Add Wound Bonus Tracking to Unit

```protobuf
message Unit {
  int32 q = 1;
  int32 r = 2;
  int32 player = 3;
  int32 unit_type = 4;
  string shortcut = 5;

  int32 available_health = 6;
  double distance_left = 7;
  int32 last_acted_turn = 8;
  int32 last_toppedup_turn = 9;

  // NEW: Wound bonus tracking for this turn
  int32 attacks_received_this_turn = 10;  // Total number of attacks received
  repeated AttackRecord attack_history = 11;  // Detailed attack history for wound bonus calculation
}

// NEW: Track individual attacks for wound bonus calculation
message AttackRecord {
  int32 q = 1;              // Attacker's Q coordinate
  int32 r = 2;              // Attacker's R coordinate
  bool is_ranged = 3;       // Was the attack from 2+ tiles away?
  int32 turn_number = 4;    // Which turn this attack occurred
}
```

**Rationale:**
- The wound bonus depends on HOW MANY attacks and FROM WHERE they came
- Need to track attack positions to calculate geometric relationships (adjacent, opposite side, etc.)
- `AttackRecord` stores enough info to calculate the wound bonus correctly
- These reset at turn end via the lazy top-up pattern

### 3. Keep DamageDistribution in UnitUnitProperties

```protobuf
message UnitUnitProperties {
  int32 attacker_id = 1;
  int32 defender_id = 2;
  DamageDistribution damage = 3;  // Keep for UI display and backward compatibility

  // Optional: Can be used to override base attack/defense if needed
  int32 attack_override = 4;   // Optional: specific A value for this matchup
  int32 defense_override = 5;  // Optional: specific D value for this matchup
}
```

**Note:** The `DamageDistribution` can indeed be calculated dynamically from:
- Attacking unit definition (attack, health)
- Attacker's terrain (attack bonus)
- Defending unit definition (defense)
- Defender's terrain (defense bonus)
- Wound bonus (from attack history)
- Splash damage settings

However, we'll keep it for now for:
1. **UI Display** - showing expected damage ranges in tooltips
2. **Backward Compatibility** - existing code relies on it
3. **Performance** - can cache common matchup distributions
4. **Future** - may want special matchup rules that override the formula

The actual combat calculation will use the formula, but the distribution can be pre-calculated or cached for display purposes.

## Implementation Plan

### Phase 1: Verify Proto Changes ✅
- ✅ Verify wound bonus fields added to `Unit`
- ✅ Verify `AttackRecord` message exists
- ✅ Verify `unit_class` and `unit_terrain` fields added
- ✅ Verify `attack_vs_class` map added

### Phase 2: Update Rules Data ✅
- ✅ Extended extraction script with XPath queries
- ✅ Added unit classification extraction
- ✅ Added attack table parsing
- ✅ Generated new `assets/weewar-rules.json` with all new fields
- ✅ Verified all units have classification and attack tables

### Phase 3: Implement Formula-Based Combat
1. Create `CalculateCombatDamageWithFormula()` function
2. Implement: `p = 0.05 * ( ( ( A + Ta ) - ( D + Td ) ) + B ) + 0.5`
3. Roll dice: for each health point, roll 6 dice, count hits where random < p
4. Keep old `CalculateCombatDamage()` as fallback

### Phase 4: Implement Wound Bonus Calculation
1. Add `CalculateWoundBonus()` function
2. Parse `attack_history` to determine bonus based on attack geometry
3. Track attacks during combat in `ProcessAttackUnit()`
4. Reset `attack_history` during lazy top-up

### Phase 5: Implement Splash Damage
1. Find adjacent units to defender
2. For each adjacent unit (except air units):
   - Calculate attack WITHOUT wound bonus
   - If attack value > 4, apply splash damage
3. Splash damage doesn't add to wound bonus

### Phase 6: Update Combat Diagnostics
1. Show formula breakdown (A, Ta, D, Td, B)
2. Show calculated hit probability (p)
3. Show expected hits (Ha * 6 * p / 6)
4. Show wound bonus calculation
5. Show splash damage targets

## Migration Strategy

1. **Proto regeneration** - Run `buf generate` to update generated code
2. **Rules data update** - Add attack values to all units
3. **Backward compatibility** - Keep both old and new combat functions
4. **Feature flag** - Optional flag to switch between old/new systems
5. **Testing** - Comprehensive tests with various scenarios

## Questions Answered

**Q: Can DamageDistribution be calculated dynamically?**
A: Yes! It can be derived from:
```
Inputs: attackerUnit, attackerHealth, attackerTile, defenderUnit, defenderTile, woundBonus
Process:
  1. Get A from attackerUnit.attack
  2. Get Ta from attackerTile's terrain bonus for attacker unit type
  3. Get D from defenderUnit.defense
  4. Get Td from defenderTile's terrain bonus for defender unit type
  5. Get B from defender's attack_history
  6. Calculate p = 0.05 * (((A + Ta) - (D + Td)) + B) + 0.5
  7. Simulate damage distribution by running formula multiple times with RNG
  8. Generate DamageDistribution with ranges and probabilities
```

However, this is expensive to do on-the-fly for UI tooltips, so we can:
- Pre-calculate for common scenarios (no wound bonus, standard terrain)
- Cache based on input parameters
- Use simplified expected value for quick estimates

## Next Steps

1. ✅ Create this proposal document
2. ✅ Verify proto changes are applied
3. ✅ Update `assets/weewar-rules.json` with attack values
4. ⏳ Implement formula-based combat functions
5. ⏳ Implement wound bonus calculation
6. ⏳ Implement splash damage
7. ⏳ Update combat diagnostics
8. ⏳ Test thoroughly with various scenarios
