import { BaseComponent, DOMValidation } from './Component';
import { EventBus, EventPayload, EventTypes, MapDataLoadedPayload } from './EventBus';
import { PhaserViewer } from './PhaserViewer';
import { Map } from './Map';

/**
 * MapViewer Component - Manages Phaser-based map visualization
 * Responsible for:
 * - Phaser initialization and lifecycle management
 * - Map data rendering (tiles and units)
 * - Camera controls and viewport management
 * - Theme and display options
 * 
 * Layout and styling are handled by parent container and CSS classes.
 */
export class MapViewer extends BaseComponent {
    private phaserViewer: PhaserViewer | null = null;
    private loadedMapData: MapDataLoadedPayload | null = null;
    private viewerContainer: HTMLElement | null = null;
    
    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('map-viewer', rootElement, eventBus, debugMode);
    }
    
    protected initializeComponent(): void {
        this.log('Initializing MapViewer component');
        
        // Subscribe to map data events
        this.subscribe<MapDataLoadedPayload>(EventTypes.MAP_DATA_LOADED, (payload) => {
            this.handleMapDataLoaded(payload);
        });
        
        this.log('MapViewer component initialized');
    }
    
    protected bindToDOM(): void {
        try {
            this.log('Binding MapViewer to DOM');
            
            // Find or create the Phaser container within our root element
            this.viewerContainer = this.findElement('#phaser-viewer-container');
            if (!this.viewerContainer) {
                // Create the container if it doesn't exist
                this.viewerContainer = document.createElement('div');
                this.viewerContainer.id = 'phaser-viewer-container';
                this.viewerContainer.className = 'w-full h-full min-h-96';
                this.rootElement.appendChild(this.viewerContainer);
            }
            
            // Initialize Phaser viewer with delay to ensure container is ready
            setTimeout(() => {
                this.initializePhaserViewer();
            }, 100);
            
            this.log('MapViewer bound to DOM');
            
        } catch (error) {
            this.handleError('Failed to bind MapViewer to DOM', error);
        }
    }
    
    protected destroyComponent(): void {
        this.log('Destroying MapViewer component');
        
        // Clean up Phaser viewer
        if (this.phaserViewer) {
            this.phaserViewer.destroy();
            this.phaserViewer = null;
        }
        
        this.loadedMapData = null;
        this.viewerContainer = null;
    }
    
    protected async hydrateExistingDOM(validation: DOMValidation): Promise<boolean> {
        try {
            this.log('Hydrating MapViewer with existing DOM');
            
            // Find existing Phaser container
            this.viewerContainer = this.findElement('#phaser-viewer-container');
            if (!this.viewerContainer) {
                throw new Error('Phaser viewer container not found during hydration');
            }
            
            // Initialize Phaser with delay
            setTimeout(() => {
                this.initializePhaserViewer();
            }, 100);
            
            // Subscribe to events
            this.subscribe<MapDataLoadedPayload>(EventTypes.MAP_DATA_LOADED, (payload) => {
                this.handleMapDataLoaded(payload);
            });
            
            this.log('MapViewer hydrated successfully');
            return true;
            
        } catch (error) {
            this.handleError('Failed to hydrate MapViewer', error);
            return false;
        }
    }
    
    protected async createMissingDOM(validation: DOMValidation): Promise<boolean> {
        try {
            this.log('Creating missing DOM for MapViewer');
            
            // Create Phaser container if missing
            if (validation.missingElements.includes('phaser-viewer-container')) {
                const container = document.createElement('div');
                container.id = 'phaser-viewer-container';
                container.className = 'w-full h-full min-h-96';
                this.rootElement.appendChild(container);
                this.viewerContainer = container;
            }
            
            setTimeout(() => {
                this.initializePhaserViewer();
            }, 100);
            
            // Subscribe to events
            this.subscribe<MapDataLoadedPayload>(EventTypes.MAP_DATA_LOADED, (payload) => {
                this.handleMapDataLoaded(payload);
            });
            
            this.log('MapViewer DOM created successfully');
            return true;
            
        } catch (error) {
            this.handleError('Failed to create MapViewer DOM', error);
            return false;
        }
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
     * Initialize the Phaser viewer
     */
    private initializePhaserViewer(): void {
        try {
            if (!this.viewerContainer) {
                throw new Error('Viewer container not available');
            }
            
            this.log('Initializing Phaser viewer');
            
            // Create new PhaserViewer instance
            this.phaserViewer = new PhaserViewer();
            
            // Set up logging
            this.phaserViewer.onLog((message: string) => {
                this.log(`PhaserViewer: ${message}`);
            });
            
            // Initialize with container ID - Phaser will adapt to whatever size the parent provides
            const success = this.phaserViewer.initialize(this.viewerContainer.id);
            if (!success) {
                throw new Error('Failed to initialize Phaser viewer');
            }
            
            // Emit ready event
            this.emit(EventTypes.MAP_VIEWER_READY, {
                componentId: this.componentId,
                success: true
            });
            
            // Load map data if we have it
            if (this.loadedMapData) {
                this.loadMapIntoViewer(this.loadedMapData);
            }
            
            this.log('Phaser viewer initialized successfully');
            
        } catch (error) {
            this.handleError('Failed to initialize Phaser viewer', error);
            
            // Emit error event
            this.emit(EventTypes.MAP_VIEWER_ERROR, {
                componentId: this.componentId,
                error: error,
            });
        }
    }
    
    /**
     * Handle map data loaded event
     */
    private handleMapDataLoaded(payload: EventPayload<MapDataLoadedPayload>): void {
        this.log(`Received map data for map: ${payload.data.mapId}`);
        this.loadedMapData = payload.data;
        
        // Load into Phaser if viewer is ready
        if (this.phaserViewer && this.phaserViewer.getIsInitialized()) {
            this.loadMapIntoViewer(payload.data);
        }
    }
    
    /**
     * Load map data into the Phaser viewer
     */
    private async loadMapIntoViewer(mapData: MapDataLoadedPayload): Promise<void> {
        if (!this.phaserViewer || !this.phaserViewer.getIsInitialized()) {
            this.log('Phaser viewer not ready, deferring map load');
            return;
        }
        
        try {
            this.log('Loading map data into Phaser viewer');
            
            // Convert map data to Phaser format
            const tilesArray: Array<{ q: number; r: number; terrain: number; color: number }> = [];
            const unitsArray: Array<{ q: number; r: number; unitType: number; playerId: number }> = [];
            
            // Process tiles from bounds
            if (mapData.bounds) {
                for (let q = mapData.bounds.minQ; q <= mapData.bounds.maxQ; q++) {
                    for (let r = mapData.bounds.minR; r <= mapData.bounds.maxR; r++) {
                        // This would need to be coordinated with the map data structure
                        // For now, create placeholder logic
                        tilesArray.push({
                            q: q,
                            r: r,
                            terrain: 1, // Default grass
                            color: 0
                        });
                    }
                }
            }
            
            // Load into Phaser viewer
            await this.phaserViewer.loadMapData(tilesArray, unitsArray);
            
            this.log(`Loaded ${tilesArray.length} tiles and ${unitsArray.length} units into viewer`);
            
        } catch (error) {
            this.handleError('Failed to load map into viewer', error);
        }
    }
    
    /**
     * Public API for loading map data
     */
    public async loadMap(mapData: any): Promise<void> {
        try {
            if (!mapData) {
                throw new Error('No map data provided');
            }
            
            this.log('Loading map data');
            
            // Process map data
            const map = Map.deserialize(mapData);
            const allTiles = map.getAllTiles();
            const allUnits = map.getAllUnits();
            
            // Convert to arrays
            const tilesArray: Array<{ q: number; r: number; terrain: number; color: number }> = [];
            const unitsArray: Array<{ q: number; r: number; unitType: number; playerId: number }> = [];
            
            allTiles.forEach(tile => {
                tilesArray.push({
                    q: tile.q,
                    r: tile.r,
                    terrain: tile.tileType,
                    color: tile.playerId || 0
                });
            });
            
            allUnits.forEach(unit => {
                unitsArray.push({
                    q: unit.q,
                    r: unit.r,
                    unitType: unit.unitType,
                    playerId: unit.playerId
                });
            });
            
            // Calculate bounds and stats
            const bounds = map.getBounds();
            
            // Store map data
            this.loadedMapData = {
                mapId: mapData.id || 'unknown',
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
            if (this.phaserViewer && this.phaserViewer.getIsInitialized()) {
                await this.phaserViewer.loadMapData(tilesArray, unitsArray);
            }
            
            // Emit data loaded event for other components
            this.emit(EventTypes.MAP_DATA_LOADED, this.loadedMapData);
            
            this.log('Map loaded successfully');
            
        } catch (error) {
            this.handleError('Failed to load map', error);
            throw error;
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
        if (this.phaserViewer) {
            this.phaserViewer.setShowGrid(show);
        }
    }
    
    public setShowCoordinates(show: boolean): void {
        if (this.phaserViewer) {
            this.phaserViewer.setShowCoordinates(show);
        }
    }
    
    public setTheme(isDark: boolean): void {
        if (this.phaserViewer) {
            this.phaserViewer.setTheme(isDark);
        }
    }
    
    /**
     * Camera controls
     */
    public getZoom(): number {
        return this.phaserViewer?.getZoom() || 1;
    }
    
    public setZoom(zoom: number): void {
        if (this.phaserViewer) {
            this.phaserViewer.setZoom(zoom);
        }
    }
    
    /**
     * Resize the viewer
     */
    public resize(width?: number, height?: number): void {
        if (this.phaserViewer && this.viewerContainer) {
            const w = width || this.viewerContainer.clientWidth;
            const h = height || this.viewerContainer.clientHeight;
            this.phaserViewer.resize(w, h);
        }
    }
    
    /**
     * Check if viewer is ready
     */
    public isPhaserReady(): boolean {
        return this.phaserViewer?.getIsInitialized() || false;
    }
}
