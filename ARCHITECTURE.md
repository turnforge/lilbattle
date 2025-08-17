# WeeWar Architecture

## System Overview

WeeWar implements a clean separation between game logic (Go), state management (WASM), and presentation (TypeScript/Phaser). The architecture emphasizes transaction safety, type safety, and maintainable separation of concerns.

## Core Design Principles

### 1. Transaction-Safe State Management

**Copy-on-Write Semantics**: The World system implements parent-child transaction layers where child layers create copies of parent objects before modification, preventing corruption of parent state.

```go
// Transaction layer creates copies to avoid parent mutation
if w.parent != nil {
    if _, existsInCurrentLayer := w.unitsByCoord[currentCoord]; !existsInCurrentLayer {
        // Unit comes from parent layer - make a copy
        unitToMove = &v1.Unit{
            Q: unit.Q, R: unit.R, Player: unit.Player,
            // ... copy all fields
        }
    }
}
```

**Rollback Safety**: ProcessMoves creates transaction snapshots, processes moves, then rolls back to original state for ordered change application:

```go
originalWorld := rtGame.World
rtGame.World = originalWorld.Push()  // Create transaction
// Process moves on transaction layer
rtGame.World = originalWorld        // Rollback for ordered apply
```

### 2. Service Layer Abstraction

**Transport Independence**: Same service implementations work across HTTP, gRPC, and WASM through interface abstraction:

```go
type GamesServiceImpl interface {
    v1.GamesServiceServer
    GetRuntimeGame(game *v1.Game, gameState *v1.GameState) (*Game, error)
}
```

**Base Implementation**: `BaseGamesServiceImpl` provides core logic that concrete implementations extend for specific transports.

### 3. Type-Safe WASM Integration

**Generated Client**: WASM client provides compile-time type checking with protobuf integration:

```typescript
// Type-safe API calls
const response = await this.wasmService.ProcessMoves({
    gameId: this.gameId,
    moves: [moveAction]
});
```

**Direct Property Access**: Simplified protobuf handling with direct field access (`change.unitMoved`) instead of complex oneof handling.

## Component Architecture

### Game Engine Layer (`lib/`)

**World**: Pure game state container implementing hex coordinate system with transaction support.
- Immutable parent-child relationships for transactions
- Efficient merged iteration across transaction layers
- Copy-on-write semantics for state safety

**Game**: Runtime game logic integrating World with Rules Engine.
- Current player/turn state management
- Move validation and processing coordination
- Rules engine integration for game mechanics

**Move Processor**: Validates and processes moves with full transaction support.
- Transaction-aware move validation
- Change result generation for state updates
- Error handling and rollback coordination

### Service Layer (`services/`)

**ProcessMoves Pipeline**:
1. Create transaction snapshot of game state
2. Process moves on transaction layer
3. Generate change results from transaction
4. Rollback to original state
5. Apply changes in ordered sequence
6. Update persistent state

**Change Application**: `ApplyChangeResults` ensures ordered, atomic application of move results to maintain state consistency.

### Frontend Layer (`web/`)

**GameState Controller**: Lightweight wrapper managing WASM service interactions.
- Move execution coordination
- State synchronization with server
- Event emission for UI updates

**GameViewer Renderer**: Phaser-based hex map rendering with unit display.
- Event-driven updates from GameState
- Hex coordinate conversion for display
- User interaction handling (clicks, selections)

## Transaction System Deep Dive

### Problem Solved

**Unit Duplication Bug**: Units appearing at both old and new positions after moves due to shared object references between transaction and parent layers.

### Solution Architecture

**Copy-on-Write in MoveUnit**:
```go
func (w *World) MoveUnit(unit *v1.Unit, newCoord AxialCoord) error {
    unitToMove := unit
    if w.parent != nil {
        // Check if unit comes from parent layer
        if _, existsInCurrentLayer := w.unitsByCoord[currentCoord]; !existsInCurrentLayer {
            // Make copy to avoid modifying parent objects
            unitToMove = &v1.Unit{/* copy all fields */}
        }
    }
    // Safe to modify copy
    UnitSetCoord(unitToMove, newCoord)
}
```

**Transaction Counter Optimization**:
```go
func (w *World) NumUnits() int32 {
    if w.parent != nil {
        return w.parent.NumUnits() + w.unitsAdded - w.unitsDeleted
    }
    return int32(len(w.unitsByCoord))
}
```

### Benefits

- **Data Integrity**: Parent layers remain immutable during transaction processing
- **Performance**: Efficient counting without expensive iteration
- **Rollback Safety**: Clean rollback to known good state
- **Test Coverage**: Comprehensive validation of transaction semantics

## Event System Architecture

### Clean Separation

**Game Logic → State → UI**: Unidirectional data flow prevents circular dependencies.

**Event Bus Pattern**: Loose coupling between components through event emission:
```typescript
// GameState emits events
this.emit('world-updated', updatedWorld);

// GameViewer subscribes to events  
gameState.on('world-updated', (world) => this.updateScene(world));
```

### State Synchronization

**Server as Source of Truth**: All state changes validated server-side with client updates.

**Optimistic UI**: Client shows immediate feedback while server processes moves.

**Conflict Resolution**: Server state takes precedence in case of discrepancies.

## Testing Architecture

### Comprehensive Coverage

**Unit Tests**: World operations with/without transactions
- Basic move operations
- Unit replacement scenarios  
- Transaction layer isolation
- Copy-on-write semantics

**Integration Tests**: End-to-end ProcessMoves pipeline
- Real WasmGamesService usage
- Transaction flow validation
- State consistency verification

**Transaction Flow Tests**: Exact simulation of ProcessMoves behavior
- Transaction creation and rollback
- Change application ordering
- Unit object sharing prevention

## Performance Considerations

### Efficient Operations

**Transaction Counters**: O(1) unit counting instead of O(n) iteration
**Copy-on-Write**: Only copies objects when actually modified
**Merged Iteration**: Lazy evaluation of parent/child object combination

### Memory Management

**Shallow Copies**: Unit objects use shallow copying for performance
**Transaction Cleanup**: Automatic cleanup when transactions complete
**Object Reuse**: Minimize allocation/deallocation in hot paths

## Future Architecture Considerations

### Scalability

**Multi-Player Support**: Transaction system ready for concurrent player actions
**State Partitioning**: World architecture supports regional game state management
**Caching Layer**: Service layer ready for redis/memcache integration

### Extension Points

**Rules Engine**: JSON-configurable game mechanics
**Transport Layer**: Easy addition of new client protocols
**Rendering Backend**: Phaser abstraction allows for alternative renderers

This architecture provides a solid foundation for turn-based strategy games with excellent separation of concerns, comprehensive testing, and robust state management.