
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
 * Schema for IndexInfoDatastore message
 */
export const IndexInfoDatastoreSchema: MessageSchema = {
  name: "IndexInfoDatastore",
  fields: [
  ],
};


/**
 * Schema for TileDatastore message
 */
export const TileDatastoreSchema: MessageSchema = {
  name: "TileDatastore",
  fields: [
  ],
};


/**
 * Schema for CrossingDatastore message
 */
export const CrossingDatastoreSchema: MessageSchema = {
  name: "CrossingDatastore",
  fields: [
  ],
};


/**
 * Schema for UnitDatastore message
 */
export const UnitDatastoreSchema: MessageSchema = {
  name: "UnitDatastore",
  fields: [
    {
      name: "attackHistory",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.AttackRecordDatastore",
      repeated: true,
    },
  ],
};


/**
 * Schema for AttackRecordDatastore message
 */
export const AttackRecordDatastoreSchema: MessageSchema = {
  name: "AttackRecordDatastore",
  fields: [
  ],
};


/**
 * Schema for WorldDatastore message
 */
export const WorldDatastoreSchema: MessageSchema = {
  name: "WorldDatastore",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 2,
      repeated: true,
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
    {
      name: "defaultGameConfig",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameConfigurationDatastore",
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.IndexInfoDatastore",
    },
  ],
};


/**
 * Schema for WorldDataDatastore message
 */
export const WorldDataDatastoreSchema: MessageSchema = {
  name: "WorldDataDatastore",
  fields: [
    {
      name: "worldId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "tilesMap",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "unitsMap",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "crossings",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "screenshotIndexInfo",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.IndexInfoDatastore",
    },
  ],
};


/**
 * Schema for GameDatastore message
 */
export const GameDatastoreSchema: MessageSchema = {
  name: "GameDatastore",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "worldId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 4,
      repeated: true,
    },
    {
      name: "config",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.GameConfigurationDatastore",
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.IndexInfoDatastore",
    },
  ],
};


/**
 * Schema for GameStateDatastore message
 */
export const GameStateDatastoreSchema: MessageSchema = {
  name: "GameStateDatastore",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldDataDatastore",
    },
    {
      name: "playerStates",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for GameConfigurationDatastore message
 */
export const GameConfigurationDatastoreSchema: MessageSchema = {
  name: "GameConfigurationDatastore",
  fields: [
    {
      name: "players",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.GamePlayerDatastore",
      repeated: true,
    },
    {
      name: "teams",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameTeamDatastore",
      repeated: true,
    },
    {
      name: "incomeConfigs",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.IncomeConfigDatastore",
    },
    {
      name: "settings",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameSettingsDatastore",
    },
  ],
};


/**
 * Schema for IncomeConfigDatastore message
 */
export const IncomeConfigDatastoreSchema: MessageSchema = {
  name: "IncomeConfigDatastore",
  fields: [
  ],
};


/**
 * Schema for GamePlayerDatastore message
 */
export const GamePlayerDatastoreSchema: MessageSchema = {
  name: "GamePlayerDatastore",
  fields: [
  ],
};


/**
 * Schema for GameTeamDatastore message
 */
export const GameTeamDatastoreSchema: MessageSchema = {
  name: "GameTeamDatastore",
  fields: [
  ],
};


/**
 * Schema for GameSettingsDatastore message
 */
export const GameSettingsDatastoreSchema: MessageSchema = {
  name: "GameSettingsDatastore",
  fields: [
    {
      name: "allowedUnits",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for PlayerStateDatastore message
 */
export const PlayerStateDatastoreSchema: MessageSchema = {
  name: "PlayerStateDatastore",
  fields: [
  ],
};


/**
 * Schema for GameMoveDatastore message
 */
export const GameMoveDatastoreSchema: MessageSchema = {
  name: "GameMoveDatastore",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "groupNumber",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "moveNumber",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "moveType",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Any",
    },
    {
      name: "changes",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Any",
      repeated: true,
    },
  ],
};



/**
 * Package-scoped schema registry for lilbattle.v1
 */
export const lilbattle_v1SchemaRegistry: Record<string, MessageSchema> = {
  "lilbattle.v1.IndexStateDatastore": IndexStateDatastoreSchema,
  "lilbattle.v1.IndexRecordDatastore": IndexRecordDatastoreSchema,
  "lilbattle.v1.IndexRecordsLRODatastore": IndexRecordsLRODatastoreSchema,
  "lilbattle.v1.IndexInfoDatastore": IndexInfoDatastoreSchema,
  "lilbattle.v1.TileDatastore": TileDatastoreSchema,
  "lilbattle.v1.CrossingDatastore": CrossingDatastoreSchema,
  "lilbattle.v1.UnitDatastore": UnitDatastoreSchema,
  "lilbattle.v1.AttackRecordDatastore": AttackRecordDatastoreSchema,
  "lilbattle.v1.WorldDatastore": WorldDatastoreSchema,
  "lilbattle.v1.WorldDataDatastore": WorldDataDatastoreSchema,
  "lilbattle.v1.GameDatastore": GameDatastoreSchema,
  "lilbattle.v1.GameStateDatastore": GameStateDatastoreSchema,
  "lilbattle.v1.GameConfigurationDatastore": GameConfigurationDatastoreSchema,
  "lilbattle.v1.IncomeConfigDatastore": IncomeConfigDatastoreSchema,
  "lilbattle.v1.GamePlayerDatastore": GamePlayerDatastoreSchema,
  "lilbattle.v1.GameTeamDatastore": GameTeamDatastoreSchema,
  "lilbattle.v1.GameSettingsDatastore": GameSettingsDatastoreSchema,
  "lilbattle.v1.PlayerStateDatastore": PlayerStateDatastoreSchema,
  "lilbattle.v1.GameMoveDatastore": GameMoveDatastoreSchema,
};

/**
 * Schema registry instance for lilbattle.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(lilbattle_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);