#!/bin/bash
set -e  # Exit on first error

# Error handler - show which command failed
trap 'echo "FAILED at line $LINENO: $BASH_COMMAND" >&2' ERR

# Weewar Game Replay Script
# Generated from game history
# Players: Computer (A), goweewar

# Give user the chance to set the world id here
WORLD_ID="7e5016a4"
# Look for the line that has "export WEEWAR_GAME_ID=...." here so we can extract gameid from it
gameIdLine=$(ww new $WORLD_ID | grep "export WEEWAR_GAME_ID")
gameId=$(echo $gameIdLine | sed -e "s/.*export.WEEWAR_GAME_ID=//g")
export WEEWAR_GAME_ID=$gameId
export WEEWAR_CONFIRM=false
echo Created game for testing: $WEEWAR_GAME_ID

# ============================================================
# Round 1
# ============================================================

# --- Player 0 (Computer (A)) ---
ww build t:3,1 2
ww assert exists unit 3,1
ww assert unit 3,1 [type eq 2]
ww build t:1,0 2
ww assert exists unit 1,0
ww assert unit 1,0 [type eq 2]
ww endturn

# --- Player 1 (goweewar) ---
ww build t:0,5 8
ww assert exists unit 0,5
ww assert unit 0,5 [type eq 8]
ww build t:2,6 1
ww assert exists unit 2,6
ww assert unit 2,6 [type eq 1]
ww move 3,5 2,4
ww assert exists unit 2,4
ww assert notexists unit 3,5
ww endturn

# ============================================================
# Round 2
# ============================================================

# --- Player 0 (Computer (A)) ---
ww move 3,1 3,0
ww assert exists unit 3,0
ww assert notexists unit 3,1
ww build t:3,1 1
ww assert exists unit 3,1
ww assert unit 3,1 [type eq 1]
ww move 1,0 2,0
ww assert exists unit 2,0
ww assert notexists unit 1,0
ww build t:1,0 1
ww assert exists unit 1,0
ww assert unit 1,0 [type eq 1]
ww endturn

# --- Player 1 (goweewar) ---
ww move 0,5 0,4
ww assert exists unit 0,4
ww assert notexists unit 0,5
ww move 2,6 1,4
ww assert exists unit 1,4
ww assert notexists unit 2,6
ww build t:0,5 2
ww assert exists unit 0,5
ww assert unit 0,5 [type eq 2]
ww build t:2,6 1
ww assert exists unit 2,6
ww assert unit 2,6 [type eq 1]
ww endturn

# ============================================================
# Round 3
# ============================================================

# --- Player 0 (Computer (A)) ---
ww move 1,0 1,1
ww assert exists unit 1,1
ww assert notexists unit 1,0
ww build t:1,0 1
ww assert exists unit 1,0
ww assert unit 1,0 [type eq 1]
ww move 3,0 2,1
ww assert exists unit 2,1
ww assert notexists unit 3,0
ww move 3,1 3,0
ww assert exists unit 3,0
ww assert notexists unit 3,1
ww build t:3,1 2
ww assert exists unit 3,1
ww assert unit 3,1 [type eq 2]
ww move 2,0 0,1
ww assert exists unit 0,1
ww assert notexists unit 2,0
ww endturn

# --- Player 1 (goweewar) ---
ww attack 0,4 2,1
ww assert unit 2,1 [health lte 5]
ww move 2,4 3,3
ww assert exists unit 3,3
ww assert notexists unit 2,4
ww move 0,5 2,3
ww assert exists unit 2,3
ww assert notexists unit 0,5
ww move 1,4 1,3
ww assert exists unit 1,3
ww assert notexists unit 1,4
ww move 2,6 2,4
ww assert exists unit 2,4
ww assert notexists unit 2,6
ww build t:0,5 1
ww assert exists unit 0,5
ww assert unit 0,5 [type eq 1]
ww build t:2,6 1
ww assert exists unit 2,6
ww assert unit 2,6 [type eq 1]
ww endturn

# ============================================================
# Round 4
# ============================================================

# --- Player 0 (Computer (A)) ---
ww move 3,0 2,0
ww assert exists unit 2,0
ww assert notexists unit 3,0
ww move 0,1 0,2
ww assert exists unit 0,2
ww assert notexists unit 0,1
# Heal/Hold at 3,1
ww move 1,0 0,1
ww assert exists unit 0,1
ww assert notexists unit 1,0
ww build t:1,0 1
ww assert exists unit 1,0
ww assert unit 1,0 [type eq 1]
ww move 2,1 3,0
ww assert exists unit 3,0
ww assert notexists unit 2,1
ww move 1,1 2,1
ww assert exists unit 2,1
ww assert notexists unit 1,1
ww endturn

# --- Player 1 (goweewar) ---
# No actions
ww endturn

# ============================================================
# Final Board State (from Weewar coordinate section)
# ============================================================

# Weewar (col,row) -> Our (Q,R) mapping:
# Weewar (2,4)    -> Q,R=(  0,  4) : Unit type 8:Artillery (Basic), goweewar
# Weewar (3,3)    -> Q,R=(  2,  3) : Unit type 2:Soldier (Advanced), goweewar
# Weewar (4,3)    -> Q,R=(  3,  3) : Unit type 1:Soldier (Basic), goweewar
# Weewar (2,5)    -> Q,R=(  0,  5) : Unit type 1:Soldier (Basic), goweewar
# Weewar (3,1)    -> Q,R=(  3,  1) : Unit type 2:Soldier (Advanced), Computer (A)
# Weewar (1,0)    -> Q,R=(  1,  0) : Unit type 1:Soldier (Basic), Computer (A)
# Weewar (2,3)    -> Q,R=(  1,  3) : Unit type 1:Soldier (Basic), goweewar
# Weewar (5,6)    -> Q,R=(  2,  6) : Unit type 1:Soldier (Basic), goweewar
# Weewar (4,4)    -> Q,R=(  2,  4) : Unit type 1:Soldier (Basic), goweewar
# Weewar (2,0)    -> Q,R=(  2,  0) : Unit type 1:Soldier (Basic), Computer (A)
# Weewar (1,2)    -> Q,R=(  0,  2) : Unit type 2:Soldier (Advanced), Computer (A)
# Weewar (0,1)    -> Q,R=(  0,  1) : Unit type 1:Soldier (Basic), Computer (A)
# Weewar (3,0)    -> Q,R=(  3,  0) : Unit type 2:Soldier (Advanced), Computer (A)
# Weewar (2,1)    -> Q,R=(  2,  1) : Unit type 1:Soldier (Basic), Computer (A)

# Final board state assertions:
ww assert unit 0,4 [type eq 8, player eq 2, health eq 10]
ww assert unit 2,3 [type eq 2, player eq 2, health eq 10]
ww assert unit 3,3 [type eq 1, player eq 2, health eq 10]
ww assert unit 0,5 [type eq 1, player eq 2, health eq 10]
ww assert unit 3,1 [type eq 2, player eq 1, health eq 10]
ww assert unit 1,0 [type eq 1, player eq 1, health eq 10]
ww assert unit 1,3 [type eq 1, player eq 2, health eq 10]
ww assert unit 2,6 [type eq 1, player eq 2, health eq 10]
ww assert unit 2,4 [type eq 1, player eq 2, health eq 10]
ww assert unit 2,0 [type eq 1, player eq 1, health eq 9]
ww assert unit 0,2 [type eq 2, player eq 1, health eq 10]
ww assert unit 0,1 [type eq 1, player eq 1, health eq 10]
ww assert unit 3,0 [type eq 2, player eq 1, health eq 5]
ww assert unit 2,1 [type eq 1, player eq 1, health eq 9]

echo 'All assertions passed!'
