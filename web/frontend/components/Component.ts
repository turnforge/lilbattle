import { EventBus, EventHandler } from './EventBus';

/**
 * Standard component state interface
 */
export interface ComponentState {
    isInitialized: boolean;
    isReady: boolean;
    hasError: boolean;
    errorMessage?: string;
    lastUpdated: number;
}

/**
 * Component initialization configuration
 */
export interface ComponentConfig {
    componentId: string;
    rootElement: HTMLElement;
    eventBus: EventBus;
    debugMode?: boolean;
}

/**
 * Base interface for all UI components
 * Enforces separation of concerns and standard lifecycle
 */
export interface Component {
    /**
     * Unique identifier for this component instance
     */
    readonly componentId: string;
    
    /**
     * Root DOM element that this component owns and manages
     */
    readonly rootElement: HTMLElement;
    
    /**
     * Initialize the component with its root element and dependencies
     * @param config - Component configuration including root element and event bus
     * @returns Promise that resolves to true if initialization succeeded
     */
    initialize(config: ComponentConfig): Promise<boolean>;
    
    /**
     * Clean up the component and release resources
     * Should unsubscribe from events, clean up DOM, and release memory
     */
    destroy(): void;
    
    /**
     * Get current component state
     * @returns Current state of the component
     */
    getState(): ComponentState;
    
    /**
     * Check if component is ready for use
     * @returns True if component is initialized and ready
     */
    isReady(): boolean;
}

/**
 * Abstract base class implementing common component functionality
 * Provides standard lifecycle management and event bus integration
 */
export abstract class BaseComponent implements Component {
    public readonly componentId: string;
    public readonly rootElement: HTMLElement;
    
    protected eventBus: EventBus;
    protected debugMode: boolean = false;
    protected state: ComponentState;
    protected eventUnsubscribers: (() => void)[] = [];
    
    constructor() {
        // These will be set during initialize()
        this.componentId = '';
        this.rootElement = null as any;
        this.eventBus = null as any;
        
        this.state = {
            isInitialized: false,
            isReady: false,
            hasError: false,
            lastUpdated: Date.now()
        };
    }
    
    public async initialize(config: ComponentConfig): Promise<boolean> {
        try {
            // Set up component properties
            (this as any).componentId = config.componentId;
            (this as any).rootElement = config.rootElement;
            this.eventBus = config.eventBus;
            this.debugMode = config.debugMode || false;
            
            this.log('Initializing component...');
            
            // Validate root element
            if (!this.rootElement) {
                throw new Error('Root element is required');
            }
            
            // Mark as component in DOM for debugging
            this.rootElement.setAttribute('data-component', this.componentId);
            
            // Call component-specific initialization
            const success = await this.initializeComponent();
            
            if (success) {
                this.state.isInitialized = true;
                this.state.isReady = true;
                this.state.lastUpdated = Date.now();
                
                // Emit initialization event
                this.eventBus.emit('component-initialized', {
                    componentId: this.componentId,
                    success: true
                }, this.componentId);
                
                this.log('Component initialized successfully');
                return true;
            } else {
                throw new Error('Component-specific initialization failed');
            }
            
        } catch (error) {
            this.handleError('Initialization failed', error);
            return false;
        }
    }
    
    public destroy(): void {
        this.log('Destroying component...');
        
        try {
            // Unsubscribe from all events
            this.eventUnsubscribers.forEach(unsubscribe => unsubscribe());
            this.eventUnsubscribers = [];
            
            // Call component-specific cleanup
            this.destroyComponent();
            
            // Remove component marker from DOM
            this.rootElement?.removeAttribute('data-component');
            
            // Reset state
            this.state.isInitialized = false;
            this.state.isReady = false;
            this.state.lastUpdated = Date.now();
            
            this.log('Component destroyed successfully');
            
        } catch (error) {
            console.error(`[${this.componentId}] Error during destroy:`, error);
        }
    }
    
    public getState(): ComponentState {
        return { ...this.state };
    }
    
    public isReady(): boolean {
        return this.state.isReady && !this.state.hasError;
    }
    
    /**
     * Subscribe to an event with automatic cleanup on destroy
     */
    protected subscribe<T = any>(eventType: string, handler: EventHandler<T>): void {
        const unsubscribe = this.eventBus.subscribe(eventType, handler, this.componentId);
        this.eventUnsubscribers.push(unsubscribe);
    }
    
    /**
     * Emit an event from this component
     */
    protected emit<T = any>(eventType: string, data: T): void {
        this.eventBus.emit(eventType, data, this.componentId);
    }
    
    /**
     * Find elements within this component's root element only
     * Enforces separation of concerns - no cross-component DOM access
     */
    protected findElement<T extends HTMLElement = HTMLElement>(selector: string): T | null {
        return this.rootElement.querySelector<T>(selector);
    }
    
    /**
     * Find multiple elements within this component's root element only
     */
    protected findElements<T extends HTMLElement = HTMLElement>(selector: string): T[] {
        return Array.from(this.rootElement.querySelectorAll<T>(selector));
    }
    
    /**
     * Handle component errors consistently
     */
    protected handleError(message: string, error: any): void {
        this.state.hasError = true;
        this.state.errorMessage = message;
        this.state.lastUpdated = Date.now();
        
        console.error(`[${this.componentId}] ${message}:`, error);
        
        // Emit error event for parent components to handle
        this.eventBus.emit('component-error', {
            componentId: this.componentId,
            error: message,
            details: error
        }, this.componentId);
    }
    
    /**
     * Log messages with component identification
     */
    protected log(message: string): void {
        if (this.debugMode) {
            console.log(`[${this.componentId}] ${message}`);
        }
    }
    
    // Abstract methods that components must implement
    
    /**
     * Component-specific initialization logic
     * Called after base initialization is complete
     */
    protected abstract initializeComponent(): Promise<boolean>;
    
    /**
     * Component-specific cleanup logic
     * Called during destroy before base cleanup
     */
    protected abstract destroyComponent(): void;
}

/**
 * Utility function to find component root elements by data attribute
 */
export function findComponentRoots(container: HTMLElement, componentType: string): HTMLElement[] {
    return Array.from(container.querySelectorAll(`[data-component-type="${componentType}"]`));
}

/**
 * Utility function to validate component root element
 */
export function validateComponentRoot(element: HTMLElement, expectedType?: string): boolean {
    if (!element) return false;
    
    if (expectedType) {
        const actualType = element.getAttribute('data-component-type');
        return actualType === expectedType;
    }
    
    return true;
}