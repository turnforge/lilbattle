
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for EntityIndexStateDatastore message
 */
export const EntityIndexStateDatastoreSchema: MessageSchema = {
  name: "EntityIndexStateDatastore",
  fields: [
    {
      name: "entityType",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "entityId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "indexType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "lastQueuedAt",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "lastIndexedAt",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "status",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "lastError",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "lastContentHash",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "retryCount",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "currentLroId",
      type: FieldType.STRING,
      id: 10,
    },
  ],
};


/**
 * Schema for IndexRecordDatastore message
 */
export const IndexRecordDatastoreSchema: MessageSchema = {
  name: "IndexRecordDatastore",
  fields: [
    {
      name: "entityId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "updatedAt",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "indexerTypes",
      type: FieldType.REPEATED,
      id: 5,
      repeated: true,
    },
  ],
};


/**
 * Schema for IndexRecordsLRODatastore message
 */
export const IndexRecordsLRODatastoreSchema: MessageSchema = {
  name: "IndexRecordsLRODatastore",
  fields: [
    {
      name: "lroId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "entityType",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "createdAt",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "updatedAt",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "callbackUrl",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "records",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "weewar.v1.IndexRecordDatastore",
      repeated: true,
    },
  ],
};



/**
 * Package-scoped schema registry for weewar.v1
 */
export const weewar_v1SchemaRegistry: Record<string, MessageSchema> = {
  "weewar.v1.EntityIndexStateDatastore": EntityIndexStateDatastoreSchema,
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