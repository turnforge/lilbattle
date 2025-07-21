import { BasePage } from './BasePage';
import { EventBus, EventTypes } from './EventBus';
import { MapViewer } from './MapViewer';
import { Map } from './Map';

/**
 * Start Game Page - Orchestrator for game configuration functionality
 * Responsible for:
 * - Map data loading and preview coordination
 * - Game configuration management
 * - Player configuration handling
 * - Game creation workflow
 * 
 * Does NOT handle:
 * - Direct DOM manipulation (delegated to components)
 * - Phaser management (delegated to MapViewer)
 * - Game logic (delegated to game engine)
 */
class StartGamePage extends BasePage {
    private currentMapId: string | null;
    private isLoadingMap: boolean = false;
    private map: Map | null = null;
    private gameConfig: GameConfiguration = {
        players: [
            { id: 1, color: 'red', type: 'human', team: 1 },
            { id: 2, color: 'blue', type: 'ai', team: 2 }
        ],
        allowedUnits: ['infantry', 'tank', 'helicopter', 'destroyer'],
        turnTimeLimit: 0,
        teamMode: 'ffa'
    };
    
    // Component instances
    private mapViewer: MapViewer | null = null;

    constructor() {
        super();
        this.loadInitialState();
        this.initializeSpecificComponents();
        this.bindSpecificEvents();
    }

    protected initializeSpecificComponents(): void {
        // Initialize components immediately
        this.initializeComponents();
    }
    
    /**
     * Initialize page components using the established component architecture
     */
    private initializeComponents(): void {
        try {
            console.log('Initializing StartGamePage components');
            
            // Subscribe to MapViewer ready event BEFORE creating the component
            console.log('StartGamePage: Subscribing to map-viewer-ready event');
            this.eventBus.subscribe('map-viewer-ready', () => {
                console.log('StartGamePage: MapViewer is ready, loading map data...');
                if (this.currentMapId) {
                  // Give Phaser time to fully initialize webgl context and scene
                  setTimeout(async () => {
                    await this.loadMapData()
                  }, 10)
                }
            }, 'start-game-page');
            
            // Create MapViewer component for preview
            const mapViewerRoot = this.ensureElement('[data-component="map-viewer"]', 'map-viewer-root');
            console.log('StartGamePage: Creating MapViewer with eventBus:', this.eventBus);
            this.mapViewer = new MapViewer(mapViewerRoot, this.eventBus, true);
            
            console.log('StartGamePage components initialized');
            
        } catch (error) {
            console.error('Failed to initialize components:', error);
            this.showToast('Error', 'Failed to initialize page components', 'error');
        }
    }
    
    /**
     * Ensure an element exists, create if missing
     */
    private ensureElement(selector: string, fallbackId: string): HTMLElement {
        let element = document.querySelector(selector) as HTMLElement;
        if (!element) {
            console.warn(`Element not found: ${selector}, creating fallback`);
            element = document.createElement('div');
            element.id = fallbackId;
            element.className = 'w-full h-full';
            const mainContainer = document.querySelector('main') || document.body;
            mainContainer.appendChild(element);
        }
        return element;
    }

    protected bindSpecificEvents(): void {
        // Bind start game button
        const startGameButton = document.querySelector('[data-action="start-game"]');
        if (startGameButton) {
            startGameButton.addEventListener('click', this.startGame.bind(this));
        }

        // Bind player type selectors
        const playerSelects = document.querySelectorAll('[data-player]');
        playerSelects.forEach(select => {
            select.addEventListener('change', this.handlePlayerConfigChange.bind(this));
        });

        // Bind unit restriction checkboxes
        const unitCheckboxes = document.querySelectorAll('[data-unit]');
        unitCheckboxes.forEach(checkbox => {
            checkbox.addEventListener('change', this.handleUnitRestrictionChange.bind(this));
        });

        // Bind turn limit selector
        const turnLimitSelect = document.querySelector('[data-config="turn-limit"]');
        if (turnLimitSelect) {
            turnLimitSelect.addEventListener('change', this.handleTurnLimitChange.bind(this));
        }

        // Bind team mode radio buttons
        const teamModeRadios = document.querySelectorAll('[name="team-mode"]');
        teamModeRadios.forEach(radio => {
            radio.addEventListener('change', this.handleTeamModeChange.bind(this));
        });
    }

    /** Load document data and set initial UI states */
    private loadInitialState(): void {
        const mapIdInput = document.getElementById("mapIdInput") as HTMLInputElement | null;
        const mapId = mapIdInput?.value.trim() || null;

        if (mapId) {
            this.currentMapId = mapId;
            console.log(`Found Map ID: ${this.currentMapId}. Will load data after Phaser initialization.`);
        } else {
            console.error("Map ID input element not found or has no value. Cannot load map.");
            this.showToast("Error", "Could not load map: Map ID missing.", "error");
        }
    }

    /**
     * Load map data and coordinate between components
     */
    private async loadMapData(): Promise<void> {
        try {
            console.log(`StartGamePage: Loading map data...`);
            
            // Load map data from the hidden JSON element
            const mapData = this.loadMapDataFromElement();
            
            if (mapData) {
                this.map = Map.deserialize(mapData);
                console.log('Map data loaded successfully');
                
                // Use MapViewer component to load the map
                if (this.mapViewer) {
                    await this.mapViewer.loadMap(mapData);
                    this.showToast('Success', 'Map loaded successfully', 'success');
                } else {
                    console.warn('MapViewer component not available');
                }
                
            } else {
                console.error('No map data found');
                this.showToast('Error', 'No map data found', 'error');
            }
            
        } catch (error) {
            console.error('Failed to load map data:', error);
            this.showToast('Error', 'Failed to load map data', 'error');
        }
    }
    
    /**
     * Load map data from the hidden JSON element in the page
     */
    private loadMapDataFromElement(): any {
        try {
            const mapDataElement = document.getElementById('map-data-json');
            console.log(`Map data element found: ${mapDataElement ? 'YES' : 'NO'}`);
            
            if (mapDataElement && mapDataElement.textContent) {
                console.log(`Raw map data content: ${mapDataElement.textContent.substring(0, 200)}...`);
                const mapData = JSON.parse(mapDataElement.textContent);
                
                if (mapData && mapData !== null) {
                    console.log('Map data found in page element');
                    return mapData;
                }
            }
            console.log('No map data found in page element');
            return null;
        } catch (error) {
            console.error('Error parsing map data from page element:', error);
            return null;
        }
    }

    private handlePlayerConfigChange(event: Event): void {
        const select = event.target as HTMLSelectElement;
        const playerId = parseInt(select.dataset.player || '0');
        const playerType = select.value;
        
        const player = this.gameConfig.players.find(p => p.id === playerId);
        if (player) {
            player.type = playerType as PlayerType;
            console.log(`Player ${playerId} type changed to: ${playerType}`);
        }
        
        this.validateGameConfiguration();
    }

    private handleUnitRestrictionChange(event: Event): void {
        const checkbox = event.target as HTMLInputElement;
        const unitType = checkbox.dataset.unit || '';
        
        if (checkbox.checked) {
            if (!this.gameConfig.allowedUnits.includes(unitType)) {
                this.gameConfig.allowedUnits.push(unitType);
            }
        } else {
            this.gameConfig.allowedUnits = this.gameConfig.allowedUnits.filter(unit => unit !== unitType);
        }
        
        console.log('Allowed units updated:', this.gameConfig.allowedUnits);
        this.validateGameConfiguration();
    }

    private handleTurnLimitChange(event: Event): void {
        const select = event.target as HTMLSelectElement;
        this.gameConfig.turnTimeLimit = parseInt(select.value);
        console.log('Turn time limit changed to:', this.gameConfig.turnTimeLimit);
        this.validateGameConfiguration();
    }

    private handleTeamModeChange(event: Event): void {
        const radio = event.target as HTMLInputElement;
        this.gameConfig.teamMode = radio.value as 'ffa' | 'teams';
        console.log('Team mode changed to:', this.gameConfig.teamMode);
        this.validateGameConfiguration();
    }

    private validateGameConfiguration(): boolean {
        const startButton = document.querySelector('[data-action="start-game"]') as HTMLButtonElement;
        let isValid = true;
        let errors: string[] = [];

        // Check if at least one unit type is allowed
        if (this.gameConfig.allowedUnits.length === 0) {
            isValid = false;
            errors.push('At least one unit type must be allowed');
        }

        // Check if we have at least 2 active players
        const activePlayers = this.gameConfig.players.filter(p => p.type !== 'none');
        if (activePlayers.length < 2) {
            isValid = false;
            errors.push('At least 2 players are required');
        }

        if (startButton) {
            startButton.disabled = !isValid;
            startButton.title = errors.length > 0 ? errors.join('; ') : '';
        }

        return isValid;
    }

    private async startGame(): Promise<void> {
        if (!this.validateGameConfiguration()) {
            this.showToast('Error', 'Please fix configuration errors before starting the game', 'error');
            return;
        }

        if (!this.currentMapId) {
            this.showToast('Error', 'No map selected', 'error');
            return;
        }

        try {
            console.log('Starting game with configuration:', this.gameConfig);
            console.log('Will call CreateGame RPC with:', {
                mapId: this.currentMapId,
                players: this.gameConfig.players.filter(p => p.type !== 'none'),
                allowedUnits: this.gameConfig.allowedUnits,
                turnTimeLimit: this.gameConfig.turnTimeLimit,
                teamMode: this.gameConfig.teamMode
            });
            
            // TODO: Call CreateGame RPC endpoint here
            // const response = await this.callCreateGameAPI();
            
            // For now, just show what would be sent and a placeholder message
            this.showToast('Info', 'Game configuration ready! (CreateGame RPC call would happen here)', 'info');
            
        } catch (error) {
            console.error('Failed to start game:', error);
            this.showToast('Error', 'Failed to start game', 'error');
        }
    }

    // Placeholder for future CreateGame API call
    private async callCreateGameAPI(): Promise<any> {
        // This will eventually call the CreateGame RPC endpoint
        const gameRequest = {
            mapId: this.currentMapId,
            players: this.gameConfig.players.filter(p => p.type !== 'none').map(p => ({
                playerId: p.id,
                playerType: p.type,
                color: p.color,
                teamId: p.team
            })),
            gameSettings: {
                allowedUnits: this.gameConfig.allowedUnits,
                turnTimeLimit: this.gameConfig.turnTimeLimit,
                teamMode: this.gameConfig.teamMode
            }
        };

        console.log('CreateGame RPC request would be:', gameRequest);
        
        // TODO: Replace with actual gRPC call
        // return await grpcClient.createGame(gameRequest);
        
        return { gameId: 'placeholder-game-id' };
    }

    public destroy(): void {
        // Clean up components
        if (this.mapViewer) {
            this.mapViewer.destroy();
            this.mapViewer = null;
        }
        
        // Clean up map data
        this.map = null;
        this.currentMapId = null;
    }
}

// Type definitions for game configuration
interface GameConfiguration {
    players: Player[];
    allowedUnits: string[];
    turnTimeLimit: number; // seconds, 0 = no limit
    teamMode: 'ffa' | 'teams';
}

interface Player {
    id: number;
    color: string;
    type: PlayerType;
    team: number;
}

type PlayerType = 'human' | 'ai' | 'open' | 'none';

document.addEventListener('DOMContentLoaded', () => {
    const startGamePage = new StartGamePage();
});