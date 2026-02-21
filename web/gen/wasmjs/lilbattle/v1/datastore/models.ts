import { Any } from "@bufbuild/protobuf/wkt";


import { IndexStateDatastore as IndexStateDatastoreInterface, IndexRecordDatastore as IndexRecordDatastoreInterface, IndexRecordsLRODatastore as IndexRecordsLRODatastoreInterface, IndexInfoDatastore as IndexInfoDatastoreInterface, TileDatastore as TileDatastoreInterface, CrossingDatastore as CrossingDatastoreInterface, UnitDatastore as UnitDatastoreInterface, AttackRecordDatastore as AttackRecordDatastoreInterface, WorldDatastore as WorldDatastoreInterface, WorldDataDatastore as WorldDataDatastoreInterface, GameDatastore as GameDatastoreInterface, GameStateDatastore as GameStateDatastoreInterface, GameConfigurationDatastore as GameConfigurationDatastoreInterface, IncomeConfigDatastore as IncomeConfigDatastoreInterface, GamePlayerDatastore as GamePlayerDatastoreInterface, GameTeamDatastore as GameTeamDatastoreInterface, GameSettingsDatastore as GameSettingsDatastoreInterface, PlayerStateDatastore as PlayerStateDatastoreInterface, GameMoveDatastore as GameMoveDatastoreInterface } from "./interfaces";




/**
 * IndexStateDatastore is the Datastore representation for IndexState
 */
export class IndexStateDatastore implements IndexStateDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.IndexStateDatastore";
  readonly __MESSAGE_TYPE = IndexStateDatastore.MESSAGE_TYPE;


  
}


/**
 * IndexRecordDatastore is the Datastore representation for IndexRecord
 */
export class IndexRecordDatastore implements IndexRecordDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.IndexRecordDatastore";
  readonly __MESSAGE_TYPE = IndexRecordDatastore.MESSAGE_TYPE;


  
}


/**
 * IndexRecordsLRODatastore is the Datastore representation for IndexRecordsLRO
 */
export class IndexRecordsLRODatastore implements IndexRecordsLRODatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.IndexRecordsLRODatastore";
  readonly __MESSAGE_TYPE = IndexRecordsLRODatastore.MESSAGE_TYPE;


  
}



export class IndexInfoDatastore implements IndexInfoDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.IndexInfoDatastore";
  readonly __MESSAGE_TYPE = IndexInfoDatastore.MESSAGE_TYPE;


  
}



export class TileDatastore implements TileDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.TileDatastore";
  readonly __MESSAGE_TYPE = TileDatastore.MESSAGE_TYPE;


  
}



export class CrossingDatastore implements CrossingDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.CrossingDatastore";
  readonly __MESSAGE_TYPE = CrossingDatastore.MESSAGE_TYPE;


  
}



export class UnitDatastore implements UnitDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.UnitDatastore";
  readonly __MESSAGE_TYPE = UnitDatastore.MESSAGE_TYPE;

  /** Attack history as nested entities (noindex for large arrays) */
  attackHistory: AttackRecordDatastore[] = [];

  
}



export class AttackRecordDatastore implements AttackRecordDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.AttackRecordDatastore";
  readonly __MESSAGE_TYPE = AttackRecordDatastore.MESSAGE_TYPE;


  
}


/**
 * WorldDatastore is the Datastore representation for World
 */
export class WorldDatastore implements WorldDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.WorldDatastore";
  readonly __MESSAGE_TYPE = WorldDatastore.MESSAGE_TYPE;

  id: string = "";
  /** Tags as noindex (not queryable) */
  tags: string[] = [];
  /** PreviewUrls as noindex */
  previewUrls: string[] = [];
  /** DefaultGameConfig as flattened embedded struct */
  defaultGameConfig?: GameConfigurationDatastore;
  /** SearchIndexInfo - needs_indexing should be indexed for worker queries */
  searchIndexInfo?: IndexInfoDatastore;

  
}


/**
 * WorldDataDatastore stores the actual world map data
 */
export class WorldDataDatastore implements WorldDataDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.WorldDataDatastore";
  readonly __MESSAGE_TYPE = WorldDataDatastore.MESSAGE_TYPE;

  /** Primary key - matches World.id (excluded from properties, used for key) */
  worldId: string = "";
  /** Map of tiles - large, no index needed */
  tilesMap: Record<string, TileDatastore> = {};
  /** Map of units - large, no index needed */
  unitsMap: Record<string, UnitDatastore> = {};
  /** Map of crossings - large, no index needed */
  crossings: Record<string, CrossingDatastore> = {};
  /** ScreenshotIndexInfo - flatten so needs_indexing is queryable */
  screenshotIndexInfo?: IndexInfoDatastore;

  
}


/**
 * GameDatastore is the Datastore representation for Game
 */
export class GameDatastore implements GameDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameDatastore";
  readonly __MESSAGE_TYPE = GameDatastore.MESSAGE_TYPE;

  /** ID - used for key, not stored as property */
  id: string = "";
  /** World ID - indexed for filtering games by world */
  worldId: string = "";
  /** Tags as noindex (not queryable) */
  tags: string[] = [];
  /** PreviewUrls as noindex */
  previewUrls: string[] = [];
  /** Config as noindex (large nested struct) */
  config?: GameConfigurationDatastore;
  /** SearchIndexInfo - flatten so needs_indexing is queryable */
  searchIndexInfo?: IndexInfoDatastore;

  
}


/**
 * GameStateDatastore stores the active game state
 */
export class GameStateDatastore implements GameStateDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameStateDatastore";
  readonly __MESSAGE_TYPE = GameStateDatastore.MESSAGE_TYPE;

  /** GameId - used for key, not stored as property */
  gameId: string = "";
  /** WorldData - large embedded struct, noindex */
  worldData?: WorldDataDatastore;
  /** Per-player runtime state - noindex */
  playerStates: Record<number, PlayerStateDatastore> = {};

  
}



export class GameConfigurationDatastore implements GameConfigurationDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameConfigurationDatastore";
  readonly __MESSAGE_TYPE = GameConfigurationDatastore.MESSAGE_TYPE;

  /** Players as noindex nested array */
  players: GamePlayerDatastore[] = [];
  /** Teams as noindex nested array */
  teams: GameTeamDatastore[] = [];
  /** IncomeConfigs embedded */
  incomeConfigs?: IncomeConfigDatastore;
  /** Settings embedded */
  settings?: GameSettingsDatastore;

  
}



export class IncomeConfigDatastore implements IncomeConfigDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.IncomeConfigDatastore";
  readonly __MESSAGE_TYPE = IncomeConfigDatastore.MESSAGE_TYPE;


  
}



export class GamePlayerDatastore implements GamePlayerDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GamePlayerDatastore";
  readonly __MESSAGE_TYPE = GamePlayerDatastore.MESSAGE_TYPE;


  
}



export class GameTeamDatastore implements GameTeamDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameTeamDatastore";
  readonly __MESSAGE_TYPE = GameTeamDatastore.MESSAGE_TYPE;


  
}



export class GameSettingsDatastore implements GameSettingsDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameSettingsDatastore";
  readonly __MESSAGE_TYPE = GameSettingsDatastore.MESSAGE_TYPE;

  /** AllowedUnits as noindex (array of ints) */
  allowedUnits: number[] = [];

  
}



export class PlayerStateDatastore implements PlayerStateDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.PlayerStateDatastore";
  readonly __MESSAGE_TYPE = PlayerStateDatastore.MESSAGE_TYPE;


  
}


/**
 * GameMoveDatastore stores individual moves
 */
export class GameMoveDatastore implements GameMoveDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "lilbattle.v1.GameMoveDatastore";
  readonly __MESSAGE_TYPE = GameMoveDatastore.MESSAGE_TYPE;

  /** Fields needed for composite key and querying */
  gameId: string = "";
  /** Group number - indexed for range queries */
  groupNumber: number = 0;
  /** Move number within group - indexed for ordering */
  moveNumber: number = 0;
  /** Field named "move_type" matches the oneof name in source
 This automatically skips all oneof members (move_unit, attack_unit, etc.)
 Stored as bytes, noindex (too large) */
  moveType?: Any;
  /** Changes - stored as bytes array, noindex (too large) */
  changes: Any[] = [];

  
}


