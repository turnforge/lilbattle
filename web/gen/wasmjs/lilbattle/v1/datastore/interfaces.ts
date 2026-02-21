// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Any } from "@bufbuild/protobuf/wkt";



/**
 * IndexStateDatastore is the Datastore representation for IndexState
 */
export interface IndexStateDatastore {
}


/**
 * IndexRecordDatastore is the Datastore representation for IndexRecord
 */
export interface IndexRecordDatastore {
}


/**
 * IndexRecordsLRODatastore is the Datastore representation for IndexRecordsLRO
 */
export interface IndexRecordsLRODatastore {
}



export interface IndexInfoDatastore {
}



export interface TileDatastore {
}



export interface CrossingDatastore {
}



export interface UnitDatastore {
  /** Attack history as nested entities (noindex for large arrays) */
  attackHistory?: AttackRecordDatastore[];
}



export interface AttackRecordDatastore {
}


/**
 * WorldDatastore is the Datastore representation for World
 */
export interface WorldDatastore {
  id: string;
  /** Tags as noindex (not queryable) */
  tags: string[];
  /** PreviewUrls as noindex */
  previewUrls: string[];
  /** DefaultGameConfig as flattened embedded struct */
  defaultGameConfig?: GameConfigurationDatastore;
  /** SearchIndexInfo - needs_indexing should be indexed for worker queries */
  searchIndexInfo?: IndexInfoDatastore;
}


/**
 * WorldDataDatastore stores the actual world map data
 */
export interface WorldDataDatastore {
  /** Primary key - matches World.id (excluded from properties, used for key) */
  worldId: string;
  /** Map of tiles - large, no index needed */
  tilesMap: Record<string, TileDatastore>;
  /** Map of units - large, no index needed */
  unitsMap: Record<string, UnitDatastore>;
  /** Map of crossings - large, no index needed */
  crossings: Record<string, CrossingDatastore>;
  /** ScreenshotIndexInfo - flatten so needs_indexing is queryable */
  screenshotIndexInfo?: IndexInfoDatastore;
}


/**
 * GameDatastore is the Datastore representation for Game
 */
export interface GameDatastore {
  /** ID - used for key, not stored as property */
  id: string;
  /** World ID - indexed for filtering games by world */
  worldId: string;
  /** Tags as noindex (not queryable) */
  tags: string[];
  /** PreviewUrls as noindex */
  previewUrls: string[];
  /** Config as noindex (large nested struct) */
  config?: GameConfigurationDatastore;
  /** SearchIndexInfo - flatten so needs_indexing is queryable */
  searchIndexInfo?: IndexInfoDatastore;
}


/**
 * GameStateDatastore stores the active game state
 */
export interface GameStateDatastore {
  /** GameId - used for key, not stored as property */
  gameId: string;
  /** WorldData - large embedded struct, noindex */
  worldData?: WorldDataDatastore;
  /** Per-player runtime state - noindex */
  playerStates: Record<number, PlayerStateDatastore>;
}



export interface GameConfigurationDatastore {
  /** Players as noindex nested array */
  players?: GamePlayerDatastore[];
  /** Teams as noindex nested array */
  teams?: GameTeamDatastore[];
  /** IncomeConfigs embedded */
  incomeConfigs?: IncomeConfigDatastore;
  /** Settings embedded */
  settings?: GameSettingsDatastore;
}



export interface IncomeConfigDatastore {
}



export interface GamePlayerDatastore {
}



export interface GameTeamDatastore {
}



export interface GameSettingsDatastore {
  /** AllowedUnits as noindex (array of ints) */
  allowedUnits: number[];
}



export interface PlayerStateDatastore {
}


/**
 * GameMoveDatastore stores individual moves
 */
export interface GameMoveDatastore {
  /** Fields needed for composite key and querying */
  gameId: string;
  /** Group number - indexed for range queries */
  groupNumber: number;
  /** Move number within group - indexed for ordering */
  moveNumber: number;
  /** Field named "move_type" matches the oneof name in source
 This automatically skips all oneof members (move_unit, attack_unit, etc.)
 Stored as bytes, noindex (too large) */
  moveType?: Any;
  /** Changes - stored as bytes array, noindex (too large) */
  changes?: Any[];
}

