import { Map } from './Map';

export interface SaveResult {
    success: boolean;
    mapId?: string;
    error?: string;
}

export interface MapSaveData {
    tiles: { [key: string]: any };
    mapUnits: any[];
}

/**
 * MapSaveManager - Handles all map persistence operations
 * 
 * Responsibilities:
 * - Build save data format from Map and PhaserEditorComponent
 * - Handle CREATE (new maps) and UPDATE (existing maps) operations
 * - Manage API communication for map persistence
 * - Format data for backend API requirements
 */
export class MapSaveManager {
    private baseUrl: string = '/api/v1/maps';

    /**
     * Save a map to the backend
     */
    async saveMap(
        map: Map,
        currentMapId: string | null,
        isNewMap: boolean,
        tilesData: Array<{ q: number; r: number; terrain: number; color: number }>,
        onProgress?: (message: string) => void
    ): Promise<SaveResult> {
        if (!map) {
            return { success: false, error: 'No map data to save' };
        }

        try {
            onProgress?.('Building save data...');
            
            const saveData = this.buildSaveData(map, tilesData);
            
            onProgress?.('Sending to server...');
            
            const result = await this.sendToServer(saveData, currentMapId, isNewMap);
            
            return result;
            
        } catch (error) {
            console.error('Save failed:', error);
            return { 
                success: false, 
                error: error instanceof Error ? error.message : 'Unknown save error' 
            };
        }
    }

    /**
     * Build the save data structure from map and tiles data
     */
    private buildSaveData(
        map: Map, 
        tilesData: Array<{ q: number; r: number; terrain: number; color: number }>
    ): MapSaveData {
        // Build tiles data in the correct format for CreateMap API
        const tiles: { [key: string]: any } = {};
        
        tilesData.forEach((tile: any) => {
            const key = `${tile.q},${tile.r}`;
            tiles[key] = {
                q: tile.q,
                r: tile.r,
                tile_type: tile.terrain,
                player: tile.color
            };
        });

        // Build units data in the correct format for CreateMap API
        const mapUnits: any[] = [];
        const allUnits = map.getAllUnits();
        
        allUnits.forEach((unit) => {
            mapUnits.push({
                q: unit.q,
                r: unit.r,
                player: unit.playerId,
                unit_type: unit.unitType
            });
        });

        return { tiles, mapUnits };
    }

    /**
     * Send save data to server
     */
    private async sendToServer(
        saveData: MapSaveData,
        currentMapId: string | null,
        isNewMap: boolean
    ): Promise<SaveResult> {
        // Build the CreateMapRequest structure
        const createMapRequest = {
            map: {
                id: currentMapId || 'new-map',
                name: 'Untitled Map', // TODO: Get from map metadata
                description: '',
                tags: [],
                difficulty: 'medium',
                creator_id: 'editor-user', // TODO: Get actual user ID
                tiles: saveData.tiles,
                map_units: saveData.mapUnits
            }
        };

        const url = isNewMap ? this.baseUrl : `${this.baseUrl}/${currentMapId}`;
        const method = isNewMap ? 'POST' : 'PATCH';

        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(createMapRequest),
        });

        if (response.ok) {
            const result = await response.json();
            const mapId = result.map?.id || result.id;
            
            return {
                success: true,
                mapId: mapId
            };
        } else {
            const errorText = await response.text();
            throw new Error(`Save failed: ${response.status} ${response.statusText} - ${errorText}`);
        }
    }

    /**
     * Export map data as downloadable file
     */
    async exportMap(
        map: Map,
        tilesData: Array<{ q: number; r: number; terrain: number; color: number }>,
        format: 'json' | 'png' = 'json'
    ): Promise<boolean> {
        try {
            if (format === 'json') {
                const saveData = this.buildSaveData(map, tilesData);
                const dataStr = JSON.stringify(saveData, null, 2);
                const dataBlob = new Blob([dataStr], { type: 'application/json' });
                
                const link = document.createElement('a');
                link.href = URL.createObjectURL(dataBlob);
                link.download = `${map.getName() || 'map'}.json`;
                link.click();
                
                return true;
            }
            
            // TODO: Implement PNG export
            return false;
            
        } catch (error) {
            console.error('Export failed:', error);
            return false;
        }
    }
}