
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for TableOptions message
 */
export const TableOptionsSchema: MessageSchema = {
  name: "TableOptions",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "schema",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "comment",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "source",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for ColumnOptions message
 */
export const ColumnOptionsSchema: MessageSchema = {
  name: "ColumnOptions",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "toFunc",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "dal.v1.ConverterFunc",
    },
    {
      name: "fromFunc",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "dal.v1.ConverterFunc",
    },
    {
      name: "gormTags",
      type: FieldType.REPEATED,
      id: 10,
      repeated: true,
    },
    {
      name: "sqlTags",
      type: FieldType.REPEATED,
      id: 11,
      repeated: true,
    },
    {
      name: "firestoreTags",
      type: FieldType.REPEATED,
      id: 12,
      repeated: true,
    },
    {
      name: "mongodbTags",
      type: FieldType.REPEATED,
      id: 13,
      repeated: true,
    },
  ],
};


/**
 * Schema for ConverterFunc message
 */
export const ConverterFuncSchema: MessageSchema = {
  name: "ConverterFunc",
  fields: [
    {
      name: "package",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "alias",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "function",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for IndexOptions message
 */
export const IndexOptionsSchema: MessageSchema = {
  name: "IndexOptions",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "fields",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "unique",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "type",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "where",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for ForeignKeyOptions message
 */
export const ForeignKeyOptionsSchema: MessageSchema = {
  name: "ForeignKeyOptions",
  fields: [
    {
      name: "references",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "onDelete",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "onUpdate",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "constraintName",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for GormOptions message
 */
export const GormOptionsSchema: MessageSchema = {
  name: "GormOptions",
  fields: [
    {
      name: "source",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "table",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "embedded",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
  ],
};


/**
 * Schema for PostgresOptions message
 */
export const PostgresOptionsSchema: MessageSchema = {
  name: "PostgresOptions",
  fields: [
    {
      name: "source",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "table",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "schema",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for DatastoreOptions message
 */
export const DatastoreOptionsSchema: MessageSchema = {
  name: "DatastoreOptions",
  fields: [
    {
      name: "kind",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "namespace",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "incompleteKey",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "ancestor",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "source",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for FirestoreOptions message
 */
export const FirestoreOptionsSchema: MessageSchema = {
  name: "FirestoreOptions",
  fields: [
    {
      name: "source",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "collection",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for MongoDBOptions message
 */
export const MongoDBOptionsSchema: MessageSchema = {
  name: "MongoDBOptions",
  fields: [
    {
      name: "source",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "collection",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "database",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};



/**
 * Package-scoped schema registry for dal.v1
 */
export const dal_v1SchemaRegistry: Record<string, MessageSchema> = {
  "dal.v1.TableOptions": TableOptionsSchema,
  "dal.v1.ColumnOptions": ColumnOptionsSchema,
  "dal.v1.ConverterFunc": ConverterFuncSchema,
  "dal.v1.IndexOptions": IndexOptionsSchema,
  "dal.v1.ForeignKeyOptions": ForeignKeyOptionsSchema,
  "dal.v1.GormOptions": GormOptionsSchema,
  "dal.v1.PostgresOptions": PostgresOptionsSchema,
  "dal.v1.DatastoreOptions": DatastoreOptionsSchema,
  "dal.v1.FirestoreOptions": FirestoreOptionsSchema,
  "dal.v1.MongoDBOptions": MongoDBOptionsSchema,
};

/**
 * Schema registry instance for dal.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(dal_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);