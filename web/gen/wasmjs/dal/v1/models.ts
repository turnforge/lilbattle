import { TableOptions as TableOptionsInterface, ColumnOptions as ColumnOptionsInterface, ConverterFunc as ConverterFuncInterface, IndexOptions as IndexOptionsInterface, ForeignKeyOptions as ForeignKeyOptionsInterface, GormOptions as GormOptionsInterface, PostgresOptions as PostgresOptionsInterface, DatastoreOptions as DatastoreOptionsInterface, FirestoreOptions as FirestoreOptionsInterface, MongoDBOptions as MongoDBOptionsInterface, ReferentialAction } from "./interfaces";




/**
 * Configuration for table mapping
 */
export class TableOptions implements TableOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.TableOptions";
  readonly __MESSAGE_TYPE = TableOptions.MESSAGE_TYPE;

  /** Table name in the database */
  name: string = "";
  /** Database schema/namespace */
  schema: string = "";
  /** Custom table comment/description */
  comment: string = "";
  /** Source message fully qualified name (e.g., "library.v1.Book")
 This links the DAL schema message to the API proto message */
  source: string = "";

  
}


/**
 * Configuration for column mapping
 */
export class ColumnOptions implements ColumnOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.ColumnOptions";
  readonly __MESSAGE_TYPE = ColumnOptions.MESSAGE_TYPE;

  /** Column name override (optional - defaults to field name) */
  name: string = "";
  /** Custom conversion function for API -> Target conversion
 Overrides built-in converters (e.g., Timestamp -> int64)
 Example:
   to_func: {
     package: "github.com/myapp/converters"
     alias: "myconv"
     function: "TimestampToMillis"
   }
 Generates: import myconv "github.com/myapp/converters"
            gorm.Field = myconv.TimestampToMillis(api.Field) */
  toFunc?: ConverterFunc;
  /** Custom conversion function for Target -> API conversion
 Example:
   from_func: {
     package: "github.com/myapp/converters"
     function: "MillisToTimestamp"
   }
 Generates: api.Field = converters.MillisToTimestamp(gorm.Field) */
  fromFunc?: ConverterFunc;
  /** GORM-specific tags (for GORM target)
 Example: ["primaryKey", "type:uuid", "default:gen_random_uuid()"]
 Generates: `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` */
  gormTags: string[] = [];
  /** Raw SQL column definition tags (for raw SQL targets)
 Target-specific usage */
  sqlTags: string[] = [];
  /** Firestore tags (for Firestore target) */
  firestoreTags: string[] = [];
  /** MongoDB tags (for MongoDB target) */
  mongodbTags: string[] = [];

  
}


/**
 * Specification for a custom converter function
 */
export class ConverterFunc implements ConverterFuncInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.ConverterFunc";
  readonly __MESSAGE_TYPE = ConverterFunc.MESSAGE_TYPE;

  /** Go package import path
 Example: "github.com/myapp/converters" */
  package: string = "";
  /** Optional import alias
 If not specified, uses the last segment of package path
 Example: "myconv" -> import myconv "github.com/myapp/converters" */
  alias: string = "";
  /** Function name to call
 Example: "TimestampToMillis"
 Generates call: alias.TimestampToMillis(value) */
  function: string = "";

  
}


/**
 * Configuration for indexes
 */
export class IndexOptions implements IndexOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.IndexOptions";
  readonly __MESSAGE_TYPE = IndexOptions.MESSAGE_TYPE;

  /** Index name */
  name: string = "";
  /** Comma-separated list of field names (for composite indexes) */
  fields: string = "";
  /** Whether this is a unique index */
  unique: boolean = false;
  /** Index type (e.g., "BTREE", "HASH", "GIN", "GIST") */
  type: string = "";
  /** Partial index condition (WHERE clause) */
  where: string = "";

  
}


/**
 * Configuration for foreign keys
 */
export class ForeignKeyOptions implements ForeignKeyOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.ForeignKeyOptions";
  readonly __MESSAGE_TYPE = ForeignKeyOptions.MESSAGE_TYPE;

  /** Referenced table and column (e.g., "authors.id") */
  references: string = "";
  /** Action on delete */
  onDelete: ReferentialAction = ReferentialAction.NO_ACTION;
  /** Action on update */
  onUpdate: ReferentialAction = ReferentialAction.NO_ACTION;
  /** Foreign key constraint name */
  constraintName: string = "";

  
}


/**
 * GORM target options (database-agnostic ORM)
 */
export class GormOptions implements GormOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.GormOptions";
  readonly __MESSAGE_TYPE = GormOptions.MESSAGE_TYPE;

  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string = "";
  /** Table name */
  table: string = "";
  /** Embedded field names (for GORM embedded structs) */
  embedded: string[] = [];

  
}


/**
 * PostgreSQL target options (raw SQL)
 */
export class PostgresOptions implements PostgresOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.PostgresOptions";
  readonly __MESSAGE_TYPE = PostgresOptions.MESSAGE_TYPE;

  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string = "";
  /** Table name in PostgreSQL */
  table: string = "";
  /** Schema name (default: "public") */
  schema: string = "";

  
}


/**
 * Google Cloud Datastore options
 */
export class DatastoreOptions implements DatastoreOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.DatastoreOptions";
  readonly __MESSAGE_TYPE = DatastoreOptions.MESSAGE_TYPE;

  /** Datastore kind name */
  kind: string = "";
  /** Datastore namespace */
  namespace: string = "";
  /** Whether to use incomplete keys (auto-generated) */
  incompleteKey: boolean = false;
  /** Ancestor path */
  ancestor: string = "";
  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string = "";

  
}


/**
 * Firestore target options
 */
export class FirestoreOptions implements FirestoreOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.FirestoreOptions";
  readonly __MESSAGE_TYPE = FirestoreOptions.MESSAGE_TYPE;

  /** Source message fully qualified name */
  source: string = "";
  /** Collection name */
  collection: string = "";

  
}


/**
 * MongoDB target options
 */
export class MongoDBOptions implements MongoDBOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "dal.v1.MongoDBOptions";
  readonly __MESSAGE_TYPE = MongoDBOptions.MESSAGE_TYPE;

  /** Source message fully qualified name */
  source: string = "";
  /** Collection name */
  collection: string = "";
  /** Database name */
  database: string = "";

  
}


