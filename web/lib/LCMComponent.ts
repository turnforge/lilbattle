/**
 * LCMComponent Interface - Defines multi-phase component initialization
 * 
 * This interface implements a breadth-first lifecycle pattern that eliminates
 * initialization order dependencies and race conditions through synchronization barriers.
 * 
 * Lifecycle Phases:
 * 1. Construction - Components are created but not initialized
 * 2. bindToDOM() - Basic DOM setup, discover child components  
 * 3. injectDependencies() - Receive references to other components
 * 4. activate() - Final setup when all dependencies are ready
 * 
 * Key Benefits:
 * - Order Independence: Components can be created in any sequence
 * - Async Safety: Each phase waits for all components before proceeding
 * - Clear Dependencies: Explicit injection points prevent race conditions
 * - Error Isolation: Component failures don't cascade to others
 */

/**
 * LCM Component - Short for LifeCycle Managed Component 
 * enable components to be declared first and their loading be managed by a
 * LifecycleController so that we have layered creation, dependency injection and setup
 * in a breadth first way.
 *
 * A key constraint on LCMComponents are that they should not perform any initialization
 * in the constructor.  This is because a LifecycleController should be used to load/setup
 * these components and they will follow a layered approach.  Performing these actions in
 * the constructor could violate the idempotency guarantees.
 */
export interface LCMComponent {
    /**
     * Phase 1: The "local" initialization of the component.
     *
     * In this phase the component initializes itself and returns any children it might
     * want initialized as part of the lifecycled loading.
     * 
     * This phase should:
     * - Set up basic DOM elements and event listeners
     * - Create child components (but don't initialize them) and return them
     * - Return array of child components for lifecycle controller discovery
     * 
     * This phase must be synchronous and should not:
     * - Access other components or external dependencies
     * - Perform async operations
     * - Emit events or notifications
     * 
     * @returns Array of child components to be managed by lifecycle controller
     */
    performLocalInit(): Promise<LCMComponent[]> | LCMComponent[];
    
    /**
     * Phase 2: Inject dependencies from parent/siblings
     * 
     * This phase should:
     * - Receive and store references to required dependencies
     * - Validate that required dependencies are provided
     * - Set up internal state based on dependencies
     * 
     * This phase can be async and may:
     * - Load external data or resources
     * - Perform validation or setup operations
     * - Initialize internal components that depend on injected references
     * 
     * @param deps Record of dependency name to dependency instance
     * @returns Promise<void> or void - can be async
     */
    setupDependencies(): Promise<void> | void;
    
    /**
     * Phase 3: Activate component when all dependencies are ready
     * 
     * This phase should:
     * - Complete final initialization
     * - Enable component functionality
     * - Start listening for external events
     * - Begin normal operation
     * 
     * This phase can be async and may:
     * - Connect to external services
     * - Load initial data
     * - Emit ready notifications
     * 
     * @returns Promise<void> or void - can be async
     */
    activate(): Promise<void> | void;
    
    /**
     * Cleanup phase: Deactivate component and clean up resources
     * 
     * This should:
     * - Stop all ongoing operations
     * - Remove event listeners
     * - Clean up external connections
     * - Dispose of child components
     * 
     * @returns Promise<void> or void - can be async
     */
    deactivate(): Promise<void> | void;
}

/**
 * Configuration for component lifecycle behavior
 */
export interface LCMComponentConfig {
    /**
     * Maximum time to wait for a lifecycle phase to complete (ms)
     * Default: 10000 (10 seconds)
     */
    phaseTimeoutMs?: number;
    
    /**
     * Whether to continue if individual components fail during a phase
     * Default: false (fail fast)
     */
    continueOnError?: boolean;
    
    /**
     * Whether to validate dependencies against declared requirements
     * Default: true
     */
    validateDependencies?: boolean;
    
    /**
     * Enable debug logging for lifecycle phases
     * Default: false
     */
    enableDebugLogging?: boolean;
}

/**
 * Event emitted during component lifecycle transitions
 */
export interface LCMComponentEvent {
    type: 'phase-start' | 'phase-complete' | 'phase-error' | 'component-ready';
    componentName: string;
    timestamp: number;
    error?: Error;
    metadata?: Record<string, any>;
}
