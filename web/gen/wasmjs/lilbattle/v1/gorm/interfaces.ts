// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Any, Timestamp } from "@bufbuild/protobuf/wkt";



/**
 * IndexStateGORM is the GORM representation for IndexState
 */
export interface IndexStateGORM {
  /** Composite index for entity lookups: WHERE entity_type = ? AND entity_id = ? */
  entityType: string;
  entityId: string;
  /** Index for range queries: WHERE indexed_at <= ?
 Used by indexer to find stale records */
  indexedAt?: Timestamp;
  /** Partial index for finding records that need indexing
 Only indexes rows where needs_indexing = true (much smaller index) */
  needsIndexing: boolean;
}


/**
 * IndexRecordsLROGORM is the GORM representation for IndexRecordsLRO
 */
export interface IndexRecordsLROGORM {
}



export interface IndexInfoGORM {
}



export interface TileGORM {
}



export interface CrossingGORM {
}



export interface UnitGORM {
}



export interface AttackRecordGORM {
}



export interface WorldGORM {
  id: string;
  /** Tags as JSON for cross-DB compatibility */
  tags: string[];
  /** PreviewUrls as JSON for cross-DB compatibility */
  previewUrls: string[];
  /** DefaultGameConfig as JSON for cross-DB compatibility
 SearchIndexInfo embedded */
  searchIndexInfo?: IndexInfoGORM;
}



export interface WorldDataGORM {
  worldId: string;
  /** Units as JSON for cross-DB compatibility */
  crossings: Record<string, CrossingGORM>;
  /** ScreenshotIndexInfo embedded */
  screenshotIndexInfo?: IndexInfoGORM;
  /** Tiles as JSON for cross-DB compatibility */
  tilesMap: Record<string, TileGORM>;
  /** Units as JSON for cross-DB compatibility */
  unitsMap: Record<string, UnitGORM>;
}


/**
 * Describes a game and its metadata
 */
export interface GameGORM {
  id: string;
  /** world_id indexed for queries filtering games by world */
  worldId: string;
  /** Tags as JSON for cross-DB compatibility */
  tags: string[];
  /** PreviewUrls as JSON for cross-DB compatibility */
  previewUrls: string[];
  /** SearchIndexInfo embedded */
  searchIndexInfo?: IndexInfoGORM;
}


/**
 * Holds the game's Active/Current state (eg world state)
 */
export interface GameStateGORM {
  gameId: string;
  /** ScreenshotIndexInfo embedded */
  worldData?: GameWorldDataGORM;
  /** Per-player runtime state as JSON for cross-DB compatibility */
  playerStates: Record<number, PlayerStateGORM>;
}



export interface GameConfigurationGORM {
  /** IncomeConfigs embedded */
  incomeConfigs?: IncomeConfigGORM;
  /** Settings as foreign key relationship */
  settings?: GameSettingsGORM;
}



export interface IncomeConfigGORM {
}



export interface GamePlayerGORM {
}



export interface GameTeamGORM {
}



export interface GameSettingsGORM {
  /** AllowedUnits as JSON for cross-DB compatibility */
  allowedUnits: number[];
}



export interface PlayerStateGORM {
}


/**
 * GameWorldDataGORM is same as WorldDataGORM but without the
 primary key so it can be embedded
 */
export interface GameWorldDataGORM {
  /** ScreenshotIndexInfo embedded */
  screenshotIndexInfo?: IndexInfoGORM;
  /** Units as JSON for cross-DB compatibility */
  crossings: Record<string, CrossingGORM>;
  /** Tiles as JSON for cross-DB compatibility */
  tilesMap: Record<string, TileGORM>;
  /** Units as JSON for cross-DB compatibility */
  unitsMap: Record<string, UnitGORM>;
}


/**
 * Holds the game's move history (can be used as a replay log)
 */
export interface GameMoveHistoryGORM {
}


/**
 * A move group - we can allow X moves in one "tick"
 */
export interface GameMoveGroupGORM {
}


/**
 * *
 Represents a single move which can be one of many actions in the game
 */
export interface GameMoveGORM {
  /** game_id is indexed for queries like WHERE game_id = ? ORDER BY group_number
 Also part of composite index idx_game_moves_lookup for range queries */
  gameId: string;
  /** group_number is part of composite index for queries like WHERE game_id = ? AND group_number >= ? */
  groupNumber: number;
  moveNumber: number;
  /** Version number for optimistic locking */
  version: number;
  /** Field named "move_type" matches the oneof name in source
 This automatically skips all oneof members (move_unit, attack_unit, end_turn, build_unit) */
  moveType?: Any;
  changes?: Any[];
}

