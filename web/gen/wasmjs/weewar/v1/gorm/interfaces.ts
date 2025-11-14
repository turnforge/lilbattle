// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * EntityIndexStateGORM is the GORM representation for EntityIndexState
 */
export interface EntityIndexStateGORM {
  /** Composite primary key: entity_type + entity_id + index_type */
  entityType: string;
  entityId: string;
  indexType: string;
  /** Timestamps stored as int64 (Unix time) */
  lastQueuedAt: number;
  lastIndexedAt: number;
  /** Status field */
  status: string;
  /** Error tracking */
  lastError: string;
  /** Content hash for change detection */
  lastContentHash: string;
  /** Retry tracking */
  retryCount: number;
  /** LRO tracking */
  currentLroId: string;
}


/**
 * IndexRecordsLROGORM is the GORM representation for IndexRecordsLRO
 */
export interface IndexRecordsLROGORM {
  /** Primary key */
  lroId: string;
  /** Entity type for all records in this LRO */
  entityType: string;
  /** Timestamps */
  createdAt: number;
  updatedAt: number;
  /** Callback URL */
  callbackUrl: string;
  /** Records stored as JSONB (PostgreSQL-specific)
 The repeated IndexRecord will be serialized to JSON
 Note: In the proto, this is `repeated IndexRecord` but GORM can't handle
 repeated message types directly, so we'll store as JSONB */
  records: Uint8Array;
}

