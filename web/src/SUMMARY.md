**Purpose:**

This folder contains the core client-side TypeScript logic for the webapp, managing UI state, user events, API interactions, and DOM manipulation using a modern component-based architecture with strict separation of concerns and event-driven communication.

**Core Architecture Components:**

## Modern Component System (New)

*   **`WorldViewer.ts`**: Phaser-based world visualization component with proper DOM scoping and event-driven initialization  
*   **`WorldStatsPanel.ts`**: Statistics display component with safe DOM selectors and event-driven updates
*   **`WorldDetailsPage.ts`**: Orchestrator page following new architecture - handles data loading and component coordination only

## Key Architecture Principles

*   **Separation of Concerns**: Clear boundaries between layout, behavior, and communication responsibilities
*   **Event-Driven**: Components communicate through EventBus events, never direct method calls  
*   **DOM Isolation**: Components only access DOM within their assigned root elements
*   **Error Resilience**: Component failures are isolated and don't affect other components
*   **Timing Awareness**: Proper handling of initialization order, race conditions, and async operations
*   **WebGL Integration**: Specialized patterns for graphics libraries like Phaser with timing considerations

## Critical Timing Patterns Learned

*   **TypeScript Field Initializers**: Avoid explicit `= null` for constructor-set fields
*   **Event Subscription Order**: Subscribe to events BEFORE creating components that emit them
*   **WebGL Context Readiness**: Use small setTimeout for graphics library initialization completion
*   **State → Subscribe → Create**: Strict three-phase initialization order
*   **Async in Handlers**: EventBus stays synchronous, handlers use `.then()/.catch()` for async operations

## Integration Capabilities

*   **Phaser.js**: WebGL-based world rendering with proper timing handling
*   **HTMX**: Component hydration support for server-driven UI updates  
*   **Canvas/WebGL**: Specialized initialization patterns for graphics contexts
*   **Toast/Modal Systems**: User feedback and interaction patterns
*   **Theme Management**: Coordinated theming across component boundaries

## Recent Session Work (2025-01-24)

### Layer System Architecture Complete ✅
*   **Generic WorldViewer**: `WorldViewer<TScene>` with template parameter for proper typing
*   **GameViewer Specialization**: `GameViewer extends WorldViewer<PhaserGameScene>` with game-specific layer access
*   **Layer-Based Interaction**: Direct layer manipulation (`getSelectionHighlightLayer()`, `getMovementHighlightLayer()`, etc.)
*   **Editor Integration**: PhaserEditorComponent uses layer callbacks for painting logic
*   **Callback Architecture**: Click handling through BaseMapLayer callbacks with validation in components
*   **Brush Size Support**: Multi-tile painting with hex distance calculations in component layer

### Architecture Improvements ✅
*   **Scene Separation**: PhaserWorldScene for rendering, components for business logic
*   **Single Source of Truth**: World model updates trigger observer pattern for visual updates
*   **Type Safety**: Proper TypeScript generics eliminate casting and improve developer experience
*   **Clean Separation**: UI logic in components, rendering logic in scenes, interaction through layers

## Recent Session Work (2025-01-22)

### Interactive Game Viewer Foundation ✅
*   **GameViewerPage Architecture**: Complete interactive game interface with lifecycle controller integration
*   **External Orchestration Pattern**: LifecycleController with breadth-first component initialization eliminates race conditions
*   **LCMComponent Interface**: Multi-phase initialization (performLocalInit, setupDependencies, activate, deactivate)
*   **WASM Bridge Architecture**: GameState component with async loading and synchronous gameplay operations
*   **Synchronous UI Pattern**: Immediate UI feedback with notification events for coordination only

### Component Communication Architecture ✅  
*   **Event-Driven Coordination**: Components communicate via EventBus without tight coupling
*   **Source Filtering**: Components ignore events they originate to prevent feedback loops
*   **Error Isolation**: Component failures don't cascade through event system
*   **Debug Support**: Comprehensive logging and lifecycle event callbacks
*   **Notification Events**: System coordination (`game-created`, `unit-moved`, `turn-ended`) for logging, animations

### Previous Session Work (2025-01-20)

#### Component Architecture Cleanup ✅
*   **WorldEditorPage Streamlining**: Removed dead code and consolidated component management patterns
*   **Panel Integration Optimization**: Improved coordination between EditorToolsPanel, TileStatsPanel, and PhaserEditor
*   **Reference Management**: Cleaner component initialization and lifecycle patterns
*   **State Management Consolidation**: Reduced complexity in page-level state handling
