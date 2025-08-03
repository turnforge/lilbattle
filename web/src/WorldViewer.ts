import { BaseComponent } from '../lib/Component';
import { EventBus } from '../lib/EventBus';
import { WorldEventTypes, WorldDataLoadedPayload } from './events';
import { PhaserWorldScene } from './phaser/PhaserWorldScene';
import { PhaserGameScene } from './phaser/PhaserGameScene';
import { Unit, Tile, World } from './World';
import { LCMComponent } from '../lib/LCMComponent';

/**
 * WorldViewer Component - Manages Phaser-based world visualization
 * Responsible for:
 * - Phaser initialization and lifecycle management
 * - World data rendering (tiles and units)
 * - Camera controls and viewport management
 * - Theme and display options
 * 
 * Layout and styling are handled by parent container and CSS classes.
 * 
 * @template TScene - The type of Phaser scene to use (defaults to PhaserWorldScene)
 */
export class WorldViewer<TScene extends PhaserWorldScene = PhaserWorldScene> extends BaseComponent implements LCMComponent {
    protected scene: TScene | null = null;
    private world: World | null = null; // Store canonical World object
    private viewerContainer: HTMLElement | null;
    
    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('world-viewer', rootElement, eventBus, debugMode);
    }

    /**
     * Factory method for creating the scene - can be overridden by subclasses
     */
    protected createScene(): TScene {
        return new PhaserWorldScene() as TScene;
    }
    
    protected destroyComponent(): void {
        this.log('Destroying WorldViewer component');
        
        // Clean up Phaser scene (it manages its own game instance)
        if (this.scene) {
            this.scene.destroy();
            this.scene = null;
        }
        
        this.world = null;
        this.viewerContainer = null;
    }
    
    
    // validateDOM method removed - not needed in pure LCMComponent approach
    
    /**
     * Initialize the appropriate Phaser scene (PhaserWorldScene or PhaserGameScene)
     */
    protected async initializePhaserScene(): Promise<void> {
        // Guard against multiple initialization
        if (this.scene) {
            console.warn('WorldViewer: Phaser scene already initialized, skipping');
            return;
        }
        
        if (!this.viewerContainer) {
            throw new Error('Viewer container not available');
        }
        
        // Create scene using factory method
        this.log('Creating Phaser scene using factory method');
        this.scene = this.createScene();
        
        // Initialize it with the container
        await this.scene.initialize(this.viewerContainer.id);
        
        this.log(`Phaser scene initialized successfully`);
        
        // Emit ready event
        this.emit(WorldEventTypes.WORLD_VIEWER_READY, {
            componentId: this.componentId,
            success: true
        }, this, this);
        
        // Load world data if we have it
        if (this.world) {
            await this.loadWorldIntoScene();
        }
    }
    
    /**
     * Load World object into the PhaserWorldScene
     */
    private async loadWorldIntoScene(): Promise<void> {
        if (!this.scene || !this.scene.getIsInitialized()) {
            this.log('Phaser scene not ready, deferring world load');
            return;
        }
        
        if (!this.world) {
            throw new Error('No World object available to load');
        }
        
        this.log('Loading World object into Phaser scene');
        
        // Load the canonical World object directly into Phaser
        await this.scene.loadWorldData(this.world);
        
        // Emit stats for other components
        const allTiles = this.world.getAllTiles();
        const allUnits = this.world.getAllUnits();
        const bounds = this.world.getBounds();
        
        const worldStats = {
            worldId: this.world.id || 'unknown',
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
        
        this.emit(WorldEventTypes.WORLD_DATA_LOADED, worldStats, this, this);
    }
    
    /**
     * Public API for loading canonical World object
     */
    public async loadWorld(world: World): Promise<void> {
        this.log('Loading canonical World object');
        
        // Store the canonical World object
        this.world = world;
        
        // Load into Phaser if ready
        if (this.scene && this.scene.getIsInitialized()) {
            await this.loadWorldIntoScene();
        }
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
        } else {
            console.error('[WorldViewer] No scene available to set callbacks on');
        }
    }

    // =============================================================================
    // LCMComponent Interface Implementation
    // =============================================================================

    /**
     * Phase 1: Initialize DOM and discover child components
     */
    performLocalInit(): LCMComponent[] {
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
        
        return [];
    }

    /**
     * Phase 2: Inject dependencies (none needed for WorldViewer)
     */
    setupDependencies(): void {
        // WorldViewer doesn't need external dependencies
    }

    /**
     * Phase 3: Activate component - Initialize Phaser here
     */
    async activate(): Promise<void> {
        // Check if already initialized
        if (this.scene) {
            throw new Error('WorldViewer: Already activated and scene exists, skipping');
            return;
        }
        
        // Subscribe to world data events about any world
        this.addSubscription(WorldEventTypes.WORLD_DATA_LOADED, null);
        
        // Now initialize PhaserWorldScene in the proper lifecycle phase
        await this.initializePhaserScene();
    }

    /**
     * Handle incoming events from the EventBus
     */
    public handleBusEvent(eventType: string, data: any, target: any, emitter: any): void {
        switch(eventType) {
            case WorldEventTypes.WORLD_DATA_LOADED:
                // this.handleWorldDataLoaded(data);
                break;
                
            default:
                // Call parent implementation for unhandled events
                super.handleBusEvent(eventType, data, target, emitter);
        }
    }

    /**
     * Cleanup phase (called by lifecycle controller if needed)
     */
    deactivate(): void {
        // Remove event subscriptions
        this.removeSubscription(WorldEventTypes.WORLD_DATA_LOADED, null);
        
        this.destroyComponent();
    }
}
