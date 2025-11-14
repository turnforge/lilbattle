// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

/**
 * Referential actions for foreign keys
 */
export enum ReferentialAction {
  NO_ACTION = 0,
  RESTRICT = 1,
  CASCADE = 2,
  SET_NULL = 3,
  SET_DEFAULT = 4,
}


/**
 * Configuration for table mapping
 */
export interface TableOptions {
  /** Table name in the database */
  name: string;
  /** Database schema/namespace */
  schema: string;
  /** Custom table comment/description */
  comment: string;
  /** Source message fully qualified name (e.g., "library.v1.Book")
 This links the DAL schema message to the API proto message */
  source: string;
}


/**
 * Configuration for column mapping
 */
export interface ColumnOptions {
  /** Column name override (optional - defaults to field name) */
  name: string;
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
  gormTags: string[];
  /** Raw SQL column definition tags (for raw SQL targets)
 Target-specific usage */
  sqlTags: string[];
  /** Firestore tags (for Firestore target) */
  firestoreTags: string[];
  /** MongoDB tags (for MongoDB target) */
  mongodbTags: string[];
}


/**
 * Specification for a custom converter function
 */
export interface ConverterFunc {
  /** Go package import path
 Example: "github.com/myapp/converters" */
  package: string;
  /** Optional import alias
 If not specified, uses the last segment of package path
 Example: "myconv" -> import myconv "github.com/myapp/converters" */
  alias: string;
  /** Function name to call
 Example: "TimestampToMillis"
 Generates call: alias.TimestampToMillis(value) */
  function: string;
}


/**
 * Configuration for indexes
 */
export interface IndexOptions {
  /** Index name */
  name: string;
  /** Comma-separated list of field names (for composite indexes) */
  fields: string;
  /** Whether this is a unique index */
  unique: boolean;
  /** Index type (e.g., "BTREE", "HASH", "GIN", "GIST") */
  type: string;
  /** Partial index condition (WHERE clause) */
  where: string;
}


/**
 * Configuration for foreign keys
 */
export interface ForeignKeyOptions {
  /** Referenced table and column (e.g., "authors.id") */
  references: string;
  /** Action on delete */
  onDelete: ReferentialAction;
  /** Action on update */
  onUpdate: ReferentialAction;
  /** Foreign key constraint name */
  constraintName: string;
}


/**
 * GORM target options (database-agnostic ORM)
 */
export interface GormOptions {
  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string;
  /** Table name */
  table: string;
  /** Embedded field names (for GORM embedded structs) */
  embedded: string[];
}


/**
 * PostgreSQL target options (raw SQL)
 */
export interface PostgresOptions {
  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string;
  /** Table name in PostgreSQL */
  table: string;
  /** Schema name (default: "public") */
  schema: string;
}


/**
 * Google Cloud Datastore options
 */
export interface DatastoreOptions {
  /** Datastore kind name */
  kind: string;
  /** Datastore namespace */
  namespace: string;
  /** Whether to use incomplete keys (auto-generated) */
  incompleteKey: boolean;
  /** Ancestor path */
  ancestor: string;
  /** Source message fully qualified name (e.g., "library.v1.Book") */
  source: string;
}


/**
 * Firestore target options
 */
export interface FirestoreOptions {
  /** Source message fully qualified name */
  source: string;
  /** Collection name */
  collection: string;
}


/**
 * MongoDB target options
 */
export interface MongoDBOptions {
  /** Source message fully qualified name */
  source: string;
  /** Collection name */
  collection: string;
  /** Database name */
  database: string;
}

