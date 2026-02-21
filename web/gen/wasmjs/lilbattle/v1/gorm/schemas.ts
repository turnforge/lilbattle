
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for IndexStateGORM message
 */
export const IndexStateGORMSchema: MessageSchema = {
  name: "IndexStateGORM",
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
      name: "indexedAt",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "needsIndexing",
      type: FieldType.BOOLEAN,
      id: 4,
    },
  ],
};


/**
 * Schema for IndexRecordsLROGORM message
 */
export const IndexRecordsLROGORMSchema: MessageSchema = {
  name: "IndexRecordsLROGORM",
  fields: [
  ],
};


/**
 * Schema for IndexInfoGORM message
 */
export const IndexInfoGORMSchema: MessageSchema = {
  name: "IndexInfoGORM",
  fields: [
  ],
};


/**
 * Schema for TileGORM message
 */
export const TileGORMSchema: MessageSchema = {
  name: "TileGORM",
  fields: [
  ],
};


/**
 * Schema for CrossingGORM message
 */
export const CrossingGORMSchema: MessageSchema = {
  name: "CrossingGORM",
  fields: [
  ],
};


/**
 * Schema for UnitGORM message
 */
export const UnitGORMSchema: MessageSchema = {
  name: "UnitGORM",
  fields: [
  ],
};


/**
 * Schema for AttackRecordGORM message
 */
export const AttackRecordGORMSchema: MessageSchema = {
  name: "AttackRecordGORM",
  fields: [
  ],
};


/**
 * Schema for WorldGORM message
 */
export const WorldGORMSchema: MessageSchema = {
  name: "WorldGORM",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 7,
      repeated: true,
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 11,
      repeated: true,
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 13,
      messageType: "lilbattle.v1.IndexInfoGORM",
    },
  ],
};


/**
 * Schema for WorldDataGORM message
 */
export const WorldDataGORMSchema: MessageSchema = {
  name: "WorldDataGORM",
  fields: [
    {
      name: "worldId",
      type: FieldType.STRING,
      id: 1,
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
      messageType: "lilbattle.v1.IndexInfoGORM",
    },
    {
      name: "tilesMap",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "unitsMap",
      type: FieldType.STRING,
      id: 7,
    },
  ],
};


/**
 * Schema for GameGORM message
 */
export const GameGORMSchema: MessageSchema = {
  name: "GameGORM",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "worldId",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 7,
      repeated: true,
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 11,
      repeated: true,
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 13,
      messageType: "lilbattle.v1.IndexInfoGORM",
    },
  ],
};


/**
 * Schema for GameStateGORM message
 */
export const GameStateGORMSchema: MessageSchema = {
  name: "GameStateGORM",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameWorldDataGORM",
    },
    {
      name: "playerStates",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for GameConfigurationGORM message
 */
export const GameConfigurationGORMSchema: MessageSchema = {
  name: "GameConfigurationGORM",
  fields: [
    {
      name: "incomeConfigs",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.IncomeConfigGORM",
    },
    {
      name: "settings",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameSettingsGORM",
    },
  ],
};


/**
 * Schema for IncomeConfigGORM message
 */
export const IncomeConfigGORMSchema: MessageSchema = {
  name: "IncomeConfigGORM",
  fields: [
  ],
};


/**
 * Schema for GamePlayerGORM message
 */
export const GamePlayerGORMSchema: MessageSchema = {
  name: "GamePlayerGORM",
  fields: [
  ],
};


/**
 * Schema for GameTeamGORM message
 */
export const GameTeamGORMSchema: MessageSchema = {
  name: "GameTeamGORM",
  fields: [
  ],
};


/**
 * Schema for GameSettingsGORM message
 */
export const GameSettingsGORMSchema: MessageSchema = {
  name: "GameSettingsGORM",
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
 * Schema for PlayerStateGORM message
 */
export const PlayerStateGORMSchema: MessageSchema = {
  name: "PlayerStateGORM",
  fields: [
  ],
};


/**
 * Schema for GameWorldDataGORM message
 */
export const GameWorldDataGORMSchema: MessageSchema = {
  name: "GameWorldDataGORM",
  fields: [
    {
      name: "screenshotIndexInfo",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.IndexInfoGORM",
    },
    {
      name: "crossings",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "tilesMap",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "unitsMap",
      type: FieldType.STRING,
      id: 7,
    },
  ],
};


/**
 * Schema for GameMoveHistoryGORM message
 */
export const GameMoveHistoryGORMSchema: MessageSchema = {
  name: "GameMoveHistoryGORM",
  fields: [
  ],
};


/**
 * Schema for GameMoveGroupGORM message
 */
export const GameMoveGroupGORMSchema: MessageSchema = {
  name: "GameMoveGroupGORM",
  fields: [
  ],
};


/**
 * Schema for GameMoveGORM message
 */
export const GameMoveGORMSchema: MessageSchema = {
  name: "GameMoveGORM",
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
      name: "version",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "moveType",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Any",
    },
    {
      name: "changes",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "google.protobuf.Any",
      repeated: true,
    },
  ],
};



/**
 * Package-scoped schema registry for lilbattle.v1
 */
export const lilbattle_v1SchemaRegistry: Record<string, MessageSchema> = {
  "lilbattle.v1.IndexStateGORM": IndexStateGORMSchema,
  "lilbattle.v1.IndexRecordsLROGORM": IndexRecordsLROGORMSchema,
  "lilbattle.v1.IndexInfoGORM": IndexInfoGORMSchema,
  "lilbattle.v1.TileGORM": TileGORMSchema,
  "lilbattle.v1.CrossingGORM": CrossingGORMSchema,
  "lilbattle.v1.UnitGORM": UnitGORMSchema,
  "lilbattle.v1.AttackRecordGORM": AttackRecordGORMSchema,
  "lilbattle.v1.WorldGORM": WorldGORMSchema,
  "lilbattle.v1.WorldDataGORM": WorldDataGORMSchema,
  "lilbattle.v1.GameGORM": GameGORMSchema,
  "lilbattle.v1.GameStateGORM": GameStateGORMSchema,
  "lilbattle.v1.GameConfigurationGORM": GameConfigurationGORMSchema,
  "lilbattle.v1.IncomeConfigGORM": IncomeConfigGORMSchema,
  "lilbattle.v1.GamePlayerGORM": GamePlayerGORMSchema,
  "lilbattle.v1.GameTeamGORM": GameTeamGORMSchema,
  "lilbattle.v1.GameSettingsGORM": GameSettingsGORMSchema,
  "lilbattle.v1.PlayerStateGORM": PlayerStateGORMSchema,
  "lilbattle.v1.GameWorldDataGORM": GameWorldDataGORMSchema,
  "lilbattle.v1.GameMoveHistoryGORM": GameMoveHistoryGORMSchema,
  "lilbattle.v1.GameMoveGroupGORM": GameMoveGroupGORMSchema,
  "lilbattle.v1.GameMoveGORM": GameMoveGORMSchema,
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