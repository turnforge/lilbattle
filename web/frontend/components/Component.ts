import { EventBus, EventHandler, EventTypes } from './EventBus';

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
 * DOM validation result for hydration scenarios
 */
export interface DOMValidation {
    isValid: boolean;
    missingElements: string[];
    invalidElements: string[];
    warnings: string[];
    existingData?: any;
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
     * Handle dynamic content updates (e.g., from HTMX or server responses)
     * @param newHTML - New HTML content to replace current content
     */
    contentUpdated(newHTML: string): void;
    
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
    
    constructor(componentId: string, rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        this.componentId = componentId;
        this.rootElement = rootElement;
        this.eventBus = eventBus;
        this.debugMode = debugMode;
        
        this.state = {
            isInitialized: false,
            isReady: false,
            hasError: false,
            lastUpdated: Date.now()
        };
        
        // Mark as component in DOM for debugging
        this.rootElement.setAttribute('data-component', this.componentId);
        
        // Initialize the component
        try {
            this.initializeComponent();
            this.bindToDOM();
            
            this.state.isInitialized = true;
            this.state.isReady = true;
            this.state.lastUpdated = Date.now();
            
            // Emit initialization event
            this.eventBus.emit(EventTypes.COMPONENT_INITIALIZED, {
                componentId: this.componentId,
                success: true
            }, this.componentId);
            
            this.log('Component initialized successfully');
            
        } catch (error) {
            this.handleError('Component initialization failed', error);
        }
    }
    
    public contentUpdated(newHTML: string): void {
        try {
            this.log('Content updated, re-binding to DOM');
            
            // Update the DOM
            this.rootElement.innerHTML = newHTML;
            
            // Re-bind to the new DOM structure
            this.bindToDOM();
            
            this.state.lastUpdated = Date.now();
            
        } catch (error) {
            this.handleError('Content update failed', error);
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
     * Called once during construction to set up the component
     */
    protected abstract initializeComponent(): void;
    
    /**
     * Bind to DOM elements (handles both empty and pre-populated root elements)
     * Called during initialization and after content updates
     */
    protected abstract bindToDOM(): void;
    
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
