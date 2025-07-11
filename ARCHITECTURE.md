# WeeWar Game Architecture

## Overview

The WeeWar game demonstrates a sophisticated, layered architecture that has evolved through multiple iterations. Starting with a TurnEngine framework approach, it has evolved into a unified game implementation with comprehensive interface definitions, extensive testing, and multiple frontend interfaces (CLI, PNG rendering, web). This architecture showcases modern game development patterns with strong separation of concerns, comprehensive testing, and data-driven design.

## Current Architecture (2024)

```
┌─────────────────────────────────────────────────────────────┐
│                   WeeWar Core Game System                   │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐│
│ │  Game Interface │  │  Map Interface  │  │  Unit Interface ││
│ │                 │  │                 │  │                 ││
│ │• GameController │  │• MapInterface   │  │• UnitInterface  ││
│ │• Turn Management│  │• Coordinate Ops │  │• Unit Actions   ││
│ │• Save/Load      │  │• Pathfinding    │  │• Combat System  ││
│ │• Game State     │  │• Hex Navigation │  │• Unit Queries   ││
│ └─────────────────┘  └─────────────────┘  └─────────────────┘│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Unified Game Struct                      │
├─────────────────────────────────────────────────────────────┤
│ • Implements all interfaces (GameInterface)                 │
│ • Comprehensive state management                            │
│ • Integrated combat, movement, and map systems             │
│ • Real WeeWar data integration                             │
│ • PNG rendering with sophisticated graphics                │
│ • JSON save/load with full state persistence              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Multiple Interfaces                      │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐│
│ │  CLI Interface  │  │  PNG Renderer   │  │  Web Interface  ││
│ │                 │  │                 │  │                 ││
│ │• REPL Loop      │  │• Hex Rendering  │  │• HTTP Server    ││
│ │• Chess Notation │  │• Unit Graphics  │  │• JSON API       ││
│ │• Interactive    │  │• Map Rendering  │  │• Browser UI     ││
│ │• Batch Mode     │  │• Auto-render    │  │• Real-time      ││
│ └─────────────────┘  └─────────────────┘  └─────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Core Architecture Components

### 1. Game Interface System
**Location**: `game_interface.go`
- **Purpose**: Defines comprehensive contracts for all game operations
- **Design**: Interface segregation with focused responsibilities
- **Reusability**: 100% - Clean contracts enable multiple implementations

```go
// Core interface combining all game operations
type GameInterface interface {
    GameController  // Game lifecycle, turns, state
    MapInterface    // Map queries, pathfinding, coordinates
    UnitInterface   // Unit actions, queries, management
}
```

**Key Design Decisions:**
- **Interface Segregation**: Separate interfaces for different concerns
- **Comprehensive Coverage**: All game operations have defined contracts
- **Implementation Flexibility**: Interfaces enable multiple implementations
- **Testing Support**: Interfaces enable easy mocking and testing

### 2. Unified Game Implementation
**Location**: `game.go`
- **Purpose**: Single struct implementing all game interfaces
- **Design**: Comprehensive state management with integrated systems
- **Evolution**: Replaced fragmented ECS approach with unified design

```go
type Game struct {
    // Core game state
    PlayerCount   int
    CurrentPlayer int
    TurnNumber    int
    GameStatus    GameStatus
    
    // Game systems
    gameMap     *Map
    units       []*Unit
    
    // Integrated systems
    rng         *rand.Rand
    pathfinder  *HexPathfinder
}
```

**Key Design Decisions:**
- **Single Source of Truth**: All game state in one place
- **Integrated Systems**: Combat, movement, pathfinding built-in
- **Performance**: Direct access without ECS overhead
- **Simplicity**: Easier to understand and maintain

### 3. Hexagonal Map System
**Location**: `map.go`, `tile.go`
- **Purpose**: Sophisticated hex-based map with pathfinding
- **Design**: Hex coordinate system with neighbor connectivity
- **Data Integration**: Real WeeWar terrain and movement costs

```go
type Map struct {
    Width, Height int
    IsHex         bool
    tiles         [][]*Tile
    
    // Hex-specific pathfinding
    pathfinder *HexPathfinder
}

type Tile struct {
    Row, Col     int
    TileType     int
    Unit         *Unit
    Neighbors    [6]*Tile  // Hex neighbors
}
```

**Key Design Decisions:**
- **Hex-First Design**: Built specifically for hexagonal grids
- **Neighbor Connectivity**: Pre-computed neighbors for efficiency
- **Integrated Pathfinding**: A* pathfinding built into map system
- **Real Data**: Terrain types match original WeeWar

### 4. Unit and Combat System
**Location**: `unit.go`, `combat.go`
- **Purpose**: Comprehensive unit management and combat resolution
- **Design**: Data-driven combat with real WeeWar mechanics
- **Data Source**: Extracted from original WeeWar game

```go
type Unit struct {
    Row, Col          int
    PlayerID          int
    UnitType          int
    AvailableHealth   int
    DistanceLeft      int
    HasAttacked       bool
}

type CombatResult struct {
    AttackerDamage, DefenderDamage int
    AttackerKilled, DefenderKilled bool
    AttackerHealth, DefenderHealth int
}
```

**Key Design Decisions:**
- **Data-Driven Combat**: Real damage matrices from WeeWar
- **Comprehensive State**: All unit state tracked precisely
- **Combat Results**: Detailed result information for interfaces
- **Movement Integration**: Movement and combat state linked

### 5. Multiple Interface Support

#### CLI Interface (`cli_impl.go`, `cli_formatter.go`)
- **REPL Loop**: Interactive command-line gameplay
- **Chess Notation**: A1, B2, C3 position system
- **Rich Formatting**: Colors, tables, status displays
- **Multiple Modes**: Interactive, batch, single commands

```go
type WeeWarCLI struct {
    game         *Game
    displayMode  CLIDisplayMode
    formatter    CLIFormatter
    interactive  bool
}
```

#### PNG Renderer (`rendering.go`, `buffer.go`)
- **Hex Rendering**: Sophisticated hexagonal grid rendering
- **Visual Elements**: Terrain, units, health bars, borders
- **Auto-Scaling**: Responsive to different map sizes
- **High Quality**: Anti-aliased graphics with rich visuals

```go
type Buffer struct {
    Width, Height int
    img           *image.RGBA
    gc            *draw2dimg.GraphicContext
}
```

#### Web Interface (Future)
- **HTTP Server**: RESTful API for game operations
- **JSON Protocol**: Structured data exchange
- **Real-time Updates**: WebSocket support for live games
- **Browser UI**: Modern web interface

### 6. Data Integration System
**Location**: `weewar_data.go`, `cmd/extract-data/`
- **Purpose**: Integration with real WeeWar game data
- **Design**: Extracted and parsed from HTML source
- **Authenticity**: Ensures accurate game mechanics

```go
type WeeWarData struct {
    Units    map[string]UnitData
    Terrain  map[string]TerrainData
    Combat   map[string]map[string]DamageDistribution
}
```

**Key Design Decisions:**
- **Real Data**: Extracted from tinyattack.com HTML
- **Comprehensive**: All unit stats, terrain, combat data
- **Validation**: Ensures calculations match original game
- **Extensibility**: Easy to add new units and mechanics

## Testing Architecture

### 1. Comprehensive Test Suite
**Location**: `*_test.go` files
- **Unit Tests**: Individual component testing
- **Integration Tests**: Full game scenario testing
- **Interface Tests**: Contract compliance testing
- **Performance Tests**: Load and stress testing

```go
// Test categories
func TestGameBasicOperations(t *testing.T)     // Core game functions
func TestCombatSystem(t *testing.T)            // Combat mechanics
func TestMapNavigation(t *testing.T)           // Map and pathfinding
func TestCLIInterface(t *testing.T)            // CLI functionality
func TestPNGRendering(t *testing.T)            // Visual rendering
func TestSaveLoad(t *testing.T)                // Persistence
```

### 2. Test Data Management
- **Deterministic**: Fixed seeds for reproducible tests
- **Comprehensive**: Tests cover all game scenarios
- **Validation**: Tests verify against known correct results
- **Coverage**: High test coverage across all systems

## Key Design Principles

### 1. Interface-Driven Design
- **Clean Contracts**: Well-defined interfaces for all operations
- **Implementation Flexibility**: Multiple implementations possible
- **Testing Support**: Interfaces enable comprehensive testing
- **Future Evolution**: Easy to add new interfaces and implementations

### 2. Unified State Management
- **Single Source of Truth**: All game state in Game struct
- **Consistency**: No state synchronization issues
- **Performance**: Direct access without indirection
- **Simplicity**: Easy to understand and debug

### 3. Data-Driven Authenticity
- **Real Data**: Extracted from original WeeWar game
- **Accuracy**: Calculations match original game exactly
- **Validation**: Comprehensive testing against known results
- **Extensibility**: Easy to add new data and mechanics

### 4. Multiple Interface Support
- **CLI**: Interactive command-line interface
- **PNG**: Visual rendering for analysis and debugging
- **Web**: HTTP API for browser-based gameplay
- **Batch**: Automated testing and AI development

### 5. Comprehensive Testing
- **Unit Tests**: Individual component validation
- **Integration Tests**: Full game scenario testing
- **Interface Tests**: Contract compliance verification
- **Performance Tests**: Load and stress testing

## Evolution and Learnings

### 1. Architecture Evolution
- **Started**: Complex ECS framework approach
- **Evolved**: Unified game implementation
- **Learned**: Simplicity often beats complexity
- **Result**: Cleaner, faster, more maintainable code

### 2. Interface Design
- **Started**: Monolithic interfaces
- **Evolved**: Segregated, focused interfaces
- **Learned**: Interface segregation principle crucial
- **Result**: Clean contracts enabling multiple implementations

### 3. Testing Strategy
- **Started**: Basic unit tests
- **Evolved**: Comprehensive test suite
- **Learned**: Testing game logic requires careful design
- **Result**: High confidence in game correctness

### 4. Data Integration
- **Started**: Hardcoded game data
- **Evolved**: Real data extraction and integration
- **Learned**: Authenticity requires real data
- **Result**: Accurate WeeWar game mechanics

## Performance Characteristics

### 1. Game Operations
- **Turn Processing**: O(1) - Direct state access
- **Pathfinding**: O(V log V) - A* with efficient heuristics
- **Combat Resolution**: O(1) - Direct lookup in damage matrices
- **State Persistence**: O(n) - Linear in game state size

### 2. Rendering Performance
- **PNG Generation**: O(n) - Linear in map size
- **Memory Usage**: Efficient buffer management
- **Image Quality**: High-quality anti-aliased graphics
- **Scalability**: Responsive to different map sizes

### 3. CLI Performance
- **Command Processing**: O(1) - Direct command dispatch
- **Display Updates**: O(n) - Linear in visible elements
- **Interactive Response**: Sub-millisecond command processing
- **Memory Usage**: Minimal overhead for CLI operations

## Future Directions

### 1. Planned Enhancements
- **AI Integration**: Support for AI players
- **Multiplayer**: Network multiplayer support
- **Map Editor**: Visual map creation tools
- **Statistics**: Comprehensive game statistics

### 2. Architecture Improvements
- **Plugin System**: Support for game modifications
- **Event System**: Decoupled event handling
- **Scripting**: Lua or similar for game customization
- **Optimization**: Further performance improvements

### 3. Interface Expansion
- **Mobile**: Native mobile app interfaces
- **VR**: Virtual reality game experience
- **API**: Comprehensive REST API
- **Streaming**: Game streaming and spectator modes

## Summary

The WeeWar implementation demonstrates a mature, well-architected game system that balances complexity and simplicity. The interface-driven design enables multiple implementations while maintaining clean separation of concerns. The unified game implementation provides performance and simplicity while the comprehensive testing ensures correctness and reliability. The data-driven approach ensures authenticity to the original game while enabling future extensions and improvements.

The architecture successfully supports multiple interfaces (CLI, PNG, web) while maintaining a single source of truth for game state. The evolution from a complex ECS framework to a unified implementation demonstrates the value of simplicity and pragmatism in software design.