// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Any } from "@bufbuild/protobuf/wkt";



/**
 * IndexStateGORM is the GORM representation for IndexState
 */
export interface IndexStateGORM {
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



export interface UnitGORM {
}



export interface AttackRecordGORM {
}



export interface WorldGORM {
}



export interface WorldDataGORM {
}


/**
 * Describes a game and its metadata
 */
export interface GameGORM {
}



export interface GameConfigurationGORM {
}



export interface IncomeConfigGORM {
}



export interface GamePlayerGORM {
}



export interface GameTeamGORM {
}



export interface GameSettingsGORM {
}


/**
 * Holds the game's Active/Current state (eg world state)
 */
export interface GameStateGORM {
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
  /** Store the oneof move_type as serialized Any */
  moveType?: Any;
  /** Skip the individual oneof fields from the source (field names, not oneof name!) */
  moveUnit: boolean;
  attackUnit: boolean;
  endTurn: boolean;
  buildUnit: boolean;
  changes?: Any;
}

