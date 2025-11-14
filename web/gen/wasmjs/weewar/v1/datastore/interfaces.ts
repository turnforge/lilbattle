// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * EntityIndexStateDatastore is the Datastore representation for EntityIndexState
 */
export interface EntityIndexStateDatastore {
  /** Composite key in Datastore: entity_type + entity_id + index_type
 Datastore uses string IDs, so we'll concatenate them */
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
 * IndexRecordDatastore is the Datastore representation for IndexRecord
 */
export interface IndexRecordDatastore {
  entityId: string;
  updatedAt: number;
  /** Note: google.protobuf.Any requires special handling - typically stored as bytes
 We'll skip entity_data for now as it needs custom serialization
 bytes entity_data = 4; */
  indexerTypes: string[];
}


/**
 * IndexRecordsLRODatastore is the Datastore representation for IndexRecordsLRO
 */
export interface IndexRecordsLRODatastore {
  /** Primary key */
  lroId: string;
  /** Entity type for all records in this LRO */
  entityType: string;
  /** Timestamps */
  createdAt: number;
  updatedAt: number;
  /** Callback URL */
  callbackUrl: string;
  /** Records - Datastore supports nested repeated message types natively */
  records?: IndexRecordDatastore[];
}

