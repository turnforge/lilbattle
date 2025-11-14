import { EntityIndexStateDatastore as EntityIndexStateDatastoreInterface, IndexRecordDatastore as IndexRecordDatastoreInterface, IndexRecordsLRODatastore as IndexRecordsLRODatastoreInterface } from "./interfaces";




/**
 * EntityIndexStateDatastore is the Datastore representation for EntityIndexState
 */
export class EntityIndexStateDatastore implements EntityIndexStateDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.EntityIndexStateDatastore";
  readonly __MESSAGE_TYPE = EntityIndexStateDatastore.MESSAGE_TYPE;

  /** Composite key in Datastore: entity_type + entity_id + index_type
 Datastore uses string IDs, so we'll concatenate them */
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
 * IndexRecordDatastore is the Datastore representation for IndexRecord
 */
export class IndexRecordDatastore implements IndexRecordDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.IndexRecordDatastore";
  readonly __MESSAGE_TYPE = IndexRecordDatastore.MESSAGE_TYPE;

  entityId: string = "";
  updatedAt: number = 0;
  /** Note: google.protobuf.Any requires special handling - typically stored as bytes
 We'll skip entity_data for now as it needs custom serialization
 bytes entity_data = 4; */
  indexerTypes: string[] = [];

  
}


/**
 * IndexRecordsLRODatastore is the Datastore representation for IndexRecordsLRO
 */
export class IndexRecordsLRODatastore implements IndexRecordsLRODatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.IndexRecordsLRODatastore";
  readonly __MESSAGE_TYPE = IndexRecordsLRODatastore.MESSAGE_TYPE;

  /** Primary key */
  lroId: string = "";
  /** Entity type for all records in this LRO */
  entityType: string = "";
  /** Timestamps */
  createdAt: number = 0;
  updatedAt: number = 0;
  /** Callback URL */
  callbackUrl: string = "";
  /** Records - Datastore supports nested repeated message types natively */
  records: IndexRecordDatastore[] = [];

  
}


