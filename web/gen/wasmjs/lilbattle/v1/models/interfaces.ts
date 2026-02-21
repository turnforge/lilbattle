// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Any, FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


/**
 * Crossing types for terrain improvements (roads on land, bridges on water)
 */
export enum CrossingType {
  CROSSING_TYPE_UNSPECIFIED = 0,
  CROSSING_TYPE_ROAD = 1,
  CROSSING_TYPE_BRIDGE = 2,
}

/**
 * Terrain type classification - used for gameplay logic
 */
export enum TerrainType {
  TERRAIN_TYPE_UNSPECIFIED = 0,
  TERRAIN_TYPE_CITY = 1,
  TERRAIN_TYPE_NATURE = 2,
  TERRAIN_TYPE_BRIDGE = 3,
  TERRAIN_TYPE_WATER = 4,
  TERRAIN_TYPE_ROAD = 5,
}

/**
 * /////// Game related models
 */
export enum GameStatus {
  GAME_STATUS_UNSPECIFIED = 0,
  GAME_STATUS_PLAYING = 1,
  GAME_STATUS_PAUSED = 2,
  GAME_STATUS_ENDED = 3,
}


export enum PathDirection {
  PATH_DIRECTION_UNSPECIFIED = 0,
  PATH_DIRECTION_LEFT = 1,
  PATH_DIRECTION_TOP_LEFT = 2,
  PATH_DIRECTION_TOP_RIGHT = 3,
  PATH_DIRECTION_RIGHT = 4,
  PATH_DIRECTION_BOTTOM_RIGHT = 5,
  PATH_DIRECTION_BOTTOM_LEFT = 6,
}


export enum IndexStatus {
  INDEX_STATUS_UNSPECIFIED = 0,
  INDEX_STATUS_PENDING = 1,
  INDEX_STATUS_INDEXING = 2,
  INDEX_STATUS_COMPLETED = 3,
  INDEX_STATUS_FAILED = 4,
}


export enum RunState {
  RUN_STATE_UNSPECIFIED = 0,
  RUN_STATE_STARTED = 1,
  RUN_STATE_FINISHED = 2,
}


export enum Type {
  TYPE_UNSPECIFIED = 0,
  TYPE_PATH = 1,
  TYPE_SVG = 2,
  TYPE_DATA_URL = 3,
}



export interface IndexInfo {
  /** We maintain an IndexInfo for each type of "indexing" operation needed
 For example one update may change the keywords (so we need to update indexes for search)
 Another might update the "units" so we may need a new screenshot
 Each one's updates - updated and indexed timestamps separately so they can be tracked sepately */
  lastUpdatedAt?: Timestamp;
  lastIndexedAt?: Timestamp;
  needsIndexing: boolean;
}



export interface Pagination {
  /** *
 Instead of an offset an abstract  "page" key is provided that offers
 an opaque "pointer" into some offset in a result set. */
  pageKey: string;
  /** *
 If a pagekey is not supported we can also support a direct integer offset
 for cases where it makes sense. */
  pageOffset: number;
  /** *
 Number of results to return. */
  pageSize: number;
}



export interface PaginationResponse {
  /** *
 The key/pointer string that subsequent List requests should pass to
 continue the pagination. */
  nextPageKey: string;
  /** *
 Also support an integer offset if possible */
  nextPageOffset: number;
  /** *
 Whether theere are more results. */
  hasMore: boolean;
  /** *
 Total number of results. */
  totalResults: number;
}



export interface World {
  createdAt?: Timestamp;
  updatedAt?: Timestamp;
  /** Version for Optimistic concurrent locking */
  version: number;
  /** Unique ID for the world */
  id: string;
  /** User that created the world */
  creatorId: string;
  /** Name if items have names */
  name: string;
  /** Description if world has a description */
  description: string;
  /** Some tags */
  tags: string[];
  /** A possible image url */
  imageUrl: string;
  /** Difficulty - example attribute */
  difficulty: string;
  /** URL to screenshot/preview image (defaults to /worlds/{id}/screenshots/{screenshotName})
 Can be overridden to point to CDN or external hosting */
  previewUrls: string[];
  /** Default game configs */
  defaultGameConfig?: GameConfiguration;
  searchIndexInfo?: IndexInfo;
}



export interface WorldData {
  /** New map-based storage (key = "q,r" coordinate string) */
  tilesMap: Record<string, Tile>;
  unitsMap: Record<string, Unit>;
  /** When this world data was updated (may have happened without world updating) */
  screenshotIndexInfo?: IndexInfo;
  /** We will only update if hash's are different */
  contentHash: string;
  /** Version for Optimistic concurrent locking */
  version: number;
  /** Improvement layer - crossings (roads on land, bridges on water) */
  crossings: Record<string, Crossing>;
}


/**
 * Crossing with explicit connectivity data
 Each crossing stores which of its 6 hex neighbors it connects to
 */
export interface Crossing {
  type: CrossingType;
  /** 6 booleans for hex neighbors in order matching AxialNeighborDeltas:
 0: LEFT (-1,0), 1: TOP_LEFT (0,-1), 2: TOP_RIGHT (1,-1),
 3: RIGHT (1,0), 4: BOTTOM_RIGHT (0,1), 5: BOTTOM_LEFT (-1,1) */
  connectsTo: boolean[];
}



export interface Tile {
  /** Q and R in Cubed coordinates */
  q: number;
  r: number;
  tileType: number;
  /** Whether the tile itself belongs to a player */
  player: number;
  shortcut: string;
  /** Keep track of turns when the move was last made and when a "top up" was last done on/for this tile.
 This helps us not having to "top up" or "reset" the stats at the end
 of each turn.  Instead as the game turn is incremented we can do a 
 lazy reset for any unit or tile where unit_or_tile.last_toppedup_turn < game.curent_turn
 So we just have to increment the game_turn and the unit is automaticaly flagged as
 needing a top up of its health/balance/movement etc */
  lastActedTurn: number;
  lastToppedupTurn: number;
}



export interface Unit {
  /** Q and R in Cubed coordinates */
  q: number;
  r: number;
  player: number;
  unitType: number;
  shortcut: string;
  /** Runtime state fields */
  availableHealth: number;
  distanceLeft: number;
  /** Keep track of turns when the move was last made and when a "top up" was last done on/for this tile.
 This helps us not having to "top up" or "reset" the stats at the end
 of each turn.  Instead as the game turn is incremented we can do a 
 lazy reset for any unit or tile where unit_or_tile.last_toppedup_turn < game.curent_turn
 So we just have to increment the game_turn and the unit is automaticaly flagged as
 needing a top up of its health/balance/movement etc */
  lastActedTurn: number;
  lastToppedupTurn: number;
  /** Details around wound bonus tracking for this turn */
  attacksReceivedThisTurn: number;
  attackHistory?: AttackRecord[];
  /** Action progression tracking - index into UnitDefinition.action_order
 Indicates which step in the action sequence the unit is currently on
 Reset to 0 at turn start via TopUpUnitIfNeeded() */
  progressionStep: number;
  /** When current step has pipe-separated alternatives (e.g., "attack|capture"),
 this tracks which alternative the user chose, preventing switching mid-step
 Cleared when advancing to next step */
  chosenAlternative: string;
  /** Turn when this unit started capturing a building (0 = not capturing)
 Capture completes at the start of the capturing player's next turn
 if the unit is still alive on the tile */
  captureStartedTurn: number;
}



export interface AttackRecord {
  q: number;
  r: number;
  isRanged: boolean;
  turnNumber: number;
}


/**
 * Rules engine terrain definition
 */
export interface TerrainDefinition {
  id: number;
  name: string;
  /** double base_move_cost = 3;     // Base movement cost
 double defense_bonus = 4;      // Defense bonus multiplier (0.0 to 1.0) */
  type: number;
  description: string;
  /** How this terrain impacts */
  unitProperties: Record<number, TerrainUnitProperties>;
  /** List of units that can be built on this terrain */
  buildableUnitIds: number[];
  incomePerTurn: number;
}


/**
 * Rules engine unit definition
 */
export interface UnitDefinition {
  id: number;
  name: string;
  description: string;
  health: number;
  coins: number;
  movementPoints: number;
  retreatPoints: number;
  defense: number;
  attackRange: number;
  minAttackRange: number;
  splashDamage: number;
  terrainProperties: Record<number, TerrainUnitProperties>;
  properties: string[];
  /** Unit classification for attack calculations */
  unitClass: string;
  unitTerrain: string;
  /** Attack table: base attack values against different unit classes
 Key format: "Light:Air", "Heavy:Land", "Stealth:Water", etc.
 Value 0 or missing key means "n/a" (cannot attack) */
  attackVsClass: Record<string, number>;
  /** Ordered list of allowed actions this turn
 Examples:
   ["move", "attack"] - can move then attack
   ["move", "attack|capture"] - can move then either attack or capture
   ["attack"] - can only attack (no movement)
 Default if empty: ["move", "attack|capture"] */
  actionOrder: string[];
  /** How many times each action type can be performed per turn
 Key: action name, Value: max count
 Example: {"attack": 2} means can attack twice
 Default if not specified: 1 per action type */
  actionLimits: Record<string, number>;
  /** Fix value for units that can repair other units (Medic, Engineer, etc.)
 Used in fix calculation: p = 0.05 * fix_value
 Default 0 means unit cannot fix */
  fixValue: number;
}


/**
 * Properties that are specific to unit on a particular terrain
 */
export interface TerrainUnitProperties {
  terrainId: number;
  unitId: number;
  movementCost: number;
  healingBonus: number;
  canBuild: boolean;
  canCapture: boolean;
  attackBonus: number;
  defenseBonus: number;
  attackRange: number;
  minAttackRange: number;
}


/**
 * Properties for unit-vs-unit combat interactions
 */
export interface UnitUnitProperties {
  attackerId: number;
  defenderId: number;
  attackOverride?: number | undefined;
  defenseOverride?: number | undefined;
  damage?: DamageDistribution;
}


/**
 * Damage distribution for combat calculations
 */
export interface DamageDistribution {
  minDamage: number;
  maxDamage: number;
  expectedDamage: number;
  ranges?: DamageRange[];
}



export interface DamageRange {
  minValue: number;
  maxValue: number;
  probability: number;
}


/**
 * Main rules engine definition - centralized source of truth
 */
export interface RulesEngine {
  /** Core entity definitions */
  units: Record<number, UnitDefinition>;
  terrains: Record<number, TerrainDefinition>;
  /** Centralized property definitions (source of truth)
 Key format: "terrain_id:unit_id" (e.g., "1:3" for terrain 1, unit 3) */
  terrainUnitProperties: Record<string, TerrainUnitProperties>;
  /** Key format: "attacker_id:defender_id" (e.g., "1:2" for unit 1 attacking unit 2) */
  unitUnitProperties: Record<string, UnitUnitProperties>;
  /** Terrain type classifications (terrain_id -> TerrainType)
 Used to determine if a terrain is city, nature, bridge, water, or road */
  terrainTypes: Record<number, any>;
}


/**
 * Describes a game and its metadata
 */
export interface Game {
  createdAt?: Timestamp;
  updatedAt?: Timestamp;
  /** Version number for optimistic locking */
  version: number;
  /** Unique ID for the game */
  id: string;
  /** User who started/created the game */
  creatorId: string;
  /** The world this game was created from */
  worldId: string;
  /** Name if items have names */
  name: string;
  /** Description if game has a description */
  description: string;
  /** Some tags */
  tags: string[];
  /** A possible image url */
  imageUrl: string;
  /** Difficulty - example attribute */
  difficulty: string;
  /** Game configuration */
  config?: GameConfiguration;
  /** URL to screenshot/preview image (defaults to /games/{id}/screenshots/{screenshotName})
 Can be overridden to point to CDN or external hosting */
  previewUrls: string[];
  searchIndexInfo?: IndexInfo;
}



export interface GameConfiguration {
  /** Player configuration */
  players?: GamePlayer[];
  /** Team configuration */
  teams?: GameTeam[];
  /** Various kinds of per turn income configs */
  incomeConfigs?: IncomeConfig;
  /** Game settings */
  settings?: GameSettings;
}



export interface IncomeConfig {
  /** How much starting coins to give each player at the start of the agme */
  startingCoins: number;
  /** Income each player just for being in the game each turn */
  gameIncome: number;
  /** Income from each landbase per turn */
  landbaseIncome: number;
  /** Income from each navalbase per turn */
  navalbaseIncome: number;
  /** Income from each airport base per turn */
  airportbaseIncome: number;
  /** Income from each missile silo per turn */
  missilesiloIncome: number;
  /** Income from each mine per turn */
  minesIncome: number;
}



export interface GamePlayer {
  /** Player ID (1-based) */
  playerId: number;
  /** ID of the system user that is assigned to this game player.  This is the "auth" user */
  userId: string;
  /** Player type */
  playerType: string;
  /** Player color */
  color: string;
  /** Team ID (0 = no team, 1+ = team number) */
  teamId: number;
  /** Nickname for the player in this game */
  name: string;
  /** Whether play is still in the game - can this just be inferred? */
  isActive: boolean;
  /** How many coins the player started off with */
  startingCoins: number;
}



export interface GameTeam {
  /** ID of the team within the game (unique to the game) */
  teamId: number;
  /** Name of the team - in a game */
  name: string;
  /** Just a color for this team */
  color: string;
  /** Whether team has active players - can also be inferred */
  isActive: boolean;
}



export interface GameSettings {
  /** List of allowed unit type IDs */
  allowedUnits: number[];
  /** Turn time limit in seconds (0 = no limit) */
  turnTimeLimit: number;
  /** Team mode */
  teamMode: string;
  /** Maximum number of turns (0 = unlimited) */
  maxTurns: number;
}


/**
 * Runtime state for a player during the game
 This is separate from GamePlayer (which is player configuration)
 PlayerState is indexed by player_id in the player_states map
 */
export interface PlayerState {
  /** Current coin balance (changes during gameplay via building, income, etc.) */
  coins: number;
  /** Whether player is still active in the game (not eliminated) */
  isActive: boolean;
}


/**
 * Holds the game's Active/Current state (eg world state)
 */
export interface GameState {
  updatedAt?: Timestamp;
  /** ID of the game whos state is being tracked */
  gameId: string;
  turnCounter: number;
  currentPlayer: number;
  /** Current world state */
  worldData?: WorldData;
  /** Current state hash for validation */
  stateHash: string;
  /** Version number for optimistic locking */
  version: number;
  status: GameStatus;
  /** Only set after a win has been possible */
  finished: boolean;
  winningPlayer: number;
  winningTeam: number;
  currentGroupNumber: number;
  /** Per-player runtime state, keyed by player_id (1-based)
 This holds mutable player state like coins that changes during gameplay */
  playerStates: Record<number, PlayerState>;
}


/**
 * Holds the game's move history (can be used as a replay log)
 */
export interface GameMoveHistory {
  /** Move history for the game */
  gameId: string;
  /** Each entry in our history is a "group" of moves */
  groups?: GameMoveGroup[];
}


/**
 * A move group - we can allow X moves in one "tick"
 */
export interface GameMoveGroup {
  /** When the moves happened (or were submitted) */
  startedAt?: Timestamp;
  endedAt?: Timestamp;
  /** Group number within the game - will be monotonically increasing */
  groupNumber: number;
  /** *
 List of moves to add - */
  moves?: GameMove[];
}


/**
 * *
 Represents a single move which can be one of many actions in the game
 */
export interface GameMove {
  player: number;
  /** Generated by the server */
  groupNumber: number;
  /** Generated by the server - will be monotonically increasing within the group */
  moveNumber: number;
  timestamp?: Timestamp;
  moveUnit?: MoveUnitAction;
  attackUnit?: AttackUnitAction;
  endTurn?: EndTurnAction;
  buildUnit?: BuildUnitAction;
  captureBuilding?: CaptureBuildingAction;
  healUnit?: HealUnitAction;
  fixUnit?: FixUnitAction;
  /** A monotonically increasing and unique (within the game) sequence number for the move
 This is generated by the server */
  sequenceNum: number;
  /** Whether the result is permenant and can be undone.
 Just moving a unit for example is not permanent, but attacking a unit
 would be (ie a player cannot undo it).  This is also determined by the server/validator */
  isPermanent: boolean;
  /** The corresponding "result" for the move.  This can be "proposed" or can be evaluated.
 Keeping this colocated with the Move for consistency and simplicity */
  changes?: WorldChange[];
  /** Human redable description for say recording "commands" if any */
  description: string;
}


/**
 * A unified "Position" type that can be used to 
 specify locations via "string shortcuts" like A1, "3,2", "r2,4" (for row/col)
 or even "relative" positions like "L,TL,TR,R"  in the shortcut field.
 Or string q/r coordinates in the q and r fields.  This can also be used 
 in the "response" to resolve a shortcut -> q,r
 */
export interface Position {
  label: string;
  q: number;
  r: number;
}


/**
 * *
 Move unit from one position to another
 */
export interface MoveUnitAction {
  from?: Position;
  to?: Position;
  /** Optional fields that can be used for showing move options as well as debugging */
  movementCost: number;
  /** Debug fields */
  reconstructedPath?: Path;
}


/**
 * *
 Attack with one unit against another
 */
export interface AttackUnitAction {
  attacker?: Position;
  defender?: Position;
  /** Optional fields for presenting during "options" and debugging */
  targetUnitType: number;
  targetUnitHealth: number;
  canAttack: boolean;
  damageEstimate: number;
}


/**
 * *
 An action to build a unit (at a city tile)
 */
export interface BuildUnitAction {
  pos?: Position;
  unitType: number;
  cost: number;
}


/**
 * *
 A move where a unit can capture a building
 */
export interface CaptureBuildingAction {
  pos?: Position;
  tileType: number;
}


/**
 * *
 End current player's turn
 */
export interface EndTurnAction {
}


/**
 * *
 Heal a unit - player manually chooses to heal instead of attacking/moving
 Auto-healing at turn start is handled separately in TopUpUnitIfNeeded
 */
export interface HealUnitAction {
  pos?: Position;
  healAmount: number;
}


/**
 * *
 Fix (repair) another friendly unit - used by Medic, Engineer, Stratotanker, Tugboat, Aircraft Carrier
 The fixer must be adjacent to the target unit
 */
export interface FixUnitAction {
  fixer?: Position;
  target?: Position;
  fixAmount: number;
}


/**
 * *
 Represents a change to the game world
 */
export interface WorldChange {
  unitMoved?: UnitMovedChange;
  unitDamaged?: UnitDamagedChange;
  unitKilled?: UnitKilledChange;
  playerChanged?: PlayerChangedChange;
  unitBuilt?: UnitBuiltChange;
  coinsChanged?: CoinsChangedChange;
  tileCaptured?: TileCapturedChange;
  captureStarted?: CaptureStartedChange;
  unitHealed?: UnitHealedChange;
  unitFixed?: UnitFixedChange;
}


/**
 * *
 A unit was healed
 */
export interface UnitHealedChange {
  previousUnit?: Unit;
  updatedUnit?: Unit;
  healAmount: number;
}


/**
 * *
 A unit was fixed (repaired) by another unit
 */
export interface UnitFixedChange {
  fixerUnit?: Unit;
  previousTarget?: Unit;
  updatedTarget?: Unit;
  fixAmount: number;
}


/**
 * *
 A unit moved from one position to another
 */
export interface UnitMovedChange {
  /** Complete unit state before the move */
  previousUnit?: Unit;
  /** Complete unit state after the move (includes updated position, distanceLeft, etc.) */
  updatedUnit?: Unit;
}


/**
 * *
 A unit took damage
 */
export interface UnitDamagedChange {
  /** Complete unit state before taking damage */
  previousUnit?: Unit;
  /** Complete unit state after taking damage */
  updatedUnit?: Unit;
}


/**
 * *
 A unit was killed
 */
export interface UnitKilledChange {
  /** Complete unit state before being killed */
  previousUnit?: Unit;
}


/**
 * *
 Active player changed
 */
export interface PlayerChangedChange {
  previousPlayer: number;
  newPlayer: number;
  previousTurn: number;
  newTurn: number;
  /** Units that had their movement/health reset for the new turn */
  resetUnits?: Unit[];
}


/**
 * *
 A new unit was built at a tile
 */
export interface UnitBuiltChange {
  /** The newly created unit */
  unit?: Unit;
  /** Tile coordinates where unit was built */
  tileQ: number;
  tileR: number;
  /** Cost in coins */
  coinsCost: number;
  /** Player's remaining coins after build */
  playerCoins: number;
}


/**
 * *
 A player's coin balance changed
 */
export interface CoinsChangedChange {
  /** Which player's coins changed */
  playerId: number;
  /** Previous coin balance */
  previousCoins: number;
  /** New coin balance */
  newCoins: number;
  /** Reason for change: "build", "income", "repair", etc. */
  reason: string;
}


/**
 * *
 A tile was captured by a unit
 */
export interface TileCapturedChange {
  /** The unit that captured the tile */
  capturingUnit?: Unit;
  /** Tile coordinates */
  tileQ: number;
  tileR: number;
  /** Tile type */
  tileType: number;
  /** Previous owner (0 for neutral) */
  previousOwner: number;
  /** New owner */
  newOwner: number;
}


/**
 * *
 A unit started capturing a building (capture not yet complete)
 */
export interface CaptureStartedChange {
  /** The unit starting the capture */
  capturingUnit?: Unit;
  /** Tile coordinates */
  tileQ: number;
  tileR: number;
  /** Tile type */
  tileType: number;
  /** Current owner (0 for neutral) */
  currentOwner: number;
}


/**
 * Compact representation of all reachable paths from a source
 */
export interface AllPaths {
  /** Starting coordinate for all paths */
  sourceQ: number;
  sourceR: number;
  /** Map of edges: key is "toQ,toR" for quick parent lookup
 Each edge represents the optimal way to reach 'to' from its parent */
  edges: Record<string, PathEdge>;
}


/**
 * A single edge in a path with movement details
 */
export interface PathEdge {
  fromQ: number;
  fromR: number;
  toQ: number;
  toR: number;
  movementCost: number;
  totalCost: number;
  terrainType: string;
  explanation: string;
  isOccupied: boolean;
}


/**
 * Full path from source to destination (constructed on-demand from AllPaths)
 */
export interface Path {
  /** Edges in order from source to destination */
  edges?: PathEdge[];
  /** len(directions) = len(edges) - 1
 and directions[i] = direction from edge[i - 1] -> edge[i] */
  directions: PathDirection[];
  /** Sum of all edge costs */
  totalCost: number;
}



export interface File {
  /** Unique path of the file */
  path: string;
  /** Type of file */
  contentType: string;
  fileSize: number;
  /** Public visibility or not */
  isPublic: boolean;
  /** When the last indexing was queued */
  createdAt?: Timestamp;
  /** when the last time the file was updated */
  updatedAt?: Timestamp;
  /** The download/get URL for this file (public URL or default signed URL) */
  downloadUrl: string;
  /** Signed URLs with different expiries (e.g., "15m", "1h", "24h")
 Only populated when requested */
  signedUrls: Record<string, string>;
}



export interface PutFileRequest {
  file?: File;
  content: Uint8Array;
}



export interface PutFileResponse {
  file?: File;
}



export interface GetFileRequest {
  path: string;
  /** If true, populate signed_urls in the File response */
  includeSignedUrls: boolean;
}



export interface GetFileResponse {
  file?: File;
}



export interface DeleteFileRequest {
  path: string;
}



export interface DeleteFileResponse {
  file?: File;
}



export interface ListFilesRequest {
  /** Directory path to list files from (relative to base path) */
  path: string;
  /** Pagination info */
  pagination?: Pagination;
  /** If true, populate signed_urls in each File response */
  includeSignedUrls: boolean;
}



export interface ListFilesResponse {
  items?: File[];
  pagination?: PaginationResponse;
}


/**
 * Request messages
 */
export interface ListGamesRequest {
  /** Pagination info */
  pagination?: Pagination;
  /** May be filter by owner id */
  ownerId: string;
}



export interface ListGamesResponse {
  items?: Game[];
  pagination?: PaginationResponse;
}



export interface GetGameRequest {
  id: string;
  version: string;
}



export interface GetGameResponse {
  game?: Game;
  state?: GameState;
  history?: GameMoveHistory;
}



export interface GetGameContentRequest {
  id: string;
  version: string;
}



export interface GetGameContentResponse {
  lilbattleContent: string;
  recipeContent: string;
  readmeContent: string;
}



export interface UpdateGameRequest {
  /** Game id to modify */
  gameId: string;
  /** *
 Game being updated */
  newGame?: Game;
  /** New world state to save */
  newState?: GameState;
  /** History to save */
  newHistory?: GameMoveHistory;
  /** *
 Mask of fields being updated in this Game to make partial changes. */
  updateMask?: FieldMask;
}


/**
 * *
 The request for (partially) updating an Game.
 */
export interface UpdateGameResponse {
  /** *
 Game being updated */
  game?: Game;
}


/**
 * *
 Request to delete an game.
 */
export interface DeleteGameRequest {
  /** *
 ID of the game to be deleted. */
  id: string;
}


/**
 * *
 Game deletion response
 */
export interface DeleteGameResponse {
}


/**
 * *
 Request to batch get games
 */
export interface GetGamesRequest {
  /** *
 IDs of the game to be fetched */
  ids: string[];
}


/**
 * *
 Game batch-get response
 */
export interface GetGamesResponse {
  games: Record<string, Game>;
}


/**
 * *
 Game creation request object
 */
export interface CreateGameRequest {
  /** *
 Game being updated */
  game?: Game;
}


/**
 * *
 Response of an game creation.
 */
export interface CreateGameResponse {
  /** *
 Game being created */
  game?: Game;
  /** The starting game state */
  gameState?: GameState;
  /** *
 Error specific to a field if there are any errors. */
  fieldErrors: Record<string, string>;
}


/**
 * *
 Request to add moves to a game
 The model is that a game in each "tick" can handle multiple moves (by possibly various players).
 It is upto the move manager/processor in the game to ensure the "transaction" of moves is handled
 atomically.

 For example we may have 3 moves where first two units are moved to a common location
 and then they attack another unit.  Here If we treat it as a single unit attacking it
 will have different outcomes than a "combined" attack.
 */
export interface ProcessMovesRequest {
  /** *
 Game ID to add moves to */
  gameId: string;
  /** *
 List of moves to add */
  moves?: GameMove[];
  /** *
 The player can submit a list of "Expected" changes when in local-first mode
 If this is list provided the server will validate it - either via the coordinator
 or by itself.  If it is not provided then the server will validate it and return
 the changes. */
  expectedResponse?: ProcessMovesResponse;
  /** Whether to only perform a dryrun and return results instead of comitting it */
  dryRun: boolean;
}


/**
 * *
 Response after adding moves to game.
 */
export interface ProcessMovesResponse {
  /** *
 Returns the moves that were passed in along wth changes and other data filled in. */
  moves?: GameMove[];
}


/**
 * *
 Request to get the game's latest state
 */
export interface GetGameStateRequest {
  /** *
 Game ID to add moves to */
  gameId: string;
}


/**
 * *
 Response holding latest game state
 */
export interface GetGameStateResponse {
  state?: GameState;
}


/**
 * *
 Request to list moves for a game
 */
export interface ListMovesRequest {
  /** *
 Game ID to add moves to */
  gameId: string;
  /** Gets moves >= from_group */
  fromGroup: number;
  /** Gets moves <= to_group */
  toGroup: number;
}


/**
 * *
 Response after adding moves to game.
 */
export interface ListMovesResponse {
  /** Whether there are more moves before this */
  hasMore: boolean;
  moveGroups?: GameMoveGroup[];
}


/**
 * *
 Request to get all available options at a position
 */
export interface GetOptionsAtRequest {
  gameId: string;
  pos?: Position;
}


/**
 * *
 Response with all available options at a position
 */
export interface GetOptionsAtResponse {
  options?: GameOption[];
  currentPlayer: number;
  gameInitialized: boolean;
  /** A Path from source to dest along with cost on each tile for tracking */
  allPaths?: AllPaths;
}


/**
 * *
 A single game option available at a position
 */
export interface GameOption {
  move?: MoveUnitAction;
  attack?: AttackUnitAction;
  build?: BuildUnitAction;
  capture?: CaptureBuildingAction;
  endTurn?: EndTurnAction;
  heal?: HealUnitAction;
}


/**
 * *
 Request for simulating combat between two units
 */
export interface SimulateAttackRequest {
  attackerUnitType: number;
  attackerTerrain: number;
  attackerHealth: number;
  defenderUnitType: number;
  defenderTerrain: number;
  defenderHealth: number;
  woundBonus: number;
  numSimulations: number;
}


/**
 * *
 Response containing damage distribution statistics
 */
export interface SimulateAttackResponse {
  /** Damage distributions: damage_value -> number_of_occurrences */
  attackerDamageDistribution: Record<number, number>;
  defenderDamageDistribution: Record<number, number>;
  /** Statistical summary */
  attackerMeanDamage: number;
  defenderMeanDamage: number;
  attackerKillProbability: number;
  defenderKillProbability: number;
}


/**
 * *
 Request for simulating fix (repair) between two units
 */
export interface SimulateFixRequest {
  fixingUnitType: number;
  fixingUnitHealth: number;
  injuredUnitType: number;
  numSimulations: number;
}


/**
 * *
 Response containing health restoration distribution statistics
 */
export interface SimulateFixResponse {
  /** Health restoration distribution: health_restored -> number_of_occurrences */
  healingDistribution: Record<number, number>;
  /** Statistical summary */
  meanHealing: number;
  /** The fix value (F) of the fixing unit type */
  fixValue: number;
}


/**
 * *
 Request to join a game as a player
 User must be authenticated to join a game.
 The player slot must be "open" (player_type = "open") to be joinable.
 */
export interface JoinGameRequest {
  /** ID of the game to join */
  gameId: string;
  /** The player ID (slot) to join as (1-based)
 Must be an open slot (player_type = "open") */
  playerId: number;
}


/**
 * *
 Response after joining a game
 */
export interface JoinGameResponse {
  /** The updated game with the user assigned to the player slot */
  game?: Game;
  /** The player ID that was joined */
  playerId: number;
}



export interface EmptyRequest {
}



export interface EmptyResponse {
}


/**
 * Request to fetch data from a URL
 */
export interface SetContentRequest {
  innerHtml: string;
}


/**
 * Response from fetch
 */
export interface SetContentResponse {
}



export interface ShowBuildOptionsRequest {
  innerHtml: string;
  hide: boolean;
  q: number;
  r: number;
}



export interface ShowBuildOptionsResponse {
}


/**
 * Request to fetch data from a URL
 */
export interface LogMessageRequest {
  message: string;
}


/**
 * Response from fetch
 */
export interface LogMessageResponse {
}


/**
 * Request to fetch data from a URL
 */
export interface SetGameStateRequest {
  game?: Game;
  state?: GameState;
}


/**
 * Response from fetch
 */
export interface SetGameStateResponse {
}


/**
 * Request to update game UI status (current player, turn counter)
 */
export interface UpdateGameStatusRequest {
  currentPlayer: number;
  turnCounter: number;
}



export interface UpdateGameStatusResponse {
}


/**
 * Request to set a tile at a specific coordinate
 */
export interface SetTileAtRequest {
  q: number;
  r: number;
  tile?: Tile;
}



export interface SetTileAtResponse {
}


/**
 * Request to set a unit at a specific coordinate
 */
export interface SetUnitAtRequest {
  q: number;
  r: number;
  unit?: Unit;
  flash: boolean;
  appear: boolean;
}



export interface SetUnitAtResponse {
}


/**
 * Request to remove a tile at a specific coordinate
 */
export interface RemoveTileAtRequest {
  q: number;
  r: number;
  animate: boolean;
}



export interface RemoveTileAtResponse {
}


/**
 * Request to remove a unit at a specific coordinate
 */
export interface RemoveUnitAtRequest {
  q: number;
  r: number;
  animate: boolean;
}



export interface RemoveUnitAtResponse {
}


/**
 * Request to show highlights on the game board
 */
export interface ShowHighlightsRequest {
  highlights?: HighlightSpec[];
}



export interface ShowHighlightsResponse {
}


/**
 * Specification for a single highlight
 */
export interface HighlightSpec {
  q: number;
  r: number;
  type: string;
  move?: MoveUnitAction;
  attack?: AttackUnitAction;
  build?: BuildUnitAction;
  capture?: CaptureBuildingAction;
  player: number;
}


/**
 * Request to clear highlights
 */
export interface ClearHighlightsRequest {
  types: string[];
}



export interface ClearHighlightsResponse {
}


/**
 * Request to show a path on the game board
 */
export interface ShowPathRequest {
  coords: number[];
  color: number;
  thickness: number;
}



export interface ShowPathResponse {
}


/**
 * Request to clear paths
 */
export interface ClearPathsRequest {
}



export interface ClearPathsResponse {
}


/**
 * Request to animate unit movement along a path
 */
export interface MoveUnitRequest {
  unit?: Unit;
  path?: HexCoord[];
}



export interface MoveUnitResponse {
}


/**
 * Hex coordinate for paths
 */
export interface HexCoord {
  q: number;
  r: number;
}


/**
 * Request to show attack effect animation
 */
export interface ShowAttackEffectRequest {
  fromQ: number;
  fromR: number;
  toQ: number;
  toR: number;
  damage: number;
  splashTargets?: SplashTarget[];
}



export interface SplashTarget {
  q: number;
  r: number;
  damage: number;
}



export interface ShowAttackEffectResponse {
}


/**
 * Request to show heal effect animation
 */
export interface ShowHealEffectRequest {
  q: number;
  r: number;
  amount: number;
}



export interface ShowHealEffectResponse {
}


/**
 * Request to show capture effect animation
 */
export interface ShowCaptureEffectRequest {
  q: number;
  r: number;
}



export interface ShowCaptureEffectResponse {
}


/**
 * Request to set allowed panels and their order
 */
export interface SetAllowedPanelsRequest {
  panelIds: string[];
}



export interface SetAllowedPanelsResponse {
}



export interface IndexState {
  entityType: string;
  entityId: string;
  /** eg "screenshots", "keywords" etc
 EntityType + EntityId + IndexType should be  unique */
  indexType: string;
  /** When the last indexing was queued */
  createdAt?: Timestamp;
  /** when the last time the entity was recorded for an update (means it is eligible for a re-indexing) */
  updatedAt?: Timestamp;
  /** When did the last indexing finish */
  indexedAt?: Timestamp;
  /** Whether indexing is needed or not */
  needsIndexing: boolean;
  /** "queued/pending", "indexing", "completed", "failed" */
  status: IndexStatus;
  /** If there was an error in the last indexing */
  lastError: string;
  /** A way to ignore multiple requests if they are updates but
 nothing has changed */
  idempotencyKey: string;
  retryCount: number;
}



export interface EnsureIndexStateRequest {
  indexState?: IndexState;
  /** *
 Mask of fields being updated in this Game to make partial changes. */
  updateMask?: FieldMask;
}



export interface EnsureIndexStateResponse {
  indexState?: IndexState;
}



export interface GetIndexStatesRequest {
  entityType: string;
  entityId: string;
  /** Optional - can be used to get "all" indexer states or just once specified here */
  indexTypes: string[];
}



export interface IndexStateList {
  states?: IndexState[];
}



export interface GetIndexStatesResponse {
  states: Record<string, IndexState>;
}



export interface ListIndexStatesRequest {
  entityType: string;
  /** Get records indexed "before" this time */
  updatedBefore?: Timestamp | undefined;
  /** Get records updated "after" this time */
  updatedAfter?: Timestamp | undefined;
  /** Filter by index types or get all */
  indexTypes: string[];
  /** "id" or "indexed_at" */
  orderBy: string;
  /** limit to max items */
  count: number;
}



export interface ListIndexStatesResponse {
  items?: IndexState[];
  /** How to identify the next "page" in this list */
  nextPageKey: string;
}



export interface DeleteIndexStatesRequest {
  entityType: string;
  entityId: string;
  /** Optional - can be used to get "all" indexer states or just once specified here */
  indexTypes: string[];
}



export interface DeleteIndexStatesResponse {
}


/**
 * Request messages
 */
export interface IndexRecord {
  entityId: string;
  updatedAt?: Timestamp;
  entityData?: Any;
  indexerTypes: string[];
}


/**
 * Each IndexRecords gets its own "long running operation" 
 so we can track how things are
 */
export interface IndexRecordsLRO {
  lroId: string;
  /** Single entity type in a index request */
  entityType: string;
  /** When this request was created */
  createdAt?: Timestamp;
  updatedAt?: Timestamp;
  /** Callback url for this request (if any)
 All records indexed in this operation will be notified via this url */
  callbackUrl: string;
  /** Records to enqueue - note that indexing by its nature is asynchronous
 so we wont wait for all the indexing to finish */
  records?: IndexRecord[];
}



export interface CreateIndexRecordsLRORequest {
  lro?: IndexRecordsLRO;
}



export interface CreateIndexRecordsLROResponse {
  lro?: IndexRecordsLRO;
}



export interface UpdateIndexRecordsLRORequest {
  lro?: IndexRecordsLRO;
  /** *
 Mask of fields being updated in this Game to make partial changes. */
  updateMask?: FieldMask;
}



export interface UpdateIndexRecordsLROResponse {
  lro?: IndexRecordsLRO;
}



export interface GetIndexRecordsLRORequest {
  lroId: string;
}



export interface GetIndexRecordsLROResponse {
  lro?: IndexRecordsLRO;
}


/**
 * Job describes the work that needs to be done.
 */
export interface Job {
  entityType: string;
  entityId: string;
  jobType: string;
  /** When the last indexing was queued */
  createdAt?: Timestamp;
  /** when the last indexing was completed */
  updatedAt?: Timestamp;
  /** Job specific data */
  jobData?: Any;
  /** Debounce so we dont run it too many time within this many seconds */
  debounceWindowSeconds: number;
  /** Whether the job is a oneoff or can repeat */
  repeatInfo?: RepeatInfo;
}



export interface RepeatInfo {
}



export interface Run {
  jobId: string;
  /** A unique run_id */
  runId: string;
  createdAt?: Timestamp;
  startedAt?: Timestamp;
  updatedAt?: Timestamp;
  state: RunState;
  /** Run specific data */
  runData?: Any;
  /** If there was an error in the last indexing */
  lastError: string;
  /** Keep a hash of the contents for quick check to check updated
 (not sure if needed) - This should be provided by the source */
  lastContentHash: string;
  /** If there were retries */
  retryCount: number;
}


/**
 * Called when the singleton presenter is initialized
 */
export interface InitializeSingletonRequest {
  gameId: string;
  gameData: string;
  gameState: string;
  moveHistory: string;
  /** User ID of the viewer (logged-in user viewing the game, for Join button visibility) */
  viewerUserId: string;
}


/**
 * Response of a turn option click
 */
export interface InitializeSingletonResponse {
  response?: InitializeGameResponse;
}


/**
 * Called when a turn option is clicked in TurnOptionsPanel
 */
export interface TurnOptionClickedRequest {
  gameId: string;
  optionIndex: number;
  optionType: string;
  pos?: Position;
}


/**
 * Response of a turn option click
 */
export interface TurnOptionClickedResponse {
  gameId: string;
}


/**
 * Called when the scene was clicked
 */
export interface SceneClickedRequest {
  gameId: string;
  pos?: Position;
  layer: string;
}


/**
 * Response of a turn option click
 */
export interface SceneClickedResponse {
  gameId: string;
}


/**
 * Called when the end turn button was clicked
 */
export interface EndTurnButtonClickedRequest {
  gameId: string;
}


/**
 * Response of a turn option click
 */
export interface EndTurnButtonClickedResponse {
  gameId: string;
}


/**
 * Called when a build option is clicked in BuildOptionsModal
 */
export interface BuildOptionClickedRequest {
  gameId: string;
  pos?: Position;
  unitType: number;
}


/**
 * Response of a build option click
 */
export interface BuildOptionClickedResponse {
}


/**
 * Called when the end turn button was clicked
 */
export interface InitializeGameRequest {
  gameId: string;
}


/**
 * Response of a turn option click
 */
export interface InitializeGameResponse {
  success: boolean;
  error: string;
  /** Initial UI state information */
  currentPlayer: number;
  turnCounter: number;
  gameName: string;
}


/**
 * Called by browser after UI/scene is fully initialized and ready for visual updates
 */
export interface ClientReadyRequest {
  gameId: string;
}


/**
 * Response for ClientReady
 */
export interface ClientReadyResponse {
  success: boolean;
}


/**
 * Request to apply changes from remote players
 Called when SyncService streams moves made by other players
 */
export interface ApplyRemoteChangesRequest {
  gameId: string;
  /** The moves containing WorldChanges to apply */
  moves?: GameMove[];
}


/**
 * Response after applying remote changes
 */
export interface ApplyRemoteChangesResponse {
  /** Whether the changes were successfully applied */
  success: boolean;
  /** Error message if application failed (e.g., state desync) */
  error: string;
  /** If state desync detected, client should reload game */
  requiresReload: boolean;
}


/**
 * SubscribeRequest to start receiving game updates
 */
export interface SubscribeRequest {
  /** Game ID to subscribe to */
  gameId: string;
  /** Player ID this client represents (for filtering own moves) */
  playerId: string;
  /** Resume from this sequence number (for reconnection).
 Server will send any missed updates since this sequence.
 Use 0 to start from current state. */
  fromSequence: number;
}


/**
 * SubscribeResponse sent once at the start of the subscription
 */
export interface SubscribeResponse {
  /** Current sequence number at time of subscription */
  currentSequence: number;
  /** Current game state (for initial load or catchup) */
  gameState?: GameState;
  /** Game metadata */
  game?: Game;
}


/**
 * GameUpdate is streamed to subscribers when game state changes
 */
export interface GameUpdate {
  /** Monotonically increasing sequence number for ordering and reconnection.
 Clients should track this to resume from correct position. */
  sequence: number;
  /** Moves were published by a player (contains WorldChanges) */
  movesPublished?: MovesPublished;
  /** A player connected to the game */
  playerJoined?: PlayerJoined;
  /** A player disconnected from the game */
  playerLeft?: PlayerLeft;
  /** Game has ended (victory, surrender, etc.) */
  gameEnded?: GameEnded;
  /** Initial state sent at subscription start */
  initialState?: SubscribeResponse;
}


/**
 * MovesPublished indicates a player made moves
 */
export interface MovesPublished {
  /** Which player made these moves */
  player: number;
  /** The moves with their WorldChanges populated */
  moves?: GameMove[];
  /** Group number for this batch of moves */
  groupNumber: number;
}


/**
 * PlayerJoined indicates a player connected
 */
export interface PlayerJoined {
  playerId: string;
  playerNumber: number;
}


/**
 * PlayerLeft indicates a player disconnected
 */
export interface PlayerLeft {
  playerId: string;
  playerNumber: number;
}


/**
 * GameEnded indicates the game has concluded
 */
export interface GameEnded {
  /** Winning player (0 if draw or N/A) */
  winner: number;
  /** Reason for game ending */
  reason: string;
}


/**
 * BroadcastRequest to send a GameUpdate to all subscribers
 Called internally by GamesService after ProcessMoves succeeds
 */
export interface BroadcastRequest {
  /** Game ID to broadcast to */
  gameId: string;
  /** The update to broadcast */
  update?: GameUpdate;
}


/**
 * BroadcastResponse after broadcasting
 */
export interface BroadcastResponse {
  /** Number of subscribers that received the update */
  subscriberCount: number;
  /** The sequence number assigned to this update */
  sequence: number;
}


/**
 * ThemeInfo contains metadata about a theme
 */
export interface ThemeInfo {
  name: string;
  version: string;
  basePath: string;
  assetType: string;
  needsPostProcessing: boolean;
}


/**
 * UnitMapping maps a unit ID to its theme-specific representation
 */
export interface UnitMapping {
  old: string;
  name: string;
  image: string;
  description: string;
}


/**
 * TerrainMapping maps a terrain ID to its theme-specific representation
 */
export interface TerrainMapping {
  old: string;
  name: string;
  image: string;
  description: string;
}


/**
 * ThemeManifest represents the full theme configuration
 This matches the structure of mapping.json files
 */
export interface ThemeManifest {
  themeInfo?: ThemeInfo;
  units: Record<number, UnitMapping>;
  terrains: Record<number, TerrainMapping>;
  playerColors: Record<number, PlayerColor>;
}


/**
 * PlayerColor defines the color scheme for a player
 */
export interface PlayerColor {
  primary: string;
  secondary: string;
  name: string;
}


/**
 * AssetResult represents a rendered asset
 */
export interface AssetResult {
  type: Type;
  data: string;
}


/**
 * WorldInfo represents a world in the catalog
 */
export interface WorldInfo {
  id: string;
  name: string;
  description: string;
  category: string;
  difficulty: string;
  tags: string[];
  icon: string;
  lastUpdated: string;
}


/**
 * Request messages
 */
export interface ListWorldsRequest {
  /** Pagination info */
  pagination?: Pagination;
  /** May be filter by owner id */
  ownerId: string;
}



export interface ListWorldsResponse {
  items?: World[];
  pagination?: PaginationResponse;
}



export interface GetWorldRequest {
  id: string;
  version: string;
}



export interface GetWorldResponse {
  world?: World;
  worldData?: WorldData;
}



export interface UpdateWorldRequest {
  /** *
 World being updated */
  world?: World;
  worldData?: WorldData;
  clearWorld: boolean;
  /** *
 Mask of fields being updated in this World to make partial changes. */
  updateMask?: FieldMask;
}


/**
 * *
 The request for (partially) updating an World.
 */
export interface UpdateWorldResponse {
  /** *
 World being updated */
  world?: World;
  worldData?: WorldData;
}


/**
 * *
 Request to delete an world.
 */
export interface DeleteWorldRequest {
  /** *
 ID of the world to be deleted. */
  id: string;
}


/**
 * *
 World deletion response
 */
export interface DeleteWorldResponse {
}


/**
 * *
 Request to batch get worlds
 */
export interface GetWorldsRequest {
  /** *
 IDs of the world to be fetched */
  ids: string[];
}


/**
 * *
 World batch-get response
 */
export interface GetWorldsResponse {
  worlds: Record<string, World>;
}


/**
 * *
 World creation request object
 */
export interface CreateWorldRequest {
  /** *
 World being updated */
  world?: World;
  worldData?: WorldData;
}


/**
 * *
 Response of an world creation.
 */
export interface CreateWorldResponse {
  /** *
 World being created */
  world?: World;
  worldData?: WorldData;
  /** *
 Error specific to a field if there are any errors. */
  fieldErrors: Record<string, string>;
}

