/**
 * Default Theme Provider
 * Extends BaseTheme with PNG-specific asset loading
 * Uses pre-colored PNG assets that don't need post-processing
 */

import { BaseTheme, ThemeInfo, ThemeMapping, PlayerColor } from './BaseTheme';
import mappingData from './default/mapping.json';

// Type for the extended mapping with playerColors
interface DefaultMappingData extends ThemeMapping {
    themeInfo: {
        name: string;
        version: string;
        base_path: string;
        asset_type: string;
        needs_post_processing: boolean;
    };
    playerColors: { [key: string]: PlayerColor };
}

const mapping = mappingData as DefaultMappingData;

/**
 * Default Theme Implementation
 * Uses PNG assets with pre-colored player variants
 */
export default class DefaultTheme extends BaseTheme {
    protected basePath = '/static/assets/themes/default';
    protected themeName = 'Default (PNG)';
    protected themeVersion = '1.0.0';

    constructor() {
        super(mapping);
    }

    /**
     * Override getThemeInfo for PNG-specific settings
     */
    getThemeInfo(): ThemeInfo {
        return {
            name: this.themeName,
            version: this.themeVersion,
            basePath: this.basePath,
            supportsTinting: false,
            needsPostProcessing: false,
            assetType: 'png',
            playerColors: this.playerColors as any,
        };
    }

    /**
     * Override getUnitPath - PNG assets don't use template paths
     */
    getUnitPath(unitId: number): string | undefined {
        // PNG assets use getUnitAssetPath instead
        return undefined;
    }

    /**
     * Override getTilePath - PNG assets don't use template paths
     */
    getTilePath(terrainId: number): string | undefined {
        // PNG assets use getTileAssetPath instead
        return undefined;
    }

    /**
     * Returns the direct path to a pre-colored unit asset
     */
    getUnitAssetPath(unitId: number, playerId: number): string | undefined {
        const unit = this.unitMapping[unitId.toString()];
        if (!unit) {
            return undefined;
        }
        // PNG assets are stored as: {basePath}/{image}/{playerId}.png
        // where image is a directory like "Units/15"
        return `${this.basePath}/${unit.image}/${playerId}.png`;
    }

    /**
     * Returns the direct path to a pre-colored terrain asset
     */
    getTileAssetPath(terrainId: number, playerId: number): string | undefined {
        const terrain = this.terrainMapping[terrainId.toString()];
        if (!terrain) {
            return undefined;
        }

        // Only city terrains use player colors; all others use player 0 (neutral)
        const effectivePlayer = this.isCityTile(terrainId) ? playerId : 0;
        // PNG assets are stored as: {basePath}/{image}/{playerId}.png
        // where image is a directory like "Tiles/5"
        return `${this.basePath}/${terrain.image}/${effectivePlayer}.png`;
    }

    /**
     * Override loadUnit - PNG theme uses pre-colored assets
     */
    async loadUnit(unitId: number, playerId: number): Promise<string> {
        throw new Error('Default theme uses pre-colored PNG assets. Use getUnitAssetPath() instead.');
    }

    /**
     * Override loadTile - PNG theme uses pre-colored assets
     */
    async loadTile(terrainId: number, playerId?: number): Promise<string> {
        throw new Error('Default theme uses pre-colored PNG assets. Use getTileAssetPath() instead.');
    }

    /**
     * Override setUnitImage for PNG assets
     */
    async setUnitImage(unitId: number, playerId: number, targetElement: HTMLElement): Promise<void> {
        const path = this.getUnitAssetPath(unitId, playerId);
        if (!path) {
            targetElement.innerHTML = '⚔️';
            return;
        }

        targetElement.innerHTML = '';
        const img = document.createElement('img');
        img.src = path;
        img.alt = this.getUnitName(unitId) || `Unit ${unitId}`;
        img.className = 'w-full h-full object-contain';
        img.style.imageRendering = 'pixelated';

        img.onerror = () => {
            targetElement.innerHTML = '⚔️';
        };

        targetElement.appendChild(img);
    }

    /**
     * Override setTileImage for PNG assets
     */
    async setTileImage(tileId: number, playerId: number, targetElement: HTMLElement): Promise<void> {
        const path = this.getTileAssetPath(tileId, playerId);
        if (!path) {
            targetElement.innerHTML = '🏞️';
            return;
        }

        targetElement.innerHTML = '';
        const img = document.createElement('img');
        img.src = path;
        img.alt = this.getTerrainName(tileId) || `Terrain ${tileId}`;
        img.className = 'w-full h-full object-contain';
        img.style.imageRendering = 'pixelated';

        img.onerror = () => {
            targetElement.innerHTML = '🏞️';
        };

        targetElement.appendChild(img);
    }
}
