import { IndexStateDatastore as IndexStateDatastoreInterface, IndexRecordDatastore as IndexRecordDatastoreInterface, IndexRecordsLRODatastore as IndexRecordsLRODatastoreInterface } from "./interfaces";




/**
 * IndexStateDatastore is the Datastore representation for IndexState
 */
export class IndexStateDatastore implements IndexStateDatastoreInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "weewar.v1.IndexStateDatastore";
  readonly __MESSAGE_TYPE = IndexStateDatastore.MESSAGE_TYPE;


  
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


  
}


