# Attack Mathematics
 
## Attack Values

Attack values are determined by the following:

* The amount of health points (Ha) of the attacking unit.
* The attack (A) of the attacking unit.
* The defense (D) of the defending unit.
* The terrain bonuses both the attacker (Ta) and defender (Td) are on during the attack.
* The wound bonus (B).

## Wound Bonus

A unit that is attacked multiple times exerts a negative wound bonus (B) for all future attacks in the same turn:

* +1 for each previous attack if the current attacker is ranged (attacks from 2+ tiles away).

- OR -

* +1 for each previous attack from a distance (e.g. by an artillery, anti-air, destroyer, etc.).
* +1 for each previous attack from a hex adjacent to the attacker and the defender.
* +1 for each previous attack from the same location (e.g. the unit died and another unit stepped into the same tile).
* +2 for each previous attack from any other hex adjacent to the defender.
* +3 for each previous attack from a hex on the opposite side of (across from) the defender.

## Attack Formula

p = 0.05 * ( ( ( A + Ta ) - ( D + Td ) ) + B ) + 0.5

if p < 0, set p to 0
if p > 1, set p to 1

For each health unit (Ha) of the attacking unit, six random numbers (r) between 0 and 1 are generated. Each time r < p, a hit is counted. The number of hits divided by 6 is the number of health points the defender loses from the attack.

If the defender is able to hit the attacker (based on its type and range), the above calculations are performed for the defender attacking the attacker. Damage isn't subtracted until the end so that both units' attacks are calculated with their initial health numbers.

## Splash Damage

* Units cause splash damage (S) to the surrounding/adjancent units of the defender.
* Only certain units can cause splash damage (e.g. artillery, missiles, sea mines, etc.).
* Splash damage is dealt to both friendly and enemy units.
* Air units are unaffected by splash damage.

Splash damage is determined by using the same formula as above, but without the wound bonus (B). For each splash damage (S) of the attacking unit, an attack value is determined by the attack formula. If the resulting attack value is greater than 4, the splash damage is dealt. Damage dealt by splash damage does not contribute to any wound bonus.

## Notes

* It is not possible for a unit to deal more damage than it has health points.
* It is usually better to attack the same unit multiple times in the same turn to take advantage of the wound bonus.
* The health of the defender is not part of the attack formula, but does play a role if the defender can attack back.
* Be aware that splash damage can hurt your own units.
