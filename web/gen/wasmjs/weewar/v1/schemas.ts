
// Generated TypeScript package-level schema registry
// DO NOT EDIT - This file is auto-generated
//
// This file consolidates schema registries from all subdirectories in this package.
// It imports directory-level schema registries and merges them into a single package-level registry.

import { MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";

// Import schema registries from all subdirectories
import { weewar_v1SchemaRegistry as modelsSchemas } from './models/schemas';

/**
 * Consolidated package-level schema registry for weewar.v1
 * Merges all directory-level schema registries into a single registry
 */
export const weewar_v1SchemaRegistry: Record<string, MessageSchema> = {
  ...modelsSchemas,
};

/**
 * Package-level schema registry instance with utility methods
 * Use this for deserializers and other schema-aware operations
 */
const registryInstance = new BaseSchemaRegistry(weewar_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);
