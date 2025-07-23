import { BaseComponent } from './Component';
import { EventBus } from './EventBus';
import { ComponentLifecycle } from './ComponentLifecycle';
import { TERRAIN_NAMES } from './ColorsAndNames';

/**
 * Terrain information for a specific tile
 */
export interface TerrainInfo {
    name: string;
    tileType: number;
    movementCost: number;
    defenseBonus: number;
    description: string;
    q: number;
    r: number;
    player?: number;
}

/**
 * Rules engine data structures (matching Go structs)
 */
interface TerrainData {
    ID: number;
    Name: string;
    BaseMoveCost: number;
    DefenseBonus: number;
    Type: number;
    Description?: string;
}

interface UnitData {
    ID: number;
    Name: string;
    MovementPoints: number;
    AttackRange: number;
    Health: number;
    Properties: string[];
}

interface MovementMatrix {
    Costs: { [unitID: number]: { [terrainID: number]: number } };
}

/**
 * TerrainStatsPanel displays detailed information about a selected terrain tile
 * 
 * This component shows:
 * - Terrain type and visual representation from rules engine
 * - Movement costs for different unit types from movement matrix
 * - Defense bonuses from terrain data
 * - Coordinate information
 * - Player ownership (if applicable)
 * 
 * The panel remains hidden until terrain is selected, then displays relevant info.
 * Uses the terrain-stats-panel-template from TerrainStatsPanel.html
 * Gets terrain data from rules engine JSON embedded in page by Go backend
 */
export class TerrainStatsPanel extends BaseComponent implements ComponentLifecycle {
    private isUIBound = false;
    private isActivated = false;
    private currentTerrain: TerrainInfo | null = null;
    
    // Rules engine data (loaded from page)
    private terrainData: { [id: number]: TerrainData } = {};
    private unitData: { [id: number]: UnitData } = {};
    private movementMatrix: MovementMatrix | null = null;

    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('terrain-stats-panel', rootElement, eventBus, debugMode);
    }

    // ComponentLifecycle Phase 1: Initialize DOM structure
    public initializeDOM(): ComponentLifecycle[] {
        if (this.isUIBound) {
            this.log('Already bound to DOM, skipping');
            return [];
        }

        try {
            this.log('Binding TerrainStatsPanel to DOM using template');
            this.bindToTemplate();
            this.loadRulesEngineData();
            this.isUIBound = true;
            this.log('TerrainStatsPanel bound to DOM successfully');
            
            // This is a leaf component - no children
            return [];
            
        } catch (error) {
            this.handleError('Failed to bind TerrainStatsPanel to DOM', error);
            throw error;
        }
    }

    // Phase 2: No external dependencies needed
    public injectDependencies(deps: Record<string, any>): void {
        this.log('TerrainStatsPanel: No dependencies required');
    }

    // Phase 3: Activate component
    public activate(): void {
        if (this.isActivated) {
            this.log('Already activated, skipping');
            return;
        }

        this.log('Activating TerrainStatsPanel');
        this.isActivated = true;
        this.log('TerrainStatsPanel activated successfully');
    }

    // Phase 4: Deactivate component
    public deactivate(): void {
        this.log('Deactivating TerrainStatsPanel');
        this.currentTerrain = null;
        this.isActivated = false;
        this.log('TerrainStatsPanel deactivated');
    }

    /**
     * Bind component to the template from TerrainStatsPanel.html
     */
    private bindToTemplate(): void {
        const template = document.getElementById('terrain-stats-panel-template') as HTMLTemplateElement;
        if (!template) {
            throw new Error('terrain-stats-panel-template not found. Make sure TerrainStatsPanel.html is included.');
        }

        // Clone the template content and append to root element
        const templateContent = template.content.cloneNode(true) as DocumentFragment;
        this.rootElement.appendChild(templateContent);

        this.log('Template bound successfully');
    }

    /**
     * Load rules engine data from embedded JSON in page
     */
    private loadRulesEngineData(): void {
        try {
            // Load terrain data
            const terrainElement = document.getElementById('terrain-data-json');
            if (terrainElement && terrainElement.textContent) {
                this.terrainData = JSON.parse(terrainElement.textContent);
                this.log('Loaded terrain data:', { count: Object.keys(this.terrainData).length });
            }

            // Load unit data
            const unitElement = document.getElementById('unit-data-json');
            if (unitElement && unitElement.textContent) {
                this.unitData = JSON.parse(unitElement.textContent);
                this.log('Loaded unit data:', { count: Object.keys(this.unitData).length });
            }

            // Load movement matrix
            const movementElement = document.getElementById('movement-matrix-json');
            if (movementElement && movementElement.textContent) {
                this.movementMatrix = JSON.parse(movementElement.textContent);
                this.log('Loaded movement matrix with', { unitTypes: Object.keys(this.movementMatrix?.Costs || {}).length });
            }

        } catch (error) {
            this.log('Error loading rules engine data:', error);
            // Continue with empty data - component will still work with fallbacks
        }
    }

    /**
     * Update the panel with information about a selected terrain tile
     */
    public updateTerrainInfo(terrainInfo: TerrainInfo): void {
        if (!this.isActivated) {
            this.log('Component not activated, cannot update terrain info');
            return;
        }

        this.currentTerrain = terrainInfo;
        this.log('Updating terrain info for tile:', terrainInfo);

        // Hide no-selection state and show terrain details
        const noSelectionDiv = this.findElement('#no-terrain-selected');
        const terrainDetailsDiv = this.findElement('#terrain-details');
        
        if (noSelectionDiv) noSelectionDiv.classList.add('hidden');
        if (terrainDetailsDiv) terrainDetailsDiv.classList.remove('hidden');

        // Get terrain data from rules engine
        const rulesTerrainData = this.terrainData[terrainInfo.tileType];
        
        // Update terrain header information
        this.updateTerrainHeader(terrainInfo, rulesTerrainData);
        
        // Update movement cost (use rules engine data if available)
        const movementCost = rulesTerrainData?.BaseMoveCost || terrainInfo.movementCost;
        this.updateMovementCost(movementCost);
        
        // Update defense bonus (use rules engine data if available)
        const defenseBonus = rulesTerrainData?.DefenseBonus || terrainInfo.defenseBonus;
        this.updateDefenseBonus(defenseBonus);
        
        // Update player ownership if applicable
        this.updatePlayerOwnership(terrainInfo.player);
        
        // Update terrain properties using rules engine data
        this.updateTerrainProperties(terrainInfo, rulesTerrainData);
    }

    /**
     * Clear terrain selection and show empty state
     */
    public clearTerrainInfo(): void {
        if (!this.isActivated) {
            return;
        }

        this.currentTerrain = null;
        this.log('Clearing terrain info');

        // Show no-selection state and hide terrain details
        const noSelectionDiv = this.findElement('#no-terrain-selected');
        const terrainDetailsDiv = this.findElement('#terrain-details');
        
        if (noSelectionDiv) noSelectionDiv.classList.remove('hidden');
        if (terrainDetailsDiv) terrainDetailsDiv.classList.add('hidden');
    }

    /**
     * Update the terrain header (icon, name, coordinates, description)
     */
    private updateTerrainHeader(terrainInfo: TerrainInfo, rulesData?: TerrainData): void {
        const iconElement = this.findElement('#terrain-icon');
        const nameElement = this.findElement('#terrain-name');
        const coordsElement = this.findElement('#terrain-coordinates');
        const descElement = this.findElement('#terrain-description');

        if (iconElement) {
            const terrainData = TERRAIN_NAMES[terrainInfo.tileType] || { icon: '🎨' };
            iconElement.textContent = terrainData.icon;
        }

        if (nameElement) {
            // Use rules engine name if available, fallback to terrainInfo name
            const displayName = rulesData?.Name || terrainInfo.name;
            nameElement.textContent = displayName;
        }

        if (coordsElement) {
            coordsElement.textContent = `(${terrainInfo.q}, ${terrainInfo.r})`;
        }

        if (descElement) {
            // Use rules engine description if available, fallback to terrainInfo description
            const description = rulesData?.Description || terrainInfo.description;
            descElement.textContent = description;
        }
    }

    /**
     * Update the movement cost display
     */
    private updateMovementCost(cost: number): void {
        const costElement = this.findElement('#movement-cost');
        if (costElement) {
            costElement.textContent = cost.toFixed(1);
        }
    }

    /**
     * Update the defense bonus display
     */
    private updateDefenseBonus(bonus: number): void {
        const bonusElement = this.findElement('#defense-bonus');
        if (bonusElement) {
            const sign = bonus >= 0 ? '+' : '';
            bonusElement.textContent = `${sign}${(bonus * 100).toFixed(0)}%`;
        }
    }

    /**
     * Update player ownership display
     */
    private updatePlayerOwnership(player?: number): void {
        const ownershipDiv = this.findElement('#player-ownership');
        const playerElement = this.findElement('#owner-player');

        if (player !== undefined && player > 0) {
            if (ownershipDiv) ownershipDiv.classList.remove('hidden');
            if (playerElement) playerElement.textContent = `Player ${player}`;
        } else {
            if (ownershipDiv) ownershipDiv.classList.add('hidden');
        }
    }

    /**
     * Update terrain properties list using rules engine data
     */
    private updateTerrainProperties(terrainInfo: TerrainInfo, rulesData?: TerrainData): void {
        const propertiesList = this.findElement('#properties-list');
        if (!propertiesList) return;

        const properties: Array<{name: string, value: string}> = [];

        // Add basic properties
        properties.push({
            name: 'Type ID',
            value: terrainInfo.tileType.toString()
        });

        properties.push({
            name: 'Hex Coordinate',
            value: `Q:${terrainInfo.q}, R:${terrainInfo.r}`
        });

        // Add rules engine data if available
        if (rulesData) {
            properties.push({
                name: 'Base Move Cost',
                value: rulesData.BaseMoveCost.toFixed(1)
            });

            if (rulesData.DefenseBonus !== 0) {
                const sign = rulesData.DefenseBonus >= 0 ? '+' : '';
                properties.push({
                    name: 'Defense Bonus',
                    value: `${sign}${(rulesData.DefenseBonus * 100).toFixed(0)}%`
                });
            }

            properties.push({
                name: 'Terrain Type',
                value: rulesData.Type.toString()
            });
        }

        // Add movement costs for different unit types if available
        if (this.movementMatrix && this.movementMatrix.Costs) {
            const unitMovements: string[] = [];
            
            // Show movement costs for first few unit types as examples
            Object.entries(this.movementMatrix.Costs).slice(0, 3).forEach(([unitId, terrainCosts]) => {
                const cost = terrainCosts[terrainInfo.tileType];
                if (cost !== undefined) {
                    const unitName = this.unitData[parseInt(unitId)]?.Name || `Unit ${unitId}`;
                    unitMovements.push(`${unitName}: ${cost.toFixed(1)}`);
                }
            });
            
            if (unitMovements.length > 0) {
                properties.push({
                    name: 'Unit Movement Costs',
                    value: unitMovements.join(', ')
                });
            }
        }

        // Generate HTML
        let propertiesHTML = '';
        properties.forEach(property => {
            propertiesHTML += `
                <div class="text-sm text-gray-600 dark:text-gray-300">
                    <span class="font-medium">${property.name}:</span> ${property.value}
                </div>
            `;
        });

        propertiesList.innerHTML = propertiesHTML || 
            '<div class="text-sm text-gray-500 dark:text-gray-400 italic">No properties available</div>';
    }

    /**
     * Get current terrain info (for external access)
     */
    public getCurrentTerrain(): TerrainInfo | null {
        return this.currentTerrain;
    }

    /**
     * Check if terrain is currently selected
     */
    public hasTerrainSelected(): boolean {
        return this.currentTerrain !== null;
    }

    /**
     * Get terrain data from rules engine (for external access)
     */
    public getTerrainData(tileType: number): TerrainData | null {
        return this.terrainData[tileType] || null;
    }

    // BaseComponent lifecycle compatibility
    protected initializeComponent(): void {
        // Handled by the new lifecycle system
    }

    protected bindToDOM(): void {
        // Handled by the new lifecycle system
    }

    protected destroyComponent(): void {
        this.deactivate();
    }
}