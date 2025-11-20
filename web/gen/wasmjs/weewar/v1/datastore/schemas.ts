
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for IndexStateDatastore message
 */
export const IndexStateDatastoreSchema: MessageSchema = {
  name: "IndexStateDatastore",
  fields: [
  ],
};


/**
 * Schema for IndexRecordDatastore message
 */
export const IndexRecordDatastoreSchema: MessageSchema = {
  name: "IndexRecordDatastore",
  fields: [
  ],
};


/**
 * Schema for IndexRecordsLRODatastore message
 */
export const IndexRecordsLRODatastoreSchema: MessageSchema = {
  name: "IndexRecordsLRODatastore",
  fields: [
  ],
};



/**
 * Package-scoped schema registry for weewar.v1
 */
export const weewar_v1SchemaRegistry: Record<string, MessageSchema> = {
  "weewar.v1.IndexStateDatastore": IndexStateDatastoreSchema,
  "weewar.v1.IndexRecordDatastore": IndexRecordDatastoreSchema,
  "weewar.v1.IndexRecordsLRODatastore": IndexRecordsLRODatastoreSchema,
};

/**
 * Schema registry instance for weewar.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(weewar_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);