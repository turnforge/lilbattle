
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for IndexInfo message
 */
export const IndexInfoSchema: MessageSchema = {
  name: "IndexInfo",
  fields: [
    {
      name: "lastUpdatedAt",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "lastIndexedAt",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "needsIndexing",
      type: FieldType.BOOLEAN,
      id: 3,
    },
  ],
};


/**
 * Schema for Pagination message
 */
export const PaginationSchema: MessageSchema = {
  name: "Pagination",
  fields: [
    {
      name: "pageKey",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pageOffset",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "pageSize",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for PaginationResponse message
 */
export const PaginationResponseSchema: MessageSchema = {
  name: "PaginationResponse",
  fields: [
    {
      name: "nextPageKey",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "nextPageOffset",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "hasMore",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "totalResults",
      type: FieldType.NUMBER,
      id: 5,
    },
  ],
};


/**
 * Schema for World message
 */
export const WorldSchema: MessageSchema = {
  name: "World",
  fields: [
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "version",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "id",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "creatorId",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 8,
      repeated: true,
    },
    {
      name: "imageUrl",
      type: FieldType.STRING,
      id: 9,
    },
    {
      name: "difficulty",
      type: FieldType.STRING,
      id: 10,
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 11,
      repeated: true,
    },
    {
      name: "defaultGameConfig",
      type: FieldType.MESSAGE,
      id: 12,
      messageType: "lilbattle.v1.GameConfiguration",
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 13,
      messageType: "lilbattle.v1.IndexInfo",
    },
  ],
};


/**
 * Schema for WorldData message
 */
export const WorldDataSchema: MessageSchema = {
  name: "WorldData",
  fields: [
    {
      name: "tilesMap",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "unitsMap",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "screenshotIndexInfo",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.IndexInfo",
    },
    {
      name: "contentHash",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "version",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "crossings",
      type: FieldType.STRING,
      id: 8,
    },
  ],
};


/**
 * Schema for Crossing message
 */
export const CrossingSchema: MessageSchema = {
  name: "Crossing",
  fields: [
    {
      name: "type",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "connectsTo",
      type: FieldType.REPEATED,
      id: 2,
      repeated: true,
    },
  ],
};


/**
 * Schema for Tile message
 */
export const TileSchema: MessageSchema = {
  name: "Tile",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "tileType",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "player",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "shortcut",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "lastActedTurn",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "lastToppedupTurn",
      type: FieldType.NUMBER,
      id: 7,
    },
  ],
};


/**
 * Schema for Unit message
 */
export const UnitSchema: MessageSchema = {
  name: "Unit",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "player",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "unitType",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "shortcut",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "availableHealth",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "distanceLeft",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "lastActedTurn",
      type: FieldType.NUMBER,
      id: 8,
    },
    {
      name: "lastToppedupTurn",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "attacksReceivedThisTurn",
      type: FieldType.NUMBER,
      id: 10,
    },
    {
      name: "attackHistory",
      type: FieldType.MESSAGE,
      id: 11,
      messageType: "lilbattle.v1.AttackRecord",
      repeated: true,
    },
    {
      name: "progressionStep",
      type: FieldType.NUMBER,
      id: 12,
    },
    {
      name: "chosenAlternative",
      type: FieldType.STRING,
      id: 13,
    },
    {
      name: "captureStartedTurn",
      type: FieldType.NUMBER,
      id: 14,
    },
  ],
};


/**
 * Schema for AttackRecord message
 */
export const AttackRecordSchema: MessageSchema = {
  name: "AttackRecord",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "isRanged",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "turnNumber",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for TerrainDefinition message
 */
export const TerrainDefinitionSchema: MessageSchema = {
  name: "TerrainDefinition",
  fields: [
    {
      name: "id",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "type",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "unitProperties",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "buildableUnitIds",
      type: FieldType.REPEATED,
      id: 8,
      repeated: true,
    },
    {
      name: "incomePerTurn",
      type: FieldType.NUMBER,
      id: 9,
    },
  ],
};


/**
 * Schema for UnitDefinition message
 */
export const UnitDefinitionSchema: MessageSchema = {
  name: "UnitDefinition",
  fields: [
    {
      name: "id",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "health",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "coins",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "movementPoints",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "retreatPoints",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "defense",
      type: FieldType.NUMBER,
      id: 8,
    },
    {
      name: "attackRange",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "minAttackRange",
      type: FieldType.NUMBER,
      id: 10,
    },
    {
      name: "splashDamage",
      type: FieldType.NUMBER,
      id: 11,
    },
    {
      name: "terrainProperties",
      type: FieldType.STRING,
      id: 12,
    },
    {
      name: "properties",
      type: FieldType.REPEATED,
      id: 13,
      repeated: true,
    },
    {
      name: "unitClass",
      type: FieldType.STRING,
      id: 14,
    },
    {
      name: "unitTerrain",
      type: FieldType.STRING,
      id: 15,
    },
    {
      name: "attackVsClass",
      type: FieldType.STRING,
      id: 16,
    },
    {
      name: "actionOrder",
      type: FieldType.REPEATED,
      id: 17,
      repeated: true,
    },
    {
      name: "actionLimits",
      type: FieldType.STRING,
      id: 18,
    },
    {
      name: "fixValue",
      type: FieldType.NUMBER,
      id: 19,
    },
  ],
};


/**
 * Schema for TerrainUnitProperties message
 */
export const TerrainUnitPropertiesSchema: MessageSchema = {
  name: "TerrainUnitProperties",
  fields: [
    {
      name: "terrainId",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "unitId",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "movementCost",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "healingBonus",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "canBuild",
      type: FieldType.BOOLEAN,
      id: 5,
    },
    {
      name: "canCapture",
      type: FieldType.BOOLEAN,
      id: 6,
    },
    {
      name: "attackBonus",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "defenseBonus",
      type: FieldType.NUMBER,
      id: 8,
    },
    {
      name: "attackRange",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "minAttackRange",
      type: FieldType.NUMBER,
      id: 10,
    },
  ],
};


/**
 * Schema for UnitUnitProperties message
 */
export const UnitUnitPropertiesSchema: MessageSchema = {
  name: "UnitUnitProperties",
  fields: [
    {
      name: "attackerId",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "defenderId",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "attackOverride",
      type: FieldType.STRING,
      id: 3,
      oneofGroup: "_attack_override",
      optional: true,
    },
    {
      name: "defenseOverride",
      type: FieldType.STRING,
      id: 4,
      oneofGroup: "_defense_override",
      optional: true,
    },
    {
      name: "damage",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.DamageDistribution",
    },
  ],
  oneofGroups: ["_attack_override", "_defense_override"],
};


/**
 * Schema for DamageDistribution message
 */
export const DamageDistributionSchema: MessageSchema = {
  name: "DamageDistribution",
  fields: [
    {
      name: "minDamage",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "maxDamage",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "expectedDamage",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "ranges",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.DamageRange",
      repeated: true,
    },
  ],
};


/**
 * Schema for DamageRange message
 */
export const DamageRangeSchema: MessageSchema = {
  name: "DamageRange",
  fields: [
    {
      name: "minValue",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "maxValue",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "probability",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for RulesEngine message
 */
export const RulesEngineSchema: MessageSchema = {
  name: "RulesEngine",
  fields: [
    {
      name: "units",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "terrains",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "terrainUnitProperties",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "unitUnitProperties",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "terrainTypes",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for Game message
 */
export const GameSchema: MessageSchema = {
  name: "Game",
  fields: [
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "version",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "id",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "creatorId",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "worldId",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 9,
      repeated: true,
    },
    {
      name: "imageUrl",
      type: FieldType.STRING,
      id: 10,
    },
    {
      name: "difficulty",
      type: FieldType.STRING,
      id: 11,
    },
    {
      name: "config",
      type: FieldType.MESSAGE,
      id: 12,
      messageType: "lilbattle.v1.GameConfiguration",
    },
    {
      name: "previewUrls",
      type: FieldType.REPEATED,
      id: 13,
      repeated: true,
    },
    {
      name: "searchIndexInfo",
      type: FieldType.MESSAGE,
      id: 15,
      messageType: "lilbattle.v1.IndexInfo",
    },
  ],
};


/**
 * Schema for GameConfiguration message
 */
export const GameConfigurationSchema: MessageSchema = {
  name: "GameConfiguration",
  fields: [
    {
      name: "players",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.GamePlayer",
      repeated: true,
    },
    {
      name: "teams",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameTeam",
      repeated: true,
    },
    {
      name: "incomeConfigs",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.IncomeConfig",
    },
    {
      name: "settings",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameSettings",
    },
  ],
};


/**
 * Schema for IncomeConfig message
 */
export const IncomeConfigSchema: MessageSchema = {
  name: "IncomeConfig",
  fields: [
    {
      name: "startingCoins",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "gameIncome",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "landbaseIncome",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "navalbaseIncome",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "airportbaseIncome",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "missilesiloIncome",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "minesIncome",
      type: FieldType.NUMBER,
      id: 7,
    },
  ],
};


/**
 * Schema for GamePlayer message
 */
export const GamePlayerSchema: MessageSchema = {
  name: "GamePlayer",
  fields: [
    {
      name: "playerId",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "playerType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "color",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "teamId",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "isActive",
      type: FieldType.BOOLEAN,
      id: 7,
    },
    {
      name: "startingCoins",
      type: FieldType.NUMBER,
      id: 8,
    },
  ],
};


/**
 * Schema for GameTeam message
 */
export const GameTeamSchema: MessageSchema = {
  name: "GameTeam",
  fields: [
    {
      name: "teamId",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "color",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "isActive",
      type: FieldType.BOOLEAN,
      id: 4,
    },
  ],
};


/**
 * Schema for GameSettings message
 */
export const GameSettingsSchema: MessageSchema = {
  name: "GameSettings",
  fields: [
    {
      name: "allowedUnits",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
    {
      name: "turnTimeLimit",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "teamMode",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "maxTurns",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for PlayerState message
 */
export const PlayerStateSchema: MessageSchema = {
  name: "PlayerState",
  fields: [
    {
      name: "coins",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "isActive",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for GameState message
 */
export const GameStateSchema: MessageSchema = {
  name: "GameState",
  fields: [
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "turnCounter",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "currentPlayer",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.WorldData",
    },
    {
      name: "stateHash",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "version",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "status",
      type: FieldType.STRING,
      id: 10,
    },
    {
      name: "finished",
      type: FieldType.BOOLEAN,
      id: 11,
    },
    {
      name: "winningPlayer",
      type: FieldType.NUMBER,
      id: 12,
    },
    {
      name: "winningTeam",
      type: FieldType.NUMBER,
      id: 13,
    },
    {
      name: "currentGroupNumber",
      type: FieldType.NUMBER,
      id: 14,
    },
    {
      name: "playerStates",
      type: FieldType.STRING,
      id: 15,
    },
  ],
};


/**
 * Schema for GameMoveHistory message
 */
export const GameMoveHistorySchema: MessageSchema = {
  name: "GameMoveHistory",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "groups",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameMoveGroup",
      repeated: true,
    },
  ],
};


/**
 * Schema for GameMoveGroup message
 */
export const GameMoveGroupSchema: MessageSchema = {
  name: "GameMoveGroup",
  fields: [
    {
      name: "startedAt",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "endedAt",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "groupNumber",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "moves",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.GameMove",
      repeated: true,
    },
  ],
};


/**
 * Schema for GameMove message
 */
export const GameMoveSchema: MessageSchema = {
  name: "GameMove",
  fields: [
    {
      name: "player",
      type: FieldType.NUMBER,
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
      name: "timestamp",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "moveUnit",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.MoveUnitAction",
      oneofGroup: "move_type",
    },
    {
      name: "attackUnit",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.AttackUnitAction",
      oneofGroup: "move_type",
    },
    {
      name: "endTurn",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "lilbattle.v1.EndTurnAction",
      oneofGroup: "move_type",
    },
    {
      name: "buildUnit",
      type: FieldType.MESSAGE,
      id: 8,
      messageType: "lilbattle.v1.BuildUnitAction",
      oneofGroup: "move_type",
    },
    {
      name: "captureBuilding",
      type: FieldType.MESSAGE,
      id: 13,
      messageType: "lilbattle.v1.CaptureBuildingAction",
      oneofGroup: "move_type",
    },
    {
      name: "healUnit",
      type: FieldType.MESSAGE,
      id: 14,
      messageType: "lilbattle.v1.HealUnitAction",
      oneofGroup: "move_type",
    },
    {
      name: "fixUnit",
      type: FieldType.MESSAGE,
      id: 15,
      messageType: "lilbattle.v1.FixUnitAction",
      oneofGroup: "move_type",
    },
    {
      name: "sequenceNum",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "isPermanent",
      type: FieldType.BOOLEAN,
      id: 10,
    },
    {
      name: "changes",
      type: FieldType.MESSAGE,
      id: 11,
      messageType: "lilbattle.v1.WorldChange",
      repeated: true,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 12,
    },
  ],
  oneofGroups: ["move_type"],
};


/**
 * Schema for Position message
 */
export const PositionSchema: MessageSchema = {
  name: "Position",
  fields: [
    {
      name: "label",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for MoveUnitAction message
 */
export const MoveUnitActionSchema: MessageSchema = {
  name: "MoveUnitAction",
  fields: [
    {
      name: "from",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "to",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "movementCost",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "reconstructedPath",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.Path",
    },
  ],
};


/**
 * Schema for AttackUnitAction message
 */
export const AttackUnitActionSchema: MessageSchema = {
  name: "AttackUnitAction",
  fields: [
    {
      name: "attacker",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "defender",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "targetUnitType",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "targetUnitHealth",
      type: FieldType.NUMBER,
      id: 8,
    },
    {
      name: "canAttack",
      type: FieldType.BOOLEAN,
      id: 9,
    },
    {
      name: "damageEstimate",
      type: FieldType.NUMBER,
      id: 10,
    },
  ],
};


/**
 * Schema for BuildUnitAction message
 */
export const BuildUnitActionSchema: MessageSchema = {
  name: "BuildUnitAction",
  fields: [
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "unitType",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "cost",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for CaptureBuildingAction message
 */
export const CaptureBuildingActionSchema: MessageSchema = {
  name: "CaptureBuildingAction",
  fields: [
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "tileType",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for EndTurnAction message
 */
export const EndTurnActionSchema: MessageSchema = {
  name: "EndTurnAction",
  fields: [
  ],
};


/**
 * Schema for HealUnitAction message
 */
export const HealUnitActionSchema: MessageSchema = {
  name: "HealUnitAction",
  fields: [
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "healAmount",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for FixUnitAction message
 */
export const FixUnitActionSchema: MessageSchema = {
  name: "FixUnitAction",
  fields: [
    {
      name: "fixer",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "target",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "fixAmount",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for WorldChange message
 */
export const WorldChangeSchema: MessageSchema = {
  name: "WorldChange",
  fields: [
    {
      name: "unitMoved",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.UnitMovedChange",
      oneofGroup: "change_type",
    },
    {
      name: "unitDamaged",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.UnitDamagedChange",
      oneofGroup: "change_type",
    },
    {
      name: "unitKilled",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.UnitKilledChange",
      oneofGroup: "change_type",
    },
    {
      name: "playerChanged",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.PlayerChangedChange",
      oneofGroup: "change_type",
    },
    {
      name: "unitBuilt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.UnitBuiltChange",
      oneofGroup: "change_type",
    },
    {
      name: "coinsChanged",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.CoinsChangedChange",
      oneofGroup: "change_type",
    },
    {
      name: "tileCaptured",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "lilbattle.v1.TileCapturedChange",
      oneofGroup: "change_type",
    },
    {
      name: "captureStarted",
      type: FieldType.MESSAGE,
      id: 8,
      messageType: "lilbattle.v1.CaptureStartedChange",
      oneofGroup: "change_type",
    },
    {
      name: "unitHealed",
      type: FieldType.MESSAGE,
      id: 9,
      messageType: "lilbattle.v1.UnitHealedChange",
      oneofGroup: "change_type",
    },
    {
      name: "unitFixed",
      type: FieldType.MESSAGE,
      id: 10,
      messageType: "lilbattle.v1.UnitFixedChange",
      oneofGroup: "change_type",
    },
  ],
  oneofGroups: ["change_type"],
};


/**
 * Schema for UnitHealedChange message
 */
export const UnitHealedChangeSchema: MessageSchema = {
  name: "UnitHealedChange",
  fields: [
    {
      name: "previousUnit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "updatedUnit",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "healAmount",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for UnitFixedChange message
 */
export const UnitFixedChangeSchema: MessageSchema = {
  name: "UnitFixedChange",
  fields: [
    {
      name: "fixerUnit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "previousTarget",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "updatedTarget",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "fixAmount",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for UnitMovedChange message
 */
export const UnitMovedChangeSchema: MessageSchema = {
  name: "UnitMovedChange",
  fields: [
    {
      name: "previousUnit",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "updatedUnit",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "lilbattle.v1.Unit",
    },
  ],
};


/**
 * Schema for UnitDamagedChange message
 */
export const UnitDamagedChangeSchema: MessageSchema = {
  name: "UnitDamagedChange",
  fields: [
    {
      name: "previousUnit",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "updatedUnit",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "lilbattle.v1.Unit",
    },
  ],
};


/**
 * Schema for UnitKilledChange message
 */
export const UnitKilledChangeSchema: MessageSchema = {
  name: "UnitKilledChange",
  fields: [
    {
      name: "previousUnit",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.Unit",
    },
  ],
};


/**
 * Schema for PlayerChangedChange message
 */
export const PlayerChangedChangeSchema: MessageSchema = {
  name: "PlayerChangedChange",
  fields: [
    {
      name: "previousPlayer",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "newPlayer",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "previousTurn",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "newTurn",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "resetUnits",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.Unit",
      repeated: true,
    },
  ],
};


/**
 * Schema for UnitBuiltChange message
 */
export const UnitBuiltChangeSchema: MessageSchema = {
  name: "UnitBuiltChange",
  fields: [
    {
      name: "unit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "tileQ",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "tileR",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "coinsCost",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "playerCoins",
      type: FieldType.NUMBER,
      id: 5,
    },
  ],
};


/**
 * Schema for CoinsChangedChange message
 */
export const CoinsChangedChangeSchema: MessageSchema = {
  name: "CoinsChangedChange",
  fields: [
    {
      name: "playerId",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "previousCoins",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "newCoins",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "reason",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for TileCapturedChange message
 */
export const TileCapturedChangeSchema: MessageSchema = {
  name: "TileCapturedChange",
  fields: [
    {
      name: "capturingUnit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "tileQ",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "tileR",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "tileType",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "previousOwner",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "newOwner",
      type: FieldType.NUMBER,
      id: 6,
    },
  ],
};


/**
 * Schema for CaptureStartedChange message
 */
export const CaptureStartedChangeSchema: MessageSchema = {
  name: "CaptureStartedChange",
  fields: [
    {
      name: "capturingUnit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "tileQ",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "tileR",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "tileType",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "currentOwner",
      type: FieldType.NUMBER,
      id: 5,
    },
  ],
};


/**
 * Schema for AllPaths message
 */
export const AllPathsSchema: MessageSchema = {
  name: "AllPaths",
  fields: [
    {
      name: "sourceQ",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "sourceR",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "edges",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for PathEdge message
 */
export const PathEdgeSchema: MessageSchema = {
  name: "PathEdge",
  fields: [
    {
      name: "fromQ",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "fromR",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "toQ",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "toR",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "movementCost",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "totalCost",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "terrainType",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "explanation",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "isOccupied",
      type: FieldType.BOOLEAN,
      id: 9,
    },
  ],
};


/**
 * Schema for Path message
 */
export const PathSchema: MessageSchema = {
  name: "Path",
  fields: [
    {
      name: "edges",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.PathEdge",
      repeated: true,
    },
    {
      name: "directions",
      type: FieldType.REPEATED,
      id: 2,
      repeated: true,
    },
    {
      name: "totalCost",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for File message
 */
export const FileSchema: MessageSchema = {
  name: "File",
  fields: [
    {
      name: "path",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "contentType",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "fileSize",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "isPublic",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "downloadUrl",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "signedUrls",
      type: FieldType.STRING,
      id: 8,
    },
  ],
};


/**
 * Schema for PutFileRequest message
 */
export const PutFileRequestSchema: MessageSchema = {
  name: "PutFileRequest",
  fields: [
    {
      name: "file",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.File",
    },
    {
      name: "content",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for PutFileResponse message
 */
export const PutFileResponseSchema: MessageSchema = {
  name: "PutFileResponse",
  fields: [
    {
      name: "file",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.File",
    },
  ],
};


/**
 * Schema for GetFileRequest message
 */
export const GetFileRequestSchema: MessageSchema = {
  name: "GetFileRequest",
  fields: [
    {
      name: "path",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "includeSignedUrls",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for GetFileResponse message
 */
export const GetFileResponseSchema: MessageSchema = {
  name: "GetFileResponse",
  fields: [
    {
      name: "file",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.File",
    },
  ],
};


/**
 * Schema for DeleteFileRequest message
 */
export const DeleteFileRequestSchema: MessageSchema = {
  name: "DeleteFileRequest",
  fields: [
    {
      name: "path",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for DeleteFileResponse message
 */
export const DeleteFileResponseSchema: MessageSchema = {
  name: "DeleteFileResponse",
  fields: [
    {
      name: "file",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.File",
    },
  ],
};


/**
 * Schema for ListFilesRequest message
 */
export const ListFilesRequestSchema: MessageSchema = {
  name: "ListFilesRequest",
  fields: [
    {
      name: "path",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Pagination",
    },
    {
      name: "includeSignedUrls",
      type: FieldType.BOOLEAN,
      id: 3,
    },
  ],
};


/**
 * Schema for ListFilesResponse message
 */
export const ListFilesResponseSchema: MessageSchema = {
  name: "ListFilesResponse",
  fields: [
    {
      name: "items",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.File",
      repeated: true,
    },
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.PaginationResponse",
    },
  ],
};


/**
 * Schema for ListGamesRequest message
 */
export const ListGamesRequestSchema: MessageSchema = {
  name: "ListGamesRequest",
  fields: [
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Pagination",
    },
    {
      name: "ownerId",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for ListGamesResponse message
 */
export const ListGamesResponseSchema: MessageSchema = {
  name: "ListGamesResponse",
  fields: [
    {
      name: "items",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
      repeated: true,
    },
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.PaginationResponse",
    },
  ],
};


/**
 * Schema for GetGameRequest message
 */
export const GetGameRequestSchema: MessageSchema = {
  name: "GetGameRequest",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "version",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for GetGameResponse message
 */
export const GetGameResponseSchema: MessageSchema = {
  name: "GetGameResponse",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
    {
      name: "state",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameState",
    },
    {
      name: "history",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.GameMoveHistory",
    },
  ],
};


/**
 * Schema for GetGameContentRequest message
 */
export const GetGameContentRequestSchema: MessageSchema = {
  name: "GetGameContentRequest",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "version",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for GetGameContentResponse message
 */
export const GetGameContentResponseSchema: MessageSchema = {
  name: "GetGameContentResponse",
  fields: [
    {
      name: "lilbattleContent",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "recipeContent",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "readmeContent",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for UpdateGameRequest message
 */
export const UpdateGameRequestSchema: MessageSchema = {
  name: "UpdateGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "newGame",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Game",
    },
    {
      name: "newState",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.GameState",
    },
    {
      name: "newHistory",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.GameMoveHistory",
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.FieldMask",
    },
  ],
};


/**
 * Schema for UpdateGameResponse message
 */
export const UpdateGameResponseSchema: MessageSchema = {
  name: "UpdateGameResponse",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
  ],
};


/**
 * Schema for DeleteGameRequest message
 */
export const DeleteGameRequestSchema: MessageSchema = {
  name: "DeleteGameRequest",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for DeleteGameResponse message
 */
export const DeleteGameResponseSchema: MessageSchema = {
  name: "DeleteGameResponse",
  fields: [
  ],
};


/**
 * Schema for GetGamesRequest message
 */
export const GetGamesRequestSchema: MessageSchema = {
  name: "GetGamesRequest",
  fields: [
    {
      name: "ids",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for GetGamesResponse message
 */
export const GetGamesResponseSchema: MessageSchema = {
  name: "GetGamesResponse",
  fields: [
    {
      name: "games",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for CreateGameRequest message
 */
export const CreateGameRequestSchema: MessageSchema = {
  name: "CreateGameRequest",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
  ],
};


/**
 * Schema for CreateGameResponse message
 */
export const CreateGameResponseSchema: MessageSchema = {
  name: "CreateGameResponse",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
    {
      name: "gameState",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameState",
    },
    {
      name: "fieldErrors",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for ProcessMovesRequest message
 */
export const ProcessMovesRequestSchema: MessageSchema = {
  name: "ProcessMovesRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "moves",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameMove",
      repeated: true,
    },
    {
      name: "expectedResponse",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.ProcessMovesResponse",
    },
    {
      name: "dryRun",
      type: FieldType.BOOLEAN,
      id: 4,
    },
  ],
};


/**
 * Schema for ProcessMovesResponse message
 */
export const ProcessMovesResponseSchema: MessageSchema = {
  name: "ProcessMovesResponse",
  fields: [
    {
      name: "moves",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.GameMove",
      repeated: true,
    },
  ],
};


/**
 * Schema for GetGameStateRequest message
 */
export const GetGameStateRequestSchema: MessageSchema = {
  name: "GetGameStateRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for GetGameStateResponse message
 */
export const GetGameStateResponseSchema: MessageSchema = {
  name: "GetGameStateResponse",
  fields: [
    {
      name: "state",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.GameState",
    },
  ],
};


/**
 * Schema for ListMovesRequest message
 */
export const ListMovesRequestSchema: MessageSchema = {
  name: "ListMovesRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "fromGroup",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "toGroup",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for ListMovesResponse message
 */
export const ListMovesResponseSchema: MessageSchema = {
  name: "ListMovesResponse",
  fields: [
    {
      name: "hasMore",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "moveGroups",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameMoveGroup",
      repeated: true,
    },
  ],
};


/**
 * Schema for GetOptionsAtRequest message
 */
export const GetOptionsAtRequestSchema: MessageSchema = {
  name: "GetOptionsAtRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
  ],
};


/**
 * Schema for GetOptionsAtResponse message
 */
export const GetOptionsAtResponseSchema: MessageSchema = {
  name: "GetOptionsAtResponse",
  fields: [
    {
      name: "options",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.GameOption",
      repeated: true,
    },
    {
      name: "currentPlayer",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "gameInitialized",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "allPaths",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.AllPaths",
    },
  ],
};


/**
 * Schema for GameOption message
 */
export const GameOptionSchema: MessageSchema = {
  name: "GameOption",
  fields: [
    {
      name: "move",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.MoveUnitAction",
      oneofGroup: "option_type",
    },
    {
      name: "attack",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.AttackUnitAction",
      oneofGroup: "option_type",
    },
    {
      name: "build",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.BuildUnitAction",
      oneofGroup: "option_type",
    },
    {
      name: "capture",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.CaptureBuildingAction",
      oneofGroup: "option_type",
    },
    {
      name: "endTurn",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.EndTurnAction",
      oneofGroup: "option_type",
    },
    {
      name: "heal",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.HealUnitAction",
      oneofGroup: "option_type",
    },
  ],
  oneofGroups: ["option_type"],
};


/**
 * Schema for SimulateAttackRequest message
 */
export const SimulateAttackRequestSchema: MessageSchema = {
  name: "SimulateAttackRequest",
  fields: [
    {
      name: "attackerUnitType",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "attackerTerrain",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "attackerHealth",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "defenderUnitType",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "defenderTerrain",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "defenderHealth",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "woundBonus",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "numSimulations",
      type: FieldType.NUMBER,
      id: 8,
    },
  ],
};


/**
 * Schema for SimulateAttackResponse message
 */
export const SimulateAttackResponseSchema: MessageSchema = {
  name: "SimulateAttackResponse",
  fields: [
    {
      name: "attackerDamageDistribution",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "defenderDamageDistribution",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "attackerMeanDamage",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "defenderMeanDamage",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "attackerKillProbability",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "defenderKillProbability",
      type: FieldType.NUMBER,
      id: 6,
    },
  ],
};


/**
 * Schema for SimulateFixRequest message
 */
export const SimulateFixRequestSchema: MessageSchema = {
  name: "SimulateFixRequest",
  fields: [
    {
      name: "fixingUnitType",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "fixingUnitHealth",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "injuredUnitType",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "numSimulations",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for SimulateFixResponse message
 */
export const SimulateFixResponseSchema: MessageSchema = {
  name: "SimulateFixResponse",
  fields: [
    {
      name: "healingDistribution",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "meanHealing",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "fixValue",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for JoinGameRequest message
 */
export const JoinGameRequestSchema: MessageSchema = {
  name: "JoinGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerId",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for JoinGameResponse message
 */
export const JoinGameResponseSchema: MessageSchema = {
  name: "JoinGameResponse",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
    {
      name: "playerId",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for EmptyRequest message
 */
export const EmptyRequestSchema: MessageSchema = {
  name: "EmptyRequest",
  fields: [
  ],
};


/**
 * Schema for EmptyResponse message
 */
export const EmptyResponseSchema: MessageSchema = {
  name: "EmptyResponse",
  fields: [
  ],
};


/**
 * Schema for SetContentRequest message
 */
export const SetContentRequestSchema: MessageSchema = {
  name: "SetContentRequest",
  fields: [
    {
      name: "innerHtml",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for SetContentResponse message
 */
export const SetContentResponseSchema: MessageSchema = {
  name: "SetContentResponse",
  fields: [
  ],
};


/**
 * Schema for ShowBuildOptionsRequest message
 */
export const ShowBuildOptionsRequestSchema: MessageSchema = {
  name: "ShowBuildOptionsRequest",
  fields: [
    {
      name: "innerHtml",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "hide",
      type: FieldType.BOOLEAN,
      id: 2,
    },
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for ShowBuildOptionsResponse message
 */
export const ShowBuildOptionsResponseSchema: MessageSchema = {
  name: "ShowBuildOptionsResponse",
  fields: [
  ],
};


/**
 * Schema for LogMessageRequest message
 */
export const LogMessageRequestSchema: MessageSchema = {
  name: "LogMessageRequest",
  fields: [
    {
      name: "message",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for LogMessageResponse message
 */
export const LogMessageResponseSchema: MessageSchema = {
  name: "LogMessageResponse",
  fields: [
  ],
};


/**
 * Schema for SetGameStateRequest message
 */
export const SetGameStateRequestSchema: MessageSchema = {
  name: "SetGameStateRequest",
  fields: [
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Game",
    },
    {
      name: "state",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameState",
    },
  ],
};


/**
 * Schema for SetGameStateResponse message
 */
export const SetGameStateResponseSchema: MessageSchema = {
  name: "SetGameStateResponse",
  fields: [
  ],
};


/**
 * Schema for UpdateGameStatusRequest message
 */
export const UpdateGameStatusRequestSchema: MessageSchema = {
  name: "UpdateGameStatusRequest",
  fields: [
    {
      name: "currentPlayer",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "turnCounter",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for UpdateGameStatusResponse message
 */
export const UpdateGameStatusResponseSchema: MessageSchema = {
  name: "UpdateGameStatusResponse",
  fields: [
  ],
};


/**
 * Schema for SetTileAtRequest message
 */
export const SetTileAtRequestSchema: MessageSchema = {
  name: "SetTileAtRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "tile",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.Tile",
    },
  ],
};


/**
 * Schema for SetTileAtResponse message
 */
export const SetTileAtResponseSchema: MessageSchema = {
  name: "SetTileAtResponse",
  fields: [
  ],
};


/**
 * Schema for SetUnitAtRequest message
 */
export const SetUnitAtRequestSchema: MessageSchema = {
  name: "SetUnitAtRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "unit",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "flash",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "appear",
      type: FieldType.BOOLEAN,
      id: 5,
    },
  ],
};


/**
 * Schema for SetUnitAtResponse message
 */
export const SetUnitAtResponseSchema: MessageSchema = {
  name: "SetUnitAtResponse",
  fields: [
  ],
};


/**
 * Schema for RemoveTileAtRequest message
 */
export const RemoveTileAtRequestSchema: MessageSchema = {
  name: "RemoveTileAtRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "animate",
      type: FieldType.BOOLEAN,
      id: 3,
    },
  ],
};


/**
 * Schema for RemoveTileAtResponse message
 */
export const RemoveTileAtResponseSchema: MessageSchema = {
  name: "RemoveTileAtResponse",
  fields: [
  ],
};


/**
 * Schema for RemoveUnitAtRequest message
 */
export const RemoveUnitAtRequestSchema: MessageSchema = {
  name: "RemoveUnitAtRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "animate",
      type: FieldType.BOOLEAN,
      id: 3,
    },
  ],
};


/**
 * Schema for RemoveUnitAtResponse message
 */
export const RemoveUnitAtResponseSchema: MessageSchema = {
  name: "RemoveUnitAtResponse",
  fields: [
  ],
};


/**
 * Schema for ShowHighlightsRequest message
 */
export const ShowHighlightsRequestSchema: MessageSchema = {
  name: "ShowHighlightsRequest",
  fields: [
    {
      name: "highlights",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.HighlightSpec",
      repeated: true,
    },
  ],
};


/**
 * Schema for ShowHighlightsResponse message
 */
export const ShowHighlightsResponseSchema: MessageSchema = {
  name: "ShowHighlightsResponse",
  fields: [
  ],
};


/**
 * Schema for HighlightSpec message
 */
export const HighlightSpecSchema: MessageSchema = {
  name: "HighlightSpec",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "type",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "move",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.MoveUnitAction",
      oneofGroup: "action",
    },
    {
      name: "attack",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.AttackUnitAction",
      oneofGroup: "action",
    },
    {
      name: "build",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.BuildUnitAction",
      oneofGroup: "action",
    },
    {
      name: "capture",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "lilbattle.v1.CaptureBuildingAction",
      oneofGroup: "action",
    },
    {
      name: "player",
      type: FieldType.NUMBER,
      id: 8,
    },
  ],
  oneofGroups: ["action"],
};


/**
 * Schema for ClearHighlightsRequest message
 */
export const ClearHighlightsRequestSchema: MessageSchema = {
  name: "ClearHighlightsRequest",
  fields: [
    {
      name: "types",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for ClearHighlightsResponse message
 */
export const ClearHighlightsResponseSchema: MessageSchema = {
  name: "ClearHighlightsResponse",
  fields: [
  ],
};


/**
 * Schema for ShowPathRequest message
 */
export const ShowPathRequestSchema: MessageSchema = {
  name: "ShowPathRequest",
  fields: [
    {
      name: "coords",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
    {
      name: "color",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "thickness",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for ShowPathResponse message
 */
export const ShowPathResponseSchema: MessageSchema = {
  name: "ShowPathResponse",
  fields: [
  ],
};


/**
 * Schema for ClearPathsRequest message
 */
export const ClearPathsRequestSchema: MessageSchema = {
  name: "ClearPathsRequest",
  fields: [
  ],
};


/**
 * Schema for ClearPathsResponse message
 */
export const ClearPathsResponseSchema: MessageSchema = {
  name: "ClearPathsResponse",
  fields: [
  ],
};


/**
 * Schema for MoveUnitRequest message
 */
export const MoveUnitRequestSchema: MessageSchema = {
  name: "MoveUnitRequest",
  fields: [
    {
      name: "unit",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Unit",
    },
    {
      name: "path",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.HexCoord",
      repeated: true,
    },
  ],
};


/**
 * Schema for MoveUnitResponse message
 */
export const MoveUnitResponseSchema: MessageSchema = {
  name: "MoveUnitResponse",
  fields: [
  ],
};


/**
 * Schema for HexCoord message
 */
export const HexCoordSchema: MessageSchema = {
  name: "HexCoord",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for ShowAttackEffectRequest message
 */
export const ShowAttackEffectRequestSchema: MessageSchema = {
  name: "ShowAttackEffectRequest",
  fields: [
    {
      name: "fromQ",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "fromR",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "toQ",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "toR",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "damage",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "splashTargets",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.SplashTarget",
      repeated: true,
    },
  ],
};


/**
 * Schema for SplashTarget message
 */
export const SplashTargetSchema: MessageSchema = {
  name: "SplashTarget",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "damage",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for ShowAttackEffectResponse message
 */
export const ShowAttackEffectResponseSchema: MessageSchema = {
  name: "ShowAttackEffectResponse",
  fields: [
  ],
};


/**
 * Schema for ShowHealEffectRequest message
 */
export const ShowHealEffectRequestSchema: MessageSchema = {
  name: "ShowHealEffectRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "amount",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for ShowHealEffectResponse message
 */
export const ShowHealEffectResponseSchema: MessageSchema = {
  name: "ShowHealEffectResponse",
  fields: [
  ],
};


/**
 * Schema for ShowCaptureEffectRequest message
 */
export const ShowCaptureEffectRequestSchema: MessageSchema = {
  name: "ShowCaptureEffectRequest",
  fields: [
    {
      name: "q",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "r",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for ShowCaptureEffectResponse message
 */
export const ShowCaptureEffectResponseSchema: MessageSchema = {
  name: "ShowCaptureEffectResponse",
  fields: [
  ],
};


/**
 * Schema for SetAllowedPanelsRequest message
 */
export const SetAllowedPanelsRequestSchema: MessageSchema = {
  name: "SetAllowedPanelsRequest",
  fields: [
    {
      name: "panelIds",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for SetAllowedPanelsResponse message
 */
export const SetAllowedPanelsResponseSchema: MessageSchema = {
  name: "SetAllowedPanelsResponse",
  fields: [
  ],
};


/**
 * Schema for IndexState message
 */
export const IndexStateSchema: MessageSchema = {
  name: "IndexState",
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
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "indexedAt",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "needsIndexing",
      type: FieldType.BOOLEAN,
      id: 7,
    },
    {
      name: "status",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "lastError",
      type: FieldType.STRING,
      id: 9,
    },
    {
      name: "idempotencyKey",
      type: FieldType.STRING,
      id: 10,
    },
    {
      name: "retryCount",
      type: FieldType.NUMBER,
      id: 11,
    },
  ],
};


/**
 * Schema for EnsureIndexStateRequest message
 */
export const EnsureIndexStateRequestSchema: MessageSchema = {
  name: "EnsureIndexStateRequest",
  fields: [
    {
      name: "indexState",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexState",
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.FieldMask",
    },
  ],
};


/**
 * Schema for EnsureIndexStateResponse message
 */
export const EnsureIndexStateResponseSchema: MessageSchema = {
  name: "EnsureIndexStateResponse",
  fields: [
    {
      name: "indexState",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexState",
    },
  ],
};


/**
 * Schema for GetIndexStatesRequest message
 */
export const GetIndexStatesRequestSchema: MessageSchema = {
  name: "GetIndexStatesRequest",
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
      name: "indexTypes",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
  ],
};


/**
 * Schema for IndexStateList message
 */
export const IndexStateListSchema: MessageSchema = {
  name: "IndexStateList",
  fields: [
    {
      name: "states",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexState",
      repeated: true,
    },
  ],
};


/**
 * Schema for GetIndexStatesResponse message
 */
export const GetIndexStatesResponseSchema: MessageSchema = {
  name: "GetIndexStatesResponse",
  fields: [
    {
      name: "states",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for ListIndexStatesRequest message
 */
export const ListIndexStatesRequestSchema: MessageSchema = {
  name: "ListIndexStatesRequest",
  fields: [
    {
      name: "entityType",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "updatedBefore",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.Timestamp",
      oneofGroup: "_updated_before",
      optional: true,
    },
    {
      name: "updatedAfter",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
      oneofGroup: "_updated_after",
      optional: true,
    },
    {
      name: "indexTypes",
      type: FieldType.REPEATED,
      id: 4,
      repeated: true,
    },
    {
      name: "orderBy",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "count",
      type: FieldType.NUMBER,
      id: 6,
    },
  ],
  oneofGroups: ["_updated_before", "_updated_after"],
};


/**
 * Schema for ListIndexStatesResponse message
 */
export const ListIndexStatesResponseSchema: MessageSchema = {
  name: "ListIndexStatesResponse",
  fields: [
    {
      name: "items",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexState",
      repeated: true,
    },
    {
      name: "nextPageKey",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for DeleteIndexStatesRequest message
 */
export const DeleteIndexStatesRequestSchema: MessageSchema = {
  name: "DeleteIndexStatesRequest",
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
      name: "indexTypes",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
  ],
};


/**
 * Schema for DeleteIndexStatesResponse message
 */
export const DeleteIndexStatesResponseSchema: MessageSchema = {
  name: "DeleteIndexStatesResponse",
  fields: [
  ],
};


/**
 * Schema for IndexRecord message
 */
export const IndexRecordSchema: MessageSchema = {
  name: "IndexRecord",
  fields: [
    {
      name: "entityId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "entityData",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Any",
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
 * Schema for IndexRecordsLRO message
 */
export const IndexRecordsLROSchema: MessageSchema = {
  name: "IndexRecordsLRO",
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
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Timestamp",
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
      messageType: "lilbattle.v1.IndexRecord",
      repeated: true,
    },
  ],
};


/**
 * Schema for CreateIndexRecordsLRORequest message
 */
export const CreateIndexRecordsLRORequestSchema: MessageSchema = {
  name: "CreateIndexRecordsLRORequest",
  fields: [
    {
      name: "lro",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexRecordsLRO",
    },
  ],
};


/**
 * Schema for CreateIndexRecordsLROResponse message
 */
export const CreateIndexRecordsLROResponseSchema: MessageSchema = {
  name: "CreateIndexRecordsLROResponse",
  fields: [
    {
      name: "lro",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexRecordsLRO",
    },
  ],
};


/**
 * Schema for UpdateIndexRecordsLRORequest message
 */
export const UpdateIndexRecordsLRORequestSchema: MessageSchema = {
  name: "UpdateIndexRecordsLRORequest",
  fields: [
    {
      name: "lro",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexRecordsLRO",
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.FieldMask",
    },
  ],
};


/**
 * Schema for UpdateIndexRecordsLROResponse message
 */
export const UpdateIndexRecordsLROResponseSchema: MessageSchema = {
  name: "UpdateIndexRecordsLROResponse",
  fields: [
    {
      name: "lro",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexRecordsLRO",
    },
  ],
};


/**
 * Schema for GetIndexRecordsLRORequest message
 */
export const GetIndexRecordsLRORequestSchema: MessageSchema = {
  name: "GetIndexRecordsLRORequest",
  fields: [
    {
      name: "lroId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for GetIndexRecordsLROResponse message
 */
export const GetIndexRecordsLROResponseSchema: MessageSchema = {
  name: "GetIndexRecordsLROResponse",
  fields: [
    {
      name: "lro",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.IndexRecordsLRO",
    },
  ],
};


/**
 * Schema for Job message
 */
export const JobSchema: MessageSchema = {
  name: "Job",
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
      name: "jobType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "jobData",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "google.protobuf.Any",
    },
    {
      name: "debounceWindowSeconds",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "repeatInfo",
      type: FieldType.MESSAGE,
      id: 8,
      messageType: "lilbattle.v1.RepeatInfo",
    },
  ],
};


/**
 * Schema for RepeatInfo message
 */
export const RepeatInfoSchema: MessageSchema = {
  name: "RepeatInfo",
  fields: [
  ],
};


/**
 * Schema for Run message
 */
export const RunSchema: MessageSchema = {
  name: "Run",
  fields: [
    {
      name: "jobId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "runId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "startedAt",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "updatedAt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Timestamp",
    },
    {
      name: "state",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "runData",
      type: FieldType.MESSAGE,
      id: 7,
      messageType: "google.protobuf.Any",
    },
    {
      name: "lastError",
      type: FieldType.STRING,
      id: 8,
    },
    {
      name: "lastContentHash",
      type: FieldType.STRING,
      id: 9,
    },
    {
      name: "retryCount",
      type: FieldType.NUMBER,
      id: 10,
    },
  ],
};


/**
 * Schema for InitializeSingletonRequest message
 */
export const InitializeSingletonRequestSchema: MessageSchema = {
  name: "InitializeSingletonRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "gameData",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "gameState",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "moveHistory",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "viewerUserId",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for InitializeSingletonResponse message
 */
export const InitializeSingletonResponseSchema: MessageSchema = {
  name: "InitializeSingletonResponse",
  fields: [
    {
      name: "response",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.InitializeGameResponse",
    },
  ],
};


/**
 * Schema for TurnOptionClickedRequest message
 */
export const TurnOptionClickedRequestSchema: MessageSchema = {
  name: "TurnOptionClickedRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "optionIndex",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "optionType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.Position",
    },
  ],
};


/**
 * Schema for TurnOptionClickedResponse message
 */
export const TurnOptionClickedResponseSchema: MessageSchema = {
  name: "TurnOptionClickedResponse",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for SceneClickedRequest message
 */
export const SceneClickedRequestSchema: MessageSchema = {
  name: "SceneClickedRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "layer",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for SceneClickedResponse message
 */
export const SceneClickedResponseSchema: MessageSchema = {
  name: "SceneClickedResponse",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for EndTurnButtonClickedRequest message
 */
export const EndTurnButtonClickedRequestSchema: MessageSchema = {
  name: "EndTurnButtonClickedRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for EndTurnButtonClickedResponse message
 */
export const EndTurnButtonClickedResponseSchema: MessageSchema = {
  name: "EndTurnButtonClickedResponse",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for BuildOptionClickedRequest message
 */
export const BuildOptionClickedRequestSchema: MessageSchema = {
  name: "BuildOptionClickedRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pos",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.Position",
    },
    {
      name: "unitType",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for BuildOptionClickedResponse message
 */
export const BuildOptionClickedResponseSchema: MessageSchema = {
  name: "BuildOptionClickedResponse",
  fields: [
  ],
};


/**
 * Schema for InitializeGameRequest message
 */
export const InitializeGameRequestSchema: MessageSchema = {
  name: "InitializeGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for InitializeGameResponse message
 */
export const InitializeGameResponseSchema: MessageSchema = {
  name: "InitializeGameResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "error",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "currentPlayer",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "turnCounter",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "gameName",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for ClientReadyRequest message
 */
export const ClientReadyRequestSchema: MessageSchema = {
  name: "ClientReadyRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for ClientReadyResponse message
 */
export const ClientReadyResponseSchema: MessageSchema = {
  name: "ClientReadyResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
  ],
};


/**
 * Schema for ApplyRemoteChangesRequest message
 */
export const ApplyRemoteChangesRequestSchema: MessageSchema = {
  name: "ApplyRemoteChangesRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "moves",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameMove",
      repeated: true,
    },
  ],
};


/**
 * Schema for ApplyRemoteChangesResponse message
 */
export const ApplyRemoteChangesResponseSchema: MessageSchema = {
  name: "ApplyRemoteChangesResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "error",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "requiresReload",
      type: FieldType.BOOLEAN,
      id: 3,
    },
  ],
};


/**
 * Schema for SubscribeRequest message
 */
export const SubscribeRequestSchema: MessageSchema = {
  name: "SubscribeRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "fromSequence",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for SubscribeResponse message
 */
export const SubscribeResponseSchema: MessageSchema = {
  name: "SubscribeResponse",
  fields: [
    {
      name: "currentSequence",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "gameState",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameState",
    },
    {
      name: "game",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.Game",
    },
  ],
};


/**
 * Schema for GameUpdate message
 */
export const GameUpdateSchema: MessageSchema = {
  name: "GameUpdate",
  fields: [
    {
      name: "sequence",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "movesPublished",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.MovesPublished",
      oneofGroup: "update_type",
    },
    {
      name: "playerJoined",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "lilbattle.v1.PlayerJoined",
      oneofGroup: "update_type",
    },
    {
      name: "playerLeft",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "lilbattle.v1.PlayerLeft",
      oneofGroup: "update_type",
    },
    {
      name: "gameEnded",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "lilbattle.v1.GameEnded",
      oneofGroup: "update_type",
    },
    {
      name: "initialState",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "lilbattle.v1.SubscribeResponse",
      oneofGroup: "update_type",
    },
  ],
  oneofGroups: ["update_type"],
};


/**
 * Schema for MovesPublished message
 */
export const MovesPublishedSchema: MessageSchema = {
  name: "MovesPublished",
  fields: [
    {
      name: "player",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "moves",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameMove",
      repeated: true,
    },
    {
      name: "groupNumber",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for PlayerJoined message
 */
export const PlayerJoinedSchema: MessageSchema = {
  name: "PlayerJoined",
  fields: [
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerNumber",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for PlayerLeft message
 */
export const PlayerLeftSchema: MessageSchema = {
  name: "PlayerLeft",
  fields: [
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerNumber",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for GameEnded message
 */
export const GameEndedSchema: MessageSchema = {
  name: "GameEnded",
  fields: [
    {
      name: "winner",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "reason",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for BroadcastRequest message
 */
export const BroadcastRequestSchema: MessageSchema = {
  name: "BroadcastRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "update",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.GameUpdate",
    },
  ],
};


/**
 * Schema for BroadcastResponse message
 */
export const BroadcastResponseSchema: MessageSchema = {
  name: "BroadcastResponse",
  fields: [
    {
      name: "subscriberCount",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "sequence",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for ThemeInfo message
 */
export const ThemeInfoSchema: MessageSchema = {
  name: "ThemeInfo",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "version",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "basePath",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "assetType",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "needsPostProcessing",
      type: FieldType.BOOLEAN,
      id: 5,
    },
  ],
};


/**
 * Schema for UnitMapping message
 */
export const UnitMappingSchema: MessageSchema = {
  name: "UnitMapping",
  fields: [
    {
      name: "old",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "image",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for TerrainMapping message
 */
export const TerrainMappingSchema: MessageSchema = {
  name: "TerrainMapping",
  fields: [
    {
      name: "old",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "image",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for ThemeManifest message
 */
export const ThemeManifestSchema: MessageSchema = {
  name: "ThemeManifest",
  fields: [
    {
      name: "themeInfo",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.ThemeInfo",
    },
    {
      name: "units",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "terrains",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "playerColors",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for PlayerColor message
 */
export const PlayerColorSchema: MessageSchema = {
  name: "PlayerColor",
  fields: [
    {
      name: "primary",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "secondary",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for AssetResult message
 */
export const AssetResultSchema: MessageSchema = {
  name: "AssetResult",
  fields: [
    {
      name: "type",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "data",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for WorldInfo message
 */
export const WorldInfoSchema: MessageSchema = {
  name: "WorldInfo",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "category",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "difficulty",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 6,
      repeated: true,
    },
    {
      name: "icon",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "lastUpdated",
      type: FieldType.STRING,
      id: 8,
    },
  ],
};


/**
 * Schema for ListWorldsRequest message
 */
export const ListWorldsRequestSchema: MessageSchema = {
  name: "ListWorldsRequest",
  fields: [
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.Pagination",
    },
    {
      name: "ownerId",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for ListWorldsResponse message
 */
export const ListWorldsResponseSchema: MessageSchema = {
  name: "ListWorldsResponse",
  fields: [
    {
      name: "items",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
      repeated: true,
    },
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.PaginationResponse",
    },
  ],
};


/**
 * Schema for GetWorldRequest message
 */
export const GetWorldRequestSchema: MessageSchema = {
  name: "GetWorldRequest",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "version",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for GetWorldResponse message
 */
export const GetWorldResponseSchema: MessageSchema = {
  name: "GetWorldResponse",
  fields: [
    {
      name: "world",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldData",
    },
  ],
};


/**
 * Schema for UpdateWorldRequest message
 */
export const UpdateWorldRequestSchema: MessageSchema = {
  name: "UpdateWorldRequest",
  fields: [
    {
      name: "world",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldData",
    },
    {
      name: "clearWorld",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "google.protobuf.FieldMask",
    },
  ],
};


/**
 * Schema for UpdateWorldResponse message
 */
export const UpdateWorldResponseSchema: MessageSchema = {
  name: "UpdateWorldResponse",
  fields: [
    {
      name: "world",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldData",
    },
  ],
};


/**
 * Schema for DeleteWorldRequest message
 */
export const DeleteWorldRequestSchema: MessageSchema = {
  name: "DeleteWorldRequest",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for DeleteWorldResponse message
 */
export const DeleteWorldResponseSchema: MessageSchema = {
  name: "DeleteWorldResponse",
  fields: [
  ],
};


/**
 * Schema for GetWorldsRequest message
 */
export const GetWorldsRequestSchema: MessageSchema = {
  name: "GetWorldsRequest",
  fields: [
    {
      name: "ids",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for GetWorldsResponse message
 */
export const GetWorldsResponseSchema: MessageSchema = {
  name: "GetWorldsResponse",
  fields: [
    {
      name: "worlds",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for CreateWorldRequest message
 */
export const CreateWorldRequestSchema: MessageSchema = {
  name: "CreateWorldRequest",
  fields: [
    {
      name: "world",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldData",
    },
  ],
};


/**
 * Schema for CreateWorldResponse message
 */
export const CreateWorldResponseSchema: MessageSchema = {
  name: "CreateWorldResponse",
  fields: [
    {
      name: "world",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "lilbattle.v1.World",
    },
    {
      name: "worldData",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "lilbattle.v1.WorldData",
    },
    {
      name: "fieldErrors",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};



/**
 * Package-scoped schema registry for lilbattle.v1
 */
export const lilbattle_v1SchemaRegistry: Record<string, MessageSchema> = {
  "lilbattle.v1.IndexInfo": IndexInfoSchema,
  "lilbattle.v1.Pagination": PaginationSchema,
  "lilbattle.v1.PaginationResponse": PaginationResponseSchema,
  "lilbattle.v1.World": WorldSchema,
  "lilbattle.v1.WorldData": WorldDataSchema,
  "lilbattle.v1.Crossing": CrossingSchema,
  "lilbattle.v1.Tile": TileSchema,
  "lilbattle.v1.Unit": UnitSchema,
  "lilbattle.v1.AttackRecord": AttackRecordSchema,
  "lilbattle.v1.TerrainDefinition": TerrainDefinitionSchema,
  "lilbattle.v1.UnitDefinition": UnitDefinitionSchema,
  "lilbattle.v1.TerrainUnitProperties": TerrainUnitPropertiesSchema,
  "lilbattle.v1.UnitUnitProperties": UnitUnitPropertiesSchema,
  "lilbattle.v1.DamageDistribution": DamageDistributionSchema,
  "lilbattle.v1.DamageRange": DamageRangeSchema,
  "lilbattle.v1.RulesEngine": RulesEngineSchema,
  "lilbattle.v1.Game": GameSchema,
  "lilbattle.v1.GameConfiguration": GameConfigurationSchema,
  "lilbattle.v1.IncomeConfig": IncomeConfigSchema,
  "lilbattle.v1.GamePlayer": GamePlayerSchema,
  "lilbattle.v1.GameTeam": GameTeamSchema,
  "lilbattle.v1.GameSettings": GameSettingsSchema,
  "lilbattle.v1.PlayerState": PlayerStateSchema,
  "lilbattle.v1.GameState": GameStateSchema,
  "lilbattle.v1.GameMoveHistory": GameMoveHistorySchema,
  "lilbattle.v1.GameMoveGroup": GameMoveGroupSchema,
  "lilbattle.v1.GameMove": GameMoveSchema,
  "lilbattle.v1.Position": PositionSchema,
  "lilbattle.v1.MoveUnitAction": MoveUnitActionSchema,
  "lilbattle.v1.AttackUnitAction": AttackUnitActionSchema,
  "lilbattle.v1.BuildUnitAction": BuildUnitActionSchema,
  "lilbattle.v1.CaptureBuildingAction": CaptureBuildingActionSchema,
  "lilbattle.v1.EndTurnAction": EndTurnActionSchema,
  "lilbattle.v1.HealUnitAction": HealUnitActionSchema,
  "lilbattle.v1.FixUnitAction": FixUnitActionSchema,
  "lilbattle.v1.WorldChange": WorldChangeSchema,
  "lilbattle.v1.UnitHealedChange": UnitHealedChangeSchema,
  "lilbattle.v1.UnitFixedChange": UnitFixedChangeSchema,
  "lilbattle.v1.UnitMovedChange": UnitMovedChangeSchema,
  "lilbattle.v1.UnitDamagedChange": UnitDamagedChangeSchema,
  "lilbattle.v1.UnitKilledChange": UnitKilledChangeSchema,
  "lilbattle.v1.PlayerChangedChange": PlayerChangedChangeSchema,
  "lilbattle.v1.UnitBuiltChange": UnitBuiltChangeSchema,
  "lilbattle.v1.CoinsChangedChange": CoinsChangedChangeSchema,
  "lilbattle.v1.TileCapturedChange": TileCapturedChangeSchema,
  "lilbattle.v1.CaptureStartedChange": CaptureStartedChangeSchema,
  "lilbattle.v1.AllPaths": AllPathsSchema,
  "lilbattle.v1.PathEdge": PathEdgeSchema,
  "lilbattle.v1.Path": PathSchema,
  "lilbattle.v1.File": FileSchema,
  "lilbattle.v1.PutFileRequest": PutFileRequestSchema,
  "lilbattle.v1.PutFileResponse": PutFileResponseSchema,
  "lilbattle.v1.GetFileRequest": GetFileRequestSchema,
  "lilbattle.v1.GetFileResponse": GetFileResponseSchema,
  "lilbattle.v1.DeleteFileRequest": DeleteFileRequestSchema,
  "lilbattle.v1.DeleteFileResponse": DeleteFileResponseSchema,
  "lilbattle.v1.ListFilesRequest": ListFilesRequestSchema,
  "lilbattle.v1.ListFilesResponse": ListFilesResponseSchema,
  "lilbattle.v1.ListGamesRequest": ListGamesRequestSchema,
  "lilbattle.v1.ListGamesResponse": ListGamesResponseSchema,
  "lilbattle.v1.GetGameRequest": GetGameRequestSchema,
  "lilbattle.v1.GetGameResponse": GetGameResponseSchema,
  "lilbattle.v1.GetGameContentRequest": GetGameContentRequestSchema,
  "lilbattle.v1.GetGameContentResponse": GetGameContentResponseSchema,
  "lilbattle.v1.UpdateGameRequest": UpdateGameRequestSchema,
  "lilbattle.v1.UpdateGameResponse": UpdateGameResponseSchema,
  "lilbattle.v1.DeleteGameRequest": DeleteGameRequestSchema,
  "lilbattle.v1.DeleteGameResponse": DeleteGameResponseSchema,
  "lilbattle.v1.GetGamesRequest": GetGamesRequestSchema,
  "lilbattle.v1.GetGamesResponse": GetGamesResponseSchema,
  "lilbattle.v1.CreateGameRequest": CreateGameRequestSchema,
  "lilbattle.v1.CreateGameResponse": CreateGameResponseSchema,
  "lilbattle.v1.ProcessMovesRequest": ProcessMovesRequestSchema,
  "lilbattle.v1.ProcessMovesResponse": ProcessMovesResponseSchema,
  "lilbattle.v1.GetGameStateRequest": GetGameStateRequestSchema,
  "lilbattle.v1.GetGameStateResponse": GetGameStateResponseSchema,
  "lilbattle.v1.ListMovesRequest": ListMovesRequestSchema,
  "lilbattle.v1.ListMovesResponse": ListMovesResponseSchema,
  "lilbattle.v1.GetOptionsAtRequest": GetOptionsAtRequestSchema,
  "lilbattle.v1.GetOptionsAtResponse": GetOptionsAtResponseSchema,
  "lilbattle.v1.GameOption": GameOptionSchema,
  "lilbattle.v1.SimulateAttackRequest": SimulateAttackRequestSchema,
  "lilbattle.v1.SimulateAttackResponse": SimulateAttackResponseSchema,
  "lilbattle.v1.SimulateFixRequest": SimulateFixRequestSchema,
  "lilbattle.v1.SimulateFixResponse": SimulateFixResponseSchema,
  "lilbattle.v1.JoinGameRequest": JoinGameRequestSchema,
  "lilbattle.v1.JoinGameResponse": JoinGameResponseSchema,
  "lilbattle.v1.EmptyRequest": EmptyRequestSchema,
  "lilbattle.v1.EmptyResponse": EmptyResponseSchema,
  "lilbattle.v1.SetContentRequest": SetContentRequestSchema,
  "lilbattle.v1.SetContentResponse": SetContentResponseSchema,
  "lilbattle.v1.ShowBuildOptionsRequest": ShowBuildOptionsRequestSchema,
  "lilbattle.v1.ShowBuildOptionsResponse": ShowBuildOptionsResponseSchema,
  "lilbattle.v1.LogMessageRequest": LogMessageRequestSchema,
  "lilbattle.v1.LogMessageResponse": LogMessageResponseSchema,
  "lilbattle.v1.SetGameStateRequest": SetGameStateRequestSchema,
  "lilbattle.v1.SetGameStateResponse": SetGameStateResponseSchema,
  "lilbattle.v1.UpdateGameStatusRequest": UpdateGameStatusRequestSchema,
  "lilbattle.v1.UpdateGameStatusResponse": UpdateGameStatusResponseSchema,
  "lilbattle.v1.SetTileAtRequest": SetTileAtRequestSchema,
  "lilbattle.v1.SetTileAtResponse": SetTileAtResponseSchema,
  "lilbattle.v1.SetUnitAtRequest": SetUnitAtRequestSchema,
  "lilbattle.v1.SetUnitAtResponse": SetUnitAtResponseSchema,
  "lilbattle.v1.RemoveTileAtRequest": RemoveTileAtRequestSchema,
  "lilbattle.v1.RemoveTileAtResponse": RemoveTileAtResponseSchema,
  "lilbattle.v1.RemoveUnitAtRequest": RemoveUnitAtRequestSchema,
  "lilbattle.v1.RemoveUnitAtResponse": RemoveUnitAtResponseSchema,
  "lilbattle.v1.ShowHighlightsRequest": ShowHighlightsRequestSchema,
  "lilbattle.v1.ShowHighlightsResponse": ShowHighlightsResponseSchema,
  "lilbattle.v1.HighlightSpec": HighlightSpecSchema,
  "lilbattle.v1.ClearHighlightsRequest": ClearHighlightsRequestSchema,
  "lilbattle.v1.ClearHighlightsResponse": ClearHighlightsResponseSchema,
  "lilbattle.v1.ShowPathRequest": ShowPathRequestSchema,
  "lilbattle.v1.ShowPathResponse": ShowPathResponseSchema,
  "lilbattle.v1.ClearPathsRequest": ClearPathsRequestSchema,
  "lilbattle.v1.ClearPathsResponse": ClearPathsResponseSchema,
  "lilbattle.v1.MoveUnitRequest": MoveUnitRequestSchema,
  "lilbattle.v1.MoveUnitResponse": MoveUnitResponseSchema,
  "lilbattle.v1.HexCoord": HexCoordSchema,
  "lilbattle.v1.ShowAttackEffectRequest": ShowAttackEffectRequestSchema,
  "lilbattle.v1.SplashTarget": SplashTargetSchema,
  "lilbattle.v1.ShowAttackEffectResponse": ShowAttackEffectResponseSchema,
  "lilbattle.v1.ShowHealEffectRequest": ShowHealEffectRequestSchema,
  "lilbattle.v1.ShowHealEffectResponse": ShowHealEffectResponseSchema,
  "lilbattle.v1.ShowCaptureEffectRequest": ShowCaptureEffectRequestSchema,
  "lilbattle.v1.ShowCaptureEffectResponse": ShowCaptureEffectResponseSchema,
  "lilbattle.v1.SetAllowedPanelsRequest": SetAllowedPanelsRequestSchema,
  "lilbattle.v1.SetAllowedPanelsResponse": SetAllowedPanelsResponseSchema,
  "lilbattle.v1.IndexState": IndexStateSchema,
  "lilbattle.v1.EnsureIndexStateRequest": EnsureIndexStateRequestSchema,
  "lilbattle.v1.EnsureIndexStateResponse": EnsureIndexStateResponseSchema,
  "lilbattle.v1.GetIndexStatesRequest": GetIndexStatesRequestSchema,
  "lilbattle.v1.IndexStateList": IndexStateListSchema,
  "lilbattle.v1.GetIndexStatesResponse": GetIndexStatesResponseSchema,
  "lilbattle.v1.ListIndexStatesRequest": ListIndexStatesRequestSchema,
  "lilbattle.v1.ListIndexStatesResponse": ListIndexStatesResponseSchema,
  "lilbattle.v1.DeleteIndexStatesRequest": DeleteIndexStatesRequestSchema,
  "lilbattle.v1.DeleteIndexStatesResponse": DeleteIndexStatesResponseSchema,
  "lilbattle.v1.IndexRecord": IndexRecordSchema,
  "lilbattle.v1.IndexRecordsLRO": IndexRecordsLROSchema,
  "lilbattle.v1.CreateIndexRecordsLRORequest": CreateIndexRecordsLRORequestSchema,
  "lilbattle.v1.CreateIndexRecordsLROResponse": CreateIndexRecordsLROResponseSchema,
  "lilbattle.v1.UpdateIndexRecordsLRORequest": UpdateIndexRecordsLRORequestSchema,
  "lilbattle.v1.UpdateIndexRecordsLROResponse": UpdateIndexRecordsLROResponseSchema,
  "lilbattle.v1.GetIndexRecordsLRORequest": GetIndexRecordsLRORequestSchema,
  "lilbattle.v1.GetIndexRecordsLROResponse": GetIndexRecordsLROResponseSchema,
  "lilbattle.v1.Job": JobSchema,
  "lilbattle.v1.RepeatInfo": RepeatInfoSchema,
  "lilbattle.v1.Run": RunSchema,
  "lilbattle.v1.InitializeSingletonRequest": InitializeSingletonRequestSchema,
  "lilbattle.v1.InitializeSingletonResponse": InitializeSingletonResponseSchema,
  "lilbattle.v1.TurnOptionClickedRequest": TurnOptionClickedRequestSchema,
  "lilbattle.v1.TurnOptionClickedResponse": TurnOptionClickedResponseSchema,
  "lilbattle.v1.SceneClickedRequest": SceneClickedRequestSchema,
  "lilbattle.v1.SceneClickedResponse": SceneClickedResponseSchema,
  "lilbattle.v1.EndTurnButtonClickedRequest": EndTurnButtonClickedRequestSchema,
  "lilbattle.v1.EndTurnButtonClickedResponse": EndTurnButtonClickedResponseSchema,
  "lilbattle.v1.BuildOptionClickedRequest": BuildOptionClickedRequestSchema,
  "lilbattle.v1.BuildOptionClickedResponse": BuildOptionClickedResponseSchema,
  "lilbattle.v1.InitializeGameRequest": InitializeGameRequestSchema,
  "lilbattle.v1.InitializeGameResponse": InitializeGameResponseSchema,
  "lilbattle.v1.ClientReadyRequest": ClientReadyRequestSchema,
  "lilbattle.v1.ClientReadyResponse": ClientReadyResponseSchema,
  "lilbattle.v1.ApplyRemoteChangesRequest": ApplyRemoteChangesRequestSchema,
  "lilbattle.v1.ApplyRemoteChangesResponse": ApplyRemoteChangesResponseSchema,
  "lilbattle.v1.SubscribeRequest": SubscribeRequestSchema,
  "lilbattle.v1.SubscribeResponse": SubscribeResponseSchema,
  "lilbattle.v1.GameUpdate": GameUpdateSchema,
  "lilbattle.v1.MovesPublished": MovesPublishedSchema,
  "lilbattle.v1.PlayerJoined": PlayerJoinedSchema,
  "lilbattle.v1.PlayerLeft": PlayerLeftSchema,
  "lilbattle.v1.GameEnded": GameEndedSchema,
  "lilbattle.v1.BroadcastRequest": BroadcastRequestSchema,
  "lilbattle.v1.BroadcastResponse": BroadcastResponseSchema,
  "lilbattle.v1.ThemeInfo": ThemeInfoSchema,
  "lilbattle.v1.UnitMapping": UnitMappingSchema,
  "lilbattle.v1.TerrainMapping": TerrainMappingSchema,
  "lilbattle.v1.ThemeManifest": ThemeManifestSchema,
  "lilbattle.v1.PlayerColor": PlayerColorSchema,
  "lilbattle.v1.AssetResult": AssetResultSchema,
  "lilbattle.v1.WorldInfo": WorldInfoSchema,
  "lilbattle.v1.ListWorldsRequest": ListWorldsRequestSchema,
  "lilbattle.v1.ListWorldsResponse": ListWorldsResponseSchema,
  "lilbattle.v1.GetWorldRequest": GetWorldRequestSchema,
  "lilbattle.v1.GetWorldResponse": GetWorldResponseSchema,
  "lilbattle.v1.UpdateWorldRequest": UpdateWorldRequestSchema,
  "lilbattle.v1.UpdateWorldResponse": UpdateWorldResponseSchema,
  "lilbattle.v1.DeleteWorldRequest": DeleteWorldRequestSchema,
  "lilbattle.v1.DeleteWorldResponse": DeleteWorldResponseSchema,
  "lilbattle.v1.GetWorldsRequest": GetWorldsRequestSchema,
  "lilbattle.v1.GetWorldsResponse": GetWorldsResponseSchema,
  "lilbattle.v1.CreateWorldRequest": CreateWorldRequestSchema,
  "lilbattle.v1.CreateWorldResponse": CreateWorldResponseSchema,
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