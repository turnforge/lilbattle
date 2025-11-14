import { EntityIndexStateGORM as EntityIndexStateGORMInterface, IndexRecordsLROGORM as IndexRecordsLROGORMInterface } from "./interfaces";




/**
 * EntityIndexStateGORM is the GORM representation for EntityIndexState
 */
export class EntityIndexStateGORM implements EntityIndexStateGORMInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.EntityIndexStateGORM";
  readonly __MESSAGE_TYPE = EntityIndexStateGORM.MESSAGE_TYPE;

  /** Composite primary key: entity_type + entity_id + index_type */
  entityType: string = "";
  entityId: string = "";
  indexType: string = "";
  /** Timestamps stored as int64 (Unix time) */
  lastQueuedAt: number = 0;
  lastIndexedAt: number = 0;
  /** Status field */
  status: string = "";
  /** Error tracking */
  lastError: string = "";
  /** Content hash for change detection */
  lastContentHash: string = "";
  /** Retry tracking */
  retryCount: number = 0;
  /** LRO tracking */
  currentLroId: string = "";

  
}


/**
 * IndexRecordsLROGORM is the GORM representation for IndexRecordsLRO
 */
export class IndexRecordsLROGORM implements IndexRecordsLROGORMInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.IndexRecordsLROGORM";
  readonly __MESSAGE_TYPE = IndexRecordsLROGORM.MESSAGE_TYPE;

  /** Primary key */
  lroId: string = "";
  /** Entity type for all records in this LRO */
  entityType: string = "";
  /** Timestamps */
  createdAt: number = 0;
  updatedAt: number = 0;
  /** Callback URL */
  callbackUrl: string = "";
  /** Records stored as JSONB (PostgreSQL-specific)
 The repeated IndexRecord will be serialized to JSON
 Note: In the proto, this is `repeated IndexRecord` but GORM can't handle
 repeated message types directly, so we'll store as JSONB */
  records: Uint8Array = new Uint8Array();

  
}


