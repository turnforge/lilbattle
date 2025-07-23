import { BaseComponent, DOMValidation } from './Component';
import { EventBus, EventPayload, EventTypes, WorldDataLoadedPayload } from './EventBus';
import { PhaserWorldScene } from './phaser/PhaserWorldScene';
import { PhaserGameScene } from './phaser/PhaserGameScene';
import { Unit, Tile, World } from './World';
import { ComponentLifecycle } from './ComponentLifecycle';

/**
 * WorldViewer Component - Manages Phaser-based world visualization
 * Responsible for:
 * - Phaser initialization and lifecycle management
 * - World data rendering (tiles and units)
 * - Camera controls and viewport management
 * - Theme and display options
 * - Scene selection (PhaserWorldScene for basic viewing, PhaserGameScene for interactive gameplay)
 * 
 * Layout and styling are handled by parent container and CSS classes.
 */
export class WorldViewer extends BaseComponent implements ComponentLifecycle {
    private scene: PhaserWorldScene | PhaserGameScene | null = null;
    private loadedWorldData: WorldDataLoadedPayload | null;
    private viewerContainer: HTMLElement | null;
    private isGameMode: boolean = false; // New: determines which scene type to use
    
    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false, isGameMode: boolean = false) {
        console.log('WorldViewer constructor: received eventBus:', eventBus, 'gameMode:', isGameMode);
        super('world-viewer', rootElement, eventBus, debugMode);
        this.isGameMode = isGameMode;
    }
    
    protected initializeComponent(): void {
        this.log('Initializing WorldViewer component');
        
        // Subscribe to world data events
        this.subscribe<WorldDataLoadedPayload>(EventTypes.WORLD_DATA_LOADED, (payload) => {
            this.handleWorldDataLoaded(payload);
        });
        
        this.log('WorldViewer component initialized');
    }
    
    protected bindToDOM(): void {
        this.log('Binding WorldViewer to DOM');
        
        // Find the phaser-viewer-container within the root element
        let phaserContainer = this.rootElement.querySelector('#phaser-viewer-container') as HTMLElement;
        
        if (!phaserContainer) {
            // If not found as child, check if root element IS the phaser container
            if (this.rootElement.id === 'phaser-viewer-container') {
                phaserContainer = this.rootElement;
            } else {
                // Create the phaser container as a child
                console.warn('phaser-viewer-container not found, creating one');
                phaserContainer = document.createElement('div');
                phaserContainer.id = 'phaser-viewer-container';
                phaserContainer.className = 'w-full h-full min-h-96';
                this.rootElement.appendChild(phaserContainer);
            }
        }
        
        this.viewerContainer = phaserContainer;
        
        // Ensure the container has the right classes
        if (!this.viewerContainer.classList.contains('w-full')) {
            this.viewerContainer.className = 'w-full h-full min-h-96';
        }
        
        this.log('WorldViewer bound to DOM, container:', this.viewerContainer);
        
        // Phaser initialization will happen in activate() phase, not here
        console.log('WorldViewer: DOM binding complete, waiting for activate() phase');
    }
    
    protected destroyComponent(): void {
        this.log('Destroying WorldViewer component');
        
        // Clean up Phaser scene (it manages its own game instance)
        if (this.scene) {
            this.scene.destroy();
            this.scene = null;
        }
        
        this.loadedWorldData = null;
        this.viewerContainer = null;
    }
    
    
    public validateDOM(rootElement: HTMLElement): DOMValidation {
        const validation: DOMValidation = {
            isValid: true,
            missingElements: [],
            invalidElements: [],
            warnings: []
        };
        
        // Check for Phaser container
        const phaserContainer = rootElement.querySelector('#phaser-viewer-container');
        if (!phaserContainer) {
            validation.isValid = false;
            validation.missingElements.push('phaser-viewer-container');
        }
        
        return validation;
    }
    
    /**
     * Initialize the appropriate Phaser scene (PhaserWorldScene or PhaserGameScene)
     */
    private async initializePhaserScene(): Promise<void> {
        console.log(`WorldViewer: initializePhaserScene() called, gameMode: ${this.isGameMode}`);
        
        // Guard against multiple initialization
        if (this.scene) {
            console.log('WorldViewer: Phaser scene already initialized, skipping');
            return;
        }
        
        if (!this.viewerContainer) {
            throw new Error('Viewer container not available');
        }
        
        // Create the appropriate scene type based on mode
        if (this.isGameMode) {
            this.log('Creating self-contained PhaserGameScene for interactive gameplay');
            this.scene = new PhaserGameScene();
            console.log('[WorldViewer] PhaserGameScene created:', this.scene);
        } else {
            this.log('Creating self-contained PhaserWorldScene for basic viewing');
            this.scene = new PhaserWorldScene();
            console.log('[WorldViewer] PhaserWorldScene created:', this.scene);
        }
        
        // Initialize it with the container
        await this.scene.initialize(this.viewerContainer.id);
        
        this.log(`${this.isGameMode ? 'PhaserGameScene' : 'PhaserWorldScene'} initialized successfully`);
        
        // Emit ready event
        console.log('WorldViewer: Emitting WORLD_VIEWER_READY event');
        this.emit(EventTypes.WORLD_VIEWER_READY, {
            componentId: this.componentId,
            success: true
        });
        console.log('WorldViewer: WORLD_VIEWER_READY event emitted');
        
        // Load world data if we have it
        if (this.loadedWorldData) {
            await this.loadWorldIntoScene();
        }
    }
    
    /**
     * Handle world data loaded event
     */
    private handleWorldDataLoaded(payload: EventPayload<WorldDataLoadedPayload>): void {
        this.log(`Received world data for world: ${payload.data.worldId}`);
        this.loadedWorldData = payload.data;
        
        // Load into Phaser if scene is ready
        if (this.scene && this.scene.getIsInitialized()) {
            this.loadWorldIntoScene();
        }
    }
    
    /**
     * Load world data into the PhaserWorldScene
     */
    private async loadWorldIntoScene(): Promise<void> {
        if (!this.scene || !this.scene.getIsInitialized()) {
            this.log('Phaser scene not ready, deferring world load');
            return;
        }
        
        if (!this.loadedWorldData) {
            this.log('No world data available to load');
            return;
        }
        
        this.log('Loading world data into Phaser scene');
        
        // This method will be called after loadWorld() sets up the world data properly
        // For now, we need to reconstruct the World from loadedWorldData
        // TODO: This is a bit awkward - we should refactor to avoid this conversion
        console.log('WorldViewer: loadWorldIntoScene called but needs world instance');
    }
    
    /**
     * Public API for loading world data
     */
    public async loadWorld(worldData: any): Promise<void> {
        if (!worldData) {
            throw new Error('No world data provided');
        }
        
        this.log('Loading world data');
        
        // Process world data
        const world = World.deserialize(worldData);
        const allTiles = world.getAllTiles();
        const allUnits = world.getAllUnits();
        
        // Calculate bounds and stats
        const bounds = world.getBounds();
        
        // Store world data
        this.loadedWorldData = {
            worldId: worldData.id || 'unknown',
            totalTiles: allTiles.length,
            totalUnits: allUnits.length,
            bounds: bounds ? {
                minQ: bounds.minQ,
                maxQ: bounds.maxQ,
                minR: bounds.minR,
                maxR: bounds.maxR
            } : { minQ: 0, maxQ: 0, minR: 0, maxR: 0 },
            terrainCounts: this.calculateTerrainCounts(allTiles)
        };
        
        // Load into Phaser if ready
        if (this.scene && this.scene.getIsInitialized()) {
            await this.scene.loadWorldData(world);
        }
        
        // Emit data loaded event for other components
        this.emit(EventTypes.WORLD_DATA_LOADED, this.loadedWorldData);
    }
    
    /**
     * Calculate terrain counts for statistics
     */
    private calculateTerrainCounts(tiles: any[]): { [terrainType: number]: number } {
        const counts: { [terrainType: number]: number } = {};
        
        tiles.forEach(tile => {
            counts[tile.tileType] = (counts[tile.tileType] || 0) + 1;
        });
        
        return counts;
    }
    
    /**
     * Set display options
     */
    public setShowGrid(show: boolean): void {
        if (this.scene) {
            this.scene.setShowGrid(show);
        }
    }
    
    public setShowCoordinates(show: boolean): void {
        if (this.scene) {
            this.scene.setShowCoordinates(show);
        }
    }
    
    public setTheme(isDark: boolean): void {
        if (this.scene) {
            this.scene.setTheme(isDark);
        }
    }
    
    /**
     * Camera controls
     */
    public getZoom(): number {
        return this.scene?.getZoom() || 1;
    }
    
    public setZoom(zoom: number): void {
        if (this.scene) {
            this.scene.setZoom(zoom);
        }
    }
    
    /**
     * Resize the viewer
     */
    public resize(width?: number, height?: number): void {
        if (this.scene && this.viewerContainer) {
            const w = width || this.viewerContainer.clientWidth;
            const h = height || this.viewerContainer.clientHeight;
            this.scene.resize(w, h);
        }
    }
    
    /**
     * Check if viewer is ready
     */
    public isPhaserReady(): boolean {
        return this.scene?.getIsInitialized() || false;
    }

    /**
     * Set interaction callbacks for game-specific functionality
     */
    public setInteractionCallbacks(
        tileCallback?: (q: number, r: number) => boolean,
        unitCallback?: (q: number, r: number) => boolean
    ): void {
        if (this.scene) {
            this.scene.setInteractionCallbacks(tileCallback, unitCallback);
        }
    }

    // =============================================================================
    // ComponentLifecycle Interface Implementation
    // =============================================================================

    /**
     * Phase 1: Initialize DOM and discover child components
     */
    initializeDOM(): ComponentLifecycle[] {
        console.log('WorldViewer: initializeDOM() - Phase 1');
        
        // DOM setup is already done in bindToDOM(), just return no child components
        return [];
    }

    /**
     * Phase 2: Inject dependencies (none needed for WorldViewer)
     */
    injectDependencies(deps: Record<string, any>): void {
        console.log('WorldViewer: injectDependencies() - Phase 2', Object.keys(deps));
        // WorldViewer doesn't need external dependencies
    }

    /**
     * Phase 3: Activate component - Initialize Phaser here
     */
    async activate(): Promise<void> {
        console.log('WorldViewer: activate() - Phase 3 - Initializing Phaser');
        
        // Now initialize PhaserWorldScene in the proper lifecycle phase
        await this.initializePhaserScene();
        
        console.log('WorldViewer: activation complete');
    }

    /**
     * Cleanup phase (called by lifecycle controller if needed)
     */
    deactivate(): void {
        console.log('WorldViewer: deactivate() - cleanup');
        this.destroyComponent();
    }
}
